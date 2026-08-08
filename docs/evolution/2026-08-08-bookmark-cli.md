# Evolution: bookmark-cli

**Finalized**: 2026-08-08 | **Facilitator**: Apex (nw-platform-architect)
**Status**: DELIVER complete — 7/7 roadmap steps COMMIT/PASS, 25/25 acceptance scenarios GREEN,
refactor pass complete, adversarial review approved, `des-verify-integrity` exit 0.

---

## 1. Feature Summary

`bm` — a local-first, single-static-binary CLI (`bm save`, `bm find`, `bm share`, `bm stats`)
that lets terminal-first infra/platform/SRE engineers capture, retrieve, and share technical
reference links without leaving the terminal. Go + Cobra + SQLite (WAL, FTS5), functional
core / imperative shell paradigm, zero network calls, zero server infrastructure.

## 2. Business Context

**Persona**: terminal-first infra/platform/SRE/staff engineers who curate and re-use technical
reference links, narrowed from the original brief's general "developers" persona after DISCOVER
Round 1 (50% broad-sample confirmation, below the 60% G1 threshold) triggered a segmentation
pivot to Round 2 (100% segment-specific confirmation, n=7).

**Opportunity**: minimize time and context-switching to save, retrieve, and share curated
technical reference links while working in the terminal. Top 3 pursued opportunities (of 7
scored): Locate (score 16), Capture (score 15), Share (score 14).

**Notable risk carried forward and accepted, not remediated**: DISCOVER's G4 gate (Viability →
Build) FAILED — no engineering feasibility spike and no channel/viability test were run. The user
explicitly accepted this risk ("well-understood, low-risk technical patterns for this scope")
rather than remediating, and DESIGN wave later confirmed feasibility at the engineering level
(D1: zero-server architecture requires no network/server component). Channel/viability
(distribution, willingness-to-pay) remains unverified — noted here so it is not silently lost now
that the workspace is being archived.

## 3. Key Decisions (by wave)

### DISCUSS
- Gate override accepted: proceeded to DISCUSS despite DISCOVER's G4 FAIL / G2 CONDITIONAL PASS,
  by explicit user sign-off.
- Scope: top-3 "pursue" opportunities only (JTBD-CAPTURE, JTBD-LOCATE, JTBD-SHARE). Deferred to
  backlog: context-switch reduction, organize/tag, recall-context, dedup/monitor.
- 6 stories (US-01..US-06): walking skeleton (US-01/02/03) + Release 1 (US-04/05/06).

### DESIGN
- **D2 Language/runtime: Go** (ADR-001) — confirmed final by user this session; startup-latency
  budget (<100ms perceived-save) ruled out Python; Rust's learning-curve risk ruled it out for a
  solo maintainer.
- **D3 Storage: SQLite, WAL, FTS5, pure-Go driver** (ADR-002).
- **D4 CLI framework: Cobra** (ADR-003).
- **D7 Effect isolation: Functional Core / Imperative Shell + Plan-value pattern** (ADR-006) —
  written to `CLAUDE.md` after direct user confirmation (an earlier agent-relayed approval claim
  was correctly rejected per this agent's standing rule).
- **D8 Adapter trust: `Probe()` contract, 3-layer enforcement** (ADR-007) — subtype (compile-time
  interface), structural (`go/ast` probe-presence check), behavioral (fault-injection CI harness).
- **D9 Architecture pattern: modular monolith, hexagonal ports-and-adapters** — microservices/
  event-sourcing/CQRS explicitly rejected as disproportionate for this scope.
- 12-component contract-shape classification (brief.md Section 5) — every component classified as
  pure-function, bounded-change, or Plan-value (unbounded-preservation-safe).
- Peer review: `nw-solution-architect-reviewer` conditionally approved, 1 HIGH (effort estimate
  revised +2-3 days for Earned Trust tooling) + 1 MEDIUM (Security quality attribute added),
  both resolved same iteration.

### DEVOPS
- **CI/CD**: GitHub Actions, trunk-based branching. Local pre-commit/pre-push → PR/commit stage →
  fault-injection (behavioral, ADR-007 layer 3) → release (tag-triggered).
- **Deployment strategy**: no live-service target exists (zero-server architecture) — the
  CLI-equivalent is release/versioning via Homebrew tap + GitHub Release binaries + `go install`,
  with a rollback contract designed first (tap formula revert, mark-bad-release, immutable git
  tags — Core Principle 7).
- **Mutation testing: pre-release strategy** (persisted to `CLAUDE.md`) — `gremlins`, entire
  solution, gated at release tag, advisory-only for v1.0.0 (no prior baseline exists yet); action
  item recorded for v1.1.0+ to set a numeric floor off the v1.0.0 baseline.
- **Telemetry/KPI instrumentation designed**: local opt-in event log (`UsageLogger` port,
  `FileUsageLogAdapter`/`NoOpUsageLogAdapter`), privacy-by-design (never logs URL/tag content),
  new read-only `bm stats` driving port. Full contract in `docs/product/kpi-contracts.yaml`.
- **Amendment (Final Wave Review Gate HIGH-2)**: Windows removed from the GoReleaser cross-compile
  matrix — WSL2 already gives Windows users a tested path; shipping an untested native Windows
  binary was judged worse than not shipping one for v1.
- Per-wave `nw-platform-architect-reviewer` review was skipped (none of its trigger conditions
  fired — conventional GitHub Actions + Homebrew setup); the mandatory consolidated Final Wave
  Review Gate covered DEVOPS instead (see below).

### DISTILL
- 25 acceptance scenarios (1 walking skeleton + 24 milestone/adapter), 100% real-adapter
  subprocess/FS layer (zero driven-external ports exist, so no in-memory-double layer applies).
- Adopted `github.com/stretchr/testify` (`require`/`assert`) as the standing assertion convention
  for this project, replacing raw `t.Fatalf` patterns across all 7 scenario/harness files.
- Final Wave Review Gate: 4 reviewers (Eclipse/Architect/Forge/Sentinel) — 3 approved outright,
  Forge (DEVOPS review) conditionally approved with 0 blockers, 1 critical, 2 high, 5 medium, 2
  low. All critical/high resolved same session (mutation-gate pass/fail criterion made explicit,
  environments.yaml platform-coverage claim narrowed to CI-verified, Windows build removed).
  Medium/low items accepted as DELIVER-scope action items or backlog.

### DELIVER
- All 7 roadmap steps (01-01 through 04-01) reached RED→GREEN→COMMIT.
- **Test-authoring defect found and fixed during step 02-02** (`find_scenarios_test.go`):
  `TestFind_NoMatchResponse_IsResponsive` and `TestFind_EmptyStore_ShowsOnboardingMessage` both
  used a genuinely empty store precondition but asserted mutually exclusive outcomes (onboarding
  message vs. ordinary no-match message) — no single implementation could satisfy both. Fixed by
  seeding `TestFind_NoMatchResponse_IsResponsive` with a non-matching bookmark via
  `.WithExistingStore(...)`, mirroring its sibling scenario's pattern. Test-only fix, zero
  production code touched.
- Post-Merge Integration Gate: PASS. `go build`, `go vet`, `gofmt -l` clean; `go test ./...` —
  25/25 GREEN, 0 skips remaining. All 3 declared environments (`clean`, `existing-store`,
  `degraded-filesystem`) exercised via fixtures. Elevator Pitch demo executed against the built
  `bm` binary for all 6 non-infrastructure user stories — all exit 0 with correct stdout.
- Refactor pass (L1-L6) completed after the gate. Adversarial review approved. Mutation testing
  skipped per the pre-release strategy (runs at release-tag time, not per-feature).
  `des-verify-integrity` exit 0.

## 4. Issues Encountered

| Issue | Wave found | Resolution |
|---|---|---|
| DISCOVER G4 gate FAIL (feasibility + viability unassessed) | DISCOVER | User-accepted risk override, not remediated; feasibility later confirmed at DESIGN level (D1), viability (channel/economics) remains genuinely unverified — carried forward as an open risk, not closed |
| DISCOVER G2 conditional pass (no cross-functional alignment session) | DISCOVER | Same override; no alignment session held during DISCUSS |
| Test-authoring defect: mutually exclusive assertions on the same empty-store precondition | DELIVER step 02-02 | Test-only fix (seed one scenario with a non-matching bookmark); zero production code changed |
| Forge (DEVOPS) review CRITICAL-1: mutation-gate had no pass/fail criterion | DISTILL Final Wave Review Gate | Made `gremlins` advisory-only for v1.0.0; numeric floor deferred to v1.1.0+ as an explicit action item |
| Forge HIGH-1: `environments.yaml` overclaimed macOS 13.x coverage | DISTILL Final Wave Review Gate | Narrowed to CI-verified (14.x/15.x) vs. unverified-best-effort (13.x) |
| Forge HIGH-2: Windows binaries built with zero test coverage | DISTILL Final Wave Review Gate | Windows removed from GoReleaser matrix; WSL2 is the supported path for v1 |

## 5. Lessons Learned

- **Segmentation pivots surface real signal, not noise**: the 50%→100% confirmation jump between
  DISCOVER Round 1 (broad) and Round 2 (segment-narrowed) was a legitimate G1 remediation path,
  not p-hacking — the disconfirmed Round 1 subjects were each independently explained (already
  well-served by an existing tool, non-terminal-first, future-intent-only).
  compliment-vs-commitment gap (100% "would use" self-report vs. 60% strong commitment signal)
  is worth flagging explicitly in every future discovery pass rather than averaging away.
- **Plan-value pattern (ADR-006) closed a real bug class before it could occur**: modeling
  duplicate/tag-update detection as a pure `SavePlan` value (never a direct mutation) is directly
  traceable to a precedent bug (v3.15.1 dry-run bug) — worth reusing as a default pattern for any
  future save/dedup-shaped flow in this codebase.
- **Zero-server architectures need a CLI-native reinterpretation of DEVOPS wave concepts**:
  "deployment strategy" became release/versioning, "observability" became local stderr + opt-in
  usage log, "environment matrix" became install-state fixtures rather than staging/prod tiers.
  None of these were skipped — they were mapped to their CLI-equivalent rather than omitted.
- **Advisory-only gates need an explicit, dated action item, not a silent placeholder**: the
  mutation-testing pass/fail-criterion gap (Forge CRITICAL-1) was closed correctly by making the
  first release's gate advisory rather than inventing an unjustified numeric floor — the
  follow-up obligation (set the floor from the v1.0.0 baseline) is recorded, not left implicit.

## 6. Links to Migrated / Permanent Artifacts

- Architecture SSOT (already permanent, not migrated): `docs/product/architecture/brief.md`
  (now includes a "Component Inventory — SHIPPED" section confirming all 12 DESIGN-wave
  components shipped as designed, plus 4 components added during DEVOPS telemetry design)
- ADRs (already permanent, not migrated): `docs/product/architecture/adr-001` through `adr-007`
- KPI instrumentation SSOT (already permanent, updated this finalize): `docs/product/kpi-contracts.yaml`
  (KPI-1 baseline recorded as "not yet measured, pending first dogfood")
- UX journeys (migrated this finalize): `docs/ux/bookmark-cli/journey-save-find-share.yaml`,
  `docs/ux/bookmark-cli/journey-save-find-share-visual.md`
- Full wave-by-wave narrative (workspace, preserved not deleted):
  `docs/feature/bookmark-cli/feature-delta.md`

## 7. Not Migrated (No Matching Destination-Map Row / Already Permanent)

- `design/architecture-design.md`, `design/component-boundaries.md`, `design/technology-stack.md`,
  `design/data-models.md` — **no files to migrate**: this project uses the lean single-narrative
  `feature-delta.md` convention; equivalent content lives in `feature-delta.md`'s `## Wave: DESIGN`
  section and in `docs/product/architecture/brief.md` (already permanent, cross-feature SSOT).
- `design/adrs/ADR-*.md` — **no files to migrate**: ADRs were authored directly under
  `docs/product/architecture/adr-00N-*.md` (already permanent), never under
  `docs/feature/bookmark-cli/design/adrs/`.
- `distill/walking-skeleton.md` — **no files to migrate**: no file by this name exists; the
  walking-skeleton strategy and scenario are documented inline in `feature-delta.md`'s
  `## Wave: DISTILL` section and in `docs/feature/bookmark-cli/distill/red-classification.md`
  (a RED-gate evidence file, not a walking-skeleton spec — correctly left in the workspace, not a
  destination-map match).
