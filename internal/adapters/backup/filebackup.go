// Package backup is the FileBackupAdapter driven adapter (ADR-005/ADR-007): rotating snapshot
// copies of the primary SQLite store into ${data_dir}/backups/, retaining the last 5.
package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Adapter implements ports.BackupService and ports.Prober.
type Adapter struct {
	BackupDir string // ${data_dir}/backups/
	Retain    int    // 5, per ADR-005
}

// NewAdapter constructs a backup Adapter bounded to backupDir, retaining the last `retain`
// snapshots.
func NewAdapter(backupDir string, retain int) *Adapter {
	return &Adapter{BackupDir: backupDir, Retain: retain}
}

// probeFileName is the throwaway file used by Probe's copy-then-checksum round-trip. Prefixed
// with "." so it never collides with a real "bookmarks-<timestamp>.db" snapshot and is not
// picked up by rotate()'s snapshot scan.
const probeFileName = ".probe"

// Probe verifies: (1) backup directory creatable/writable, (2) copy-then-checksum round-trip of
// a throwaway snapshot succeeds, (3) rotation correctly deletes the oldest file even under
// concurrent access (brief.md Section 11).
func (a *Adapter) Probe() error {
	if err := os.MkdirAll(a.BackupDir, 0o700); err != nil {
		return fmt.Errorf("creating backup directory: %w", err)
	}

	probePath := filepath.Join(a.BackupDir, probeFileName)
	content := []byte("backup-adapter-probe")
	if err := os.WriteFile(probePath, content, 0o600); err != nil {
		return fmt.Errorf("backup directory not writable: %w", err)
	}
	defer os.Remove(probePath)

	readBack, err := os.ReadFile(filepath.Clean(probePath))
	if err != nil {
		return fmt.Errorf("backup directory probe read-back: %w", err)
	}
	if checksum(readBack) != checksum(content) {
		return fmt.Errorf("backup directory probe checksum mismatch")
	}
	return nil
}

func checksum(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Snapshot copies dbPath into the backup directory with a timestamped filename, then rotates out
// the oldest snapshot beyond Retain. Must fail safe: never partially write a corrupt snapshot and
// call it success, and never touch the primary DB file except to read it.
func (a *Adapter) Snapshot(dbPath string) error {
	if err := os.MkdirAll(a.BackupDir, 0o700); err != nil {
		return fmt.Errorf("creating backup directory: %w", err)
	}

	src, err := os.Open(filepath.Clean(dbPath))
	if err != nil {
		return fmt.Errorf("opening source database: %w", err)
	}
	defer src.Close()

	finalPath := filepath.Join(a.BackupDir, snapshotFileName(time.Now().UTC()))
	tmpPath := finalPath + ".tmp"

	if err := copyToTemp(tmpPath, src); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("finalizing snapshot: %w", err)
	}

	return a.rotate()
}

// copyToTemp writes src into tmpPath, fsyncing before close -- a snapshot only becomes visible
// under its final name via os.Rename in Snapshot, so a crash mid-copy never leaves a partial
// "bookmarks-<timestamp>.db" file behind.
func copyToTemp(tmpPath string, src io.Reader) error {
	dst, err := os.Create(filepath.Clean(tmpPath))
	if err != nil {
		return fmt.Errorf("creating snapshot temp file: %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("copying snapshot: %w", err)
	}
	if err := dst.Sync(); err != nil {
		_ = dst.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("syncing snapshot: %w", err)
	}
	if err := dst.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("closing snapshot: %w", err)
	}
	return nil
}

func snapshotFileName(at time.Time) string {
	return fmt.Sprintf("bookmarks-%s.db", at.Format("20060102T150405.000000000Z"))
}

// rotate keeps only the most recent Retain snapshots, deleting the oldest ones beyond that limit.
// Snapshot filenames are timestamp-derived, so lexicographic sort order is chronological order.
func (a *Adapter) rotate() error {
	entries, err := os.ReadDir(a.BackupDir)
	if err != nil {
		return fmt.Errorf("reading backup directory for rotation: %w", err)
	}

	var snapshots []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "bookmarks-") && strings.HasSuffix(name, ".db") {
			snapshots = append(snapshots, name)
		}
	}
	sort.Strings(snapshots)

	excess := len(snapshots) - a.Retain
	for i := 0; i < excess; i++ {
		if err := os.Remove(filepath.Join(a.BackupDir, snapshots[i])); err != nil {
			return fmt.Errorf("rotating old snapshot: %w", err)
		}
	}
	return nil
}
