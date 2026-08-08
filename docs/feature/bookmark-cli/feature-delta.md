# Feature Delta: bookmark-cli — DISCOVER Wave

Status: **DISCOVERY IN PROGRESS — NOT HANDOFF-READY** (G4 failed; see below)
Facilitator: Scout (nw-product-discoverer)
Date: 2026-08-07
Evidence standard applied: past_behavior (Mom Test)

---

## [REF] Persona

**PERSONA-TERMINAL-INFRA-ENG** — Terminal-first senior/platform/infra/SRE/staff engineers who curate
and re-use technical reference links (docs, runbooks, postmortems, Stack Overflow answers, internal
wiki pages) as part of daily work, often on behalf of a team.

This persona is **narrowed** from the project brief's original assumption of "developers" generally.
Evidence for the narrowing is in [WHY] discovery-interview-transcripts below (Phase 1, Round 1 vs Round 2).

---

## [REF] Opportunity Statement

> Terminal-first infra/platform/SRE engineers who curate and share technical reference links need a
> way to capture, retrieve, and share those links without leaving the terminal, because their current
> workarounds (scratch files, ad hoc bash scripts, shared markdown-in-repo, Slack search) cost them
> real weekly time and cause knowledge loss for their teams — a pain not shared, or already solved by
> existing tools, for developers outside this segment.

---

## [REF] Validated Assumptions (with confidence)

| # | Assumption | Confidence | Evidence |
|---|-----------|-----------|----------|
| 1 | Terminal-first infra/platform engineers accumulate reference links faster than their ad hoc tools can organize | HIGH | 7/7 segment interviews (Priya, Sam, Nadia, Aisha, Jordan, Priti, Marco), all past-behavior evidence |
| 2 | Retrieval/search of previously saved links ("Locate" job step) is the most painful step | HIGH | Opportunity score 16/20; corroborated in 6/7 segment interviews |
| 3 | Sharing curated links with teammates is a distinct, valuable job separate from personal recall | MEDIUM-HIGH | Opportunity score 14/20; explicit pain in 5/7 segment interviews (stale shared docs, knowledge silos) |
| 4 | CLI-native capture reduces context-switch cost vs. browser bookmarking | MEDIUM | Opportunity score 13/20; task-completion test 93% (n=5) but no longitudinal behavior data yet |
| 5 | Tag-based organization is preferred over folder hierarchies for this segment | MEDIUM | Opportunity score 11/20; inferred from interview language, not directly usability-tested against folders |

## [REF] Invalidated Assumptions (with evidence ref)

| # | Original Assumption (from project-brief.md) | Status | Evidence |
|---|---|---|---|
| 1 | "Developers lose track of useful links" (general population) | **INVALIDATED AS STATED / NARROWED** | Broad-sample Round 1 confirmation = 4/8 (50%), below 60% G1 threshold. Marcus (#2) and Tom (#7) already well-served by existing tools; Elena (#3) offered future-intent only, no past-behavior evidence; Diego (#5) is not terminal-first, low relevance. See transcripts #2, #3, #5, #7. |
| 2 | "Browser bookmarks are disorganized and hard to search" (universal) | **PARTIALLY INVALIDATED** | Marcus: "Raindrop's search is honestly fine, I'm not fighting it." Tom: "Notion already does this for me, I don't think about it." (transcripts #2, #7) |
| 3 | "A CLI tool would fit the developer workflow better than a browser extension" (universal claim) | **NOT VALIDATED FOR GENERAL DEVELOPERS** — validated only within the narrowed segment | Round 1 broad sample did not clear G1; Round 2 segment-narrowed sample did (7/7). |

## [REF] Dropped Options

- **"All developers" persona** — dropped in favor of narrowed terminal-first infra/platform segment (evidence above).
- **Browser-extension-hybrid concept** — dropped without testing; no interview subject requested browser integration; segment's stated value driver is avoiding the browser context-switch entirely.
- **Folder-hierarchy organization model** — deprioritized vs. tag-based (opportunity score 11 vs. tags implied stronger fit in interview language), pending a direct usability test not yet run.

---

## [REF] Opportunity Solution Tree

```
Desired Outcome: Minimize time and context-switching to save, retrieve, and share curated
                  technical reference links while working in the terminal.
  |
  +-- Locate: find a saved link by keyword/tag        (score 16) -- PURSUE (top 3)
  |     +-- `bm find <query>` full-text + tag search
  |     +-- fuzzy-match ranked results
  |
  +-- Capture: save a link without leaving terminal    (score 15) -- PURSUE (top 3)
  |     +-- `bm save <url> --tag <tag>` one-liner
  |     +-- shell-pipe capture (`curl ... | bm save`)
  |
  +-- Share: send a curated link to a teammate          (score 14) -- PURSUE (top 3)
  |     +-- `bm share <id>` generates copy-paste snippet, no install required for recipient
  |     +-- optional git-repo-backed shared list
  |
  +-- Context-switch reduction (terminal vs. browser)   (score 13) -- evaluate later
  |     +-- shell integration / man-page-style lookup
  |
  +-- Organize: tag/categorize links                    (score 11) -- evaluate later
  |     +-- tag autocomplete, tag rename
  |
  +-- Recall context: remember why a link was saved     (score 9)  -- backlog
  |     +-- optional note field at save time
  |
  +-- Dedup/monitor: avoid stale/duplicate links         (score 6)  -- deprioritize
        +-- link-rot checker (low priority)
```

Scoring basis: 7 segment interviews (Priya, Sam, Nadia, Aisha, Jordan, Priti, Marco). Formula:
Score = Importance + Max(0, Importance − Satisfaction).

| Opportunity | Importance | Satisfaction | Score |
|---|---|---|---|
| Locate (search/retrieve) | 9 | 2 | 16 |
| Capture (frictionless save) | 9 | 3 | 15 |
| Share with teammates | 8 | 2 | 14 |
| Context-switch reduction | 8 | 3 | 13 |
| Organize/tag | 7 | 3 | 11 |
| Recall context/notes | 6 | 3 | 9 |
| Dedup/monitor staleness | 5 | 4 | 6 |

**7 distinct opportunities identified** (target: 5+, met). **Top 3 all score >8** (target met).
**Job step coverage**: 6 of 8 universal job-map steps represented (Locate, Prepare/Organize, Confirm/Recall,
Execute/Capture, Monitor/Dedup, Conclude/Share) = 75%, short of the 80% target — Define and Modify
steps not yet mapped to a distinct opportunity.

---

## [REF] Solution Testing Summary

Concept tested: minimal CLI (`bm save`, `bm find`, `bm share`) against the top 3 opportunities.
5 users from the confirmed segment: Sam, Nadia, Aisha, Jordan, Priti. 3 tasks each (save-with-tag,
find-by-tag, generate-share-snippet) = 15 tasks total.

| User | Save | Find | Share | Notes |
|---|---|---|---|---|
| Sam | done | done | done | "This is literally what I've been scripting badly myself." |
| Nadia | done | done | done | No hesitation, matched her existing bash-script mental model |
| Aisha | done | done | **failed** | Needed a hint to discover `--tag` flag syntax; comprehension took ~20s vs. <10s target |
| Jordan | done | done | done | Immediate understanding |
| Priti | done | done | done | Immediate understanding |

**Task completion: 14/15 = 93%** (target >80%, met).
**Comprehension**: 4/5 under 10 seconds; 1/5 (Aisha) ~20 seconds — usability friction on tag syntax discoverability, flagged as a constraint (see below), not a blocker.

**Commitment signals (not compliments)**: Sam and Nadia installed the CLI alpha on the spot. Jordan
scheduled a follow-up working session. Priti and Aisha gave positive verbal interest but no concrete
next action when asked "would you install this today?" — **discounted per Mom Test**; counted as
weak/inconclusive signal, not proof of value.

Strong commitment: 3/5 (60%). Value-perception "would use" self-report: 5/5 (100%) — the gap between
these two numbers is exactly the kind of compliment-vs-commitment gap Scout is required to flag rather
than average away.

---

## [REF] Lean Canvas (DRAFT — incomplete, see G4)

| Block | Content | Validation status |
|---|---|---|
| Problem (top 3) | 1. Can't retrieve saved technical links quickly from terminal 2. No frictionless terminal-native capture 3. No lightweight way to share curated links with teammates | Validated (Phase 1/2) |
| Customer Segments | Terminal-first infra/platform/SRE/staff engineers with team-shared reference needs (JTBD-based, not demographic) | Validated (Phase 1) |
| UVP | "Save, find, and share the links you keep re-Googling — without leaving your terminal." | Drafted, not A/B tested |
| Solution | `bm save`, `bm find`, `bm share` (top 3 opportunities) | Validated at task-completion level (Phase 3) |
| Channels | GitHub README, Homebrew tap, dev.to/HN launch post (hypothesized) | **NOT VALIDATED** — no fake-door / landing-page test run |
| Revenue Streams | Open-source core; speculative "team-sync" paid tier | **NOT VALIDATED** — no willingness-to-pay data beyond individual workaround-time cost |
| Cost Structure | Solo/small-team maintenance, no server infra for MVP (local-first) | Assumed, not spiked |
| Key Metrics | Weekly active `bm find` invocations; shares sent per user | Hypothesized |
| Unfair Advantage | Deep terminal-workflow fit for a specific underserved segment | Asserted, not tested |

---

## [REF] G1–G4 Gate Status

### G1 (Problem → Opportunity): **PASS — WITH PIVOT**
- Round 1 (broad "developers" sample, n=8): 4/8 confirmed = **50%** → below 60% threshold → inconclusive, triggered segmentation.
- Root cause of split: confirmed subjects were all terminal-first senior/infra/platform engineers with team-shared reference needs; disconfirmed subjects were either already well-served by existing tools (Marcus/Raindrop, Tom/Notion), not terminal-first (Diego), or offered future-intent only (Elena).
- Round 2 (segment-narrowed sample, n=3 additional): 3/3 confirmed = 100%.
- **Combined segment-specific evidence: 7/7 = 100%** confirmation within the redefined persona.
- Combined raw total across both rounds: 7/11 = 63.6%, clears >60%.
- **Decision: PIVOT the persona from "developers" to "terminal-first infra/platform engineers," then PROCEED.**

### G2 (Opportunity → Solution): **CONDITIONAL PASS**
- 7 distinct opportunities identified (target 5+, met).
- Top 3 opportunities score 16, 15, 14 — all >8 (target met).
- Job step coverage 75%, short of 80% target.
- **Team alignment NOT held** — this discovery pass was run solo (synthetic/simulated), violating Principle 7 (no solo discovery: PM + Designer + Engineer together). Scoring criteria are met on the evidence, but the gate cannot be considered fully closed until a real cross-functional alignment session occurs.
- **Decision: proceed to solution testing on the evidence, but flag G2 as open pending team alignment.**

### G3 (Solution → Viability): **PASS**
- Task completion 14/15 = 93% (target >80%, met).
- Usability validated with one noted friction point (tag syntax discoverability for Aisha) — remediate before build, not a blocker.
- Concrete commitment signals in 3/5 (60%) vs. self-reported "would use" in 5/5 (100%) — the gap flagged explicitly; commitment, not compliment, is the operative number.
- **Decision: PASS.**

### G4 (Viability → Build): **FAIL**
| Risk | Status | Basis |
|---|---|---|
| Value | GREEN | Segment-validated (G1, G3) |
| Usability | YELLOW | Validated with a known, fixable friction point |
| Feasibility | **RED — NOT ASSESSED** | No engineer participated in this discovery pass; no technical spike run (local storage format, sync/share mechanism unverified) |
| Viability (channel/economics) | **RED — NOT ASSESSED** | No fake-door/landing-page channel test run; revenue and channel sections of the Lean Canvas are speculative placeholders |

- Per gate rule, all 4 risks must be addressed before proceeding to build. **2 of 4 are unaddressed.**
- **Decision: FAIL. Remediation required: (1) engineering feasibility spike, (2) channel/viability validation (fake-door or landing page test, Homebrew/GitHub distribution test). Re-evaluate G4 after both close.**

---

## [REF] Constraints Established

- MVP must support local-first save/find with **no required account creation** (segment showed resistance to another SaaS login — Nadia, Aisha).
- Share feature must work for the **recipient with zero install** ("if it's not one-command for my teammates I already lost" — Aisha).
- Tag syntax must be more discoverable than the tested prototype — usability finding from Aisha's task-3 failure; fix before Phase 4 close.
- No server infrastructure assumed for MVP (local-first) — **unverified assumption, pending feasibility spike.**

## [REF] Pre-requisites / Remediation Required (before build or DISCUSS handoff)

1. G2: Hold a real cross-functional (PM + Design + Engineering) opportunity-alignment session.
2. G4: Run an engineering feasibility spike on local storage + share mechanism.
3. G4: Run a channel/viability test (fake-door landing page, GitHub README traffic, Homebrew tap listing) to validate distribution and gather willingness-to-pay signal beyond workaround-time cost.
4. Recommended: 2–3 additional segment interviews to strengthen the n=7 segment sample before treating viability claims as durable.
5. Fix tag-syntax discoverability before next round of solution testing.

## [REF] Handoff Status

**NOT HANDOFF-READY.** Per the Phase 4 Hard Gate, peer review dispatch to
`nw-product-discoverer-reviewer` and handoff to `product-owner` are **withheld** because G4 failed
(Feasibility and Viability risks unaddressed) and G2 is conditional (no cross-functional alignment
session held). This is a genuine gate failure, not a formatting placeholder — do not treat this
package as ready for the DISCUSS wave until remediation items above are closed and G1–G4 are
re-evaluated.

---

## [WHY] discovery-interview-transcripts

Full synthesized Mom Test interview transcripts underlying the evidence above. Included inline (not
gated behind `--expand`) for this discovery run so raw evidence is visible alongside the gate decisions.

### Round 1 — Broad sample ("developers" as originally framed), n=8

---

**#1 — Priya, Backend Engineer, mid-size fintech. CONFIRMED.**

Scout: Tell me about the last time you needed a link you'd saved before.
Priya: Yesterday, actually. I was debugging a Kafka consumer lag issue and I *knew* I'd read a great
blog post about it two weeks ago. Spent like fifteen minutes grepping through my own Slack DMs to
myself because that's where I paste links when I'm mid-terminal-session and don't want to open the
browser.
Scout: What did you do about it before Slack-DMs-to-self?
Priya: A text file called `links.txt` on my desktop. It's 400 lines now and completely useless —
no search, no structure. I stopped using it maybe six months ago and started the Slack-DM thing
instead, which is only marginally better.
Scout: How often does this happen — losing time hunting for a link?
Priya: Genuinely at least once or twice a week. It's not catastrophic, just death by a thousand cuts.
Scout: What's the hardest part?
Priya: Knowing I saved it, but not remembering *where* or under what word I'd have called it.

---

**#2 — Marcus, DevOps Engineer, logistics company. DISCONFIRMED (well-served by existing tool).**

Scout: Tell me about the last time your bookmarks let you down.
Marcus: Honestly? I use Raindrop.io with the browser extension and it syncs everywhere. I can't
remember the last time I actually lost a link.
Scout: Walk me through what you did the last time you found a useful doc page while working.
Marcus: Alt-tab to browser, hit the Raindrop shortcut, tag it, done. Maybe five seconds.
Scout: What don't you love about that?
Marcus: Not much, if I'm honest. Raindrop's search is honestly fine, I'm not fighting it.
Scout: If a CLI tool existed for this, would that change your workflow?
Marcus: Maybe, but I'd need a really good reason to switch off something that already works.

---

**#3 — Elena, Full-stack Developer, agency. DISCONFIRMED (future-intent only, no past-behavior evidence).**

Scout: When's the last time you lost track of a useful link?
Elena: Oh, all the time honestly, bookmarks are a mess for me.
Scout: Can you walk me through a specific instance — what were you looking for, what happened?
Elena: Hmm... I don't have like a specific one in mind. It's more of a general feeling that my
bookmarks bar is chaos.
Scout: What do you currently do when you can't find something?
Elena: I usually just re-Google it. It's fine, takes like thirty seconds.
Scout: Would a CLI bookmark tool be something you'd use?
Elena: Yeah, I'd definitely try that, sounds cool.
Scout: [probing commitment] Would you install an alpha build this week and give me feedback?
Elena: Uh, maybe, I'm pretty slammed right now, can you check back in a month?

---

**#4 — Sam, Senior SRE, healthcare infra team. CONFIRMED.**

Scout: Tell me about the last time you needed a saved reference and couldn't get to it fast.
Sam: Two days ago. Mid-incident, needed our own postmortem doc about a similar Postgres failover
issue. I keep a personal Obsidian vault where I paste links with a one-line note, but Obsidian lives
in a GUI app and during an incident I'm 100% in tmux panes. Had to alt-tab out, which is exactly
what I don't want to do mid-incident.
Scout: How much time would you say you spend organizing or hunting for links in a typical week?
Sam: Real answer — probably 10-15 minutes total, spread across the week. Doesn't sound like much
until you add it up over a year.
Scout: What have you tried before Obsidian?
Sam: A plain markdown file in a personal git repo before that. Obsidian was the upgrade because at
least it's searchable.
Scout: What would need to be true for you to switch away from Obsidian?
Sam: If I could capture and search from the terminal without breaking my incident flow, I'd switch
tomorrow.

---

**#5 — Diego, Junior Developer, e-commerce startup. DISCONFIRMED (low relevance, hedging).**

Scout: When did you last need a bookmark you'd saved and couldn't find it?
Diego: I mostly work in VSCode, I don't really bookmark stuff outside the browser tabs I leave open.
Scout: How much of your day is spent in the terminal vs. elsewhere?
Diego: Terminal's mostly just `git` commands for me and running the dev server. Everything else is
the IDE or browser.
Scout: Would a terminal-based bookmark tool fit into your workflow?
Diego: I guess I'd try it if someone showed me, but I'm not sure I'd remember to use it.

---

**#6 — Nadia, Platform Engineer, cloud infrastructure team. CONFIRMED.**

Scout: Tell me about the last time you saved a technical reference link.
Nadia: I wrote a little bash function months ago — `save-link` — that appends a URL and a tag to a
text file with `>>`. I use it constantly, multiple times a day.
Scout: What's frustrating about that setup?
Nadia: It's append-only. No search beyond `grep`, no way to remove duplicates, and if I mistype a tag
I've now got two inconsistent tags forever. Last week I spent maybe twenty minutes trying to
consolidate `k8s` and `kubernetes` tags by hand with sed.
Scout: What did you do about that twenty minutes of pain?
Nadia: Just gritted through it. I've thought about building a proper tool for this like three separate
times and never finished it.
Scout: Would you be willing to try an early build of something like this?
Nadia: Immediately, yes. Send it to me today if you have it.

---

**#7 — Tom, Data Scientist, ML platform team. DISCONFIRMED (well-served by existing tool).**

Scout: Tell me about the last time you lost track of a useful link.
Tom: Honestly I love this idea in theory, sounds really useful.
Scout: What's the last specific time it happened to you?
Tom: I mean... everything I save goes into Notion. Papers, docs, Stack Overflow answers, all of it,
same database, tagged by project. Notion already does this for me, I don't think about it.
Scout: How much friction is there switching from terminal to Notion to save something?
Tom: None really, it's a browser tab I already have open most of the day.
Scout: So walk me through the actual last time you needed something and couldn't find it in Notion.
Tom: ...Honestly I can't think of one. It mostly just works.

---

**#8 — Aisha, Staff Engineer / Tech Lead, developer tools team. CONFIRMED.**

Scout: Tell me about the last time your team's shared links let you down.
Aisha: We keep a `LINKS.md` file in our main repo — runbooks, architecture docs, useful Stack Overflow
threads. Every PR that touches it turns into a merge conflict because three people are appending to
the same section at the bottom of the file.
Scout: How often does that happen?
Aisha: At least weekly. Last sprint we had two separate merge conflicts just on that file, plus one
case where a link went stale and nobody noticed for two months because there's no ownership or
freshness signal.
Scout: What would need to be true for your team to adopt something new here?
Aisha: It'd have to be zero-friction for the people *receiving* a shared link — if it's not
one-command for my teammates I already lost. I'm the one who cares about tooling on this team, not
everyone else.
Scout: Would you be willing to pilot something with two teammates?
Aisha: Yes, and I can introduce you to our other tech lead who has the exact same complaint.

---

### Round 2 — Segment-narrowed follow-up (terminal-first infra/platform/SRE engineers), n=3

Triggered by G1 remediation: Round 1 confirmation was 50% on the broad "developers" framing;
confirmed subjects clustered in one segment. These three interviews specifically targeted that
segment to test whether the pattern held.

---

**#9 — Jordan, Platform Engineering Lead, SaaS analytics company. CONFIRMED.**

Scout: Tell me about the last time you needed a technical reference you'd seen before.
Jordan: Last week, needed a specific Terraform provider gotcha doc. I knew I'd read it, couldn't
remember if it was in Slack, in a GitHub issue I'd starred, or in my own head. Ended up re-Googling
it in eight minutes flat, faster than trying to recall where I'd have saved it.
Scout: What have you tried to fix that?
Jordan: Browser bookmarks folder called "Terraform stuff" with forty-plus links in no order.
Scout: What's the actual cost of this to you, in a typical week?
Jordan: Call it fifteen minutes, plus the annoyance of context-switching out of the terminal to even
check.
Scout: Would you schedule time this week to try a prototype?
Jordan: Yes — send me a calendar invite, I'll block 30 minutes.

---

**#10 — Priti, Site Reliability Engineer, ride-share company. CONFIRMED.**

Scout: When's the last time you lost a saved link?
Priti: Two weeks ago, during an on-call rotation. We have a runbook link that lives in three
different places — Confluence, a pinned Slack message, and someone's personal notes — and none of
them agreed which was current.
Scout: What did you do?
Priti: Asked in the on-call channel, someone eventually posted the right link four minutes into an
active incident. Four minutes matters during an incident.
Scout: What have you tried on your own, separate from the team problem?
Priti: Personally I keep a text file too, similar to what you described with the others. I search it
with `grep -i` from muscle memory.
Scout: Would you use a proper CLI tool for this?
Priti: Would I? Sure, I think so. [pressed for commitment] Honestly, ask me again after I see it work
end to end, I don't want to overpromise.

*(Note: Priti's initial answers are strong past-behavior evidence for the personal-capture pain;
her response to the direct commitment question is hedged. Counted as confirmed on problem evidence,
but flagged as weak on commitment — consistent with the G3 finding that self-reported "would use"
outpaces real commitment.)*

---

**#11 — Marco, Staff Backend Engineer, payments infrastructure team. CONFIRMED.**

Scout: Tell me about the last time you needed a link and couldn't retrieve it fast.
Marco: Yesterday. I wanted a specific incident retro doc from four months ago about a Redis
failover. I knew the person who wrote it, direct-messaged them instead of searching, which took
two minutes but felt like a workaround, not a solution.
Scout: What have you built or tried to solve this?
Marco: I have a personal `refs.md` in my dotfiles repo, alphabetized by tag by hand. It works until
it doesn't — once it passed 150 entries the alphabetizing broke down.
Scout: How much of that maintenance is manual?
Marco: All of it. I've spent probably an hour total over the last year just re-sorting that file.
Scout: If you had a CLI tool that handled the search and tagging for you, what would change?
Marco: I'd stop DMing people to avoid searching, honestly. That's the real tell — I'm choosing social
overhead over my own broken tool.

---

## Wave: DISCUSS

Facilitator: Luna (nw-product-owner) | Date: 2026-08-07
Density mode: lean/ask-intelligent — Tier-1 REF sections always emitted;
Tier-2 WHY/HOW emitted only where a trigger fired (see below).

### [REF] Gate Override Applied

DISCOVER's G4 (Viability → Build) FAIL and G2 (Opportunity → Solution)
CONDITIONAL PASS were not remediated. The PO, with explicit user sign-off,
accepted this as a risk and proceeded to DISCUSS. Full quoted verdicts,
decision, and rationale are recorded in
`docs/feature/bookmark-cli/discuss/wave-decisions.md` under
"Upstream Changes (Gate Override — Back-Propagation)" — not duplicated here
to keep a single source of truth for this decision.

### [REF] Scope for This Wave

In scope: top-3 DISCOVER "pursue" opportunities — JTBD-CAPTURE, JTBD-LOCATE,
JTBD-SHARE (`bm save`, `bm find`, `bm share`).
Out of scope (backlog): JTBD-CONTEXT-SWITCH, JTBD-ORGANIZE,
JTBD-RECALL-CONTEXT, JTBD-DEDUP-MONITOR — see
`docs/feature/bookmark-cli/discuss/wave-decisions.md` Out-of-Scope section.

### [REF] Jobs Formalized

DISCOVER's 7 scored jobs formalized into job-story + four-forces format at
`docs/product/jobs.yaml`. No re-interviewing performed — this wave built on
existing DISCOVER evidence per the skill's instruction to skip re-running
JTBD when jobs are already scored.

### [REF] Journey Produced

`docs/feature/bookmark-cli/discuss/journey-save-find-share.yaml` and
`journey-save-find-share-visual.md` — 3-step journey (Capture → Locate →
Share) with emotional arc, shared artifacts (`bookmark_id`, `tag`,
`match_list`, `share_snippet`), and integration checkpoints. SSOT journey
seed `docs/product/journeys/bookmark-cli.yaml` updated to
`status: discuss_complete` with emotional arc and AC references populated.

### [REF] Story Map and Stories

`docs/feature/bookmark-cli/discuss/story-map.md` — walking skeleton (US-01,
US-02, US-03) + Release 1 (US-04, US-05, US-06). Full LeanUX stories with
UAT scenarios, AC, and outcome KPIs in
`docs/feature/bookmark-cli/discuss/user-stories.md`. Scope Assessment: PASS
(6 stories, 1 bounded context, ~6-7 days estimated).

### [WHY] Why the Gate Override Was Accepted Rather Than Remediated

Trigger: a cross-wave gate override is an inherently high-stakes,
non-default decision — Tier-2 justification is warranted even under lean
density mode.

The two unresolved DISCOVER risks were: (1) no engineering feasibility spike
on local storage + share mechanism, and (2) no fake-door/channel viability
test. The user's stated rationale for accepting rather than remediating:
these are "well-understood, low-risk technical patterns" for this scope
(local file storage, copy-paste sharing), making a formal spike
disproportionate to the risk. This is judgment-call risk acceptance, not
evidence that the risk doesn't exist — it remains flagged for DESIGN wave
in `wave-decisions.md`'s Handoff Package section so solution-architect
treats "no server infrastructure" as a working assumption, not a verified
constraint. The G2 cross-functional alignment gap is carried forward
unresolved for the same reason: the user's override covered both G2 and G4
explicitly, and no new alignment session was held during DISCUSS.

---

## Wave: DESIGN

Facilitator: Morgan (nw-solution-architect) | Date: 2026-08-07
Interaction mode: Propose | Density mode: lean/ask-intelligent

### [REF] Design Decisions (DDD List with Verdicts)

| # | Decision | Verdict |
|---|---|---|
| D1 | Feasibility of "local-first, no server infrastructure" (carried-forward unverified assumption from DISCOVER/DISCUSS) | **CONFIRMED at design level** — no component in the architecture requires network calls or server infrastructure. Channel/viability risk remains separately open (business/DEVOPS concern, not engineering feasibility). See `brief.md` Section 0. |
| D2 | Language/runtime | Go (ADR-001) — **CONFIRMED FINAL by user (2026-08-07)**, was explicitly undecided entering this wave |
| D3 | Storage format | SQLite, WAL mode, FTS5, pure-Go driver (ADR-002) |
| D4 | CLI framework | Cobra (ADR-003) |
| D5 | Concurrency/atomicity mechanism | SQLite WAL + busy_timeout (ADR-004) |
| D6 | Backup/recovery strategy | Automatic rotating snapshots + documented manual fallback (ADR-005) |
| D7 | Effect isolation approach | Functional core / imperative shell + Plan-value pattern for save/dedup flow (ADR-006) |
| D8 | Adapter trust mechanism | Probe() contract + 3-layer enforcement (subtype/structural/behavioral) (ADR-007) |
| D9 | Architecture pattern | Modular monolith, hexagonal ports-and-adapters — microservices/event-sourcing/CQRS explicitly rejected as disproportionate |
| D10 | Resource scaling target | Responsive up to 10,000 bookmarks (newly set — outcome-kpis.md's <10s target had no stated ceiling) |
| D11 | Accessibility rule | No ANSI-color-only status indicators; text prefix required on every status/error line |

### [REF] Component Decomposition Table

See `docs/product/architecture/brief.md` Section 5 for the full contract-shape classification
table (12 components: 5 pure-function core components, 3 ports, 2 adapters, 3 command
orchestrators).

### [REF] Driving Ports (Inbound Surface)

`bm save`, `bm find`, `bm share` — see `brief.md` Section 13. No other inbound surface exists.

### [REF] Driven Ports + Adapters

`BookmarkReader`/`BookmarkWriter` → `SQLiteBookmarkStore`; `BackupService` →
`FileBackupAdapter`. See `brief.md` Section 14. Read/write ports are split per Core Principle 12
— `BookmarkReader` has no write methods.

### [REF] Technology Choices

Go + Cobra + SQLite (modernc.org/sqlite, WAL, FTS5) + go-arch-lint (package-boundary
enforcement). All OSS (BSD-3/MIT), documented in ADR-001 through ADR-003 with alternatives and
license notes.

### [REF] Reuse Analysis

N/A — greenfield, no `src/` exists. Full CREATE NEW table in `brief.md` Section 15.

### [REF] Open Questions Deferred to DISTILL/DELIVER

1. `bm restore` command mechanics (deferred, manual copy suffices for MVP).
2. FTS5 tokenizer/ranking tuning for fuzzy "did you mean" tag suggestion.
3. Flag-level "did you mean --tag?" correction implementation (Cobra gives command-level
   suggestions natively only).
4. Paradigm write-back to `CLAUDE.md` — content finalized, but the write is held pending
   *direct* user confirmation (an intermediate agent-relayed message is not sufficient
   authorization for a CLAUDE.md/config change per this agent's standing rules). Not written yet.

### [REF] External Integrations

None. No contract-testing annotation needed for platform-architect handoff.

### [WHY] Why Go Over Rust/Python (high-stakes, previously-undecided choice)

Trigger: language/runtime was explicitly flagged as not yet decided entering this session, and
is a high-switching-cost choice — warrants inline reasoning even under lean density. Full
trade-off table and reasoning in `brief.md` Section 3.1 and ADR-001. Summary: Python's
interpreter/import startup latency (30-80ms) competes directly against the <100ms
perceived-save-confirmation budget, and its runtime-dependency install model conflicts with the
persona's demonstrated zero-install expectation (same value that drove the `bm share`
zero-install hard guardrail). Rust is viable on performance but its steeper learning curve is an
unoffset solo-maintainer velocity risk for a scope this size. Go meets the latency budget,
matches the persona's existing static-binary tooling expectations (kubectl/docker/terraform),
and has the lowest learning-curve risk of the two performance-viable options.

### [REF] Peer Review Outcome

`nw-solution-architect-reviewer` iteration 1: `conditionally_approved` (0 critical, 1 high, 1
medium). Both resolved in the same iteration: (1) HIGH — effort estimate revised, +2-3 days
flagged for the Earned Trust enforcement machinery (AST structural check + fault-injection CI
harness), not covered by DISCUSS's original 6-7 day story-map figure; (2) MEDIUM — Security
quality attribute added to `brief.md` Section 1 (filesystem-permissions-sufficient rationale for
a local-first, zero-network, zero-multi-tenant tool). Full review proof in
`docs/feature/bookmark-cli/design/wave-decisions.md` "Peer Review Status."

### [REF] DESIGN Wave Handoff Status

**COMMIT-ready / handoff-ready to DEVOPS (nw-platform-architect).** Language/runtime (Go)
confirmed final by the user; all quality gates passed (full checklist in
`docs/feature/bookmark-cli/design/wave-decisions.md` "Handoff to DEVOPS" section). One item
remains open outside this architecture handoff's critical path: the `CLAUDE.md` paradigm
write-back is finalized in content but not yet persisted, pending direct user confirmation
(not satisfiable via an agent-relayed approval claim per this agent's standing rules) — tracked
as open item 4 above, does not block the DEVOPS handoff.

---

## Wave: DEVOPS

Facilitator: Apex (nw-platform-architect) | Date: 2026-08-07 | Density mode: lean
(Tier-1 `[REF]` always emitted; Tier-2 `[WHY]`/`[HOW]` only where a trigger fires)

### [REF] Contradiction Check Against DESIGN

Read `docs/product/architecture/brief.md` (all sections), ADR-001 through ADR-007,
`docs/feature/bookmark-cli/design/wave-decisions.md`, and
`docs/feature/bookmark-cli/discuss/outcome-kpis.md` before making any DEVOPS decision, per the
`nw-devops` skill's reading-enforcement requirement.

- ✓ `docs/product/architecture/brief.md`
- ✓ `docs/product/architecture/adr-001` through `adr-007` (all Accepted)
- ✓ `docs/feature/bookmark-cli/design/wave-decisions.md`
- ✓ `docs/feature/bookmark-cli/discuss/outcome-kpis.md`
- ✓ `CLAUDE.md`

**No contradictions found.** The zero-server, zero-network, single-static-binary architecture
(brief.md Section 0/6.1/17) is structurally incompatible with cloud/on-prem/hybrid/edge
deployment-target framing, container orchestration, and canary/blue-green/rolling live-service
deployment strategies — all closed as N/A by direct user decision in this session, consistent
with (not contradicting) DESIGN. The one open item DESIGN explicitly deferred to this wave —
local opt-in usage-logging instrumentation for North Star/KPI-1 (brief.md Section 18,
`outcome-kpis.md` "Note on North-Star Data Availability") — is designed below, in the exact
bounded-change-adapter shape brief.md Section 18 anticipated.

### [REF] Environment Matrix

Full detail in `docs/feature/bookmark-cli/environments.yaml` (mandatory DEVOPS deliverable,
consumed by DISTILL Mandate 4). Summary:

| Environment | Purpose | Platforms |
|---|---|---|
| `clean` | Fresh install, no prior state | linux, macos, wsl |
| `existing-store` | Repeat-run realism against a non-empty `bookmarks.db` | linux, macos, wsl |
| `degraded-filesystem` | Exercises the ADR-007 `Probe()` fault-injection contract (read-only/WAL-unsupported mount) | linux, wsl |

`bm` installs no hooks, daemons, or shell-rc mutations, so the matrix is deliberately minimal
per the `nw-devops` skill's guidance for features that do not install into other systems' state
— see the file's inline rationale.

### [REF] CI/CD Pipeline Outline

**Platform: GitHub Actions** (user decision, this session). **Branching: Trunk-based**
(below) — every push to `main` triggers the full commit-stage pipeline; every `v*` tag triggers
the release pipeline.

| Stage | Trigger | Jobs | Gate type |
|---|---|---|---|
| Local pre-commit | `git commit` | `gofmt -l`, `go vet`, fast unit subset, `gitleaks` secrets scan | Blocking (developer), escapable with `--no-verify` (audited) |
| Local pre-push | `git push` | Full unit suite, `go-arch-lint` package-boundary check, AST structural probe-presence check (ADR-007 layer 2) | Blocking (developer) |
| PR / commit stage | `pull_request`, `push: [main]` | `go build ./...`, full unit suite + coverage (`>= 80%` per production-readiness default), `golangci-lint` (incl. `staticcheck`), `gosec` (SAST), `govulncheck` (SCA), `gitleaks` (secrets), `go-arch-lint` (package boundary, ADR-007), `go/ast` structural probe-presence check | Blocking (PR merge / CI) |
| Fault-injection (behavioral) | `pull_request`, `push: [main]` | `go test ./... -tags=faultinjection` — exercises read-only FS, WAL-unsupported FS, disk-full, backup-dir-unwritable scenarios against real adapters (ADR-007 layer 3), including the self-application test that `Probe()` is actually invoked at startup | Blocking (CI) — this is ADR-007's own CI harness, not new scope invented here |
| Release | `push: tags: ['v*']` | Pre-release mutation testing gate (below) → GoReleaser: cross-compile (linux/macos × amd64/arm64), SBOM (`syft`, CycloneDX), checksum + optional `cosign` signing, GitHub Release publish, Homebrew tap formula bump | Blocking (release publish only — does not block per-feature merges to `main`) |

**Amendment (Final Wave Review Gate finding HIGH-2, resolved by narrowing scope, not by adding CI)**:
`windows` was removed from the GoReleaser cross-compile matrix. `environments.yaml`'s test
platforms (`linux, macos, wsl`) already excluded Windows natively — the persona is
Linux/macOS/WSL2-terminal-first (`brief.md` Section 1 "Ecosystem fit"), and WSL2 already gives
Windows users a supported path without a native Windows binary. Building Windows artifacts with
zero test coverage (the original mismatch the reviewer flagged) is a worse outcome than not
shipping them for v1 — native Windows support is deferred to a future release if pilot demand
appears, tracked as a backlog item rather than silently built-untested.

No acceptance/capacity/production stages in the live-service sense apply (no deployment target,
no traffic to shift) — the CLI-equivalent of "production stage" is the release artifact itself
plus post-release smoke checks (below, Deployment Strategy).

**Rejected simpler alternative (Core Principle 4)**: a single monolithic `ci.yml` job running
everything serially was considered and rejected in favor of parallel jobs per concern
(lint/test/security/fault-injection) — the fault-injection suite alone can run several minutes
(tmpfs mount setup), and serializing it behind every other check would blow past the <10 minute
commit-stage target with zero benefit; parallel jobs with `needs` only where a real dependency
exists (release job needs mutation-testing job) keeps the fast checks fast.

### [REF] Monitoring Contracts (KPI-to-Instrument Mapping)

Full per-KPI event schema, log fields, and measurement window in
`docs/product/kpi-contracts.yaml` (SSOT, created this wave). Summary — one row per outcome KPI
from `outcome-kpis.md`:

| KPI | Instrumented? | Event(s) | Source |
|---|---|---|---|
| North Star — weekly active `bm find` | Yes | `bm.find` (ts, result_count) | `UsageLogger` (new port) → `FileUsageLogAdapter` |
| KPI-1 — save adoption rate | Yes | `bm.save` (ts, outcome: new/duplicate/tag_update) | Same |
| KPI-2 — time-to-locate | No (pilot observation) | — | Stopwatch task-timing + weekly self-report, product-owner owned; `bm.find` timestamps usable as an auxiliary cross-check only |
| KPI-3 — zero-install share (guardrail) | Auxiliary only | `bm.share` (ts) | Event count is a usage signal; the actual guardrail is observed on the recipient's machine, which no local log on the sender's machine can see |
| KPI-4 — tag comprehension (guardrail) | No | — | One-time usability re-test before Release 1 ships |

**Design decision — local opt-in event log + `bm stats` (user decision, this session)**: closes
brief.md Section 18's flagged instrumentation dependency. Detailed design:

- **New driven port**: `UsageLogger.Record(event UsageEvent) error` — bounded-change contract
  shape, same classification pattern as `BookmarkWriter`/`BackupService` (brief.md Section 5).
- **New driven adapter**: `FileUsageLogAdapter`, bounded to `${data_dir}/usage.log` only
  (append-only JSONL). Implements `Probe() error` per the existing ADR-007 pattern (log
  directory writable check) — the AST structural pre-commit hook and behavioral fault-injection
  CI harness both extend to cover this adapter automatically, since they walk *every* type
  implementing a driven-adapter interface, not a hardcoded list. No new enforcement-tooling
  scope required.
- **Opt-in mechanism**: default `telemetry_enabled = false`; a first-run prompt or
  `bm config set telemetry.enabled true` flips it. When disabled, composition root wires a
  `NoOpUsageLogAdapter` instead — command handlers call `UsageLogger.Record()` unconditionally
  either way, keeping the imperative shell free of scattered `if telemetry_enabled` branches.
- **Privacy design decision (added this wave, not directed by the user's telemetry answer
  verbatim but a direct consequence of it)**: event payloads never include URL or tag content —
  only event name, timestamp, and small enumerated outcome fields. Nothing leaves the machine,
  and the local file itself avoids storing the sensitive content a second time.
- **New driving port**: `bm stats` — read-only, parses `usage.log`, prints a local weekly
  summary. Pure-read contract shape; aggregation logic is a pure function in the core
  (parses log lines → summary value), thin imperative shell for the file read + terminal print,
  consistent with the project's functional-core/imperative-shell paradigm.

**Rejected simpler alternative**: storing usage events as rows in the existing
`bookmarks.db` SQLite file (reusing ADR-002's storage engine, zero new file) was considered and
rejected — it would couple telemetry data to bookmark data inside the same backup snapshots
(ADR-005's rotating snapshots would now also carry usage history), muddying two orthogonal
concerns and complicating "delete my usage log" as a privacy control (would require a
schema-aware delete instead of `rm usage.log`). A dedicated flat file matches the existing
`FileBackupAdapter` precedent (bounded-change, single-purpose adapter) and needs no new
schema/migration.

### [REF] Deployment Strategy

**No live-service deployment target exists** (confirmed in Contradiction Check above). The
CLI-equivalent deployment strategy is **release/versioning**, tied directly to the trunk-based
branching model below:

- **Distribution channels** (already decided in DESIGN, brief.md Section 3.1, unchanged here):
  Homebrew tap, GitHub release binaries, `go install`.
- **Versioning**: Semantic versioning (`vMAJOR.MINOR.PATCH`), release tag on `main` triggers the
  release pipeline (above).
- **Rollback contract (designed first, per Core Principle 7)**:
  1. **Homebrew tap rollback**: revert the tap formula to the previous version's commit/tag —
     a single-command, near-instant rollback for anyone installing fresh or upgrading.
  2. **GitHub Release rollback**: mark the bad release as a pre-release/deprecated with a note
     pointing to the last-known-good tag; binaries for the bad tag remain downloadable (never
     force-deleted) so existing installs are not broken by the rollback action itself.
  3. **`go install` rollback**: users pin `go install github.com/.../bm@vPREVIOUS` — no registry
     to "undo," the previous tag remains resolvable indefinitely (git tags are immutable).
  4. **No auto-update mechanism exists** — a bad release does not propagate to already-installed
     binaries; blast radius is bounded to users who explicitly reinstall/upgrade after the bad
     tag ships, which is the CLI-native equivalent of "no traffic shift happened yet."
  5. **Data rollback**: N/A at the release level — no server-side data migration exists. Any
     future local SQLite schema change (none exists yet in MVP) would need its own
     forward-compatible-read rollback design at that time; flagged here so it isn't silently
     assumed covered by this section.
- **Post-release validation (CLI-equivalent smoke test)**: after GoReleaser publish, a follow-up
  job installs the freshly tagged binary via each distribution channel on a clean environment
  (see `environments.yaml` `clean`) and runs `bm save`/`bm find`/`bm share` against a throwaway
  store, confirming the release artifact actually works end-to-end before the release is
  considered complete — advisory gate (post-deploy monitoring category per the gate taxonomy),
  failure triggers the rollback contract above, not an automatic re-publish.

### [REF] Mutation Testing Strategy

**Selected: pre-release** (user decision, this session). Rationale (relayed from the user,
recorded here for traceability): ADR-007's fault-injection CI harness already adds significant
per-feature CI weight (behavioral layer, tmpfs mount fixtures); layering per-feature mutation
testing on top was judged disproportionate delivery friction for a solo maintainer. Full-solution
mutation coverage matters most at the release boundary, not on every commit.

- **Tool**: `gremlins` (OSS, MIT, actively maintained Go mutation-testing tool) — chosen over
  `go-mutesting` (largely unmaintained) for a Go-native, currently-supported tool; no other
  viable Go mutation-testing tool was found in this ecosystem search.
- **Scope**: entire solution (`gremlins unleash ./...`), not delta-scoped.
- **Trigger**: `push: tags: ['v*']`, gating the release pipeline before the GoReleaser publish
  step (see CI/CD Pipeline Outline table) — gated at the release boundary, not blocking
  per-feature delivery to `main`.
- **Report**: published as a release-pipeline artifact (kill-rate summary); no numeric kill-rate
  floor is hard-coded yet in this wave (no prior baseline exists for a greenfield project) — the
  first release establishes the baseline, subsequent releases compare against it. This is
  recorded as an open item for whoever runs the first release to close, not silently decided
  here without data.
- **Gate logic for v1.0.0, made explicit (Final Wave Review Gate finding CRITICAL-1)**: the
  `gremlins` job is **advisory/informational only for the first release** — it publishes its
  kill-rate report as a pipeline artifact but does NOT fail the release pipeline (`continue-on-error:
  true` equivalent for this one job only; all other release-gate jobs remain blocking). Rationale:
  a pass/fail floor requires a baseline that does not yet exist for a greenfield codebase, and
  inventing an arbitrary floor (e.g. "80%") without evidence would be exactly the kind of
  undocumented, unjustified number `nw-test-design-mandates`/Core Principle 13 warns against.
  **Action item for whoever prepares the v1.1.0+ release**: read the v1.0.0 kill-rate report and
  set an explicit numeric floor (e.g. "no release ships below the v1.0.0 baseline minus 5 points")
  before the second release's `gremlins` job is made blocking. Until that floor is set, the job
  stays advisory — this is a deliberate, documented deferral, not an oversight.

### [REF] Observability Stack

No traditional server-side observability applies (no logs aggregator, no metrics backend, no
distributed tracing target — confirmed zero network calls, brief.md Section 0/17). The
CLI-equivalent local stack:

| Signal class | Tool/mechanism |
|---|---|
| Logs (usage/business events) | `${data_dir}/usage.log`, JSONL, opt-in only (see Monitoring Contracts) |
| Logs (operational/error) | Structured stderr output with mandatory text-prefix per brief.md Section 9's accessibility rule (`Saved`, `already saved as`, `no matches found`, `Error:`) — doubles as the CLI's only "error tracking" surface, no separate error-tracking service exists |
| Metrics | `bm stats` — local aggregation of `usage.log`, no metrics backend |
| Traces | N/A — single-process, single-request-per-invocation CLI; no distributed call graph exists to trace |
| Health checks | ADR-007 `Probe()` contract at every startup — the CLI's equivalent of a liveness/readiness check, refusing the operation (`health.startup.refused`) rather than degrading silently |

### [REF] Branching Strategy

**Selected: Trunk-Based Development** (user decision, this session). Single `main` branch,
short-lived feature branches (<1 day). CI/CD alignment:

- **Triggers**: `push: [main]` runs the full commit + fault-injection stages (main must always
  stay releasable, per trunk-based's own requirement for robust automated gates). `tags: ['v*']`
  runs the release pipeline.
- **Branch protection on `main`**: required status checks (build, test+coverage, lint, SAST,
  SCA, secrets scan, `go-arch-lint`, AST structural probe check, fault-injection suite) all
  must pass before merge; linear history required; force-push restricted.
- **PR gates**: for a solo maintainer, self-merge is expected — the automated status checks
  above are the actual gate, not a second-approver review (no team to provide one). This is
  documented explicitly rather than silently omitting the "review approvals" cell from the
  gate taxonomy.

### [REF] Coexistence Matrix

Full detail in `environments.yaml`. Summary: **N/A / empty by design** — `bm` installs no hooks,
daemons, or shell-rc mutations and shares no external state with other tools, so there is
nothing for a deployment of `bm` to break. Recorded explicitly (not omitted) so the empty matrix
reads as a verified conclusion, not an oversight.

### [REF] Pre-requisites From DESIGN

Constraints the platform/pipeline must satisfy, carried forward from `brief.md`:

- Zero network calls anywhere in the architecture (brief.md Section 0) — any future CI/CD or
  telemetry design that introduces one would re-open the closed feasibility question; the local
  opt-in usage log above was designed specifically to preserve this constraint.
- Every driven adapter (including the new `FileUsageLogAdapter`) must implement and be probed
  via `Probe() error` at startup (ADR-007) — enforced by the existing 3-layer tooling with no
  new scope required, since the tooling walks all adapters generically.
- <100ms perceived-save-confirmation budget (brief.md Section 1) — the usage-log write for
  `bm.save` must not block the printed confirmation, mirroring the existing backup-snapshot
  ordering note in brief.md Section 8 (snapshot copy happens after confirmation is printed).
- No ANSI-color-only status indicators (brief.md Section 9) — carried into the Observability
  Stack's stderr design above.
- 10,000-bookmark resource-scaling target (brief.md Section 10) — no pipeline change required;
  noted so DISTILL doesn't need to re-derive it.

### [REF] Handoff Status

**Handoff-ready to DISTILL (nw-acceptance-designer)**, pending one open item: the
`## Mutation Testing Strategy` write to `CLAUDE.md` is finalized in content (pre-release
template text) but held pending direct user confirmation — see
`docs/feature/bookmark-cli/devops/wave-decisions.md` for the full explanation and the specific
confirmation this agent needs, following the same standing-rule precedent DESIGN wave already
established for CLAUDE.md writes (agent-relayed approval claims are not sufficient
authorization). This does not block the DISTILL handoff — acceptance-designer can proceed using
`environments.yaml` and this section's content regardless of when the CLAUDE.md write lands.

Per-wave Forge review (`nw-platform-architect-reviewer`) was evaluated against its trigger list
(novel deployment target, new CI/CD framework, observability rewrite, security posture change,
maintainer-flagged uncertainty) — **none fired**: GitHub Actions + Homebrew is a conventional,
well-understood setup for this project size, so per-wave review is skipped per the `nw-devops`
skill's default. The mandatory consolidated review (Eclipse + Architect + Forge + Sentinel) fires
at the end of DISTILL against the full `feature-delta.md`.

---

## Wave: DISTILL

Facilitator: Quinn (nw-acceptance-designer) | Date: 2026-08-07 | Density mode: lean
(Tier-1 `[REF]` always emitted; Tier-2 `[WHY]`/`[HOW]` only where a trigger fires)

### [REF] Wave-Decision Reconciliation (HARD GATE)

Read `discuss/wave-decisions.md`, `design/wave-decisions.md`, `devops/wave-decisions.md` in full
before any scenario was written. Checked every DISCUSS decision against DESIGN and DEVOPS for
contradiction (email-vs-in-app, REST-vs-gRPC-class conflicts). **Zero contradictions found** —
DESIGN's Go/SQLite/Cobra/Plan-value/Probe choices and DEVOPS's GitHub-Actions/trunk-based/
pre-release-mutation/local-telemetry choices are each consistent extensions of DISCUSS's stories,
not overrides. **Reconciliation passed — 0 contradictions.**

### [REF] Language + Infrastructure Policy + Port Bootstrap

- `[lang-mode] Go` — detected via `go.mod` bootstrap this wave (greenfield: no `go.mod` existed
  before DISTILL; module `bookmark-cli`, Go 1.25, `go get github.com/spf13/cobra`,
  `go get modernc.org/sqlite` fetched successfully — network available, no offline-mode fallback
  needed).
- `[policy-mode]` — `docs/architecture/atdd-infrastructure-policy.md` did not exist; bootstrapped
  this wave (`--policy=inherit` default, file was absent → treated as "create + populate", not a
  rewrite). All 3 driving/driven-internal ports in scope populated in the same pass (no
  driven-external ports exist in this project, so that section is explicitly empty, not omitted).
- `[port-mode]` — `tests/common/state_delta.go` did not exist; bootstrapped this wave (Go binding
  of the Polyglot Adapter Matrix's state-delta port: `AssertStateDelta`, `SetTo`, `Unchanged`,
  `AppendedWith`, `PrependedWith`, `Containing`, `NormalizedTo`, `IdempotentAfter`,
  `LegacyHealed`). First DISTILL run in this project — per-project apply-if-absent bootstrap,
  future features in this repo inherit it.

### [REF] Scenario List With Tags

25 scenarios (1 walking skeleton + 24 milestone/adapter). Go idiom per the Polyglot Adapter
Matrix: `*_scenarios_test.go` files, `testing` package, `t.Skip("pending")` one-at-a-time marker
(all scenarios skip-marked except the walking skeleton, per Mandate 5).

| Scenario (Go test name) | File | Story | Tags |
|---|---|---|---|
| `TestWalkingSkeleton_SaveFindShare` | `walking_skeleton_test.go` | US-01/02/03 | `@walking_skeleton @driving_port @us-01 @us-02 @us-03` |
| `TestSave_WithoutTag_StillSavesAndRetrievable` | `save_scenarios_test.go` | US-01 | `@us-01` |
| `TestSave_ExactDuplicate_DetectedNotDuplicated` | `save_scenarios_test.go` | US-01/US-05 | `@us-01 @us-05 @error` |
| `TestSave_MalformedURL_RejectedWithClearMessage` | `save_scenarios_test.go` | US-05 | `@us-05 @error` |
| `TestSave_SameURLNewTag_OffersTagUpdateNotDuplicate` | `save_scenarios_test.go` | US-05 | `@us-05` |
| `TestSave_URLWithQueryParams_SavesWithFullFidelity` | `save_scenarios_test.go` | US-05 | `@us-05` |
| `TestSave_WithoutTag_ShowsDiscoverabilityHint` | `save_scenarios_test.go` | US-04 | `@us-04` |
| `TestSaveHelp_ShowsConcreteExample` | `save_scenarios_test.go` | US-04 | `@us-04` |
| `TestSave_NearMissFlag_SuggestsDidYouMean` | `save_scenarios_test.go` | US-04 | `@us-04 @error` |
| `TestSave_ConfirmationIsResponsive` | `save_scenarios_test.go` | US-01 | `@us-01` |
| `TestFind_ByTagAndKeyword_ShowsMatchWithMetadata` | `find_scenarios_test.go` | US-02 | `@us-02` |
| `TestFind_ByKeywordAlone_ShowsMatch` | `find_scenarios_test.go` | US-02 | `@us-02` |
| `TestFind_NearMissTypo_StillSurfacesMatch` | `find_scenarios_test.go` | US-02 | `@us-02 @property` |
| `TestFind_MultipleMatches_RankedNotForcedToOne` | `find_scenarios_test.go` | US-02 | `@us-02` |
| `TestFind_NoMatch_SuggestsClosestTag` | `find_scenarios_test.go` | US-06 | `@us-06 @error` |
| `TestFind_NoMatch_NoCloseTag_ShowsCleanMessage` | `find_scenarios_test.go` | US-06 | `@us-06 @error` |
| `TestFind_EmptyStore_ShowsOnboardingMessage` | `find_scenarios_test.go` | US-06 | `@us-06 @error` |
| `TestFind_NoMatchResponse_IsResponsive` | `find_scenarios_test.go` | US-06 | `@us-06` |
| `TestShare_CuratedLink_ProducesZeroInstallSnippet` | `share_scenarios_test.go` | US-03 | `@us-03` |
| `TestShare_SnippetMatchesFindRecordExactly` | `share_scenarios_test.go` | US-03 | `@us-03` |
| `TestShare_LinkWithoutTag_ProducesValidSnippet` | `share_scenarios_test.go` | US-03 | `@us-03` |
| `TestShare_UnknownID_FailsClearly` | `share_scenarios_test.go` | US-03 | `@us-03 @error` |
| `TestSave_CreatesBackupSnapshot` | `adapter_integration_scenarios_test.go` | US-01/NFR | `@real-io @adapter-integration @us-01` |
| `TestSave_WithTelemetryEnabled_RecordsUsageEvent` | `adapter_integration_scenarios_test.go` | DEVOPS telemetry | `@real-io @adapter-integration` |
| `TestSave_DegradedFilesystem_RefusesCleanly` | `adapter_integration_scenarios_test.go` | NFR (ADR-007) | `@real-io @adapter-integration @error @us-07-nfr` |

**Error/edge scenario count: 11/25 = 44%** (malformed URL, near-miss flag, all 3 no-match/empty-
store scenarios, unknown share id, degraded-filesystem — exceeds the 40% target).

### [REF] Walking Skeleton Strategy

Per the Architecture of Reference (retired per-feature Strategy A/B/C/D choice): this project has
**zero driven-external/non-deterministic ports** (confirmed `brief.md` Section 0/17, zero network
calls anywhere). Every driven port is driving (CLI, real adapter) or driven-internal (SQLite
store, file backup, usage log — all real adapters per the bootstrapped
`atdd-infrastructure-policy.md`). Consequently **every scenario in this feature runs at the
subprocess/FS acceptance layer with 100% real adapters** — there is no in-memory-double layer to
choose between. Walking skeleton scenario: `TestWalkingSkeleton_SaveFindShare`, invoking the real
`bm` binary via subprocess (Pillar 3), asserting on CLI-observable output only (traditional
assertions, per the Layered Test Discipline table's WS row).

### [REF] Adapter Coverage Table (Mandate 6)

| Adapter | `@real-io` scenario | Covered by |
|---|---|---|
| `SQLiteBookmarkStore` (`BookmarkReader`+`BookmarkWriter`) | YES | Every save/find/share scenario — real SQLite file under `t.TempDir()` |
| `FileBackupAdapter` (`BackupService`) | YES | `TestSave_CreatesBackupSnapshot` |
| `FileUsageLogAdapter` (`UsageLogger`, telemetry-enabled) | YES | `TestSave_WithTelemetryEnabled_RecordsUsageEvent` |
| `NoOpUsageLogAdapter` (`UsageLogger`, telemetry-disabled) | N/A — real trivial adapter, no business logic to exercise with real I/O; implicitly exercised by every scenario that does not call `WithTelemetryEnabled()` | default composition-root path |

Zero "NO — MISSING" rows.

### [REF] Scaffolds (Mandate 7 — RED-Ready)

All scaffolds compile (`go build ./...` exit 0) and panic with `"... -- RED scaffold"` messages
(Go's assertion-class RED marker, per the skill's language mapping). `cmd/bm/main.go`'s composition
root recovers each command's panic into a controlled non-zero-exit CLI failure so acceptance
assertions fail on observable output, not a crash trace.

| Scaffold file | Scaffolds |
|---|---|
| `internal/core/types.go` | `Record`, `ValidationResult`, `NormalizedTag`, `DuplicateVerdict`, `SavePlan`(+`SavePlanKind`), `RankedMatch(es)`, `ShareSnippet` — types only, no panics |
| `internal/core/validator.go` | `ValidateURL` |
| `internal/core/normalizer.go` | `NormalizeTag` |
| `internal/core/duplicate.go` | `CheckDuplicate` |
| `internal/core/planner.go` | `PlanSave` (ADR-006 Plan-value pattern) |
| `internal/core/matcher.go` | `RankMatches` |
| `internal/core/formatter.go` | `FormatSnippet` |
| `internal/ports/ports.go` | `Prober`, `BookmarkReader`, `BookmarkWriter`, `BackupService`, `UsageLogger` — interfaces only, no panics |
| `internal/adapters/sqlitestore/store.go` | `Store.Probe/FindByID/Search/All/Execute` |
| `internal/adapters/backup/filebackup.go` | `Adapter.Probe/Snapshot` |
| `internal/adapters/usagelog/usagelog.go` | `FileUsageLogAdapter.Probe/Record` (real, not scaffolded: `NoOpUsageLogAdapter` — trivial, no logic to defer) |
| `cmd/bm/main.go` | Cobra command wiring (`save`/`find`/`share`/`stats`), composition root, panic-recovery boundary — not itself scaffolded, delegates to the above |
| `cmd/bm/render.go` | `renderSaveConfirmation`, `renderFindResult` (delegate to core scaffolds) |

### [REF] Test Placement

`tests/acceptance/bookmark_cli/` (Go idiom: acceptance tests live under a top-level `tests/`
directory outside `internal/`, since `internal/` packages restrict import visibility and
subprocess-based acceptance tests only need the built binary, not internal package access).
`tests/common/state_delta.go` hosts the project-local, per-project state-delta port (inherited by
future features in this repo). No precedent existed in this greenfield repo; this layout follows
the `nw-distill` skill's default `tests/{test-type-path}/{feature-id}/acceptance/` convention
adapted to Go's `testing`-package idiom (no separate `.feature` file — the Go test function body
+ its Given/When/Then doc comment together are the SSOT, per the Polyglot Adapter Matrix's Go row
which specifies `*_scenarios_test.go`, not a Gherkin `.feature` file).

### [REF] Assertion Convention — testify (require/assert)

**Standing convention for this project, established this session (coordinator directive,
2026-08-08), applies to every DISTILL run in `bookmark-cli` going forward** — not a one-off
choice for this feature only.

- **Library**: `github.com/stretchr/testify` (`require` + `assert` sub-packages). Added to
  `go.mod` (`go get github.com/stretchr/testify@latest` — resolved `v1.11.1`; `go mod tidy` also
  pulled the sub-package-only imports `require`/`assert` into `go.sum`, plus their transitive
  deps `davecgh/go-spew`, `pmezard/go-difflib`, `gopkg.in/yaml.v3`).
- **Convention**: `require.*` for a precondition/step whose failure makes the rest of the
  scenario meaningless to continue (e.g. `require.Equal(t, 0, result.ExitCode, ...)` before
  inspecting `result.Stdout` content — a non-zero exit means the stdout assertions that follow
  would just be testing an error message, not the intended behavior). `assert.*` for independent,
  collectible outcome checks within a single `Then` (e.g. multiple `assert.Contains(...)` calls
  checking different substrings of the same confirmation line — each one is worth reporting even
  if an earlier one already failed, since they diagnose different aspects of the same output).
- **Applied to all 7 scenario/harness files** in `tests/acceptance/bookmark_cli/`:
  `walking_skeleton_test.go`, `save_scenarios_test.go`, `find_scenarios_test.go`,
  `share_scenarios_test.go`, `adapter_integration_scenarios_test.go`, `harness_test.go` (the CLI
  composition-root harness itself — `require.NoError`/`require.NoError` replace the prior
  `c.t.Fatalf` calls in `WithReadOnlyDataDir` and `run`), and `domain_types_test.go` (no change
  needed — pure type declarations, zero assertions, explicitly noted as N/A in-file rather than
  silently skipped).
- **Zero raw `t.Fatalf`/`if`-then-`t.Fatalf` assertion patterns remain** in any scenario or
  harness file (verified: `grep -n "t\.Fatalf\|c\.t\.Fatalf" tests/acceptance/bookmark_cli/*_test.go`
  returns no matches). `statedelta.AssertStateDelta` (the project's Mandate-8 state-delta port,
  `tests/common/state_delta.go`) is unchanged — it is a purpose-built Universe-guard assertion,
  orthogonal to the general-purpose require/assert convention, and continues to take a
  `statedelta.TestingT` (a `*testing.T`-compatible minimal interface) directly.
- **Re-verification of the pre-DELIVER RED gate after the rewrite**: re-ran the full suite once
  with all `t.Skip(...)` markers temporarily lifted (backed up first, restored after). Same
  classification as before the rewrite — 24/25 scenarios FAIL with `MISSING_FUNCTIONALITY`
  (`require.Equal`/`require.NotEqual`/`assert.Contains` etc. firing against the RED-scaffold panic
  output, not a build or import error), 1 scenario (`TestSaveHelp_ShowsConcreteExample`) passes
  for the same pre-documented reason (static Cobra `Example:` metadata, not deferred logic — see
  `distill/red-classification.md` "Note on the help-example scenario"). `go build ./...` and
  `go vet ./...` both exit 0 throughout. The testify migration is a pure test-assertion-library
  swap — it changed zero scenario semantics and zero production code. Full detail appended to
  `distill/red-classification.md` (see "Testify Migration — Re-Verification" section there).

### [REF] Driving Adapter Coverage

All 4 driving-port commands (`bm save`, `bm find`, `bm share`, `bm stats`) are wired in
`cmd/bm/main.go`. `bm stats` is DEVOPS's new driving port (`kpi-contracts.yaml`) — no dedicated
acceptance scenario was added for it this wave (no user story in DISCUSS scope owns `bm stats`
directly; it is telemetry tooling, not a walking-skeleton/Release-1 story) but its RunE is
scaffolded consistently with the other three so `go build` succeeds and its `--help`/disabled-
telemetry path (`"telemetry is not enabled..."`) is real, working code, not a scaffold panic —
flagged here as an intentionally out-of-DISTILL-scope command rather than a silent omission. `bm
save`, `bm find`, `bm share` are each exercised via subprocess in ≥1 scenario (Driving Adapter
Verification mandate): exit code, stdout format, and argument handling (`--tag`, near-miss flags,
positional args) are all asserted.

### [REF] Pre-requisites

- `docs/product/architecture/brief.md` Sections 5/13/14 (component contract shapes, driving/
  driven ports) — directly drove the `internal/core`/`internal/ports`/`internal/adapters` package
  layout above.
- `docs/feature/bookmark-cli/environments.yaml` — `degraded-filesystem` environment directly
  produced `TestSave_DegradedFilesystem_RefusesCleanly`; `existing-store` environment underlies
  every scenario using `WithExistingStore(...)`; `clean` environment underlies every scenario
  using a bare `NewCLI(t)`.
- `docs/product/kpi-contracts.yaml` — `bm.save`/`bm.find`/`bm.share` event names and privacy
  contract (no URL/tag content in log payloads) directly produced
  `TestSave_WithTelemetryEnabled_RecordsUsageEvent`'s assertions.
- `docs/feature/bookmark-cli/discuss/story-map.md` — US-01→US-06 priority order drove scenario
  authorship order and the DELIVER-facing one-at-a-time sequencing note in
  `distill/red-classification.md`.

### [REF] Mandate 8/9/10/11 Application

- **Mandate 8 (Universe-bound assertion, layers 1-3)**: applied to every state-**mutating**
  scenario (save-path scenarios) via `statedelta.AssertStateDelta` against a CLI-observable
  Universe (`find.match_count`, built from re-running `bm find` before/after — never an internal
  SQLite column). Read-only scenarios (find/share) have no mutation to assert a delta on, so they
  use traditional assertions, consistent with Mandate 8's own scope ("state-mutating" steps only).
- **Mandate 9 (layer-dependent PBT mode)**: this feature's acceptance layer is subprocess/FS
  (layer 3) end to end — no in-memory-double layer exists (zero driven-external ports to fake).
  Per the Layered Test Discipline table, layer 3 is **example-only**; no `@given`/PBT machinery is
  imported anywhere in `tests/acceptance/`. `TestFind_NearMissTypo_StillSurfacesMatch` is tagged
  `@property` (fuzzy-match ranking is a universal-invariant AC) but is pinned as a single
  representative example at this layer, per Mandate 9 — true PBT exploration of the ranking
  function belongs to DELIVER's unit-layer tests against `core.RankMatches` directly (owned by
  the crafter, not DISTILL).
- **Mandate 10 (two-tier acceptance)**: **Tier B NOT added.** Evaluated explicitly: `bm save`/
  `bm find`/`bm share` are each single-command, single-outcome operations; the closest thing to a
  "journey" (walking skeleton's save→find→share) is exactly 3 steps but each step's precondition
  is fully captured by the prior step's *observable output* (the bookmark id), not by a rich,
  domain-varied input space requiring generative exploration — a single Tier A example already
  covers the space. No feature journey in this pass meets both Mandate 10 triggers
  simultaneously (≥3 chained scenarios AND domain-rich input space).
- **Mandate 11 (integration sad paths stay example-based)**: all `@error`-tagged and
  `@adapter-integration`-tagged scenarios are named, explicit `Test<Scenario>` functions (no PBT
  machinery), one example per failure mode, consistent with layer 3 example-only discipline.

### [REF] Mandate-12 Compliance Evidence (SSOT + Zero Duplication)

- **Criterion 1 (domain types module)**: `tests/acceptance/bookmark_cli/domain_types_test.go` —
  `SaveOutcome`, `StoreState`, `Bookmark`, `FilesystemCondition` typed enums/structs for every
  domain noun the scenarios use.
- **Criterion 2 (typed composition parameters)**: the CLI harness's composition-root-equivalent
  (`harness_test.go`'s `CLI` type) consumes `Bookmark` (not raw positional strings) in
  `WithExistingStore(...)`; `Save(url, tag string)` keeps `url`/`tag` as plain strings only where
  no richer domain enum exists yet (both are free-text user input at this layer, not closed-set
  enums — consistent with Mandate-12's "no raw `str` where a domain enum exists," which does not
  apply to genuinely open string domains).
- **Criterion 3 (no business logic in step bodies)**: harness methods (`Save`, `Find`, `Share`,
  `SaveHelp`, `SaveWithFlag`) are single-purpose delegations to `c.run(...)` (subprocess
  invocation) — no control flow beyond argument-slice construction. Scenario bodies themselves
  compose harness calls + `statedelta`/plain assertions; this is the acceptance-layer equivalent
  of "step methods delegate to composition-root services," adapted to Go's example-based `testing`
  idiom (no `given`/`when`/`then` decorators exist in Go — the closest equivalent, `t.Run`
  subtests, was not needed since Pillar-2 chaining is expressed via helper reuse, not nested
  subtests).
- **Criterion 4 (step-reuse-ratio, informational)**: 43 total `cli.<Method>(...)` invocations
  across 7 unique harness methods (`Save`, `Find`, `Share`, `SaveHelp`, `SaveWithFlag`,
  `BackupDir`, `UsageLogPath`) = **6.14× ratio**. Config-shaped, single-command-per-story feature
  shape (per the mandate's own natural-ceiling guidance) — no forced ratio-maximization was
  applied; Gherkin-equivalent doc-comment readability (Pillar 1) was preserved throughout.

### [REF] Pre-DELIVER Fail-for-the-Right-Reason Gate

Full detail: `docs/feature/bookmark-cli/distill/red-classification.md`. Summary: `go build ./...`
and `go vet ./...` both pass (zero BROKEN-class failures). 24/25 scenarios classify as clean
MISSING_FUNCTIONALITY; 1 (`TestSaveHelp_ShowsConcreteExample`) is an intentional non-scaffolded
exception (static CLI help metadata, not business logic). Two WRONG_ASSERTION bugs were found and
fixed during this gate run (both test-only fixes, zero production-code changes) — see the "Note on
the help-example scenario" and "Assertion bugs found and fixed" sections there. All scenarios
except the walking skeleton are `t.Skip`-marked for DELIVER's one-at-a-time cycle (Mandate 5); the
walking skeleton is the single active RED scenario at hand-off.

### [REF] DISTILL Wave Handoff Status

**Pending Final Wave Review Gate** (four reviewers dispatched next, per `application` deliverable-
type routing — no plugin/skill reviewer needed). Self-review checklist (Dimension 9 + Mandate 7)
below.

- [x] 1. WS strategy declared (Architecture of Reference — no per-feature A/B/C/D choice; 100%
      real adapters, no in-memory layer)
- [x] 2. WS scenario tagged `@walking_skeleton @driving_port`
- [x] 3. Every driven adapter has ≥1 `@real-io` scenario (adapter coverage table above)
- [x] 4. N/A — no in-memory doubles used anywhere in this feature (see Mandate 9 note)
- [x] 5. N/A — no container preference applicable (embedded SQLite, no external service)
- [x] 6. Mandate 7 — all production modules imported by tests have scaffold files
- [x] 7. Mandate 7 — all scaffolds carry a `SCAFFOLD: true`/`-- RED scaffold` marker
- [x] 8. Mandate 7 — all scaffold methods raise/panic (Go's assertion-class marker), not a
      generic `errors.New("todo")`
- [x] 9. Mandate 7 — tests are RED (not BROKEN) against scaffolds — verified in red-classification.md
- [x] 10. Driving Adapter — `bm save`/`bm find`/`bm share` each exercised via subprocess with
      exit-code + stdout-format + argument-handling assertions
- [x] 11. F-001 — every driven adapter has ≥1 `@real-io @adapter-integration` scenario
- [x] 12. F-002 — N/A (Go idiom has no `capsys`-equivalent step-scoping issue; stdout/stderr
      captured directly on the `exec.Cmd`, same scope throughout the harness method)
- [x] 13. F-005 — scenario files import ONLY the test harness + `tests/common/statedelta`; zero
      imports from `internal/adapters/*` (verified: `grep -L "internal/adapters" tests/acceptance/bookmark_cli/*_test.go` matches all scenario files, i.e. none import adapters directly)
- [x] 14. F-004 — timing assertions use a 2s budget (not a flaky sub-200ms figure)
- [x] 15. F-003 — N/A (Go has no import-time `sys.path` manipulation equivalent)

### [REF] Final Wave Review Gate Outcome

Four reviewers dispatched in parallel (Haiku) against the full `feature-delta.md`, per `application`
deliverable-type routing (no `@nw-plugin-validator`/`@nw-skill-reviewer` — confirmed N/A).

| Reviewer | Wave reviewed | Verdict | Blockers | Critical | High | Medium | Low |
|---|---|---|---|---|---|---|---|
| Eclipse (`nw-product-owner-reviewer`) | DISCUSS | **approved** | 0 | 0 | 0 | 0 | 0 |
| Architect (`nw-solution-architect-reviewer`) | DESIGN | **approved** | 0 | 0 | 0 | 0 | 0 |
| Forge (`nw-platform-architect-reviewer`) | DEVOPS | **conditionally_approved** | 0 | 1 | 2 | 5 | 2 |
| Sentinel (`nw-acceptance-designer-reviewer`) | DISTILL | **approved** | 0 | 0 | 0 | 0 | 0 |

**Cross-wave consistency check**: no contradictions surfaced between reviewers — Eclipse/
Architect/Sentinel's independent approvals are mutually consistent (no reviewer's approval
depended on a claim another reviewer's findings undermined).

**Forge's findings — resolved or accepted-with-conditions** (blocker_count was 0 throughout, so
per the gate rule "zero blockers, zero high (or accepted-with-conditions)" this satisfies handoff
without a `needs_revision` re-dispatch cycle; narrow, scoped edits applied directly, consistent
with the precedent already set in `design/wave-decisions.md`'s own peer-review remediation):

| ID | Finding | Resolution |
|---|---|---|
| CRITICAL-1 | Mutation-testing release gate had no pass/fail criterion | **Resolved this session** — `gremlins` job made explicitly advisory-only for v1.0.0 (publishes report, does not block release); action item recorded for whoever prepares v1.1.0+ to set a numeric floor off the v1.0.0 baseline. See DEVOPS "Mutation Testing Strategy" section, amended above. |
| HIGH-1 | `environments.yaml` claimed macOS 13.x/14.x/15.x coverage but CI only runs `macos-latest` | **Resolved this session** — `environments.yaml` `platform_coverage` narrowed to distinguish CI-verified (14.x/15.x via `macos-latest`) from unverified-best-effort (13.x), with an explicit backlog note rather than a silent overclaim. |
| HIGH-2 | Windows release binaries built (GoReleaser cross-compile) but zero Windows test environment exists | **Resolved this session** — `windows` removed from the GoReleaser cross-compile matrix for v1; WSL2 already gives Windows users a supported, tested path. Native Windows build deferred to backlog, not shipped untested. See DEVOPS "CI/CD Pipeline Outline" table, amended above. |
| MEDIUM-1 | AST probe-presence tooling's coverage of the new `FileUsageLogAdapter` unverified | **Accepted as DELIVER-scope action item** — DELIVER's crafter must confirm the `go/ast` structural check (ADR-007 layer 2) fires on `FileUsageLogAdapter` before considering that adapter GREEN; not a DISTILL-scope gap (DISTILL's own `TestSave_WithTelemetryEnabled_RecordsUsageEvent` already exercises this adapter with real I/O). |
| MEDIUM-2 | Fault-injection harness's auto-coverage of new telemetry adapters unverified | **Accepted as DELIVER-scope action item** — same disposition as MEDIUM-1; DELIVER's fault-injection CI job (`-tags=faultinjection`) must be confirmed to enumerate `FileUsageLogAdapter` once implemented. |
| MEDIUM-3 | Telemetry privacy constraint (no URL/tag in event payloads) not CI-enforced | **Accepted as DELIVER-scope action item** — `TestSave_WithTelemetryEnabled_RecordsUsageEvent` already asserts this at the acceptance layer (`!strings.Contains(logContent, "kube.io")`); a dedicated static-analysis/linter enforcement is a DELIVER/DEVOPS hardening task, not a DISTILL blocker. |
| MEDIUM-4 | `usage.log` has no documented rotation/retention policy | **Accepted as backlog item** — out of DISCUSS scope (no story requires log rotation); flagged for a future release, not MVP. |
| MEDIUM-5 | Fault-injection CI suite's actual runtime vs. the <10min commit-stage target unmeasured | **Accepted as DELIVER-scope action item** — cannot be measured until the fault-injection suite exists (DELIVER GREEN phase); DEVOPS's parallel-jobs design already anticipated this risk. |
| LOW-1 | No DISTILL scenario asserts the coexistence-matrix "N/A" claim | **Accepted as backlog item** — informational only, the underlying claim is correct today. |
| LOW-2 | Disabled-telemetry (`NoOpUsageLogAdapter`) path relies on implicit coverage | **Accepted as backlog item, informational** — every save/find/share scenario in this suite already exercises the default (telemetry-disabled) composition-root path; a dedicated `TestSave_WithTelemetryDisabled_DoesNotCreateLogFile` scenario is a low-cost DELIVER-phase addition, not required for handoff. |

**Result: zero blockers, zero unresolved critical/high findings.** Both HIGH findings and the one
CRITICAL finding were closed with direct, narrow, scoped documentation edits this session (no
architectural rework, no re-dispatch of `@nw-platform-architect` needed — same "straightforward,
scoped edits, no iteration-2 re-review" pattern DESIGN's own peer review already established).
Five MEDIUM/LOW findings are accepted-with-conditions as documented DELIVER-scope or backlog
action items above, none of which block scenario authorship or the DELIVER RED→GREEN cycle.

### [REF] DISTILL Wave Final Handoff Status

**HANDOFF-READY to DELIVER.** All four Final Wave Review Gate verdicts are APPROVED or
CONDITIONALLY_APPROVED with documented action items (table above) — gate condition satisfied
(zero blockers; zero unresolved high findings). Pre-DELIVER fail-for-the-right-reason gate passed
(`docs/feature/bookmark-cli/distill/red-classification.md`). Mandate compliance evidence (CM-A
through CM-H, plus Mandate-12 criteria 1-4) recorded in the `[REF]` sections above.

**Handoff package for DELIVER** (`@nw-functional-software-crafter`, per `CLAUDE.md`'s paradigm
routing):
- `tests/acceptance/bookmark_cli/*_test.go` (25 scenarios, 24 skip-marked, walking skeleton active RED)
- `tests/common/state_delta.go` (project-local state-delta port)
- `internal/core/`, `internal/ports/`, `internal/adapters/*`, `cmd/bm/` (RED scaffolds — every
  `panic("... -- RED scaffold")` call site is DELIVER's GREEN-phase worklist)
- `docs/feature/bookmark-cli/distill/red-classification.md` (RED gate evidence)
- `docs/architecture/atdd-infrastructure-policy.md` (bootstrapped this wave)
- This `feature-delta.md` (full DISCUSS→DESIGN→DEVOPS→DISTILL chain + review verdicts)

**Suggested DELIVER sequencing** (per `story-map.md` priority + one-at-a-time discipline):
1. `TestWalkingSkeleton_SaveFindShare` (already active RED)
2. `TestSaveHelp_ShowsConcreteExample` (trivially GREEN — static CLI metadata already correct)
3. Remaining US-01 scenarios, then US-02, US-03, US-04, US-05, US-06, then the 3
   `@adapter-integration` scenarios (backup snapshot, telemetry, degraded-filesystem) last, since
   they depend on `bm save` already being GREEN.
