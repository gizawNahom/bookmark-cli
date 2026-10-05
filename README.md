# bm

Save, find, and share technical reference links without leaving your terminal.

> **Practice project.** `bm` was built to practise the nWave methodology end-to-end: discovery,
> requirements, design, DevOps, acceptance-test design, and outside-in TDD delivery. The CLI works,
> but the main artefact is the design record in [`docs/`](docs/), including the decision to build
> despite a failed viability gate.

`bm` is a local-first CLI for people who live in a shell: capture a link mid-incident with one
command, get it back by keyword or tag seconds later, and hand it to a teammate in a form they can
use without installing anything. No account, no server, no network calls — a single static binary
and a SQLite file under your home directory.

## Install

```sh
go install github.com/gizawNahom/bookmark-cli/cmd/bm@latest
```

Or build from a clone:

```sh
go build -o bm ./cmd/bm
```

Tagged releases publish signed binaries, checksums, and an SBOM via GoReleaser.

## Usage

### Save a link

```sh
bm save https://kube.io/docs/failover --tag k8s
```

```
Saved [a1b2c3d4] https://kube.io/docs/failover (tag: k8s)
```

The `--tag` flag is optional. Saving without one still works, and adds a one-line nudge:

```
Saved [a1b2c3d4] https://kube.io/docs/failover
Hint: add --tag <name> next time to make this easier to find later
```

Saving a URL you already have does not create a duplicate. If the URL is already stored with the
same tag, `bm` says so and leaves the store alone; if you pass a new tag for an already-saved URL,
it offers to add that tag to the existing entry rather than creating a second record.

A near-miss flag typo is corrected rather than rejected outright — `--tags` gets you
`unknown flag: --tags -- did you mean --tag?`.

### Find a link

```sh
bm find failover
```

```
[a1b2c3d4] https://kube.io/docs/failover (tag: k8s) -- saved 3 days ago
```

Multiple words are treated as one query, so `bm find kube failover` is fine unquoted. Matching is
typo-tolerant. When nothing matches but a close tag exists, you get a `did you mean` nudge; when
nothing matches and nothing is close, you get a plain `no matches found`. A genuinely empty store
gets its own distinct message rather than being reported as a failed search.

### Share a link

```sh
bm share a1b2c3d4
```

```
https://kube.io/docs/failover (tag: k8s)
```

The snippet is exactly the stored record — copy-paste ready, with no install required on the
recipient's end. Untagged bookmarks emit the bare URL.

### Stats

`bm stats` is **not implemented yet.** It currently reports that telemetry is disabled and points
at a `bm config` command that does not exist. Usage logging is opt-in and off by default; when
`BM_TELEMETRY_ENABLED=true` is set, the summary path is still an unimplemented scaffold.

## Where your data lives

Everything is local. `bm` stores a SQLite database at `~/.bm/bookmarks.db` and rotating backup
snapshots under `~/.bm/backups` (the five most recent). Set `BM_DATA_DIR` to relocate the whole
store — the acceptance suite uses this to isolate each test run:

```sh
BM_DATA_DIR=/tmp/scratch-bm bm save https://example.com/doc --tag scratch
```

A backup snapshot is taken *after* the save confirmation is printed, so a failing snapshot warns on
stderr but never turns a successful save into a failed command. Errors exit non-zero with an
`Error:` prefix on stderr.

`BM_TELEMETRY_ENABLED=true` opts into a local usage log at `~/.bm/usage.log`. It is never
transmitted anywhere — see the stats caveat above.

## Development

```sh
go build ./...
go test ./...
go run ./tools/probecheck ./...   # driven-adapter probe-contract gate (ADR-007)
go-arch-lint check                # package boundary enforcement
```

The architecture is functional core / imperative shell. Business logic — URL validation, tag
normalization, save planning, match ranking, snippet formatting — lives in `internal/core` as pure
functions with no I/O. A thin imperative shell (`cmd/bm`, plus driven adapters in
`internal/adapters`) isolates every effect. Mutations are described by an immutable `SavePlan`
value and executed separately, so a command can never report a write that did not happen. Both
boundaries are enforced in CI, not by convention.

## Documentation

`docs/` carries the full design record: architecture brief, seven ADRs covering language, storage,
CLI framework, concurrency, backup, effect isolation, and probe enforcement, plus the per-wave
decision log under `docs/feature/bookmark-cli/`.

Two things worth knowing before reading it. **All interview evidence in the discovery documents is
AI-synthesized** — the named subjects are fictional and every quote is fabricated; those files
demonstrate discovery method, not findings about real users. And the DISCOVER wave's final
viability gate is recorded as FAILED, with the decision to build anyway documented rather than
quietly dropped.

## License

MIT — see [LICENSE](LICENSE).
