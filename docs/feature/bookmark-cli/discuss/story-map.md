# Story Map: bookmark-cli

## User: PERSONA-TERMINAL-INFRA-ENG (terminal-first infra/platform/SRE/staff engineer)
## Goal: Capture a technical reference link without leaving the terminal, retrieve it later, and share it with a teammate who needs zero install.

## Backbone

| Capture a link | Locate a link | Share a link |
|---|---|---|
| Save link + tag in one command (US-01) | Search saved links by keyword/tag (US-02) | Generate a zero-install share snippet (US-03) |
| Discover the tag flag without failing (US-04) | See a helpful message on no match (US-06) | |
| Get a clear error on a bad/duplicate URL (US-05) | | |

---

### Walking Skeleton

- **Capture**: US-01 — `bm save <url> --tag <tag>` happy path
- **Locate**: US-02 — `bm find <query>` happy path
- **Share**: US-03 — `bm share <id>` happy path

This is the exact 3-task sequence already tested in DISCOVER Phase 3 solution
testing (93% task completion, n=5) — the walking skeleton formalizes that
tested flow into shippable stories rather than re-designing it.

### Release 1: Trustworthy under real conditions

Target outcome: the walking skeleton survives the conditions DISCOVER
evidence says actually occur (tag-flag confusion, bad input, no-match
searches), not just the clean happy path.

- US-04 — Discoverable tag syntax (remediates Aisha's tested failure, DISCOVER Phase 3)
- US-05 — Clear error on invalid/duplicate URL at save time
- US-06 — Helpful "no results" guidance at find time

Outcome KPI targeted: tag-flag comprehension time (see `outcome-kpis.md`).

### Backlog (out of scope this pass — see `wave-decisions.md` Out-of-scope)

JTBD-CONTEXT-SWITCH, JTBD-ORGANIZE, JTBD-RECALL-CONTEXT, JTBD-DEDUP-MONITOR —
carried forward in `docs/product/jobs.yaml` as `scope_status: out_of_scope_backlog`,
consistent with DISCOVER's own prioritization (`evaluate_later`/`backlog`/`deprioritize`).

## Priority Rationale

1. **Walking skeleton first** (US-01, US-02, US-03): validates the end-to-end
   flow works at all, and is the minimum needed to demo the core value
   proposition ("save, find, share — without leaving your terminal").
   Riskiest-assumption-first: this is also the exact flow DISCOVER already
   task-tested, so it carries the lowest re-validation risk of the six stories.
2. **US-04 next** (Release 1): the one concrete usability failure DISCOVER
   surfaced (Aisha, tag-flag discoverability, ~20s vs <10s target) is a known,
   named defect against the walking skeleton's own Capture step — fixing it
   before adding new surface area is higher value than any net-new feature.
3. **US-05 and US-06** (Release 1): both are error-path completions of stories
   already in the walking skeleton (happy-path bias check, per
   nw-po-review-dimensions Dimension 1) — without them, US-01/US-02 pass DoR
   item 4 (UAT scenarios) only for the happy path, which is a completeness gap.
4. Backlog jobs deferred: none scored above the top-3 pursue threshold in
   DISCOVER, and none are required for the walking skeleton to hold together.

Value x Urgency / Effort scoring (1-5 scale):

| Story | Value | Urgency | Effort | Score | Release |
|---|---|---|---|---|---|
| US-01 Capture happy path | 5 | 5 | 2 | 12.5 | Walking Skeleton |
| US-02 Locate happy path | 5 | 5 | 2 | 12.5 | Walking Skeleton |
| US-03 Share happy path | 5 | 4 | 2 | 10.0 | Walking Skeleton |
| US-04 Tag discoverability | 4 | 4 | 1 | 16.0 | Release 1 |
| US-05 Save error handling | 3 | 3 | 1 | 9.0 | Release 1 |
| US-06 Find no-results | 3 | 3 | 1 | 9.0 | Release 1 |
