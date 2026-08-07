# Shared Artifacts Registry: bookmark-cli

Source journey: `docs/feature/bookmark-cli/discuss/journey-save-find-share.yaml`

```yaml
shared_artifacts:
  bookmark_id:
    source_of_truth: "local bookmark store record, assigned at `bm save` time"
    consumers: ["bm save (created)", "bm find (displayed in results)", "bm share (argument + echoed in output)"]
    owner: "bookmark-cli local storage layer"
    integration_risk: "HIGH -- if the ID shown by `bm find` doesn't match what `bm share` accepts, the walking skeleton breaks end-to-end"
    validation: "Integration checkpoint in journey Step 3: id round-trips save -> find -> share unchanged"

  tag:
    source_of_truth: "user input at `bm save` time, normalized by the CLI"
    consumers: ["bm save (echoed in confirmation)", "bm find (matched against + displayed)", "bm share (displayed in snippet)"]
    owner: "bookmark-cli local storage layer"
    integration_risk: "MEDIUM -- inconsistent normalization (e.g. `k8s` vs `kubernetes`) causes false negatives in find; DISCOVER evidence (Nadia, transcript #6) shows this already bit her ad hoc script"
    validation: "Tag normalization tested at save time; find matches against normalized form, not raw input"

  match_list:
    source_of_truth: "local bookmark store, filtered/ranked at `bm find` time against records created by `bm save`"
    consumers: ["bm find (displayed results)", "bm share (user selects an id from this list)"]
    owner: "bookmark-cli local storage layer"
    integration_risk: "MEDIUM -- ranking logic changes must not silently drop valid matches"
    validation: "Every returned record's bookmark_id must exist in the store and be shareable"

  share_snippet:
    source_of_truth: "generated at `bm share` time directly from the resolved bookmark record (url + tag)"
    consumers: ["recipient's terminal/chat -- must require zero install"]
    owner: "bookmark-cli share command"
    integration_risk: "HIGH -- any transformation between the stored record and the snippet risks drift from source, and any recipient-side dependency violates the hard zero-install constraint (Aisha, DISCOVER #8)"
    validation: "Snippet content diffed against the record shown in the preceding `bm find` step during acceptance testing"
```

## Quality Gate Validation

- **Journey completeness**: PASS — all 3 steps have goals, CLI commands, emotional annotations, shared artifacts, and integration checkpoints (see `journey-save-find-share.yaml`).
- **Emotional coherence**: PASS — arc is anxious -> anxious/time-pressured -> confident -> proud, no jarring transitions; CAPTURE's quick win buffers the higher-anxiety LOCATE step; LOCATE's resolution buffers the higher-stakes SHARE step.
- **Horizontal integration**: PASS — all 4 shared artifacts above have a single source of truth and documented consumers; `bookmark_id` integration checkpoint is explicit and testable.
- **CLI UX compliance**: PASS — `bm <verb> [args]` pattern consistent across all 3 commands; `<id>` vocabulary consistent across find output and share input (no "record"/"entry"/"key" synonyms introduced).
