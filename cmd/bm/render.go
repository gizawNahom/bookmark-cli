package main

import (
	"fmt"
	"strings"
	"time"

	"bookmark-cli/internal/core"
)

// planSaveOrFail delegates to the pure core decision (ADR-006 Plan-value pattern).
func planSaveOrFail(url, tag string, existing []core.Record) core.SavePlan {
	return core.PlanSave(url, tag, existing)
}

// rankOrFail delegates to the pure core ranking engine (ADR-006): typo-tolerant matching (US-02)
// and closest-tag suggestion on no-match (US-06) both live in core.RankMatches, implemented in
// step 02-01. This shell function exists only to keep the find-command call site symmetric with
// planSaveOrFail/formatOrFail.
func rankOrFail(query string, candidates []core.Record) core.RankedMatches {
	return core.RankMatches(query, candidates)
}

// formatOrFail is the walking-skeleton-minimal share-snippet formatter. Exact-fidelity formatting
// (US-03) is core.FormatSnippet's responsibility, implemented and wired in step 03-01.
func formatOrFail(rec core.Record) core.ShareSnippet {
	if rec.Tag == "" {
		return core.ShareSnippet{Text: rec.URL}
	}
	return core.ShareSnippet{Text: fmt.Sprintf("%s (tag: %s)", rec.URL, rec.Tag)}
}

func joinArgs(args []string) string {
	return strings.Join(args, " ")
}

// renderSaveConfirmation formats the outcome line for `bm save`. Every status/error line carries
// a text prefix independent of any future ANSI color decoration (accessibility rule, brief.md
// Section 9). Outcome varies by plan.Kind (ADR-006 Plan-value pattern: the message is derived
// from the pure decision, never inferred from a write that may or may not have happened):
//   - PlanNew: "Saved [id] url (tag: x)"
//   - PlanDuplicate: already-saved notice including the existing bookmark id
//   - PlanTagUpdate: an offer to add the new tag to the existing bookmark, mentioning the tag
func renderSaveConfirmation(rec core.Record, plan core.SavePlan) string {
	switch plan.Kind {
	case core.PlanDuplicate:
		return fmt.Sprintf("%s is already saved as [%s] -- no new entry created", rec.URL, rec.ID)
	case core.PlanTagUpdate:
		return fmt.Sprintf("[%s] %s is already saved -- add tag %q to it?", rec.ID, rec.URL, plan.Tag)
	default:
		if rec.Tag == "" {
			return fmt.Sprintf("Saved [%s] %s\n%s", rec.ID, rec.URL, discoverabilityHint())
		}
		return fmt.Sprintf("Saved [%s] %s (tag: %s)", rec.ID, rec.URL, rec.Tag)
	}
}

// discoverabilityHint nudges an untagged save toward tagging (US-04) without failing the command --
// a one-line, text-prefixed hint (accessibility rule, brief.md Section 9).
func discoverabilityHint() string {
	return "Hint: add --tag <name> next time to make this easier to find later"
}

// renderEmptyStoreMessage is the onboarding message shown when `bm find` runs against a genuinely
// empty store (US-06 AC). It is deliberately worded distinct from renderFindResult's no-match
// text -- "haven't saved" vs "no matches" -- so the two cases are never conflated. Text-prefixed,
// no ANSI-color-only status indicators (accessibility rule).
func renderEmptyStoreMessage() string {
	return "You haven't saved any links yet -- run `bm save <url>` to get started"
}

// renderFindResult formats ranked matches, or a "no matches" fallback for a non-empty store
// (US-06 AC). When RankedMatches carries a SuggestedTag (closest existing tag within edit-distance
// threshold, core.RankMatches), the message includes a "did you mean" nudge; otherwise it is a
// clean, non-blank "no matches found" message.
//
// Trailing newline is deliberately one-per-match (present only when there is >=1 match): the
// acceptance harness's state-delta Universe (captureFindUniverse, tests/acceptance/bookmark_cli/
// harness_test.go) derives find.match_count from counting "\n" in stdout, refined during DELIVER
// GREEN per that file's own comment. Keeping the no-match message newline-free is what makes that
// count land on 0 rather than a phantom 1.
func renderFindResult(query string, matches core.RankedMatches) string {
	if len(matches.Matches) == 0 {
		if matches.SuggestedTag != "" {
			return fmt.Sprintf("no matches for %q -- did you mean %q?", query, matches.SuggestedTag)
		}
		return "no matches found"
	}

	lines := make([]string, 0, len(matches.Matches))
	for _, m := range matches.Matches {
		if m.Record.Tag == "" {
			lines = append(lines, fmt.Sprintf("[%s] %s -- %s", m.Record.ID, m.Record.URL, savedAgo(m.Record.SavedAt)))
		} else {
			lines = append(lines, fmt.Sprintf("[%s] %s (tag: %s) -- %s", m.Record.ID, m.Record.URL, m.Record.Tag, savedAgo(m.Record.SavedAt)))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

// savedAgo renders the "saved N days ago" metadata (US-02 AC) from a record's saved timestamp.
// A record saved less than a day ago reads "saved today" rather than "saved 0 days ago".
func savedAgo(savedAt time.Time) string {
	days := int(time.Since(savedAt).Hours() / 24)
	if days <= 0 {
		return "saved today"
	}
	if days == 1 {
		return "saved 1 day ago"
	}
	return fmt.Sprintf("saved %d days ago", days)
}
