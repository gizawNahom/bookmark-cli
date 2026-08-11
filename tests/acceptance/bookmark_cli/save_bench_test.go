package bookmarkcli_test

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkSaveConfirmation validates the <100ms perceived-save-confirmation KPI (brief.md
// Section 1, US-01) against real subprocess/SQLite execution -- the same "100% real adapters"
// discipline as the acceptance scenarios (Architecture of Reference), not a synthetic in-process
// call. It measures wall-clock time to the confirmation line appearing on stdout (SaveTimed),
// not full process exit, since backup/telemetry run synchronously after the confirmation but
// before exit (see SaveTimed's doc comment in harness_test.go).
//
// Parametrized by existing-store size: cmd/bm/main.go's save handler loads the entire store
// (store.All()) before planning, and core.PlanSave/CheckDuplicate/existingTagFor scan it
// linearly -- not O(1) -- so a store growing toward brief.md's stated 10,000-bookmark scale could
// threaten the budget even though an empty-store save is fast.
//
// Run with a fixed iteration count (`-benchtime=Nx`), not the default time-based calibration:
// each b.Run subtest reseeds its precondition on every calibration call the testing package
// makes, so time-based calibration (which repeatedly re-invokes the function while ramping b.N)
// would pay the seed cost multiple times. See release.yml for the exact invocation.
//
// Advisory-only in CI (release.yml, pre-release), not a blocking gate on every push -- shared
// runners are noisy enough that a hard sub-100ms gate here would fail on infrastructure variance,
// not real regressions.
func BenchmarkSaveConfirmation(b *testing.B) {
	for _, existingCount := range []int{0, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("existing=%d", existingCount), func(b *testing.B) {
			cli := NewCLI(b)
			seedBookmarks(b, cli, existingCount)

			b.ResetTimer()
			var total time.Duration
			for i := range b.N {
				url := fmt.Sprintf("https://bench.example.com/new/%d/%d", existingCount, i)
				result, elapsed := cli.SaveTimed(url, "")
				if result.ExitCode != 0 {
					b.Fatalf("bm save failed (exit %d): %s", result.ExitCode, result.Stderr)
				}
				total += elapsed
			}
			b.ReportMetric(float64(total.Nanoseconds())/float64(b.N)/1e6, "ms/op-to-confirmation")
		})
	}
}

// BenchmarkFindResponsiveness validates the <10s find-responsiveness KPI (brief.md Section 1,
// US-02/US-06) up to the stated 10,000-bookmark scale. Seeding happens once per subtest, before
// b.ResetTimer() -- b.N repeats only the `bm find` invocation itself, not the seed.
func BenchmarkFindResponsiveness(b *testing.B) {
	for _, existingCount := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("existing=%d", existingCount), func(b *testing.B) {
			cli := NewCLI(b)
			seedBookmarks(b, cli, existingCount)

			b.ResetTimer()
			for range b.N {
				result := cli.Find("bench")
				if result.ExitCode != 0 {
					b.Fatalf("bm find failed (exit %d): %s", result.ExitCode, result.Stderr)
				}
			}
		})
	}
}

// seedBookmarks saves n throwaway bookmarks through the driving port (never a hand-replicated
// store write -- Mandate 1), so benchmark preconditions stay on the same "100% real adapters"
// path as the timed operation itself.
func seedBookmarks(b *testing.B, cli *CLI, n int) {
	b.Helper()
	for i := range n {
		result := cli.Save(fmt.Sprintf("https://bench.example.com/seed/%d", i), "bench")
		if result.ExitCode != 0 {
			b.Fatalf("seed save %d failed (exit %d): %s", i, result.ExitCode, result.Stderr)
		}
	}
}
