# Slice 01: Walking Skeleton — Save, Find, Share (Happy Path)

Feature: bookmark-cli | Wave: DISCUSS | Type: Walking Skeleton (thinnest end-to-end slice)

## Outcome

A terminal-first infra/platform engineer can save a link, find it again, and
share it with a teammate who needs zero install — the complete core value
proposition, happy path only.

## Stories in This Slice

| Story | Command | Job |
|---|---|---|
| US-01 | `bm save <url> --tag <tag>` | JTBD-CAPTURE |
| US-02 | `bm find <query>` | JTBD-LOCATE |
| US-03 | `bm share <id>` | JTBD-SHARE |

Full story detail: `docs/feature/bookmark-cli/discuss/user-stories.md`.

## Why This Is the Walking Skeleton

This is the exact 3-task sequence already tested in DISCOVER Phase 3
solution testing (93% task completion, n=5) — no new design risk is
introduced by sequencing it as the first slice. It touches every backbone
activity in the story map (`story-map.md`) exactly once, the minimum
required for a complete walking skeleton per `nw-user-story-mapping`.

## Demo Script

```
$ bm save https://kube.io/docs/failover --tag k8s
Saved [a1b2] "https://kube.io/docs/failover" (tag: k8s)

$ bm find k8s failover
[a1b2] https://kube.io/docs/failover  (tag: k8s, saved 12 days ago)

$ bm share a1b2
Share this with your team:
https://kube.io/docs/failover  (via bm, tag: k8s)
```

## Risk Carried Forward

Local-first storage and the share mechanism are unverified by an
engineering feasibility spike (DISCOVER G4 FAIL, accepted as risk — see
`docs/feature/bookmark-cli/discuss/wave-decisions.md`). This slice is the
first place that assumption becomes concrete; DESIGN wave should confirm it
holds before implementation.

## Dependencies

None — this is the first slice.
