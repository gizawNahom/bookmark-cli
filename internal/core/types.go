// Package core holds bookmark-cli's Functional Core: pure business-logic types and functions
// with no I/O (per CLAUDE.md "Development Paradigm" and ADR-006 Plan-value pattern).
//
// __SCAFFOLD__ = true -- created by DISTILL (nw-acceptance-designer) as a RED-ready scaffold.
// Types carry no behavior and never panic; only the functions in the sibling files (validator.go,
// normalizer.go, duplicate.go, planner.go, matcher.go, formatter.go) are scaffolded to panic.
package core

import "time"

// Record is a single saved bookmark, as read back from the store.
type Record struct {
	ID      string
	URL     string
	Tag     string // "" means untagged
	SavedAt time.Time
}

// ValidationResult is the pure outcome of URLValidator.Validate.
type ValidationResult struct {
	Valid  bool
	Reason string // human-readable reason when Valid == false
}

// NormalizedTag is the pure outcome of TagNormalizer.Normalize.
type NormalizedTag struct {
	Value string
}

// DuplicateVerdict is the pure outcome of DuplicateDetector.Check.
type DuplicateVerdict struct {
	IsDuplicate    bool
	ExistingID     string
	ExistingHasTag bool
	TagDiffers     bool // true when the incoming tag differs from (or adds to) the existing untagged record
}

// SavePlanKind enumerates the three possible outcomes of SavePlanner.Plan (ADR-006 Plan-value
// pattern) -- New | Duplicate | TagUpdate.
type SavePlanKind string

const (
	PlanNew       SavePlanKind = "new"
	PlanDuplicate SavePlanKind = "duplicate"
	PlanTagUpdate SavePlanKind = "tag_update"
)

// SavePlan is the pure decision returned by SavePlanner.Plan. BookmarkWriter.Execute(plan) is the
// only impure step allowed to act on it -- SavePlanner itself never writes.
type SavePlan struct {
	Kind       SavePlanKind
	URL        string
	Tag        string
	ExistingID string // populated for Duplicate and TagUpdate
}

// RankedMatch pairs a Record with its relevance score from Matcher.Rank.
type RankedMatch struct {
	Record Record
	Score  float64
}

// RankedMatches is the pure outcome of Matcher.Rank, including an optional closest-tag
// suggestion for the no-match case (US-06).
type RankedMatches struct {
	Matches      []RankedMatch
	SuggestedTag string // "" when no close tag exists
}

// ShareSnippet is the pure outcome of SnippetFormatter.Format.
type ShareSnippet struct {
	Text string
}
