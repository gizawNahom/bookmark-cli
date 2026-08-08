# ATDD Infrastructure Policy

Per `nw-distill` § Project Infrastructure Policy. One file per project. Apply-if-exists;
write-if-absent; rewrite with `--policy=fresh`. Git history is the audit trail.

Bootstrapped: 2026-08-07, DISTILL wave, bookmark-cli (first feature in this project).

## Driving

| Port | Mechanism | Note |
|---|---|---|
| CLI (`bm save`/`bm find`/`bm share`/`bm stats`) | subprocess — build the real `bm` binary once per test binary run (`go build -o <tmp>/bm ./cmd/bm`), invoke via `os/exec.Command` against a per-test `${data_dir}` under `t.TempDir()` | Matches Architecture of Reference "Driving → real adapter (CLI runner)"; exercises the actual Cobra wiring, not the service function directly (Driving Adapter Verification mandate) |

## Driven internal (real)

| Port | Mechanism | Note |
|---|---|---|
| `BookmarkReader` / `BookmarkWriter` (SQLite, WAL, FTS5) | Real `SQLiteBookmarkStore` adapter against a real file under `t.TempDir()/${data_dir}/bookmarks.db` — `modernc.org/sqlite` pure-Go driver, no Testcontainers needed (embedded, not a network service) | Per ADR-002/ADR-004; acceptance tests read the DB file directly (or via `bm find`/`bm stats` output) to build the state-delta `after` snapshot |
| `BackupService` (rotating file snapshots) | Real `FileBackupAdapter` against `t.TempDir()/${data_dir}/backups/` | Per ADR-005; probed at composition-root startup like the primary store |
| `UsageLogger` (opt-in local event log) | Real `FileUsageLogAdapter` against `t.TempDir()/${data_dir}/usage.log` when `telemetry_enabled=true`; `NoOpUsageLogAdapter` (real, trivial, not a test double) when disabled | Per DEVOPS `kpi-contracts.yaml`; no port is faked, both concrete adapters are the real production types |

## Driven external / non-deterministic (fake)

| Port | Fake | Note |
|---|---|---|
| *(none)* | *(none)* | This project has zero external/non-deterministic driven ports — no clock, network, email, or third-party API dependency exists anywhere in the architecture (`brief.md` Section 0/17, confirmed zero network calls). Recorded explicitly so the empty section reads as a verified conclusion, not an oversight. |

## Layer mapping (informs Mandate 8/9/11 application)

- All acceptance scenarios in this project run at the **subprocess/FS acceptance layer** (real
  CLI binary + real SQLite file + real filesystem backup directory) — there is no in-memory
  double layer in DISTILL scope for bookmark-cli, because every driven port in this project is
  driven-internal (real-by-default) or absent (no driven-external ports exist to fake). Pure-core
  functions (`URLValidator`, `TagNormalizer`, `SavePlanner`, `Matcher`, `SnippetFormatter`) get
  unit-layer property-based tests in DELIVER (owned by the crafter, not DISTILL) — DISTILL's job
  is the acceptance layer only.
- Per Mandate 10, Tier B (state-machine PBT) is **not** added for this feature: no journey in
  this feature chains ≥3 scenarios over a domain-rich input space in a way that a single Tier A
  example per journey doesn't already cover; `bm save`/`bm find`/`bm share` are single-command,
  single-outcome operations, not multi-step stateful workflows. See `feature-delta.md` DISTILL
  section for the explicit Mandate 10 evaluation.
