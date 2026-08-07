# ADR-005: Backup / Recovery Strategy

## Status
Accepted

## Context
NFR Guardrail: local storage is a single point of failure; DESIGN must decide and document a
backup/recovery approach — even a manual one — before MVP ships. No story in DISCUSS implements
backup/restore; the "right" mechanism depends on the storage format DESIGN chooses (ADR-002).

## Decision
**Automatic rotating file snapshots** (last 5 retained) of the SQLite DB file into
`${data_dir}/backups/`, triggered after each successful `bm save`, **plus** documenting
`${data_dir}` clearly so users can fold it into their own manual backup/sync tooling as a
belt-and-suspenders fallback.

## Alternatives Considered

**Manual-only** ("back up `${store_location}` yourself", explicitly offered as an acceptable
floor by the guardrail's own wording) — rejected as the *sole* mechanism: it satisfies the letter
of the guardrail but leaves data reliability entirely dependent on user diligence the DISCOVER
evidence already shows this persona doesn't reliably have (Nadia's ad hoc bash-function tooling
had zero backup and she never built the "proper tool" she'd thought about three times). An
automatic floor costs little (single-file copy, small file size) and closes a gap manual-only
leaves open.

**Continuous replication / cloud sync target** — rejected as disproportionate: reintroduces a
network dependency this design explicitly avoids (Section 0's feasibility argument depends on
zero network calls anywhere in the architecture), and no story or KPI asks for multi-device sync.

## Consequences

**Positive**: closes the SPOF gap with near-zero added latency (snapshot copy happens after the
save confirmation is already printed, staying inside the <100ms perceived-save budget) and
near-zero added complexity (file copy + rotation, no new dependency). Manual documentation adds a
second layer without extra engineering cost.

**Negative**: no dedicated `bm restore` command in this wave — recovery is "copy a snapshot file
back into place" by hand. Explicitly flagged as an open question deferred to DELIVER
(`brief.md` Section 19), not silently dropped.
