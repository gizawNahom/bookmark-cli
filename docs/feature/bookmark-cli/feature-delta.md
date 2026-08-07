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
