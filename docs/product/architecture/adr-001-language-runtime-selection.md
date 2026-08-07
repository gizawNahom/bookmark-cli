# ADR-001: Language/Runtime Selection

## Status
**Accepted** — confirmed final by user (2026-08-07). Language/paradigm was explicitly open
going into this DESIGN session; Go is now the confirmed choice for `bookmark-cli`.

## Context
`bookmark-cli` is a greenfield CLI tool for terminal-first infra/platform/SRE engineers. Hard
quality-attribute requirements: <100ms perceived save confirmation, single-command install with
no runtime dependency expected by the persona (they already use kubectl/docker/terraform-style
static binaries), and a solo/small-team maintenance model (time-to-market and learning-curve
matter more than maximal performance ceiling or maximal compile-time safety).

## Decision
**Go.**

## Alternatives Considered

**Rust** — comparable startup latency and single-binary distribution story to Go. Rejected: its
stronger compile-time safety guarantees (no unsafe memory, ownership/borrow checking) address
risk classes (memory safety, data races in unsafe code) that this design doesn't need to defend
against — there's no unsafe FFI, no exotic concurrency primitive, and SQLite access is mediated
through a driver either way. The steeper learning curve is a real solo-maintainer velocity risk
with no offsetting requirement here, violating "simplest solution first, team capability match."

**Python** — fastest to write, but interpreter/import startup (30-80ms) consumes a large
fraction of the <100ms save-confirmation budget before any actual save work happens. `pip
install` also requires a Python runtime present on the user's machine, reintroducing the kind of
recipient-side friction this exact persona already rejected for `bm share` (Aisha, DISCOVER #8:
"if it's not one-command for my teammates I already lost").

## Consequences

**Positive**: single static binary distribution (Homebrew tap, GitHub release, `go install`) —
matches the persona's existing tool expectations; ~1-5ms cold start leaves ample budget under
the 100ms save-confirmation KPI; mature, pure-Go SQLite driver available (no cgo cross-compile
pain); strong stdlib support for CLI concerns (flag parsing via Cobra, file locking).

**Negative**: Go's type system is weaker than Rust's for enforcing the read/write port split at
compile time (interfaces are structurally, not nominally, typed — a type can accidentally
satisfy `BookmarkWriter` without meaning to). Mitigated by the AST-based structural
enforcement check in ADR-007, not by the language alone.
