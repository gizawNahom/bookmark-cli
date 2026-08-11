// Package bookmarkcli_test holds bookmark-cli's Tier A acceptance scenarios. Per the Architecture
// of Reference (Driving port -> real adapter) and Pillar 3 (app as in production), every scenario
// invokes the real `bm` binary via subprocess against a real SQLite store under a per-test
// ${data_dir} -- the production composition root in cmd/bm/main.go, never a hand-replicated
// wiring. This is the subprocess/FS acceptance layer (Layered Test Discipline table): example-only
// input mode (Mandate 9), state-delta + Universe required for mutating steps (Mandate 8).
package bookmarkcli_test

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// binPath is the path to the `bm` binary, built once in TestMain and shared read-only across all
// scenario tests in this package (composition-root re-use, not re-compiled per scenario).
var binPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "bm-acceptance-bin-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	binPath = filepath.Join(tmp, "bm")
	buildArgs := []string{"build", "-o", binPath}
	// When GOCOVERDIR is set (CI's integration-coverage step), build bm as a coverage-
	// instrumented binary so each subprocess invocation below flushes counters there on exit --
	// see go.dev/blog/integration-test-coverage. Left off locally (no GOCOVERDIR) to avoid the
	// instrumentation build-time cost on every `go test` run.
	if os.Getenv("GOCOVERDIR") != "" {
		buildArgs = append(buildArgs, "-cover")
	}
	buildArgs = append(buildArgs, "./cmd/bm")
	build := exec.Command("go", buildArgs...)
	build.Dir = repoRoot()
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build bm binary for acceptance tests: " + err.Error() + "\n" + string(out))
	}

	os.Exit(m.Run())
}

// repoRoot walks up from the current working directory (tests/acceptance/bookmark_cli when `go
// test` runs this package) to the module root, identified by go.mod.
func repoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found above " + dir)
		}
		dir = parent
	}
}

// Result is the CLI harness's observation of a single `bm` invocation -- fully driving-port
// observable (stdout/stderr/exit code), never an internal field.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// CLI is the acceptance-layer composition root: every scenario drives the system exclusively
// through this type's methods, which shell out to the real `bm` binary (Mandate 1: driving port
// only, zero internal-component imports from this package).
//
// t is testing.TB (not the concrete *testing.T) so the same harness serves both scenario tests
// and the performance benchmarks in save_bench_test.go -- *testing.B satisfies testing.TB too.
type CLI struct {
	t                testing.TB
	dataDir          string
	telemetryEnabled bool
}

// NewCLI constructs a CLI harness with an isolated ${data_dir} under t.TempDir() -- each scenario
// gets a fresh store, never sharing state with another scenario (Pillar 2's chained narrative
// operates WITHIN a scenario's own Given/When steps, not by leaking store state across tests).
func NewCLI(t testing.TB) *CLI {
	t.Helper()
	return &CLI{t: t, dataDir: t.TempDir()}
}

// WithTelemetryEnabled opts this scenario's store into the local usage-log adapter
// (FileUsageLogAdapter), matching the DEVOPS wave's opt-in-only telemetry design
// (kpi-contracts.yaml) -- selects the real adapter, never a fake, per the Architecture of
// Reference (driven-internal ports get real adapters).
func (c *CLI) WithTelemetryEnabled() *CLI {
	c.telemetryEnabled = true
	return c
}

// UsageLogPath returns where the FileUsageLogAdapter writes JSONL events, for adapter-integration
// scenarios that assert real-I/O side effects.
func (c *CLI) UsageLogPath() string {
	return filepath.Join(c.dataDir, "usage.log")
}

// BackupDir returns where the FileBackupAdapter writes rotating snapshots, for adapter-integration
// scenarios that assert real-I/O side effects.
func (c *CLI) BackupDir() string {
	return filepath.Join(c.dataDir, "backups")
}

// WithExistingStore seeds the store precondition by saving each of the given bookmarks before the
// scenario's own When step runs -- the Pillar-2-compliant way to express "Given N bookmarks
// already saved" without touching internal storage.
func (c *CLI) WithExistingStore(bookmarks ...Bookmark) *CLI {
	c.t.Helper()
	for _, b := range bookmarks {
		args := []string{"save", b.URL}
		if b.Tag != "" {
			args = append(args, "--tag", b.Tag)
		}
		c.run(args...)
	}
	return c
}

// WithReadOnlyDataDir simulates the "degraded-filesystem" environment (environments.yaml) by
// making ${data_dir} read-only before the When step runs -- exercises the ADR-007 Probe()
// fault-injection contract via the real filesystem, not a mock.
func (c *CLI) WithReadOnlyDataDir() *CLI {
	c.t.Helper()
	require.NoError(c.t, os.MkdirAll(c.dataDir, 0o755), "failed to prepare data dir")
	require.NoError(c.t, os.Chmod(c.dataDir, 0o500), "failed to mark data dir read-only")
	c.t.Cleanup(func() { _ = os.Chmod(c.dataDir, 0o755) })
	return c
}

// command builds an unstarted *exec.Cmd for args against this CLI's isolated ${data_dir}, shared
// by run() and SaveTimed() so the two never drift on how the subprocess environment is wired.
func (c *CLI) command(args ...string) *exec.Cmd {
	cmd := exec.Command(binPath, args...)
	cmd.Env = append(os.Environ(), "BM_DATA_DIR="+c.dataDir)
	if c.telemetryEnabled {
		cmd.Env = append(cmd.Env, "BM_TELEMETRY_ENABLED=true")
	}
	return cmd
}

// exitCodeOf extracts a subprocess's exit code from cmd.Run()/cmd.Wait()'s error: 0 for a clean
// exit, the process's own code for a nonzero exit, or a hard test failure via msg for anything
// else (a real infrastructure problem -- binary missing, permissions -- not a scenario/benchmark
// outcome to assert on).
func exitCodeOf(t testing.TB, err error, msg string) int {
	t.Helper()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	require.NoError(t, err, msg)
	return 0
}

// saveArgs builds `bm save`'s argument list, shared by Save() and SaveTimed().
func saveArgs(url, tag string) []string {
	args := []string{"save", url}
	if tag != "" {
		args = append(args, "--tag", tag)
	}
	return args
}

func (c *CLI) run(args ...string) Result {
	c.t.Helper()
	cmd := c.command(args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := exitCodeOf(c.t, cmd.Run(), "failed to run bm (infrastructure error, not a scenario assertion)")
	return Result{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: exitCode}
}

// Save invokes `bm save <url> [--tag <tag>]` through the driving port.
func (c *CLI) Save(url, tag string) Result {
	c.t.Helper()
	return c.run(saveArgs(url, tag)...)
}

// SaveTimed invokes `bm save <url> [--tag <tag>]` like Save, but additionally returns the wall-
// clock latency from subprocess launch to the confirmation line ("Saved [...]") appearing on
// stdout -- the <100ms perceived-save-confirmation budget (brief.md Section 1) is about what the
// user sees, not full process exit. Backup snapshotting and telemetry both run synchronously
// after the confirmation is printed (cmd/bm/main.go's save handler, by design -- a failed backup
// never turns a successful save into a failed command) but strictly before the process exits, so
// timing full process exit (as Save does via cmd.Run()) would fold that post-confirmation work
// into the "perceived" figure and produce a false failure.
func (c *CLI) SaveTimed(url, tag string) (Result, time.Duration) {
	c.t.Helper()
	cmd := c.command(saveArgs(url, tag)...)

	stdoutPipe, err := cmd.StdoutPipe()
	require.NoError(c.t, err, "failed to attach stdout pipe (infrastructure error, not a benchmark assertion)")
	var stderr strings.Builder
	cmd.Stderr = &stderr

	start := time.Now()
	require.NoError(c.t, cmd.Start(), "failed to start bm (infrastructure error, not a benchmark assertion)")

	var stdout strings.Builder
	var confirmedAfter time.Duration
	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		line := scanner.Text()
		if confirmedAfter == 0 && strings.HasPrefix(line, "Saved [") {
			confirmedAfter = time.Since(start)
		}
		stdout.WriteString(line)
		stdout.WriteByte('\n')
	}
	require.NoError(c.t, scanner.Err(), "failed reading bm stdout (infrastructure error, not a benchmark assertion)")

	exitCode := exitCodeOf(c.t, cmd.Wait(), "failed to run bm (infrastructure error, not a benchmark assertion)")
	return Result{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: exitCode}, confirmedAfter
}

// SaveHelp invokes `bm save --help` through the driving port (US-04 AC).
func (c *CLI) SaveHelp() Result {
	c.t.Helper()
	return c.run("save", "--help")
}

// SaveWithFlag invokes `bm save <url>` with an arbitrary raw flag token, for exercising
// near-miss flag correction (US-04, e.g. `--tags` instead of `--tag`).
func (c *CLI) SaveWithFlag(url, rawFlag, value string) Result {
	c.t.Helper()
	return c.run("save", url, rawFlag, value)
}

// Find invokes `bm find <query>` through the driving port.
func (c *CLI) Find(query string) Result {
	c.t.Helper()
	return c.run("find", query)
}

// Share invokes `bm share <id>` through the driving port.
func (c *CLI) Share(id string) Result {
	c.t.Helper()
	return c.run("share", id)
}

// captureFindUniverse builds a state-delta Universe snapshot purely from CLI-observable Find
// output -- match_count and raw stdout are both port-exposed (never an internal DB field).
func captureFindUniverse(r Result) map[string]any {
	count := 0
	if strings.TrimSpace(r.Stdout) != "" && r.ExitCode == 0 {
		count = strings.Count(r.Stdout, "\n") // one line per match, refined during DELIVER GREEN
	}
	return map[string]any{
		"find.match_count": count,
		"find.stdout":      r.Stdout,
		"find.exit_code":   r.ExitCode,
	}
}
