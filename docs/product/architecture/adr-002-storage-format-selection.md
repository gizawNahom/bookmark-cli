# ADR-002: Storage Format Selection

## Status
Accepted

## Context
Local-first storage is carried forward from DISCOVER/DISCUSS as a working assumption (feasibility
confirmed at the design level — see `brief.md` Section 0). Requirements: atomic-write-at-minimum
concurrency guarantee (NFR Guardrail), typo-tolerant ranked search (US-02 AC), responsive up to
10,000 bookmarks (Section 10 target), and a documented backup/recovery approach (NFR Guardrail).

## Decision
**SQLite**, accessed via a pure-Go driver (`modernc.org/sqlite`, no cgo), in **WAL mode**, with
the **FTS5** extension for full-text/fuzzy search.

## Alternatives Considered

**JSON file (whole-file read/rewrite)** — zero new binary dependency, simplest mental model.
Rejected as primary choice: concurrency safety would require hand-rolled advisory locking
(`flock`) plus a write-to-temp-then-atomic-rename pattern to get atomic-write-at-minimum — this
is exactly the "convention, not contract" surface Core Principle 13 warns against (correctness
depends on every code path remembering to follow the pattern, versus SQLite's WAL mode being an
engine-level, probeable guarantee). No native ranking/fuzzy-match primitive — would need a
hand-built or additional-dependency fuzzy-match library anyway, eroding the "zero new
dependency" advantage.

**BoltDB/bbolt (embedded KV store)** — built-in ACID locking, single dependency, no cgo.
Rejected: no native full-text/fuzzy search primitive; would require hand-building an index on
top of the KV store, which is strictly more custom-built surface area than SQLite+FTS5 for the
same end result, with no offsetting benefit at this scale (10k rows is well within SQLite's
comfortable range).

## Consequences

**Positive**: WAL mode gives an engine-level, empirically probeable atomic-write/multi-reader
guarantee (Section 12's probe contract directly targets this). FTS5 provides ranked, typo-tolerant
search without hand-building ranking logic. Single-file store is trivially copyable for backup
(ADR-005). Pure-Go driver preserves the single-static-binary distribution story from ADR-001 (no
cgo cross-compilation friction).

**Negative**: adds a storage-engine dependency where a JSON file would have none. Mitigated: the
dependency is a mature, extremely widely-deployed embedded engine (SQLite is arguably the most
audited embedded database in existence), not a service dependency, and the pure-Go driver avoids
the cgo/cross-compile complexity that would otherwise be the sharpest edge of this choice.

**Known substrate risk (flagged for Earned Trust probe, see ADR-007/brief.md Section 11)**: WAL
mode's fsync/locking guarantees can silently degrade on some overlay or network filesystems. This
must be probed at startup, not assumed.
