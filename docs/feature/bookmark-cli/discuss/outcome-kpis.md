# Outcome KPIs: bookmark-cli

## Feature: bookmark-cli (walking skeleton: save / find / share)

### Objective
Within one pilot cycle, terminal-first infra/platform engineers stop reaching
for scratch files, bash scripts, and Slack-DM-to-self as their default way to
capture, retrieve, and share technical reference links — and `bm` becomes the
habitual replacement.

### Outcome KPIs

| # | Who | Does What | By How Much | Baseline | Measured By | Type |
|---|-----|-----------|-------------|----------|-------------|------|
| 1 | Terminal-first infra/platform engineers (pilot: Sam, Nadia + 3 more) | Save a reference link via `bm save` instead of their prior workaround (text file, bash script, Obsidian) | 80%+ of new saves go through `bm save` within 2 weeks of adoption | 0% (new tool; prior workaround = 100%) | Opt-in local usage count + weekly pilot self-report | Leading |
| 2 | Same pilot segment | Locate a previously saved link via `bm find` | Median time-to-locate drops from 10-15 min/week (DISCOVER baseline, Sam #4) to a single lookup under 10 seconds | 10-15 min/week spent hunting for links | Task-timing during pilot + weekly self-reported time estimate | Leading |
| 3 | Pilot team members receiving a shared link | Use a `bm share` snippet with zero additional install steps | 100% of share recipients need zero install (hard guardrail, not a stretch target) | N/A (current workaround already requires no install, e.g. Slack paste) — this KPI protects against regression as `bm share` is built | Pilot observation: recipient confirms usable link with no setup | Leading (guardrail) |
| 4 | Aisha-profile users (staff/tech-lead curating for a team) | Comprehend the `--tag` flag syntax without a failed attempt | Comprehension time stays under 10 seconds (DISCOVER Phase 3 target; Aisha's tested failure was ~20s) | ~20s / 1 failed attempt out of 5 users (DISCOVER Phase 3) | Task-timing during Release 1 usability re-test (US-04) | Leading (guardrail) |

### Metric Hierarchy

- **North Star**: Weekly active `bm find` invocations per pilot user — the
  clearest signal that retrieval habit has formed (matches DISCOVER Lean
  Canvas "Key Metrics" hypothesis).
- **Leading Indicators**: `bm save` adoption rate vs. prior workaround (KPI 1),
  time-to-locate (KPI 2).
- **Guardrail Metrics**: zero-install share compliance (KPI 3) — a hard
  constraint from DISCOVER (Aisha, #8), must never regress; tag-flag
  comprehension time (KPI 4) — the one known usability defect being fixed in
  Release 1, must not silently reappear.

### Measurement Plan

| KPI | Data Source | Collection Method | Frequency | Owner |
|-----|------------|-------------------|-----------|-------|
| 1. Save adoption | Local usage log (opt-in) | Count `bm save` invocations vs. pilot self-report of workaround usage | Weekly during pilot | product-owner + pilot facilitator |
| 2. Time-to-locate | Pilot session observation | Stopwatch task-timing + weekly self-report | Weekly during pilot | product-owner |
| 3. Zero-install share | Pilot observation | Direct confirmation from recipient (no install performed) | Per share event during pilot | product-owner |
| 4. Tag comprehension | Usability re-test | Stopwatch task-timing, same protocol as DISCOVER Phase 3 | Once, before Release 1 ships | product-owner |

### Hypothesis

We believe that a local-first, zero-account CLI (`bm save`/`bm find`/`bm share`)
for terminal-first infra/platform engineers will replace ad hoc scratch-file
and bash-script workarounds for capturing, retrieving, and sharing technical
reference links.
We will know this is true when 80%+ of pilot users' saves go through `bm save`
within 2 weeks, median time-to-locate drops to under 10 seconds, and 100% of
share recipients need zero install.

## Note on North-Star Data Availability

North-star and KPI-1 telemetry require **local usage logging** (opt-in) that
does not yet exist in the codebase (greenfield project, no `src/` yet). This
is flagged as an instrumentation dependency for the DEVOPS wave
(`platform-architect`), consistent with the Outcome KPI Framework's "Handoff
to DEVOPS" guidance — not a blocker for DISCUSS-wave completion, since the KPI
definitions themselves are complete and measurable-in-principle.
