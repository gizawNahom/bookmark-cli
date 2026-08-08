package bookmarkcli_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Adapter-coverage scenarios (Mandate 6): every driven adapter gets at least one @real-io
// @adapter-integration scenario exercising it with real I/O, not synthetic/mocked data.
//
// | Adapter               | @real-io scenario                                  |
// |------------------------|----------------------------------------------------|
// | SQLiteBookmarkStore    | every save/find/share scenario (real file, t.TempDir()) |
// | FileBackupAdapter      | TestSave_CreatesBackupSnapshot (below)              |
// | FileUsageLogAdapter    | TestSave_WithTelemetryEnabled_RecordsUsageEvent (below) |
// | NoOpUsageLogAdapter    | real trivial adapter, no business logic -- covered implicitly by every scenario that does NOT call WithTelemetryEnabled() (default composition-root path) |

// @real-io @adapter-integration @us-01
//
// Scenario: A successful save triggers a rotating backup snapshot
//
//	Given Nadia Petrova has an empty bookmark store
//	When she runs "bm save https://kube.io/docs/failover --tag k8s"
//	Then a backup snapshot file exists under the backup directory
func TestSave_CreatesBackupSnapshot(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t)

	result := cli.Save("https://kube.io/docs/failover", "k8s")
	require.Equal(t, 0, result.ExitCode, "bm save exited non-zero, stderr: %s", result.Stderr)

	entries, err := os.ReadDir(cli.BackupDir())
	require.NoError(t, err, "expected the backup directory to exist after a successful save")
	assert.NotEmpty(t, entries, "expected at least one backup snapshot file under %s", cli.BackupDir())
}

// @real-io @adapter-integration
//
// Scenario: With telemetry opted in, a save is recorded to the local usage log
//
//	Given Nadia Petrova has opted in to local telemetry
//	When she runs "bm save https://kube.io/docs/failover --tag k8s"
//	Then a "bm.save" event is appended to her local usage log
//	And the event never contains the URL or tag content
func TestSave_WithTelemetryEnabled_RecordsUsageEvent(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t).WithTelemetryEnabled()

	result := cli.Save("https://kube.io/docs/failover", "k8s")
	require.Equal(t, 0, result.ExitCode, "bm save exited non-zero, stderr: %s", result.Stderr)

	raw, err := os.ReadFile(cli.UsageLogPath())
	require.NoError(t, err, "expected a usage.log file when telemetry is enabled")
	logContent := string(raw)
	assert.Contains(t, logContent, "bm.save", "expected a bm.save event in the usage log")
	assert.NotContains(t, logContent, "kube.io", "usage log must never contain URL content (privacy-by-design)")
}

// @real-io @adapter-integration @error @us-07-nfr
//
// Scenario: A degraded (read-only) filesystem refuses the write instead of failing silently
//
//	Given ${data_dir} is mounted read-only (environments.yaml "degraded-filesystem")
//	When Nadia Petrova runs "bm save https://kube.io/docs/failover --tag k8s"
//	Then the command refuses the operation with a health.startup.refused message
//	And no partial write occurs
//
// Exercises the ADR-007 Probe() fault-injection contract against a real read-only directory, not
// a mocked filesystem -- environments.yaml explicitly names this matrix entry as the contract
// this scenario validates.
func TestSave_DegradedFilesystem_RefusesCleanly(t *testing.T) {
	t.Skip("pending -- enable one scenario at a time per DELIVER RED->GREEN cycle, ADR-025")
	cli := NewCLI(t).WithReadOnlyDataDir()

	result := cli.Save("https://kube.io/docs/failover", "k8s")

	require.NotEqual(t, 0, result.ExitCode, "expected a non-zero exit on a read-only data directory")
	combined := result.Stdout + result.Stderr
	assert.Contains(t, combined, "health.startup.refused", "expected a health.startup.refused message")
}
