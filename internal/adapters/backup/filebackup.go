// Package backup is the FileBackupAdapter driven adapter (ADR-005/ADR-007): rotating snapshot
// copies of the primary SQLite store into ${data_dir}/backups/, retaining the last 5.
//
// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer).
package backup

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

// Probe verifies: (1) backup directory creatable/writable, (2) copy-then-checksum round-trip of
// a throwaway snapshot succeeds, (3) rotation correctly deletes the oldest file even under
// concurrent access (brief.md Section 11).
func (a *Adapter) Probe() error {
	panic("backup.Adapter.Probe not yet implemented -- RED scaffold")
}

// Snapshot copies dbPath into the backup directory with a timestamped filename, then rotates out
// the oldest snapshot beyond Retain. Must fail safe: never partially write a corrupt snapshot and
// call it success, and never touch the primary DB file except to read it.
func (a *Adapter) Snapshot(dbPath string) error {
	panic("backup.Adapter.Snapshot not yet implemented -- RED scaffold")
}
