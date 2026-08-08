package main

import (
	"fmt"
	"strings"

	"bookmark-cli/internal/core"
)

// planSaveOrFail delegates to the pure core decision (ADR-006 Plan-value pattern).
func planSaveOrFail(url, tag string, existing []core.Record) core.SavePlan {
	return core.PlanSave(url, tag, existing)
}

// rankOrFail is the walking-skeleton-minimal candidate filter for `bm find`: a match requires
// every whitespace-separated query term to appear (case-insensitive) in the record's URL or tag.
// This is deliberately thin -- typo-tolerant ranking (US-02) and closest-tag suggestions (US-06)
// are core.RankMatches' responsibility, implemented and wired in step 02-01.
func rankOrFail(query string, candidates []core.Record) core.RankedMatches {
	terms := strings.Fields(strings.ToLower(query))

	var matches []core.RankedMatch
	for _, candidate := range candidates {
		if containsAllTerms(candidate, terms) {
			matches = append(matches, core.RankedMatch{Record: candidate, Score: 1})
		}
	}
	return core.RankedMatches{Matches: matches}
}

func containsAllTerms(candidate core.Record, terms []string) bool {
	haystack := strings.ToLower(candidate.URL + " " + candidate.Tag)
	for _, term := range terms {
		if !strings.Contains(haystack, term) {
			return false
		}
	}
	return true
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

// renderSaveConfirmation formats the "Saved [id] url (tag: x)" line. Every status/error line
// carries a text prefix ("Saved") independent of any future ANSI color decoration (accessibility
// rule, brief.md Section 9).
func renderSaveConfirmation(rec core.Record, plan core.SavePlan) string {
	if rec.Tag == "" {
		return fmt.Sprintf("Saved [%s] %s", rec.ID, rec.URL)
	}
	return fmt.Sprintf("Saved [%s] %s (tag: %s)", rec.ID, rec.URL, rec.Tag)
}

// renderFindResult formats ranked matches, or the "no matches found" fallback (US-06's richer
// closest-tag-suggestion / empty-store onboarding variants land in step 02-02).
//
// Trailing newline is deliberately one-per-match (present only when there is >=1 match): the
// acceptance harness's state-delta Universe (captureFindUniverse, tests/acceptance/bookmark_cli/
// harness_test.go) derives find.match_count from counting "\n" in stdout, refined during DELIVER
// GREEN per that file's own comment. Keeping "no matches found" newline-free is what makes that
// count land on 0 rather than a phantom 1.
func renderFindResult(matches core.RankedMatches) string {
	if len(matches.Matches) == 0 {
		return "no matches found"
	}

	lines := make([]string, 0, len(matches.Matches))
	for _, m := range matches.Matches {
		if m.Record.Tag == "" {
			lines = append(lines, fmt.Sprintf("[%s] %s", m.Record.ID, m.Record.URL))
		} else {
			lines = append(lines, fmt.Sprintf("[%s] %s (tag: %s)", m.Record.ID, m.Record.URL, m.Record.Tag))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}
