package core

// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer).
const matcherScaffold = true

// RankMatches is a pure function: ranks candidates against a query, typo-tolerant (US-02 AC).
// When there are zero matches, RankedMatches.SuggestedTag carries the closest existing tag, if
// any exists close enough to suggest (US-06 AC).
func RankMatches(query string, candidates []Record) RankedMatches {
	panic("core.RankMatches not yet implemented -- RED scaffold")
}
