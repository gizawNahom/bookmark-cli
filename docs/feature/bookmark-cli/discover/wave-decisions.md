# Wave Decisions — DISCOVER (bookmark-cli)

Facilitator: Scout (nw-product-discoverer) | Date: 2026-08-07
Status: **Discovery decisions recorded — handoff blocked pending G2/G4 remediation** (see feature-delta.md)

## D1-Tagged Decisions

**D1-001: Narrow target persona from "developers" (broad, per project-brief.md) to "terminal-first
infra/platform/SRE/staff engineers with shared technical-reference needs."**
- Rationale: Broad-sample Phase 1 confirmation rate was 4/8 = 50%, below the 60% G1 threshold.
  Confirmed subjects clustered entirely in one segment (senior/infra/platform roles, team-shared
  reference use); disconfirmed subjects were either already well-served by existing tools or not
  terminal-first.
- Source: Phase 1 interviews, Round 1 (transcripts #1–#8) and Round 2 (transcripts #9–#11) in
  `docs/feature/bookmark-cli/feature-delta.md`.

**D1-002: Prioritize "Locate" (search/retrieve) and "Capture" (frictionless save) as the top two
opportunities over "Organize" (tagging depth) or "Dedup/monitor" (staleness detection).**
- Rationale: Opportunity Algorithm scores — Locate 16/20, Capture 15/20, vs. Organize 11/20 and
  Dedup 6/20.
- Source: Phase 2 opportunity scoring, 7 segment interviews (Priya, Sam, Nadia, Aisha, Jordan, Priti, Marco).

**D1-003: Include team-sharing as a first-class, third top opportunity, not an afterthought.**
- Rationale: Opportunity score 14/20; explicit pain in 5/7 segment interviews describing knowledge
  silos, stale shared docs, and merge-conflict-prone shared link files.
- Source: transcripts #6 (Nadia — implicit), #8 (Aisha), #9 (Jordan), #10 (Priti), #11 (Marco).

**D1-004: Defer the browser-extension-hybrid concept; commit to CLI-only scope for MVP.**
- Rationale: Segment consistently values context-switch avoidance as the core driver; no interview
  subject requested browser integration; adding browser surface area would dilute the differentiator.
- Source: Phase 1–2 interviews, all segment transcripts.

**D1-005: Do not proceed to build or to DISCUSS-wave handoff until Feasibility and Viability risks
are independently assessed.**
- Rationale: G4 hard gate — 2 of 4 big risks (Feasibility, Viability/channel) are unaddressed in this
  discovery pass; no engineer participated, no channel/fake-door test was run.
- Source: Phase 4 risk review, `feature-delta.md` G4 section.

**D1-006: Prioritize fixing tag-syntax discoverability before the next round of solution testing.**
- Rationale: Only usability failure observed in Phase 3 testing (1 of 5 users, 1 of 15 tasks);
  comprehension time exceeded the <10s target for that user (~20s).
- Source: Phase 3 solution testing summary, `feature-delta.md`.

---

## Constraints (with evidence source)

| Constraint | Evidence source |
|---|---|
| MVP must support local-first save/find with no required account creation | Nadia (#6), Aisha (#8) — resistance to another SaaS login inferred from workaround choices (bash script, git-repo markdown) rather than any existing bookmark SaaS |
| Share feature must require zero install for the recipient | Aisha (#8): "if it's not one-command for my teammates I already lost" |
| Tag syntax must be more discoverable than the tested prototype | Phase 3 testing, Aisha task-3 failure |
| No server infrastructure assumed for MVP — **unverified, pending feasibility spike** | Inferred from local-first constraint above; not independently confirmed by an engineer |

---

## Validated Assumptions (with confidence)

| Assumption | Confidence | Evidence source |
|---|---|---|
| Terminal-first infra/platform engineers accumulate reference links faster than ad hoc tools organize them | HIGH | 7/7 segment interviews, past-behavior evidence |
| "Locate" (search/retrieve) is the most painful job step | HIGH | Opportunity score 16/20; 6/7 segment interviews |
| Team-sharing of curated links is a distinct, valuable job | MEDIUM-HIGH | Opportunity score 14/20; 5/7 segment interviews |
| CLI-native capture reduces context-switch cost vs. browser bookmarking | MEDIUM | Opportunity score 13/20; task-completion 93% (n=5), no longitudinal data |
| Tag-based organization preferred over folder hierarchies | MEDIUM | Opportunity score 11/20; inferred from interview language, not directly usability-tested |

## Invalidated Assumptions (with evidence ref)

| Assumption | Status | Evidence ref |
|---|---|---|
| "Developers [generally] lose track of useful links" | Invalidated as stated / narrowed to segment | Round 1 broad sample 4/8 (50%) confirmed, below 60% threshold — transcripts #2, #3, #5, #7 disconfirm |
| "Browser bookmarks are disorganized and hard to search" (universal) | Partially invalidated | Marcus (#2): "Raindrop's search is honestly fine"; Tom (#7): "Notion already does this for me" |
| "A CLI tool fits developer workflow better than a browser extension" (as a universal claim) | Not validated for general developers; validated only within narrowed segment | Round 1 vs. Round 2 comparison, `feature-delta.md` G1 section |

---

## Gate Status Summary (see feature-delta.md for full detail)

| Gate | Status |
|---|---|
| G1 (Problem → Opportunity) | PASS — with persona pivot |
| G2 (Opportunity → Solution) | CONDITIONAL PASS — scoring met, cross-functional alignment outstanding |
| G3 (Solution → Viability) | PASS |
| G4 (Viability → Build) | **FAIL** — Feasibility and Viability/channel risks unaddressed |

**Handoff to product-owner is blocked** pending G2 and G4 remediation. See
`docs/feature/bookmark-cli/feature-delta.md` §Pre-requisites for the remediation list.
