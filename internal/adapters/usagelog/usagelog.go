// Package usagelog implements the local opt-in telemetry driven adapter (DEVOPS wave,
// docs/product/kpi-contracts.yaml): FileUsageLogAdapter (real, bounded to
// ${data_dir}/usage.log, append-only JSONL) and NoOpUsageLogAdapter (real, trivial -- selected
// at the composition root when telemetry_enabled=false; not a test double).
package usagelog

import "bookmark-cli/internal/ports"

// FileUsageLogAdapter is the "real" telemetry adapter, selected when telemetry_enabled=true.
//
// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer) for the file-writing behavior.
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
	panic("usagelog.FileUsageLogAdapter.Probe not yet implemented -- RED scaffold")
}

// Record appends a privacy-safe usage event (event name + timestamp + small enumerated fields
// only -- never URL or tag content) to LogPath as a JSONL line.
func (a *FileUsageLogAdapter) Record(event ports.UsageEvent) error {
	panic("usagelog.FileUsageLogAdapter.Record not yet implemented -- RED scaffold")
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
