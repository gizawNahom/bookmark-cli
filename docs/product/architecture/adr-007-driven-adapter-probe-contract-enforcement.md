# ADR-007: Driven Adapter Probe Contract + 3-Layer Enforcement Tooling

## Status
Accepted

## Context
Core Principle 13 (Earned Trust) requires every driven adapter (here: `SQLiteBookmarkStore`,
`FileBackupAdapter`) to demonstrate empirically that it can honor its contract in the real
environment, via a `Probe()` method run at composition-root startup, enforced through three
semantically orthogonal layers so a single-layer bypass is still caught.

## Decision
1. Every driven adapter implements `Probe() error`, exercising the fault-injection scenarios
   catalogued in `brief.md` Section 11 (read-only filesystem, WAL-unsupported filesystem, disk
   full, backup-directory unwritable).
2. Composition root ("wire then probe then use"): `bm` runs each adapter's cheap probe subset at
   every startup; a failure emits a structured `health.startup.refused` message and the command
   aborts before attempting any write.
3. Three enforcement layers:
   - **Subtype (compile-time)**: Go's interface satisfaction check — native, no extra tooling.
   - **Structural (AST pre-commit hook)**: custom `go/ast`-based script (or `ruleguard`/custom
     `golangci-lint` rule) asserting every type implementing a driven-adapter interface also
     defines `Probe() error`.
   - **Behavioral (CI)**: `go test ./... -tags=faultinjection` exercises the catalogued substrate
     lies against real adapters (tmpfs mounted read-only, size-capped tmpfs for disk-full,
     permission-denied directories).

## Alternatives Considered

**`import-linter`** — investigated for the structural layer. Rejected: it is a Python-only tool
(this project's language is Go per ADR-001, ruling it out immediately), and independent of
language, its contracts are import-graph-only with no API for method-presence enforcement on
types — it can enforce "package A must not import package B" but cannot express "every type
satisfying interface X must also define method `Probe()`," which is exactly what this layer
needs to check.

**Single-layer enforcement (compile-time only, trusting code review for the rest)** — rejected:
a single layer can be bypassed (e.g. an adapter defines `Probe()` to satisfy the structural
check but the composition root never calls it, or `Probe()` exists but doesn't actually exercise
a real fault scenario). The mandate's own reasoning — "a single-layer bypass is caught by at
least one of the other two" — is the explicit justification for keeping three independent checks
rather than optimizing down to one.

**No probe, rely on integration tests alone** — rejected: integration tests run in CI's
controlled environment and would not catch a *production* environment silently lying about WAL
support (e.g. a user running `bm` from a network-mounted home directory) — the mandate requires
this to be checked in the real environment where the adapter runs, not only in CI.

## Consequences

**Positive**: closes the exact bug class named in the mandate (adapter claims a contract it
cannot actually honor in the deployed environment) before any write is attempted, with three
independent checks so no single oversight (forgot to call probe, forgot to define probe, defined
a probe that doesn't test anything real) goes undetected.

**Negative**: adds implementation and CI surface area (fault-injection test harness,
AST-walking pre-commit script) that a CLI this size would not otherwise need. Accepted per Core
Principle 13's explicit stance that probing is "a first-class design responsibility," not
optional, and the mandate applies regardless of project size.

**Self-application**: a fourth, narrow behavioral test verifies `Probe()` is actually invoked at
startup for each adapter (not merely defined) — this principle applies to its own enforcement.
