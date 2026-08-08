package bookmarkcli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// @walking_skeleton @driving_port @us-01 @us-02 @us-03
//
// Scenario: A terminal-first engineer saves a link, finds it again, and shares it with a
// teammate who needs zero install
//
//	Given Nadia Petrova is in her terminal mid-debugging session
//	When she runs "bm save https://kube.io/docs/failover --tag k8s"
//	Then she sees a confirmation with a bookmark id and the tag
//	When she runs "bm find k8s failover"
//	Then she sees the matching link with the same bookmark id
//	When she runs "bm share <id>"
//	Then she sees a copy-paste-ready snippet requiring zero install for the recipient
//
// This is the exact 3-task sequence DISCOVER Phase 3 already task-tested (93% completion, n=5) --
// slice-01-walking-skeleton.md formalizes it as the first slice. WS layer (Layered Test
// Discipline table): example-only, traditional assertions -- proves wiring end-to-end through the
// real `bm` binary (subprocess), not a fixture-shape shortcut.
//
// Assertion convention: testify require/assert (project standard, see feature-delta.md DISTILL
// section "Assertion Convention"). require.* halts the scenario on a precondition/step failure
// that would make later steps meaningless; assert.* is used for independent, collectible outcome
// checks within a single Then.
func TestWalkingSkeleton_SaveFindShare(t *testing.T) {
	cli := NewCLI(t)

	saveResult := cli.Save("https://kube.io/docs/failover", "k8s")
	require.Equal(t, 0, saveResult.ExitCode, "bm save exited non-zero, stderr: %s", saveResult.Stderr)
	assert.Contains(t, saveResult.Stdout, "Saved", "expected a save confirmation naming the outcome")
	assert.Contains(t, saveResult.Stdout, "k8s", "expected the tag to appear in the confirmation")

	findResult := cli.Find("k8s failover")
	require.Equal(t, 0, findResult.ExitCode, "bm find exited non-zero, stderr: %s", findResult.Stderr)
	assert.Contains(t, findResult.Stdout, "https://kube.io/docs/failover", "expected the saved link in find results")

	id := extractBookmarkID(t, saveResult.Stdout)

	shareResult := cli.Share(id)
	require.Equal(t, 0, shareResult.ExitCode, "bm share exited non-zero, stderr: %s", shareResult.Stderr)
	assert.Contains(t, shareResult.Stdout, "https://kube.io/docs/failover",
		"expected the share snippet to contain the saved URL exactly")
	// US-03 AC: zero-install for the recipient is a property of the OUTPUT SHAPE (plain text,
	// no attachment/link-to-install), asserted here as "the snippet is plain text containing the
	// URL" -- the only observable proxy the sender's own CLI invocation can assert; the actual
	// zero-install guardrail is pilot-observation-only per kpi-contracts.yaml KPI-3.
}

// extractBookmarkID pulls the bookmark id out of a save confirmation line, e.g.
// "Saved [a1b2] https://kube.io/docs/failover (tag: k8s)" -> "a1b2". Business-language helper,
// not a technical parsing detail leaked into the scenario body above.
func extractBookmarkID(t *testing.T, saveStdout string) string {
	t.Helper()
	start := strings.Index(saveStdout, "[")
	end := strings.Index(saveStdout, "]")
	require.True(t, start != -1 && end != -1 && end > start,
		"could not locate a bookmark id in save confirmation: %q", saveStdout)
	return saveStdout[start+1 : end]
}
