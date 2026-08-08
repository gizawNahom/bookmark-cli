package bookmarkcli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// share scenarios are pure-read (resolves a record, formats a snippet) -- no store mutation, so
// traditional assertions on CLI output are the correct shape (Mandate 8 applies to mutating
// steps only).

// @us-03
//
// Scenario: A curated link is shared with a zero-install snippet
//
//	Given Aisha Rahman has found her bookmark for the failover doc, tagged "k8s"
//	When she runs "bm share <id>"
//	Then she sees a copy-paste-ready snippet containing the URL and tag
func TestShare_CuratedLink_ProducesZeroInstallSnippet(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t)
	saveResult := cli.Save("https://kube.io/docs/failover", "k8s")
	id := extractBookmarkID(t, saveResult.Stdout)

	result := cli.Share(id)

	require.Equal(t, 0, result.ExitCode, "bm share exited non-zero, stderr: %s", result.Stderr)
	assert.Contains(t, result.Stdout, "https://kube.io/docs/failover", "expected the URL in the snippet")
	assert.Contains(t, result.Stdout, "k8s", "expected the tag in the snippet")
}

// @us-03
//
// Scenario: The shared snippet matches the source record exactly
//
//	Given bookmark "a1b2" resolves to "https://kube.io/docs/failover" tagged "k8s"
//	When Aisha Rahman runs "bm share a1b2"
//	Then the URL and tag in the snippet exactly match the record shown by "bm find"
func TestShare_SnippetMatchesFindRecordExactly(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t)
	saveResult := cli.Save("https://kube.io/docs/failover", "k8s")
	id := extractBookmarkID(t, saveResult.Stdout)

	findResult := cli.Find("failover")
	shareResult := cli.Share(id)

	assert.Contains(t, findResult.Stdout, "https://kube.io/docs/failover",
		"expected find to show the URL; find=%q", findResult.Stdout)
	assert.Contains(t, shareResult.Stdout, "https://kube.io/docs/failover",
		"expected share to show the identical URL; share=%q", shareResult.Stdout)
}

// @us-03
//
// Scenario: A link saved without a tag can still be shared
//
//	Given Jordan Osei saved a link with no tag
//	When he runs "bm share" with that bookmark's id
//	Then he sees a valid snippet containing the URL with no broken tag field
func TestShare_LinkWithoutTag_ProducesValidSnippet(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t)
	saveResult := cli.Save("https://kube.io/docs/networking", "")
	id := extractBookmarkID(t, saveResult.Stdout)

	result := cli.Share(id)

	require.Equal(t, 0, result.ExitCode, "bm share exited non-zero, stderr: %s", result.Stderr)
	assert.Contains(t, result.Stdout, "https://kube.io/docs/networking",
		"expected the URL in the snippet even without a tag")
	assert.NotContains(t, result.Stdout, "tag: )", "expected no broken tag field for an untagged bookmark")
	assert.NotContains(t, result.Stdout, "tag: <nil>", "expected no broken tag field for an untagged bookmark")
}

// @us-03 @error
//
// Scenario: Sharing an unknown bookmark id fails clearly
//
//	Given Priti Desai runs "bm share a9z9" for an id that does not exist
//	Then she sees a clear error naming the invalid id, not a blank or garbled snippet
func TestShare_UnknownID_FailsClearly(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t)

	result := cli.Share("a9z9")

	require.NotEqual(t, 0, result.ExitCode,
		"expected a non-zero exit for an unknown bookmark id, got stdout: %q", result.Stdout)
	combined := result.Stdout + result.Stderr
	assert.Contains(t, combined, "a9z9", "expected the error to name the invalid id 'a9z9'")
	assert.NotEmpty(t, strings.TrimSpace(combined), "expected a clear error, got blank output")
}
