# Wave Decisions — DISCUSS (bookmark-cli)

Facilitator: Luna (nw-product-owner) | Date: 2026-08-07
Status: Requirements crafted, DoR-validated, pending peer review before DESIGN handoff.

## Upstream Changes (Gate Override — Back-Propagation)

The DISCOVER wave (`docs/feature/bookmark-cli/discover/wave-decisions.md`,
`docs/feature/bookmark-cli/feature-delta.md`) explicitly failed its own
handoff gate:

> **G4 (Viability → Build): FAIL** — "Per gate rule, all 4 risks must be
> addressed before proceeding to build. 2 of 4 are unaddressed. Decision:
> FAIL. Remediation required: (1) engineering feasibility spike, (2)
> channel/viability validation... Re-evaluate G4 after both close."

> **G2 (Opportunity → Solution): CONDITIONAL PASS** — "Team alignment NOT
> held — this discovery pass was run solo (synthetic/simulated), violating
> Principle 7 (no solo discovery: PM + Designer + Engineer together)...
> the gate cannot be considered fully closed until a real cross-functional
> alignment session occurs."

> "**NOT HANDOFF-READY.** ... do not treat this package as ready for the
> DISCUSS wave until remediation items above are closed."

**Decision**: The product owner (with the user's explicit sign-off) accepts
the G2/G4 risk and proceeds to DISCUSS without the remediation items being
closed. This is a deliberate PO-level risk acceptance, not a silent skip of
the gate.

**Rationale (user-provided)**: "Low actual risk — the local-first storage and
copy-paste share mechanism are well-understood, low-risk technical patterns;
a formal spike is overkill for this scope."

**What this means downstream**:
- DESIGN wave (solution-architect) should treat "local-first storage, no
  server infrastructure for MVP" as a working assumption carried from
  DISCOVER, not an independently verified constraint — flagged again here so
  it isn't silently promoted to "verified" by omission.
- G2's cross-functional alignment gap is not remediated by this wave; DISCUSS
  proceeded on the DISCOVER evidence as recorded (solo/simulated discovery
  pass). This is noted as an open item, not resolved.
- If a future feasibility spike (still recommended, just deprioritized)
  invalidates local-first storage, only Technical Notes in `user-stories.md`
  need revision — all UAT scenarios and acceptance criteria are observable
  CLI behavior, not implementation-specific, so they should survive that
  change intact.

---

## Scope Assessment (Elephant Carpaccio Gate)

Assessed against the DISCOVER top-3 pursue opportunities (JTBD-CAPTURE,
JTBD-LOCATE, JTBD-SHARE), consistent with the lightweight-UX-depth
configuration.

- **Story count**: 6 (US-01 through US-06)
- **Bounded contexts**: 1 (local bookmark store + CLI surface — no
  multi-service split)
- **Walking skeleton integration points**: 3 (save → find → share, single
  local data store, no external systems)
- **Estimated effort**: Walking skeleton (US-01, US-02, US-03) ≈ 3 days;
  Release 1 (US-04, US-05, US-06) ≈ 3-4 days. Total ≈ 6-7 days.
- **Independent outcomes**: All 6 stories serve one coherent journey (capture
  → locate → share); none is independently shippable as a separate product.

None of the oversized signals fire (not >10 stories, not >3 bounded
contexts, walking skeleton has exactly 3 integration points, estimated
effort well under 2 weeks, no independently-shippable sub-features).

## Scope Assessment: PASS — 6 stories, 1 bounded context, estimated 6-7 days

---

## Out-of-Scope (This Pass)

Carried forward from DISCOVER's own prioritization, tracked in
`docs/product/jobs.yaml` as `scope_status: out_of_scope_backlog`:

- **JTBD-CONTEXT-SWITCH** (score 13, `evaluate_later`) — shell integration /
  man-page-style lookup to reduce terminal-vs-browser context switching.
- **JTBD-ORGANIZE** (score 11, `evaluate_later`) — tag autocomplete, tag
  rename, deeper organizational tooling beyond the flat save-with-tag model.
- **JTBD-RECALL-CONTEXT** (score 9, `backlog`) — optional note field
  capturing why a link was saved.
- **JTBD-DEDUP-MONITOR** (score 6, `deprioritize`) — link-rot / staleness
  checker.

None of these are required for the walking skeleton or Release 1 to deliver
coherent, demonstrable value.

---

## Job Discovery Formalization

DISCOVER's 7 scored jobs were not re-interviewed; they were formalized into
job-story + four-forces format in `docs/product/jobs.yaml` (per the
skill guidance to treat already-scored jobs as a starting corpus rather than
re-running JTBD discovery from scratch). No DIVERGE artifacts existed for
this feature (project went DISCOVER → DISCUSS directly), so grounding is
entirely in DISCOVER evidence — noted as the risk carried by the G2/G4
override above.

---

## DoR Validation Summary

Full per-story validation lives inline with each story's checklist in
`user-stories.md`; summary below.

| Story | Problem statement | Persona | 3+ examples | UAT (3-7) | AC from UAT | Right-sized | Tech notes | Deps tracked | job_id | Elevator Pitch |
|---|---|---|---|---|---|---|---|---|---|---|
| US-01 Capture | PASS | PASS | PASS (3) | PASS (4) | PASS (4) | PASS | PASS | PASS (none) | PASS | PASS |
| US-02 Locate | PASS | PASS | PASS (3) | PASS (5) | PASS (5) | PASS | PASS | PASS (US-01) | PASS | PASS |
| US-03 Share | PASS | PASS | PASS (3) | PASS (5) | PASS (5) | PASS | PASS | PASS (US-01, US-02) | PASS | PASS |
| US-04 Tag discoverability | PASS | PASS | PASS (3) | PASS (4) | PASS (4) | PASS | PASS | PASS (US-01) | PASS | PASS |
| US-05 Save error handling | PASS | PASS | PASS (3) | PASS (4) | PASS (4) | PASS | PASS | PASS (US-01) | PASS | PASS |
| US-06 Find no-results | PASS | PASS | PASS (3) | PASS (4) | PASS (4) | PASS | PASS | PASS (US-02) | PASS | PASS |

**DoR Status: PASSED** (all 6 stories, all 9 items — see `nw-po-review-dimensions`
Dimension 0 Elevator Pitch check folded in as item 9 alongside the 8-item
LeanUX checklist and the job_id traceability requirement).

---

## Peer Review

**Iteration 1** (`nw-product-owner-reviewer`): `conditionally_approved`.
0 critical issues, 2 high, 3 medium, 1 low, **no blocking issues**. All hard
gates passed: DoR (8/8 all stories), JTBD traceability (6/6), Elevator Pitch
(6/6), journey coherence, slice composition.

High-severity findings (both design-phase NFR gaps, not story defects):
1. Accessibility NFR not explicit.
2. Data reliability / backup strategy for local-first storage not addressed.

**Remediation applied** (same iteration, no re-review needed — findings were
advisory NFR gaps with `due_phase: DESIGN`, not structural defects in the
stories themselves): added a "NFR Guardrails Added After Peer Review"
subsection to `user-stories.md` System Constraints, explicitly carrying
data-reliability, accessibility, concurrency, and resource-scaling
constraints forward to DESIGN rather than retrofitting them as new DISCUSS
stories (would prescribe storage-format solutions, violating Core Principle
5: problem-first, solution-never).

Medium/low findings (concurrent modification, resource exhaustion, recipient
validation depth, ranking/fuzzy-match specificity) are DESIGN-phase
decisions by nature and are captured in the same NFR Guardrails subsection
or left as DESIGN's responsibility per the reviewer's own `due_phase` marks.

**Result: peer review requirement satisfied in 1 iteration** — no critical
or high-severity *blocking* issue existed, so a second review pass was not
required before handoff.

---

## Handoff Package (to DESIGN wave — solution-architect)

- `docs/feature/bookmark-cli/discuss/journey-save-find-share-visual.md`
- `docs/feature/bookmark-cli/discuss/journey-save-find-share.yaml`
- `docs/feature/bookmark-cli/discuss/shared-artifacts-registry.md`
- `docs/feature/bookmark-cli/discuss/story-map.md`
- `docs/feature/bookmark-cli/discuss/user-stories.md`
- `docs/feature/bookmark-cli/discuss/outcome-kpis.md`
- `docs/product/jobs.yaml` (SSOT)
- `docs/product/journeys/bookmark-cli.yaml` (SSOT, updated)
- `docs/feature/bookmark-cli/slices/` (slice briefs, walking skeleton + Release 1)

Risk flagged for DESIGN wave: local-first/no-server-infra assumption is
carried forward unverified (see Upstream Changes above) — solution-architect
should treat this as a design input requiring confirmation, not a closed
decision.
