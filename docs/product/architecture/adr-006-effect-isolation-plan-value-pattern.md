# ADR-006: Effect Isolation via Functional Core / Imperative Shell + Plan-Value Pattern

## Status
Accepted

## Context
Core Principle 12 (2026-05-15 mandate) requires that "side-effect-free function silently writes"
be non-representable by design, not merely testable-around. US-05's duplicate/tag-update flow
("offered to add the tag to the existing entry rather than creating a duplicate") is exactly the
kind of confirm-then-mutate flow where this bug class historically hides.

## Decision
Structure the save/find/share flows as a **Functional Core / Imperative Shell**:
- All business logic (URL validation, tag normalization, duplicate detection, save planning,
  match ranking, snippet formatting) is implemented as pure functions with no I/O.
- The duplicate/tag-update decision is modeled as a **Plan-value**: `SavePlanner.Plan(url, tag,
  existing) -> SavePlan` is pure and returns a `SavePlan{New | Duplicate | TagUpdate}` value.
  `BookmarkWriter.Execute(plan)` is the only impure function, and it is the sole place a write to
  the store can originate from.
- `BookmarkReader` (used by `bm find` and `bm share`) exposes **no write methods** — read and
  write are split into separate driving-adjacent ports, so `find`/`share` are structurally
  incapable of mutating the store, not merely conventionally expected not to.

## Alternatives Considered

**Direct imperative implementation** (validate inline, check duplicates inline, write inline
within the command handler) — the natural default for a small CLI. Rejected: this is precisely
the shape that lets "preview/check logic silently writes" bugs hide, because there is no type-level
distinction between a function that decides and a function that acts — the exact failure mode the
mandate's v3.15.1 dry-run precedent names. Also harder to unit test without a real database.

**Command objects with an internal `dryRun` boolean flag** — a common but rejected pattern per
the mandate's own reasoning: a boolean flag threaded through an otherwise-effectful function
still allows the effectful path to be reached accidentally (e.g. a default-false flag forgotten
in a new call site) — it's a runtime check, not a structural impossibility. The Plan-value
pattern instead makes the pure decision and the impure execution two different functions with
different types, so there is no path through which "just checking" can accidentally write.

## Consequences

**Positive**: the bug class "duplicate-check offered a tag-update but silently created a second
row anyway" becomes structurally unrepresentable — there is no code path where `SavePlanner`
(pure) can write, and no code path where `find`/`share` (read-only port) can write. All business
logic is testable with plain unit tests, no database or filesystem needed, satisfying the
Testability quality attribute cheaply.

**Negative**: slightly more indirection than a direct imperative implementation for a CLI this
small — three named types (`SavePlan` and its variants) where a simpler tool might inline an
if/else. Accepted: the mandate treats this as identity-essential, and the indirection cost is low
relative to the bug class it closes.
