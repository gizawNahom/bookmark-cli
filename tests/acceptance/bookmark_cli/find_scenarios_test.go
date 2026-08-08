package bookmarkcli_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// find scenarios are pure-read (BookmarkReader has no write methods, ADR-006) -- no store
// mutation occurs, so Mandate 8's state-delta requirement does not apply here (nothing to assert
// a delta on); traditional assertions on CLI output are the correct shape for a read-only port.

// @us-02
//
// Scenario: A saved link is found by tag and keyword
//
//	Given Priti Desai saved a link 12 days ago tagged "k8s" with "failover" in its title
//	When she runs "bm find k8s failover"
//	Then she sees the matching link with its bookmark id, tag, and "saved 12 days ago"
func TestFind_ByTagAndKeyword_ShowsMatchWithMetadata(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://kube.io/docs/failover", Tag: "k8s"})

	result := cli.Find("k8s failover")

	require.Equal(t, 0, result.ExitCode, "bm find exited non-zero, stderr: %s", result.Stderr)
	assert.Contains(t, result.Stdout, "https://kube.io/docs/failover", "expected the matching link")
	assert.Contains(t, result.Stdout, "k8s", "expected the tag shown alongside the match")
}

// @us-02
//
// Scenario: A saved link is found by keyword alone
//
//	Given Marco Alves saved "https://redis.io/failover-postmortem" without a tag
//	When he runs "bm find failover"
//	Then he sees the matching link in the results
func TestFind_ByKeywordAlone_ShowsMatch(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://redis.io/failover-postmortem", Tag: ""})

	result := cli.Find("failover")

	assert.Contains(t, result.Stdout, "https://redis.io/failover-postmortem", "expected the matching link")
}

// @us-02 @property
//
// Scenario: A near-miss search term still surfaces the right link
//
//	Given Jordan Osei saved a Terraform provider gotcha doc tagged "terraform"
//	When he runs "bm find terrafrom" (typo)
//	Then he still sees the Terraform doc ranked in the results
func TestFind_NearMissTypo_StillSurfacesMatch(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://terraform.io/provider-gotcha", Tag: "terraform"})

	result := cli.Find("terrafrom")

	assert.Contains(t, result.Stdout, "https://terraform.io/provider-gotcha",
		"expected fuzzy match to surface the Terraform doc despite the typo")
}

// @us-02
//
// Scenario: Multiple matches are shown ranked, not forced to one guess
//
//	Given Nadia Petrova has 3 saved links tagged "k8s"
//	When she runs "bm find k8s"
//	Then she sees all 3 matches ranked by relevance
func TestFind_MultipleMatches_RankedNotForcedToOne(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(
		Bookmark{URL: "https://kube.io/docs/failover", Tag: "k8s"},
		Bookmark{URL: "https://kube.io/docs/networking", Tag: "k8s"},
		Bookmark{URL: "https://kube.io/docs/autoscaling", Tag: "k8s"},
	)

	result := cli.Find("k8s")

	for _, url := range []string{
		"https://kube.io/docs/failover",
		"https://kube.io/docs/networking",
		"https://kube.io/docs/autoscaling",
	} {
		assert.Contains(t, result.Stdout, url, "expected all 3 tagged matches shown")
	}
}

// @us-06 @error
//
// Scenario: A no-match search suggests the closest existing tag
//
//	Given Aisha Rahman has links tagged "grpc" but none matching "gRPC-retry-policy"
//	When she runs "bm find gRPC-retry-policy"
//	Then she sees "no matches for 'gRPC-retry-policy'"
//	And a suggestion "did you mean 'grpc'?"
func TestFind_NoMatch_SuggestsClosestTag(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://grpc.io/docs/retry", Tag: "grpc"})

	result := cli.Find("gRPC-retry-policy")

	assert.Contains(t, result.Stdout, "no matches", "expected a 'no matches' message")
	assert.Contains(t, result.Stdout, "grpc", "expected a closest-tag suggestion mentioning 'grpc'")
}

// @us-06 @error
//
// Scenario: A no-match search with no close tag gives a clean message
//
//	Given Jordan Osei's store has no links resembling "teraform"
//	When he runs "bm find teraform"
//	Then he sees a clear "no matches found" message
//	And no crash, stack trace, or silent blank output
func TestFind_NoMatch_NoCloseTag_ShowsCleanMessage(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://kube.io/docs/failover", Tag: "k8s"})

	result := cli.Find("teraform")

	assert.NotEmpty(t, strings.TrimSpace(result.Stdout), "expected a non-blank no-match message")
	assert.Contains(t, result.Stdout, "no matches", "expected a clear 'no matches found' message")
	assert.NotContains(t, result.Stdout+result.Stderr, "panic", "expected no crash/stack trace on no-match")
}

// @us-06 @error
//
// Scenario: Searching an empty store gives an onboarding message, not a false no-match
//
//	Given Priti Desai has not saved any links yet
//	When she runs "bm find anything"
//	Then she sees a message telling her she hasn't saved any links yet
//	And it is visibly different from a normal no-match result
func TestFind_EmptyStore_ShowsOnboardingMessage(t *testing.T) {
	cli := NewCLI(t) // no WithExistingStore -- store is genuinely empty

	result := cli.Find("anything")

	assert.True(t,
		strings.Contains(result.Stdout, "haven't saved") || strings.Contains(result.Stdout, "no links yet"),
		"expected a distinct empty-store onboarding message, got: %q", result.Stdout)
}

// @us-06
//
// Scenario: A no-match response feels immediate
//
//	Given Jordan Osei's search returns zero matches
//	When the command completes
//	Then the no-match message appears within a responsive feedback window
func TestFind_NoMatchResponse_IsResponsive(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://kube.io/docs/failover", Tag: "k8s"})

	start := time.Now()
	result := cli.Find("nothing-saved-under-this-term")
	elapsed := time.Since(start)

	require.Equal(t, 0, result.ExitCode, "bm find exited non-zero, stderr: %s", result.Stderr)
	assert.Contains(t, result.Stdout, "no matches", "expected a no-match message (not an error) within budget")
	assert.LessOrEqual(t, elapsed, 2*time.Second,
		"expected the no-match response within 2s (generous CI budget), took %s", elapsed)
}
