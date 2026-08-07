# ADR-004: Concurrency / Atomicity Approach

## Status
Accepted

## Context
NFR Guardrail (added after DISCUSS peer review): `bm save`/`bm find`/`bm share` may be invoked
concurrently from multiple terminal panes or scripts against the same local store; no story's
acceptance criteria assume single-process access. DESIGN must specify at least an atomic-write
guarantee.

## Decision
**SQLite WAL (write-ahead log) mode with a configured `busy_timeout`** (e.g. 5000ms) as the
concurrency/atomicity mechanism, layered on top of ADR-002's storage choice.

## Alternatives Considered

**Hand-rolled advisory file locking (`flock`) around a JSON file, with write-to-temp-then-atomic-rename**
— viable and used successfully by many CLI tools. Rejected as the primary mechanism specifically
because it is a *convention* every write path must remember to follow correctly, rather than an
*engine-level contract* that can be probed once and trusted (Core Principle 13: Earned Trust —
"every dependency you don't probe is an act of faith"). SQLite's WAL mode turns this into a
property of the storage engine itself, which is both more reliably correct across future code
changes and more directly probeable (`PRAGMA journal_mode` returns `wal`, or it doesn't).

**No locking, "last write wins" full-file overwrite** — rejected outright: directly violates the
NFR guardrail's atomic-write-at-minimum requirement and risks data loss under the concurrent-pane
usage pattern the guardrail explicitly calls out.

## Consequences

**Positive**: concurrent readers proceed without blocking; a single writer at a time with other
writers queuing behind `busy_timeout` rather than failing or corrupting data; the guarantee is
empirically probeable at startup (ADR-007) rather than assumed.

**Negative**: WAL mode's guarantees depend on the underlying filesystem actually honoring
fsync/locking semantics — known to silently degrade on some overlay/network filesystems (Docker
overlayfs, WSL2 DrvFs). This is *not* treated as a closed risk by this ADR alone — it is exactly
why the probe contract in `brief.md` Section 11 exists: the design assumes the environment may
lie, and specifies the check rather than the trust.
