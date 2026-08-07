# CLAUDE.md

## Development Paradigm

**Functional-leaning imperative (Go)** — Functional Core / Imperative Shell.

Confirmed final by the user (2026-08-07), following from the DESIGN wave's component
classification (`docs/product/architecture/brief.md` Section 5):

- All business logic (URL validation, tag normalization, save planning, match ranking, snippet
  formatting) is implemented as **pure functions** — no I/O, unit-testable in isolation.
- **Composition over inheritance**; immutable value types for `SavePlan`, `ValidationResult`,
  `RankedMatches`.
- A **thin imperative shell** (Cobra command handlers + driven adapters) isolates all I/O from
  the core.
- This is not a "pure FP language" choice (Go isn't one) — it is the Functional Core /
  Imperative Shell pattern applied within Go.

**DELIVER-wave routing**: implementation work for this paradigm routes to
`@nw-functional-software-crafter`.

See `docs/product/architecture/brief.md` ("Paradigm Selection") and ADR-006
(`docs/product/architecture/adr-006-effect-isolation-plan-value-pattern.md`) for full rationale.

## Mutation Testing Strategy

This project uses **pre-release** mutation testing. Runs on entire solution before each
release. Delivery not blocked.
