package core

import (
	"sort"
	"strings"
)

// maxTermEditDistance bounds how many single-character edits (insert/delete/substitute) a query
// term may be from a candidate word and still count as a typo-tolerant match (US-02 AC).
const maxTermEditDistance = 2

// RankMatches is a pure function: ranks candidates against a query, typo-tolerant (US-02 AC).
// A candidate matches when every whitespace-separated query term matches at least one word drawn
// from the candidate's tag or URL -- exactly, as a substring, or within maxTermEditDistance edits
// (fuzzy/typo-tolerant, US-02 Scenario 3). Matches are ordered by descending relevance score, so
// multiple matches under a shared tag are all surfaced rather than collapsed to one guess
// (US-02 Scenario 4).
//
// When there are zero matches, RankedMatches.SuggestedTag carries the closest existing tag, if
// any exists close enough to suggest (US-06 AC).
func RankMatches(query string, candidates []Record) RankedMatches {
	terms := queryTerms(query)
	if len(terms) == 0 {
		return RankedMatches{}
	}

	matches := make([]RankedMatch, 0, len(candidates))
	for _, candidate := range candidates {
		if score, matched := scoreCandidate(terms, candidate); matched {
			matches = append(matches, RankedMatch{Record: candidate, Score: score})
		}
	}

	if len(matches) == 0 {
		return RankedMatches{SuggestedTag: closestTag(terms, candidates)}
	}

	sortByRelevance(matches)
	return RankedMatches{Matches: matches}
}

// queryTerms splits and lowercases the raw query into whitespace-separated terms.
func queryTerms(query string) []string {
	return strings.Fields(strings.ToLower(query))
}

// isWordSeparator splits on anything that is not a lowercase letter or digit, so both candidate
// text (URL/tag) and query terms tokenize consistently on punctuation such as hyphens.
func isWordSeparator(r rune) bool {
	return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
}

// splitWords tokenizes a lowercased string into word fragments using isWordSeparator.
func splitWords(text string) []string {
	return strings.FieldsFunc(strings.ToLower(text), isWordSeparator)
}

// candidateWords returns every matchable word from a record's tag and URL, lowercased.
func candidateWords(candidate Record) []string {
	words := splitWords(candidate.URL)
	if candidate.Tag != "" {
		words = append(words, strings.ToLower(candidate.Tag))
	}
	return words
}

// scoreCandidate requires every query term to match at least one candidate word, and returns the
// average per-term match strength as the candidate's overall relevance score.
func scoreCandidate(terms []string, candidate Record) (float64, bool) {
	words := candidateWords(candidate)

	total := 0.0
	for _, term := range terms {
		termScore, matched := bestWordScore(term, words)
		if !matched {
			return 0, false
		}
		total += termScore
	}
	return total / float64(len(terms)), true
}

// bestWordScore finds the strongest match for a single query term among a candidate's words:
// exact match scores highest, substring match next, and a fuzzy (edit-distance-bounded) match
// scores lowest but still counts -- this is what surfaces a typo like "terrafrom" against the
// stored tag "terraform".
func bestWordScore(term string, words []string) (float64, bool) {
	best := 0.0
	matched := false
	for _, word := range words {
		switch {
		case word == term:
			return 1.0, true
		case strings.Contains(word, term):
			// Only "word contains term" counts as a substring match (e.g. query "fail" against
			// word "failover"). The reverse direction ("term contains word") is deliberately
			// excluded: a short candidate word (e.g. tag "grpc") would otherwise substring-match
			// against almost any longer, unrelated compound query term (e.g. "grpc-retry-policy"),
			// producing false positives that mask genuine no-match + closest-tag-suggestion cases
			// (US-06 AC).
			if 0.8 > best {
				best = 0.8
				matched = true
			}
		default:
			distance := levenshtein(term, word)
			if distance <= maxTermEditDistance {
				score := 1.0 - float64(distance)/float64(maxRuneCount(term, word)+1)
				if score > best {
					best = score
					matched = true
				}
			}
		}
	}
	return best, matched
}

// closestTag returns the existing tag closest (by edit distance) to any query term, for the
// no-match "did you mean" suggestion (US-06). Returns "" when no tag is close enough.
//
// Compares against both whole query terms (typo case: "teraform" vs tag "terraform") and their
// word-fragments split on punctuation (compound case: "grpc-retry-policy" contains fragment
// "grpc", close to tag "grpc") -- a query term may be a hyphenated compound that embeds a
// near-exact tag name even though the term as a whole is far from that tag by edit distance.
func closestTag(terms []string, candidates []Record) string {
	const maxSuggestDistance = 2

	tokens := make([]string, 0, len(terms)*2)
	tokens = append(tokens, terms...)
	for _, term := range terms {
		tokens = append(tokens, splitWords(term)...)
	}

	best := ""
	bestDistance := maxSuggestDistance + 1
	seen := make(map[string]bool)
	for _, candidate := range candidates {
		tag := strings.ToLower(candidate.Tag)
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true

		for _, token := range tokens {
			if d := levenshtein(token, tag); d <= maxSuggestDistance && d < bestDistance {
				bestDistance = d
				best = candidate.Tag
			}
		}
	}
	return best
}

// sortByRelevance orders matches by descending score, breaking ties by original candidate order
// for deterministic, stable output across repeated calls with the same input.
func sortByRelevance(matches []RankedMatch) {
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})
}

// levenshtein computes the classic single-character-edit distance between two strings.
func levenshtein(a, b string) int {
	rowLen := len(b) + 1
	prev := make([]int, rowLen)
	curr := make([]int, rowLen)
	for j := 0; j < rowLen; j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = minInt(prev[j]+1, minInt(curr[j-1]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

func maxRuneCount(a, b string) int {
	if len(a) > len(b) {
		return len(a)
	}
	return len(b)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
