// Package usagelog implements the local opt-in telemetry driven adapter (DEVOPS wave,
// docs/product/kpi-contracts.yaml): FileUsageLogAdapter (real, bounded to
// ${data_dir}/usage.log, append-only JSONL) and NoOpUsageLogAdapter (real, trivial -- selected
// at the composition root when telemetry_enabled=false; not a test double).
package usagelog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gizawNahom/bookmark-cli/internal/ports"
)

// FileUsageLogAdapter is the "real" telemetry adapter, selected when telemetry_enabled=true.
type FileUsageLogAdapter struct {
	LogPath string // ${data_dir}/usage.log
}

// NewFileUsageLogAdapter constructs a FileUsageLogAdapter bounded to logPath.
func NewFileUsageLogAdapter(logPath string) *FileUsageLogAdapter {
	return &FileUsageLogAdapter{LogPath: logPath}
}

// Probe verifies the log directory is writable (ADR-007 pattern, extended to this adapter by the
// project's generic 3-layer enforcement tooling -- no new tooling scope required).
func (a *FileUsageLogAdapter) Probe() error {
	dir := filepath.Dir(a.LogPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating usage log directory: %w", err)
	}

	probePath := filepath.Join(dir, ".usage-log-probe")
	if err := os.WriteFile(probePath, []byte("probe"), 0o644); err != nil {
		return fmt.Errorf("usage log directory not writable: %w", err)
	}
	defer os.Remove(probePath)
	return nil
}

// usageLogEntry is the JSONL record shape written to LogPath. Deliberately excludes any field
// carrying URL or tag content -- event name, timestamp, and small enumerated outcome fields only
// (privacy-by-design, brief.md Section 8 / kpi-contracts.yaml).
type usageLogEntry struct {
	Event       string `json:"event"`
	Timestamp   string `json:"timestamp"`
	Outcome     string `json:"outcome,omitempty"`
	ResultCount int    `json:"result_count,omitempty"`
}

// Record appends a privacy-safe usage event (event name + timestamp + small enumerated fields
// only -- never URL or tag content) to LogPath as a JSONL line.
func (a *FileUsageLogAdapter) Record(event ports.UsageEvent) error {
	dir := filepath.Dir(a.LogPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating usage log directory: %w", err)
	}

	f, err := os.OpenFile(a.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening usage log: %w", err)
	}
	defer f.Close()

	entry := usageLogEntry{
		Event:       event.Event,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Outcome:     event.Outcome,
		ResultCount: event.ResultCount,
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("encoding usage event: %w", err)
	}
	line = append(line, '\n')
	if _, err := f.Write(line); err != nil {
		return fmt.Errorf("writing usage event: %w", err)
	}
	return nil
}

// NoOpUsageLogAdapter is the real (not faked) adapter selected when telemetry_enabled=false.
// Command handlers call UsageLogger.Record() unconditionally either way -- this keeps the
// imperative shell free of scattered "if telemetry_enabled" branches. Trivial enough that it is
// not scaffolded RED: there is no business logic to defer to DELIVER.
type NoOpUsageLogAdapter struct{}

// NewNoOpUsageLogAdapter constructs the no-op adapter.
func NewNoOpUsageLogAdapter() *NoOpUsageLogAdapter { return &NoOpUsageLogAdapter{} }

// Probe always succeeds -- there is no resource to verify.
func (a *NoOpUsageLogAdapter) Probe() error { return nil }

// Record discards the event.
func (a *NoOpUsageLogAdapter) Record(event ports.UsageEvent) error { return nil }
