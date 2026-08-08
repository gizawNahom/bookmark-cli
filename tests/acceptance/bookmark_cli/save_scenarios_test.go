package bookmarkcli_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	statedelta "bookmark-cli/tests/common"
)

// increasedBy is a small local Predicate constructor (state_delta.go ships the general-purpose
// predicates; this one is specific enough to this feature's "match count" universe entry that it
// stays local rather than being promoted to the shared port).
func increasedBy(n int) statedelta.Predicate {
	return func(before, after any) (bool, string) {
		b, _ := before.(int)
		a, _ := after.(int)
		if a == b+n {
			return true, ""
		}
		return false, fmt.Sprintf("expected count to increase by %d from %d, got %d", n, b, a)
	}
}

// @us-01
//
// Scenario: Link is saved without a tag and is still retrievable
//
//	Given Sam Okafor is mid-incident and wants to capture a link quickly
//	When he runs "bm save https://wiki.internal/postmortem-2024-03"
//	Then he sees a confirmation with a bookmark id and no tag shown
//	And the link is still saved and retrievable by keyword
func TestSave_WithoutTag_StillSavesAndRetrievable(t *testing.T) {
	cli := NewCLI(t)
	before := captureFindUniverse(cli.Find("postmortem"))

	result := cli.Save("https://wiki.internal/postmortem-2024-03", "")

	after := captureFindUniverse(cli.Find("postmortem"))
	statedelta.AssertStateDelta(t, before, after,
		[]string{"find.match_count"},
		map[string]statedelta.Predicate{"find.match_count": increasedBy(1)},
	)

	require.Equal(t, 0, result.ExitCode, "bm save exited non-zero, stderr: %s", result.Stderr)
	assert.Contains(t, result.Stdout, "[", "expected a bookmark id in the confirmation")
}

// @us-01 @us-05 @error
//
// Scenario: Saving the same URL twice is caught, not silently duplicated
//
//	Given Marco Alves already saved "https://kube.io/docs/failover" last month
//	When he runs "bm save https://kube.io/docs/failover --tag k8s" again
//	Then he sees a message that this URL is already saved, with its existing bookmark id
//	And no second entry is created
func TestSave_ExactDuplicate_DetectedNotDuplicated(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://kube.io/docs/failover", Tag: "k8s"})
	before := captureFindUniverse(cli.Find("failover"))

	result := cli.Save("https://kube.io/docs/failover", "k8s")

	after := captureFindUniverse(cli.Find("failover"))
	statedelta.AssertStateDelta(t, before, after,
		[]string{"find.match_count"},
		map[string]statedelta.Predicate{"find.match_count": statedelta.Unchanged()},
	)

	assert.Contains(t, result.Stdout, "already saved",
		"expected an 'already saved' message with the existing id")
}

// @us-05 @error
//
// Scenario: A malformed URL is rejected with a clear message
//
//	Given Nadia Petrova runs "bm save not-a-url --tag misc"
//	Then she sees a message that this does not look like a valid URL
//	And nothing is saved
func TestSave_MalformedURL_RejectedWithClearMessage(t *testing.T) {
	cli := NewCLI(t)
	before := captureFindUniverse(cli.Find("misc"))

	result := cli.Save("not-a-url", "misc")

	after := captureFindUniverse(cli.Find("misc"))
	statedelta.AssertStateDelta(t, before, after,
		[]string{"find.match_count"},
		map[string]statedelta.Predicate{"find.match_count": statedelta.Unchanged()},
	)

	require.NotEqual(t, 0, result.ExitCode, "expected a non-zero exit for a malformed URL, got stdout: %q", result.Stdout)
	combined := result.Stdout + result.Stderr
	assert.True(t,
		strings.Contains(combined, "valid URL") || strings.Contains(combined, "doesn't look like"),
		"expected a specific 'not a valid URL' message, got: %q", combined)
}

// @us-05
//
// Scenario: Saving the same URL with a new tag updates rather than duplicates
//
//	Given Sam Okafor already saved "https://kube.io/docs/failover" untagged
//	When he runs "bm save https://kube.io/docs/failover --tag k8s"
//	Then he is offered to add the "k8s" tag to the existing bookmark
//	And no second entry is created
func TestSave_SameURLNewTag_OffersTagUpdateNotDuplicate(t *testing.T) {
	cli := NewCLI(t).WithExistingStore(Bookmark{URL: "https://kube.io/docs/failover", Tag: ""})
	before := captureFindUniverse(cli.Find("failover"))

	result := cli.Save("https://kube.io/docs/failover", "k8s")

	after := captureFindUniverse(cli.Find("failover"))
	statedelta.AssertStateDelta(t, before, after,
		[]string{"find.match_count"},
		map[string]statedelta.Predicate{"find.match_count": statedelta.Unchanged()},
	)

	assert.Contains(t, result.Stdout, "k8s", "expected the tag-update offer to mention the new tag 'k8s'")
}

// @us-05
//
// Scenario: A URL with query parameters still saves correctly
//
//	Given Priti Desai runs "bm save https://wiki.internal/runbook?id=42&version=3 --tag oncall"
//	Then the link is saved with the full URL including query parameters intact
func TestSave_URLWithQueryParams_SavesWithFullFidelity(t *testing.T) {
	cli := NewCLI(t)
	url := "https://wiki.internal/runbook?id=42&version=3"

	saveResult := cli.Save(url, "oncall")
	require.Equal(t, 0, saveResult.ExitCode, "bm save exited non-zero, stderr: %s", saveResult.Stderr)

	findResult := cli.Find("runbook")
	assert.Contains(t, findResult.Stdout, url, "expected the full URL with query params intact in find output")
}

// @us-04
//
// Scenario: A hint nudges the user toward tagging without failing first
//
//	Given Aisha Rahman runs "bm save https://kube.io/docs/failover" with no --tag
//	When the command completes
//	Then she sees a one-line hint suggesting "--tag <name>" to make it easier to find later
func TestSave_WithoutTag_ShowsDiscoverabilityHint(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t)

	result := cli.Save("https://kube.io/docs/failover", "")

	assert.Contains(t, result.Stdout, "--tag", "expected a hint mentioning --tag")
}

// @us-04
//
// Scenario: Help text shows a concrete example, not just abstract syntax
//
//	Given Priti Desai runs "bm save --help"
//	Then she sees an example line such as "bm save <url> --tag k8s"
func TestSaveHelp_ShowsConcreteExample(t *testing.T) {
	cli := NewCLI(t)

	result := cli.SaveHelp()

	assert.Contains(t, result.Stdout, "--tag k8s", "expected a concrete example with --tag k8s in help output")
}

// @us-04 @error
//
// Scenario: A near-miss flag name is corrected with a suggestion
//
//	Given Marco Alves runs "bm save <url> --tags k8s" (typo: plural)
//	Then he sees a suggestion "did you mean --tag?"
//	And the command does not fail silently or with a bare parse error
func TestSave_NearMissFlag_SuggestsDidYouMean(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t)

	result := cli.SaveWithFlag("https://kube.io/docs/failover", "--tags", "k8s")

	combined := result.Stdout + result.Stderr
	// Precise match on "did you mean" only -- an earlier version of this assertion also accepted
	// a bare "--tag" substring, which trivially matched Cobra's default "unknown flag: --tags"
	// error (since "--tags" itself contains "--tag" as a prefix) and made this scenario pass
	// without any did-you-mean implementation. Caught during the pre-DELIVER fail-for-the-right-
	// reason gate; documented in distill/red-classification.md.
	assert.Contains(t, strings.ToLower(combined), "did you mean",
		"expected a 'did you mean --tag?' suggestion, got: %q", combined)
}

// @us-01
//
// Scenario: Confirmation appears fast enough to not break flow
//
//	Given Nadia Petrova is mid-flow in a tmux session
//	When she runs "bm save <url> --tag <tag>"
//	Then the confirmation appears within a responsive feedback window
//
// Budget set to 2s (not the 100ms perceived-save target) per the F-004 flakiness-under-load
// guardrail (timing assertions in acceptance tests must use a generous budget); the <100ms
// perceived-save target itself is validated by a benchmark in DELIVER, not this subprocess test.
func TestSave_ConfirmationIsResponsive(t *testing.T) {
	cli := NewCLI(t)

	start := time.Now()
	result := cli.Save("https://kube.io/docs/failover", "k8s")
	elapsed := time.Since(start)

	require.Equal(t, 0, result.ExitCode, "bm save exited non-zero, stderr: %s", result.Stderr)
	assert.LessOrEqual(t, elapsed, 2*time.Second,
		"expected bm save to respond within 2s (generous CI budget), took %s", elapsed)
}
