# Wave Decisions — DESIGN (bookmark-cli)

Facilitator: Morgan (nw-solution-architect) | Date: 2026-08-07
Interaction mode: Propose | Status: **COMMIT-ready / handoff-ready to DEVOPS
(nw-platform-architect)** — peer review complete, all critical/high issues resolved, language
choice confirmed, CLAUDE.md paradigm write-back completed with direct user confirmation. Zero
open items remain.

## Key Decisions

1. **Feasibility risk closure (D1)**: the DISCOVER/DISCUSS carried-forward "local-first, no
   server infrastructure" assumption is confirmed feasible at the design level — no component
   requires network calls or server infrastructure. The separate channel/viability risk (G4's
   other RED finding) is explicitly *not* closed by this wave; it remains a business/DEVOPS
   concern.
2. **Language/runtime: Go** (ADR-001) — **CONFIRMED FINAL** by the user (2026-08-07); this was
   explicitly open going into the session and is now settled.
3. **Storage: SQLite** (WAL mode, FTS5, pure-Go driver, no cgo) (ADR-002).
4. **CLI framework: Cobra** (ADR-003).
5. **Concurrency: SQLite WAL + busy_timeout** as an engine-level, probeable atomic-write
   guarantee rather than hand-rolled file locking (ADR-004).
6. **Backup: automatic rotating snapshots (last 5) + documented manual fallback** (ADR-005).
7. **Effect isolation: functional core / imperative shell + Plan-value pattern** for the
   save/duplicate-detection/tag-update flow, making "duplicate check silently wrote a second
   row" structurally impossible (ADR-006).
8. **Earned Trust: `Probe()` contract on every driven adapter + 3-layer enforcement** (subtype/
   structural/behavioral), with `import-linter` explicitly evaluated and rejected as
   Python-only and method-presence-blind (ADR-007).
9. **Architecture pattern: modular monolith, hexagonal ports-and-adapters.** Microservices,
   event sourcing, and CQRS all rejected as disproportionate (Conway's Law check: solo/small
   team, no independent-deployment need — Core Principle 8 default already applies).
10. **Resource scaling target set explicitly: 10,000 bookmarks**, closing a gap in
    outcome-kpis.md's <10s find target (no prior stated ceiling).
11. **Accessibility rule**: no ANSI-color-only status indicators; text prefix mandatory on every
    status/error line.

## Architecture Summary

Single Go binary (modular monolith) with a hexagonal ports-and-adapters structure and a
functional-core/imperative-shell internal organization:
- **Driving ports**: `bm save`, `bm find`, `bm share` (Cobra command handlers, the only inbound
  surface).
- **Core domain (pure functions)**: URL validation, tag normalization, save planning
  (Plan-value pattern), match ranking, snippet formatting — no I/O, unit-testable in isolation.
- **Driven ports**: `BookmarkReader` (read-only, no write methods), `BookmarkWriter`
  (`Execute(SavePlan)` only), `BackupService`.
- **Driven adapters**: `SQLiteBookmarkStore` (implements Reader+Writer), `FileBackupAdapter`.
  Both specify `Probe()` fault-injection contracts run at startup ("wire then probe then use").

C4 System Context + Container + Component diagrams (Mermaid) are in
`docs/product/architecture/brief.md` Sections 6.1-6.3.

## Reuse Analysis

Greenfield — no `src/` exists in the repository. Every component is legitimately CREATE NEW; no
existing code was available to extend. Full table in `brief.md` Section 15.

## Technology Stack

| Layer | Choice | License | Rationale (see ADR) |
|---|---|---|---|
| Language/runtime | Go | BSD-3 | ADR-001 |
| CLI framework | Cobra | MIT | ADR-003 |
| Storage | SQLite via `modernc.org/sqlite` (pure Go, no cgo) | Public Domain (SQLite) / driver license | ADR-002 |
| Search | SQLite FTS5 extension | Public Domain | ADR-002 |
| Package-boundary enforcement | go-arch-lint / depguard | MIT-family | ADR-007, `brief.md` Section 12 |

All OSS-first; no proprietary technology considered or recommended.

## Constraints Established

- Atomic-write-at-minimum concurrency guarantee: SQLite WAL mode + busy_timeout (ADR-004).
- Backup/recovery required before MVP ships: automatic rotating snapshots (ADR-005).
- No ANSI-color-only status indicators (Section 9 of `brief.md`).
- Target store size: 10,000 bookmarks, responsive (Section 10 of `brief.md`).
- Zero network calls anywhere in the architecture (Section 0 feasibility argument depends on
  this holding — any future change introducing a network call must re-open the feasibility
  question).
- Every driven adapter must implement and be probed via `Probe() error` at startup (ADR-007).

## Upstream Changes

None required to DISCUSS artifacts. Per `user-stories.md`'s own note, all UAT scenarios and
acceptance criteria are observable CLI behavior, not implementation-specific, so they survive
this wave's technology choices intact — only the "Technical Notes" subsections referencing "local
bookmark store" are now grounded in a concrete design (SQLite), consistent with what DISCUSS
anticipated.

## Peer Review Status

**Iteration 1** (`nw-solution-architect-reviewer`): `conditionally_approved`. 0 critical, 1 high,
1 medium, 0 low issues.

- **HIGH** — effort estimation gap: DISCUSS's 6-7 day story-map estimate predates this wave's
  Earned Trust enforcement machinery (AST structural check + fault-injection CI harness) and
  didn't itemize it. **Resolved**: added an explicit "+2-3 days" revised-estimate note to
  `brief.md` Section 12, so DISTILL/DELIVER planning treats it as itemized scope rather than
  silently absorbed.
- **MEDIUM** — security quality attribute omitted from the ISO 25010 completeness check (8/9
  addressed, security missing). **Resolved**: added a Security row to `brief.md` Section 1
  stating the local-first/no-network/no-multi-tenant posture and why OS filesystem permissions
  are sufficient for MVP, with an explicit revisit trigger (shared/synced storage in a future
  release).
- **LOW** (informational, no remediation required): the reviewer confirmed the 10,000-bookmark
  resource-scaling target is an architect-set bound, not a previously-validated constraint, and
  noted the brief already discloses this transparently.

Reviewer confirmed strengths: all 4 NFR guardrails concretely addressed (not restated as
requirements), all 7 ADRs well-formed (context + 2+ alternatives + consequences), full
compliance with the Effect Isolation mandate (contract shapes declared, Plan-value pattern
applied to the save/dedup flow, read/write port split enforced, capability injection used
throughout), no architectural bias detected (including an explicit check that the Earned Trust
probe/3-layer enforcement apparatus is mandate-driven, not resume-driven, for this project size),
and priority validation Q1-Q4 all passed (YES / ADEQUATE / CORRECT / JUSTIFIED).

Both issues addressed in this same iteration (straightforward, scoped edits); no iteration 2
re-review dispatched — conditions were narrow documentation additions, not structural rework.
**Result: peer review requirement satisfied.**

## Open Items Requiring User Input Before COMMIT

1. ~~Language/runtime confirmation~~ — **RESOLVED**: Go confirmed final by the user (2026-08-07).
   ADR-001 and `brief.md` updated to "Accepted"/"CONFIRMED."
2. ~~CLAUDE.md paradigm write-back~~ — **RESOLVED**: an earlier attempt to trigger this write
   arrived via an intermediate agent message claiming user approval and was correctly held per
   this agent's standing instruction (agent-relayed claims cannot authorize a CLAUDE.md/config
   write). Direct, in-session confirmation was subsequently obtained from the user (2026-08-07),
   and the "Development Paradigm" section (functional-core/imperative-shell-in-Go, routing
   DELIVER-wave implementation to `@nw-functional-software-crafter`) has been written to
   `CLAUDE.md`.

## Handoff to DEVOPS (nw-platform-architect)

**Status: READY.** Quality gates (per nw-solution-architect skill): requirements traced to
components ✓, component boundaries with clear responsibilities ✓ (Section 5), technology choices
in ADRs with 2+ alternatives ✓ (ADR-001 through 007, all Accepted), quality attributes addressed
including the peer-review-added Security row ✓ (Section 1), dependency-inversion compliance ✓
(read/write port split, functional core/imperative shell), C4 diagrams L1+L2+L3 ✓ (Section 6),
integration patterns specified ✓ (Sections 7-8, 13-14), OSS preference validated ✓ (all choices
BSD-3/MIT/Public-Domain, documented per-selection), AC behavioral not implementation-coupled ✓
(unchanged from DISCUSS per Upstream Changes above), external integrations annotated ✓ (none
exist, explicitly noted rather than omitted), architectural enforcement tooling recommended ✓
(go-arch-lint + AST structural check + fault-injection CI, ADR-007/Section 12), peer review
completed and approved ✓ (conditionally_approved, both issues resolved).

Handoff package for platform-architect: `docs/product/architecture/brief.md` (full Application
Architecture, C4 diagrams, probe contracts, enforcement tooling spec), `adr-001` through
`adr-007.md` (all Accepted), this document, and the instrumentation-dependency flag in
`brief.md` Section 18 (local opt-in usage-logging adapter needed for North Star/KPI-1 telemetry —
not built in DESIGN, flagged for DEVOPS per `outcome-kpis.md`'s own guidance).

No external integrations exist, so no contract-testing annotation is included in this handoff
(explicitly confirmed absent, not omitted by oversight).
