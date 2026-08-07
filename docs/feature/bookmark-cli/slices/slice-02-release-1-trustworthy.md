# Slice 02: Release 1 — Trustworthy Under Real Conditions

Feature: bookmark-cli | Wave: DISCUSS | Type: Release 1 (depends on Slice 01)

## Outcome

The walking skeleton survives the conditions DISCOVER evidence says actually
occur: tag-flag confusion (the one known usability failure), bad/duplicate
input at save time, and no-match searches at find time.

## Stories in This Slice

| Story | Fixes | Job |
|---|---|---|
| US-04 | Tag-flag discoverability (Aisha's tested failure, ~20s vs <10s target) | JTBD-CAPTURE |
| US-05 | Malformed/duplicate URL handling at save time | JTBD-CAPTURE |
| US-06 | Helpful no-results experience at find time | JTBD-LOCATE |

Full story detail: `docs/feature/bookmark-cli/discuss/user-stories.md`.

## Why This Order

US-04 fixes a *named, evidenced* usability defect against a command already
in the walking skeleton — higher value than any new surface area. US-05 and
US-06 close the happy-path bias gap the walking skeleton alone would leave
(see `story-map.md` Priority Rationale, item 3): without them, US-01/US-02
would only have UAT coverage for the happy path.

## Demo Script (each story independently demoable)

```
# US-04
$ bm save https://kube.io/docs/failover
Saved [a1b2] "https://kube.io/docs/failover"
tip: add --tag <name> to make this easier to find later

# US-05
$ bm save https://kube.io/docs/failover --tag k8s
already saved as [a1b2]

# US-06
$ bm find gRPC-retry-policy
no matches for 'gRPC-retry-policy' -- did you mean 'grpc'?
```

## Dependencies

Depends on Slice 01 (US-01, US-02) for the save/find commands being
extended. Outcome KPI 4 (tag comprehension guardrail) requires a usability
re-test using the DISCOVER Phase 3 protocol before this slice is considered
done.
