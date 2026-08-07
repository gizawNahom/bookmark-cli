<!-- markdownlint-disable MD024 -->
# User Stories: bookmark-cli

Wave: DISCUSS | Persona: PERSONA-TERMINAL-INFRA-ENG | Scope: walking skeleton
(save/find/share) + Release 1 (usability + error handling)

## System Constraints

Carried forward from DISCOVER (`docs/feature/bookmark-cli/discover/wave-decisions.md`,
constraints table) and binding on every story below:

- MVP must support local-first save/find with **no required account creation**.
- The `bm share` output must require **zero install for the recipient**.
- Tag syntax must be **more discoverable** than the DISCOVER-tested prototype
  (Aisha's task-3 comprehension took ~20s vs. a <10s target) — see US-04.
- **No server infrastructure is assumed for MVP.** This is carried forward as
  an unverified assumption per PO risk-acceptance decision — see
  `wave-decisions.md` Upstream Changes / gate-override entry. If a future
  feasibility spike invalidates local-first storage, these stories'
  acceptance criteria (observable CLI behavior) do not change, but Technical
  Notes referencing "local bookmark store" will need re-grounding in DESIGN.

### NFR Guardrails Added After Peer Review (iteration 1)

The following gaps were flagged by `nw-product-owner-reviewer` (high/medium
severity, non-blocking, `conditionally_approved`) and are recorded here as
explicit constraints for DESIGN to design against, rather than left implicit:

- **Data reliability**: local-first storage is a single point of failure —
  DESIGN must decide and document a backup/recovery approach (even a manual
  one, e.g. "back up `${store_location}` yourself") before MVP ships. No
  story in this pass implements backup/restore; this is a DESIGN-phase
  decision, not a DISCUSS-phase story, because the "right" backup mechanism
  depends on the storage format DESIGN chooses.
- **Accessibility**: this is a terminal-only CLI with no color-only signals
  in any mockup in this pass (see `journey-save-find-share-visual.md`) — text
  output is screen-reader-compatible by the nature of the medium. No
  additional accessibility story is added for this pilot-scoped MVP; DESIGN
  should confirm no ANSI-color-only status indicators are introduced when
  visual output (e.g. colored tags) is designed in later releases.
- **Concurrency**: `bm save`/`bm find`/`bm share` may be invoked from
  multiple terminal panes or scripts concurrently against the same local
  store. DESIGN must specify the store's concurrency guarantee (atomic
  writes at minimum) — no story's acceptance criteria assume single-process
  access.
- **Resource scaling**: outcome KPI 2's "<10 second locate" target has no
  stated upper bound on store size. DESIGN should define a target store size
  (e.g. "responsive up to 10,000 bookmarks") before implementation.

These are carried to DESIGN as constraints, not new DISCUSS-wave stories,
because introducing them as separate user stories now would prescribe
storage-format solutions that belong in DESIGN (Core Principle 5:
problem-first, solution-never).

---

## US-01: Save a Link Without Leaving the Terminal

### job_id
JTBD-CAPTURE

### Elevator Pitch
- **Before**: Nadia Petrova pastes URLs into a personal bash function that
  appends to a flat text file — no structure, and a mistyped tag creates two
  inconsistent tags forever.
- **After**: Nadia runs `bm save <url> --tag <tag>` from wherever she is in
  the terminal and sees an immediate confirmation with a bookmark id.
- **Decision enabled**: Nadia decides, in the moment, whether the link is
  safely captured and she can return to her work — without needing a second
  backup copy pasted into Slack-DM-to-self.

### Problem
Nadia Petrova is a platform engineer who finds a useful technical reference
mid-task and needs to capture it without alt-tabbing out of her terminal.
She finds it costly to interrupt her flow to paste into a GUI app, and her
current workaround (a hand-rolled bash function appending to a text file) has
no structure and no duplicate protection.

### Who
- Terminal-first infra/platform/SRE/staff engineer | Mid-task in a terminal
  session, often mid-debug or mid-incident | Motivated to capture a link in
  one command without breaking flow

### Solution
A `bm save <url> --tag <tag>` command that saves a link locally in one
command and confirms with a bookmark id.

### Domain Examples
1. **Happy path** — Nadia Petrova, mid-debugging a Kafka consumer lag issue,
   runs `bm save https://kube.io/docs/failover --tag k8s` and sees
   `Saved [a1b2] "https://kube.io/docs/failover" (tag: k8s)`.
2. **Edge case** — Sam Okafor, mid-incident, runs
   `bm save https://wiki.internal/postmortem-2024-03` with no `--tag` and
   still gets a confirmed, retrievable save.
3. **Error/boundary** — Marco Alves runs `bm save` on a URL he already saved
   last month; the system tells him it's already saved instead of creating a
   second entry.

### UAT Scenarios (BDD)

#### Scenario: Link is saved with a tag in one command
```gherkin
Given Nadia Petrova is in her terminal mid-debugging session
When she runs "bm save https://kube.io/docs/failover --tag k8s"
Then she sees "Saved [a1b2] https://kube.io/docs/failover (tag: k8s)"
And the bookmark is retrievable later by that id
```

#### Scenario: Link is saved without a tag
```gherkin
Given Sam Okafor is mid-incident and wants to capture a link quickly
When he runs "bm save https://wiki.internal/postmortem-2024-03"
Then he sees a confirmation with a bookmark id and no tag shown
And the link is still saved and retrievable by keyword
```

#### Scenario: Saving the same URL twice is caught, not silently duplicated
```gherkin
Given Marco Alves already saved "https://kube.io/docs/failover" last month
When he runs "bm save https://kube.io/docs/failover --tag k8s" again
Then he sees a message that this URL is already saved, with its existing bookmark id
And no second entry is created
```

#### Scenario: Confirmation appears fast enough to not break flow
```gherkin
Given Nadia Petrova is mid-flow in a tmux session
When she runs "bm save <url> --tag <tag>"
Then the confirmation appears within a responsive (<100ms perceived) feedback window
```

### Acceptance Criteria
- [ ] `bm save <url> --tag <tag>` returns a confirmation with a bookmark id and the tag
- [ ] `bm save <url>` (no `--tag`) still saves the link successfully
- [ ] Saving a URL already present is detected and reported, not silently duplicated
- [ ] Confirmation output appears within a responsive feedback window

### Outcome KPIs
- **Who**: Terminal-first infra/platform engineers (pilot segment)
- **Does what**: Save a link via `bm save` instead of their prior workaround
- **By how much**: 80%+ of new saves go through `bm save` within 2 weeks of adoption
- **Measured by**: Opt-in local usage count + weekly pilot self-report
- **Baseline**: 0% (new tool; prior workaround = 100%)

### Technical Notes (Optional)
- Local-first storage per System Constraints — no account/auth flow.
- Duplicate detection requires a lookup key (URL); deeper error-path coverage
  (malformed URLs, tag-conflict-on-duplicate) is in US-05.

---

## US-02: Locate a Saved Link by Keyword or Tag

### job_id
JTBD-LOCATE

### Elevator Pitch
- **Before**: Priti Desai searches a flat text file with `grep -i` from
  muscle memory, or during an on-call incident waits 4 minutes for a
  teammate to post the right runbook link in Slack.
- **After**: Priti runs `bm find k8s failover` from her terminal and sees the
  matching link with its tag and bookmark id in the results.
- **Decision enabled**: Priti decides whether the top match is the right
  runbook to act on immediately, without waiting on a teammate to confirm.

### Problem
Priti Desai is a site reliability engineer who knows she saved a link before
but can't remember where or under what word. She finds it costly, especially
mid-incident, to hunt through unstructured notes or ask teammates to
re-locate a link that should already be findable.

### Who
- Terminal-first infra/SRE engineer | Searching under time pressure, often
  mid-incident | Motivated to retrieve a known-saved link in seconds

### Solution
A `bm find <query>` command that searches saved links by keyword and tag,
returning ranked, retrievable results.

### Domain Examples
1. **Happy path** — Priti Desai runs `bm find k8s failover` and gets the
   bookmark she saved 12 days ago tagged "k8s".
2. **Edge case** — Jordan Osei searches with a typo, `bm find terrafrom`, and
   still sees his Terraform provider doc ranked in the results (fuzzy match).
3. **Error/boundary** — Aisha Rahman searches `bm find gRPC-retry-policy` with
   nothing saved under that term and gets a helpful "no matches" message
   instead of a blank result (full detail in US-06).

### UAT Scenarios (BDD)

#### Scenario: A saved link is found by tag and keyword
```gherkin
Given Priti Desai saved a link 12 days ago tagged "k8s" with "failover" in its title
When she runs "bm find k8s failover"
Then she sees the matching link with its bookmark id, tag, and "saved 12 days ago"
```

#### Scenario: A saved link is found by keyword alone
```gherkin
Given Marco Alves saved "https://redis.io/failover-postmortem" without a tag
When he runs "bm find failover"
Then he sees the matching link in the results
```

#### Scenario: A near-miss search term still surfaces the right link
```gherkin
Given Jordan Osei saved a Terraform provider gotcha doc
When he runs "bm find terrafrom" (typo)
Then he still sees the Terraform doc ranked in the results
```

#### Scenario: Multiple matches are shown ranked, not forced to one guess
```gherkin
Given Nadia Petrova has 3 saved links tagged "k8s"
When she runs "bm find k8s"
Then she sees all 3 matches ranked by relevance
```

#### Scenario: A search with no saved matches gives a helpful response
```gherkin
Given Aisha Rahman has no saved link matching "gRPC-retry-policy"
When she runs "bm find gRPC-retry-policy"
Then she sees a message that no matches were found, not a blank result
```

### Acceptance Criteria
- [ ] `bm find <query>` matches against both tag and keyword content
- [ ] Search is typo-tolerant for near-miss terms (fuzzy match)
- [ ] Multiple matches are shown ranked, never a single forced guess
- [ ] No-match searches return a guiding message, not a blank/silent result
- [ ] Every result displays the same bookmark id shown at save time

### Outcome KPIs
- **Who**: Terminal-first infra/platform engineers (pilot segment)
- **Does what**: Locate a previously saved link via `bm find`
- **By how much**: Median time-to-locate drops from 10-15 min/week to a single lookup under 10 seconds
- **Measured by**: Task-timing during pilot + weekly self-reported time estimate
- **Baseline**: 10-15 min/week spent hunting for links (Sam Okafor, DISCOVER #4)

### Technical Notes (Optional)
- Depends on `${bookmark_id}` integrity from US-01 (shared artifact — see
  `shared-artifacts-registry.md`).
- Full "no results" UX detail (tag suggestions, empty-store distinction) is
  in US-06.

---

## US-03: Share a Curated Link With a Teammate, Zero Install

### job_id
JTBD-SHARE

### Elevator Pitch
- **Before**: Aisha Rahman pastes links into a shared `LINKS.md` file in the
  team repo, causing weekly merge conflicts, or pastes into Slack where links
  get buried and go stale unnoticed.
- **After**: Aisha runs `bm share a1b2` and gets a copy-paste-ready snippet
  she can drop anywhere; her teammate uses it immediately with no install.
- **Decision enabled**: Aisha decides to hand off a curated link the moment
  she finds it, instead of batching it into a periodically-conflicting shared
  file.

### Problem
Aisha Rahman is a staff engineer / tech lead who curates links for her team.
She finds it costly that sharing today means either a merge-conflict-prone
shared file or an easily-buried chat message, and any solution that requires
her teammates to install something is a non-starter.

### Who
- Terminal-first staff engineer / tech lead curating links on behalf of a
  team | Sharing after having already located a valuable link | Motivated by
  zero recipient-side friction

### Solution
A `bm share <id>` command that generates a copy-paste-ready snippet from a
saved bookmark, requiring no install or account for the recipient.

### Domain Examples
1. **Happy path** — Aisha Rahman runs `bm share a1b2` for the failover doc
   and gets a snippet her teammate opens with zero install.
2. **Edge case** — Jordan Osei shares a link that has no tag; the snippet
   still contains a valid URL with no broken tag field.
3. **Error/boundary** — Priti Desai mistypes a bookmark id, `bm share a9z9`,
   and gets a clear error instead of a blank or garbled snippet.

### UAT Scenarios (BDD)

#### Scenario: A curated link is shared with a zero-install snippet
```gherkin
Given Aisha Rahman has found bookmark "a1b2" for the failover doc, tagged "k8s"
When she runs "bm share a1b2"
Then she sees a copy-paste-ready snippet containing the URL and tag
And the snippet requires no install or account for the recipient to use
```

#### Scenario: The shared snippet matches the source record exactly
```gherkin
Given bookmark "a1b2" resolves to "https://kube.io/docs/failover" tagged "k8s"
When Aisha Rahman runs "bm share a1b2"
Then the URL and tag in the snippet exactly match the record shown by "bm find"
```

#### Scenario: A link saved without a tag can still be shared
```gherkin
Given Jordan Osei saved a link with no tag
When he runs "bm share" with that bookmark's id
Then he sees a valid snippet containing the URL with no broken tag field
```

#### Scenario: Sharing an unknown bookmark id fails clearly
```gherkin
Given Priti Desai runs "bm share a9z9" for an id that does not exist
Then she sees a clear error naming the invalid id, not a blank or garbled snippet
```

#### Scenario: A teammate uses the shared link with zero install
```gherkin
Given Aisha Rahman has sent her teammate the output of "bm share a1b2"
When her teammate opens the link
Then the teammate needs no installation, account, or additional tool to use it
```

### Acceptance Criteria
- [ ] `bm share <id>` produces a copy-paste-ready snippet with URL and tag
- [ ] Snippet content exactly matches the stored record (no transformation drift)
- [ ] Links saved without a tag still produce a valid, non-broken snippet
- [ ] An invalid/unknown bookmark id produces a clear, actionable error
- [ ] The recipient requires zero install, account, or additional tooling

### Outcome KPIs
- **Who**: Pilot team members receiving a shared link
- **Does what**: Use a `bm share` snippet with zero additional install steps
- **By how much**: 100% of share recipients need zero install (hard guardrail)
- **Measured by**: Pilot observation — recipient confirms usable link with no setup
- **Baseline**: N/A (new capability); protects against regression as the feature is built

### Technical Notes (Optional)
- Hard constraint: no server round-trip may be required on the recipient's
  side for MVP (consistent with System Constraints above).
- Depends on `${bookmark_id}` resolving to the exact record from US-01/US-02.

---

## US-04: Discover the Tag Flag Syntax Without a Failed Attempt

### job_id
JTBD-CAPTURE

### Elevator Pitch
- **Before**: In DISCOVER Phase 3 solution testing, Aisha Rahman needed a
  hint to discover the `--tag` flag syntax, taking ~20 seconds versus the
  <10 second target — the one usability failure across 15 tested tasks.
- **After**: Aisha runs `bm save <url>` and sees an in-context hint about the
  `--tag` flag, without needing outside help.
- **Decision enabled**: Aisha decides whether to tag the link right now or
  later, informed by output that surfaces the tag option in the moment
  rather than hiding it behind `--help`.

### Problem
Aisha Rahman (and users like her) find the `--tag` flag on `bm save` is not
discoverable enough on first use — DISCOVER usability testing showed a
comprehension failure that this story exists to remediate before the
walking skeleton ships.

### Who
- First-time or infrequent `bm` user | Attempting to tag a link without
  prior instruction | Motivated to discover the right syntax without
  external help

### Solution
In-context hints (post-save nudge, concrete `--help` examples, "did you
mean" flag correction) that make `--tag` discoverable without a failed
attempt.

### Domain Examples
1. Aisha Rahman runs `bm save <url>` with no tag and sees a one-line hint:
   `tip: add --tag <name> to make this easier to find later.`
2. Priti Desai runs `bm save --help` and sees a concrete example line
   (`bm save <url> --tag k8s`), not just an abstract flag description.
3. Marco Alves, unfamiliar with the tool, runs `bm save <url> --tags k8s`
   (typo) and gets `did you mean --tag?` instead of a bare parse error.

### UAT Scenarios (BDD)

#### Scenario: A hint nudges the user toward tagging without failing first
```gherkin
Given Aisha Rahman runs "bm save https://kube.io/docs/failover" with no --tag
When the command completes
Then she sees a one-line hint suggesting "--tag <name>" to make it easier to find later
```

#### Scenario: Help text shows a concrete example, not just abstract syntax
```gherkin
Given Priti Desai runs "bm save --help"
Then she sees an example line such as "bm save <url> --tag k8s"
And not only an abstract "--tag <TAG>" description
```

#### Scenario: A near-miss flag name is corrected with a suggestion
```gherkin
Given Marco Alves runs "bm save <url> --tags k8s" (typo: plural)
Then he sees a suggestion "did you mean --tag?"
And the command does not fail silently or with a bare parse error
```

#### Scenario: Tag flag comprehension time improves to under target
```gherkin
Given a first-time user attempts to tag a link without prior instruction
When they run "bm save <url>" and read the resulting hint or help output
Then they successfully use "--tag" within 10 seconds, matching the DISCOVER Phase 3 target
```

### Acceptance Criteria
- [ ] Saving without a tag surfaces a one-line hint about `--tag`
- [ ] `bm save --help` shows a concrete example, not only abstract flag syntax
- [ ] A near-miss flag name (e.g. `--tags`) triggers a "did you mean" suggestion
- [ ] Re-tested comprehension time for `--tag` is under 10 seconds (was ~20s in DISCOVER Phase 3)

### Outcome KPIs
- **Who**: Aisha-profile users (staff/tech-lead curating for a team)
- **Does what**: Comprehend the `--tag` flag syntax without a failed attempt
- **By how much**: Comprehension time under 10 seconds (guardrail, not a stretch)
- **Measured by**: Task-timing during Release 1 usability re-test, same protocol as DISCOVER Phase 3
- **Baseline**: ~20s / 1 failed attempt out of 5 users (DISCOVER Phase 3)

### Technical Notes (Optional)
- Directly remediates the sole DISCOVER Phase 3 usability failure; must be
  re-tested using the same protocol before Release 1 is considered done.

---

## US-05: Get a Helpful Error When Saving an Invalid or Duplicate URL

### job_id
JTBD-CAPTURE

### Elevator Pitch
- **Before**: Nadia Petrova's bash script silently appends whatever URL she
  gives it — typos, duplicates, and malformed entries all get saved without
  warning, discovered only weeks later during a painful cleanup.
- **After**: Nadia runs `bm save <bad-or-duplicate-url>` and immediately sees
  a specific, actionable message explaining why it wasn't saved as a new
  entry.
- **Decision enabled**: Nadia decides in the moment whether to fix the URL
  and retry or move on, instead of discovering data-quality problems weeks
  later.

### Problem
Nadia Petrova's prior ad hoc tooling has no input validation or duplicate
protection, and DISCOVER evidence shows this directly causes rework (twenty
minutes spent consolidating inconsistent tags by hand). This story ensures
`bm save` doesn't repeat that failure mode.

### Who
- Terminal-first infra/platform engineer | Saving links routinely over weeks
  and months | Motivated to trust that saved data stays clean without manual
  cleanup

### Solution
Input validation and duplicate detection on `bm save`, with specific,
actionable error messages.

### Domain Examples
1. Nadia Petrova runs `bm save not-a-url --tag misc` and sees "this doesn't
   look like a valid URL" instead of a silent save.
2. Marco Alves runs `bm save https://kube.io/docs/failover --tag k8s` for a
   URL he already saved and sees "already saved as [a1b2]".
3. Sam Okafor runs `bm save https://kube.io/docs/failover` (same URL, new
   tag) and is offered to add the tag to the existing entry rather than
   creating a duplicate.

### UAT Scenarios (BDD)

#### Scenario: A malformed URL is rejected with a clear message
```gherkin
Given Nadia Petrova runs "bm save not-a-url --tag misc"
Then she sees a message that this does not look like a valid URL
And nothing is saved
```

#### Scenario: An exact duplicate URL is detected, not silently re-saved
```gherkin
Given Marco Alves already saved "https://kube.io/docs/failover" as bookmark "a1b2"
When he runs "bm save https://kube.io/docs/failover --tag k8s" again
Then he sees "already saved as [a1b2]"
And no second entry is created
```

#### Scenario: Saving the same URL with a new tag updates rather than duplicates
```gherkin
Given Sam Okafor already saved "https://kube.io/docs/failover" untagged
When he runs "bm save https://kube.io/docs/failover --tag k8s"
Then he is offered to add the "k8s" tag to the existing bookmark "a1b2"
And no second entry is created
```

#### Scenario: A URL with query parameters still saves correctly
```gherkin
Given Priti Desai runs "bm save https://wiki.internal/runbook?id=42&version=3 --tag oncall"
Then the link is saved with the full URL including query parameters intact
```

### Acceptance Criteria
- [ ] Malformed URLs are rejected with a specific, actionable message
- [ ] Exact duplicate URLs are detected and reported with the existing bookmark id
- [ ] Re-saving an existing URL with a new tag offers to update rather than duplicate
- [ ] URLs with query parameters or special characters save with full fidelity

### Outcome KPIs
- **Who**: Terminal-first infra/platform engineers (pilot segment)
- **Does what**: Trust `bm save` to keep saved data clean without manual cleanup
- **By how much**: Zero manual tag/URL cleanup sessions reported during the pilot (vs. Nadia's ~20 min/incident baseline)
- **Measured by**: Weekly pilot self-report
- **Baseline**: ~20 minutes spent hand-consolidating inconsistent tags (Nadia Petrova, DISCOVER #6)

### Technical Notes (Optional)
- Extends US-01's save path; both stories share the local bookmark store as
  source of truth.

---

## US-06: Get a Helpful "No Results" Experience When Find Has No Matches

### job_id
JTBD-LOCATE

### Elevator Pitch
- **Before**: Jordan Osei searches a 40-entry unsorted browser bookmarks
  folder, finds nothing, gives up, and re-Googles from scratch (8 minutes).
- **After**: Jordan runs `bm find <query>` with no matches and immediately
  sees whether to refine his search or that nothing was ever saved under
  that term.
- **Decision enabled**: Jordan decides within seconds whether to refine his
  search term or abandon `bm` and re-Google, instead of wondering if the
  tool is broken.

### Problem
Priti Desai and Jordan Osei both need to trust a "no results" response
enough to act on it immediately rather than assume the tool failed —
DISCOVER evidence shows the cost of a bad no-match experience is a full
re-Google fallback (8 minutes).

### Who
- Terminal-first infra/platform/SRE engineer | Searching and getting zero
  matches | Motivated to know immediately whether to refine or give up

### Solution
A distinguishable, guiding "no results" response on `bm find` — including a
closest-tag suggestion when available and a distinct empty-store message.

### Domain Examples
1. Aisha Rahman searches `bm find gRPC-retry-policy` with links tagged
   "grpc" but nothing matching exactly, and sees "did you mean 'grpc'?"
2. Jordan Osei searches `bm find teraform` with zero Terraform-tagged links
   yet saved and sees a clean "no matches found" message, not a crash.
3. Priti Desai searches `bm find anything` on day one before saving anything
   and sees a distinct "you haven't saved any links yet" message.

### UAT Scenarios (BDD)

#### Scenario: A no-match search suggests the closest existing tag
```gherkin
Given Aisha Rahman has links tagged "grpc" but none matching "gRPC-retry-policy"
When she runs "bm find gRPC-retry-policy"
Then she sees "no matches for 'gRPC-retry-policy'"
And a suggestion "did you mean 'grpc'?"
```

#### Scenario: A no-match search with no close tag gives a clean message
```gherkin
Given Jordan Osei's store has no links resembling "teraform"
When he runs "bm find teraform"
Then he sees a clear "no matches found" message
And no crash, stack trace, or silent blank output
```

#### Scenario: Searching an empty store gives an onboarding message, not a false no-match
```gherkin
Given Priti Desai has not saved any links yet
When she runs "bm find anything"
Then she sees a message telling her she hasn't saved any links yet
And it is visibly different from a normal no-match result
```

#### Scenario: A no-match response feels immediate
```gherkin
Given Jordan Osei's search returns zero matches
When the command completes
Then the no-match message appears within a responsive feedback window
```

### Acceptance Criteria
- [ ] No-match searches suggest the closest existing tag when one exists
- [ ] No-match searches with no close tag show a clean message, never a crash/blank output
- [ ] An empty store produces a distinct onboarding message, not a generic no-match message
- [ ] No-match responses appear within the same responsive-feedback target as successful searches

### Outcome KPIs
- **Who**: Terminal-first infra/platform engineers (pilot segment)
- **Does what**: Get a fast, clear "no match" signal instead of falling back to re-Googling
- **By how much**: Reduce re-Google fallback rate after a `bm find` no-match to near zero (from an observed 8-minute fallback, DISCOVER #9)
- **Measured by**: Weekly pilot self-report + qualitative observation
- **Baseline**: 8-minute re-Google fallback after failing to find a link manually (Jordan Osei, DISCOVER #9)

### Technical Notes (Optional)
- Extends US-02's find path; both stories share the local bookmark store as
  source of truth.
