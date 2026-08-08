package main

import (
	"strings"

	"bookmark-cli/internal/core"
)

// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer). These helpers delegate to the
// core (pure) scaffolds, which panic "not yet implemented -- RED scaffold"; runCommand() in
// main.go converts that panic into a controlled CLI failure. Real rendering logic (including the
// US-04 tag-discoverability hint and the accessibility text-prefix rule, brief.md Section 9)
// lands here during DELIVER GREEN.

func planSaveOrFail(url, tag string, existing []core.Record) core.SavePlan {
	return core.PlanSave(url, tag, existing)
}

func rankOrFail(query string, candidates []core.Record) core.RankedMatches {
	return core.RankMatches(query, candidates)
}

func formatOrFail(rec core.Record) core.ShareSnippet {
	return core.FormatSnippet(rec)
}

func joinArgs(args []string) string {
	return strings.Join(args, " ")
}

// renderSaveConfirmation formats the "Saved [id] url (tag: x)" line plus, when no tag was given,
// the US-04 discoverability hint. Every status line carries a text prefix (accessibility rule,
// brief.md Section 9) independent of any future ANSI color decoration.
func renderSaveConfirmation(rec core.Record, plan core.SavePlan) string {
	panic("main.renderSaveConfirmation not yet implemented -- RED scaffold")
}

// renderFindResult formats ranked matches, or the distinguishable no-match / empty-store /
// closest-tag-suggestion messages (US-06).
func renderFindResult(matches core.RankedMatches) string {
	panic("main.renderFindResult not yet implemented -- RED scaffold")
}
