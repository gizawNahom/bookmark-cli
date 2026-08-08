// Package ports declares bookmark-cli's driven-port interfaces (hexagonal boundary). Every
// driven adapter implements one of these plus Prober per ADR-007 ("wire then probe then use").
package ports

import "bookmark-cli/internal/core"

// Prober is implemented by every driven adapter (ADR-007 Earned Trust). Probe() is run at
// composition-root startup; a failure aborts the operation with health.startup.refused rather
// than attempting a write that might silently fail.
type Prober interface {
	Probe() error
}

// BookmarkReader is read-only by construction -- no write methods are present on this interface,
// per ADR-006's read/write port split. bm find and bm share are structurally incapable of
// mutating the store through this port.
type BookmarkReader interface {
	FindByID(id string) (core.Record, bool, error)
	Search(query string) ([]core.Record, error)
	All() ([]core.Record, error)
}

// BookmarkWriter exposes exactly one mutating method: Execute(plan). There is no ad hoc mutation
// method -- every write must originate from a SavePlan produced by the pure core.PlanSave.
type BookmarkWriter interface {
	Execute(plan core.SavePlan) (core.Record, error)
}

// BackupService snapshots the primary store into the backup directory (ADR-005). Bounded to
// ${data_dir}/backups/ only -- never touches the primary DB file except to read it for copying.
type BackupService interface {
	Prober
	Snapshot(dbPath string) error
}

// UsageEvent is the local opt-in telemetry record (kpi-contracts.yaml). Payloads never include
// URL or tag content -- event name, timestamp, and small enumerated outcome fields only.
type UsageEvent struct {
	Event       string // "bm.save" | "bm.find" | "bm.share"
	Outcome     string // for bm.save: "new" | "duplicate" | "tag_update"
	ResultCount int    // for bm.find: number of ranked matches returned
}

// UsageLogger is the local opt-in telemetry port (DEVOPS wave, kpi-contracts.yaml). Command
// handlers call Record() unconditionally; the composition root selects FileUsageLogAdapter or
// NoOpUsageLogAdapter once, at startup, based on the telemetry_enabled config flag.
type UsageLogger interface {
	Prober
	Record(event UsageEvent) error
}
