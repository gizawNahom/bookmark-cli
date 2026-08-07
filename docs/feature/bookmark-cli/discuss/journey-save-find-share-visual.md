# Journey (Visual): Save, Locate, and Share a Technical Reference Link

Feature: bookmark-cli | Persona: PERSONA-TERMINAL-INFRA-ENG | Wave: DISCUSS
Grounded in DISCOVER evidence (no DIVERGE wave was run for this feature — see
Upstream Changes entry in `wave-decisions.md`).

## Scope

Covers the top-3 "pursue" jobs from DISCOVER: JTBD-CAPTURE, JTBD-LOCATE,
JTBD-SHARE — the same three tasks tested in DISCOVER Phase 3 solution testing
(`bm save`, `bm find`, `bm share`, 93% task completion, n=5).

Out of scope for this journey (see Out-of-scope section of `wave-decisions.md`):
JTBD-CONTEXT-SWITCH, JTBD-ORGANIZE, JTBD-RECALL-CONTEXT, JTBD-DEDUP-MONITOR.

## ASCII Flow with Emotional Annotations

```
[Trigger: mid-task,        [Step 1: CAPTURE]        [Step 2: LOCATE]           [Step 3: SHARE]
 finds a link worth         `bm save <url>            (days/weeks later,        `bm share <id>`
 keeping]                    --tag <tag>`               often mid-incident)
                                                        `bm find <query>`
   |                             |                          |                         |
   v                             v                          v                         v
Feels: focused,           Feels: quick relief —      Feels: anxious/            Feels: proud/relieved —
mildly anxious             "one command, back to      time-pressured -> then     "my teammate got it in
("I'll lose this           work" (Sam, #4)             confident on finding      one command, no install"
 if I don't act now")                                  (Priya #1, Priti #10)     (Aisha, #8)

Sees:                      Sees:                       Sees:                     Sees:
terminal prompt,           confirmation line with      ranked results with       a copy-paste-ready
URL just retrieved          saved ID + tag echoed        matched keyword/tag      snippet or link the
 via curl/browser           back                         highlighted              recipient opens with
                                                                                   zero install

Artifacts: ${url},         Artifacts: ${bookmark_id},  Artifacts: ${query},      Artifacts: ${bookmark_id},
${tag}                     ${saved_at}                  ${match_list}             ${share_snippet}
```

## Emotional Arc

- **Start**: Focused, mid-task, mildly anxious about losing the link if not captured now.
- **Middle**: Anxiety spikes at the LOCATE step — often searching under time pressure
  (mid-incident, per Priti #10 and Sam #4) — resolving into confidence once results
  appear and the match is unambiguous.
- **End**: Proud/relieved at SHARE — validated the moment the teammate confirms they
  received something usable with zero install (Aisha #8).

No jarring transitions: CAPTURE's quick win builds trust that carries into LOCATE;
LOCATE's resolution (find succeeds) is the confidence buffer needed before the
higher-stakes SHARE step (where a teammate is watching).

## Step 1: CAPTURE — `bm save <url> --tag <tag>`

```
+-- Step 1: Save a link -------------------------------------------------+
| $ bm save https://kube.io/docs/failover --tag k8s                      |
| Saved [a1b2] "https://kube.io/docs/failover" (tag: k8s)                |
+--------------------------------------------------------------------------+
```

- Shared artifacts: `${bookmark_id}` (source: local bookmark store, single
  record created here), `${tag}` (source: user input, normalized on save).
- Emotional state: entry = focused/mildly anxious -> exit = quick relief.
- Integration checkpoint: `${bookmark_id}` generated here must be the same ID
  displayed by LOCATE and consumed by SHARE.
- Failure modes:
  - Duplicate URL already saved -> must not silently create a second entry.
  - Malformed/unreachable-looking URL -> must not save silently without feedback.
  - Tag flag syntax not discovered by the user (DISCOVER Phase 3 finding: Aisha,
    task-3, ~20s vs <10s target) -> must be fixed before this ships (see US-04).

## Step 2: LOCATE — `bm find <query>`

```
+-- Step 2: Find a saved link --------------------------------------------+
| $ bm find k8s failover                                                  |
| [a1b2] https://kube.io/docs/failover  (tag: k8s, saved 12 days ago)     |
+--------------------------------------------------------------------------+
```

- Shared artifacts: `${query}` (user input), `${match_list}` (source: local
  bookmark store, filtered/ranked by keyword+tag match against records
  created in Step 1).
- Emotional state: entry = anxious (often time-pressured) -> exit = confident.
- Integration checkpoint: every record returned must reference the same
  `${bookmark_id}` assigned at save time — no drift between save-time ID and
  find-time ID.
- Failure modes:
  - No match found -> must guide the user (e.g., suggest nearest tag) rather
    than return a bare empty result.
  - Ambiguous match (multiple close hits) -> must present ranked results, not
    force a single guess.

## Step 3: SHARE — `bm share <id>`

```
+-- Step 3: Share a link --------------------------------------------------+
| $ bm share a1b2                                                          |
| Share this with your team:                                               |
| https://kube.io/docs/failover  (via bm, tag: k8s)                        |
+--------------------------------------------------------------------------+
```

- Shared artifacts: `${bookmark_id}` (consumed from Step 1/2, must resolve to
  the exact record), `${share_snippet}` (source: generated at share time from
  the resolved record — copy-paste text or link requiring zero install for
  the recipient).
- Emotional state: entry = confident (has the right link) -> exit = proud/relieved.
- Integration checkpoint: the URL and tag shown in the share snippet must
  exactly match the record shown in Step 2 — no re-fetch or transformation
  that could drift from source.
- Failure modes:
  - Invalid/unknown `${bookmark_id}` (typo, deleted record) -> must fail with
    an actionable message, not a blank/garbled snippet.
  - Recipient-side friction (anything requiring recipient install) -> hard
    constraint violation per DISCOVER (Aisha, #8); must never happen.

## CLI Vocabulary Consistency

Command pattern: `bm <verb> [args]` (verb-first, consistent with `bm save`,
`bm find`, `bm share`). All three steps share `<id>` as the vocabulary term
for `${bookmark_id}` — no step calls it "record", "entry", or "key" instead.
