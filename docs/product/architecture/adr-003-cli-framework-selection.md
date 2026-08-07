# ADR-003: CLI Framework Selection

## Status
Accepted

## Context
US-04 requires the `--tag` flag and overall command syntax to be discoverable within 10 seconds
(remediating a tested ~20s failure), including concrete `--help` examples and "did you mean"-style
correction for near-miss input.

## Decision
**Cobra** (`spf13/cobra`).

## Alternatives Considered

**urfave/cli** — simpler API and smaller footprint than Cobra. Rejected: no built-in
command-suggestion support, and, more importantly, it doesn't carry the same convention
familiarity — the target persona already has muscle memory for Cobra's help/flag conventions
from kubectl, Docker CLI, and Terraform CLI, which directly supports the US-04 discoverability
goal in a way a less-familiar framework's conventions would not.

**Kong** — declarative, struct-tag-based configuration, appealing for its conciseness. Rejected:
smaller community, no built-in suggestion support, and no persona-familiarity advantage over
Cobra to offset the smaller ecosystem.

## Consequences

**Positive**: matches conventions the persona already knows from kubectl/docker/terraform,
directly supporting US-04's <10s comprehension target. Built-in command-suggestion support
reduces custom code for near-miss command correction. Very large community/maturity (MIT
license) reduces long-term maintenance risk for a solo maintainer.

**Negative**: Cobra's built-in suggestion support is at the command level, not the flag level —
the `--tags` → "did you mean --tag?" correction in US-04 will need a small custom layer on top
(flagged as an open implementation detail in `brief.md` Section 19, not a blocker to this ADR).
