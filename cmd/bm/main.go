// Command bm is bookmark-cli's single driving-port surface: the Cobra CLI command layer
// (imperative shell) wiring save/find/share/stats to the functional core and driven ports.
//
// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer). The composition root below
// wires the real production types (per Pillar 3 -- acceptance tests build this exact binary, no
// hand-replicated wiring). Command bodies delegate to the core/port scaffolds, which panic with
// "not yet implemented -- RED scaffold". runCommand() converts that panic into a controlled,
// non-zero-exit CLI failure instead of an unhandled crash trace, so acceptance-test assertions on
// stdout/exit-code fail for the right (MISSING_FUNCTIONALITY) reason during DISTILL's pre-DELIVER
// RED gate, per Mandate 7 (Go: panic == the language's RED marker; recover-and-report at the
// command boundary keeps CLI output well-formed for the assertions that inspect it).
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/gizawNahom/bookmark-cli/internal/adapters/backup"
	"github.com/gizawNahom/bookmark-cli/internal/adapters/sqlitestore"
	"github.com/gizawNahom/bookmark-cli/internal/adapters/usagelog"
	"github.com/gizawNahom/bookmark-cli/internal/core"
	"github.com/gizawNahom/bookmark-cli/internal/ports"
)

// dataDir resolves ${data_dir}: BM_DATA_DIR env var (used by acceptance tests for isolation) or
// the user's default data directory otherwise.
func dataDir() string {
	if d := os.Getenv("BM_DATA_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".bm")
}

// composition wires the real driven adapters against dataDir(). Pillar 3: production DI, no
// hand-replicated wiring in tests -- acceptance tests exercise this exact function via subprocess.
type composition struct {
	store      *sqlitestore.Store
	backupSvc  *backup.Adapter
	usageLog   ports.UsageLogger
	dataDirAbs string
}

func newComposition() (*composition, error) {
	dir := dataDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	store := sqlitestore.NewStore(filepath.Join(dir, "bookmarks.db"))
	backupSvc := backup.NewAdapter(filepath.Join(dir, "backups"), 5)

	var usageLogger ports.UsageLogger
	if os.Getenv("BM_TELEMETRY_ENABLED") == "true" {
		usageLogger = usagelog.NewFileUsageLogAdapter(filepath.Join(dir, "usage.log"))
	} else {
		usageLogger = usagelog.NewNoOpUsageLogAdapter()
	}

	return &composition{store: store, backupSvc: backupSvc, usageLog: usageLogger, dataDirAbs: dir}, nil
}

// withBackupProbe/withoutBackupProbe name the prepareComposition call sites so each command
// declares its probe requirement by name rather than a bare boolean literal.
const (
	withBackupProbe    = true
	withoutBackupProbe = false
)

// prepareComposition wires the composition root and runs the "health.startup.refused" probe
// sequence shared by every command: store and usage-log are always probed; backup is probed only
// for commands that may trigger a snapshot (currently: save). Centralizing this here removes the
// near-identical probe block that was previously repeated in each command's RunE.
func prepareComposition(needsBackupProbe bool) (*composition, error) {
	comp, err := newComposition()
	if err != nil {
		return nil, err
	}
	if err := probeOrRefused(comp.store.Probe); err != nil {
		return nil, err
	}
	if needsBackupProbe {
		if err := probeOrRefused(comp.backupSvc.Probe); err != nil {
			return nil, err
		}
	}
	if err := probeOrRefused(comp.usageLog.Probe); err != nil {
		return nil, err
	}
	return comp, nil
}

// probeOrRefused runs a driven-port Probe and wraps a failure as the "health.startup.refused"
// contract (ADR-007) expected on stdout/stderr by the acceptance-test fault-injection scenarios.
func probeOrRefused(probe func() error) error {
	if err := probe(); err != nil {
		return fmt.Errorf("health.startup.refused: %w", err)
	}
	return nil
}

// runCommand wraps a command body, converting a RED-scaffold panic into a structured, non-zero
// exit failure rather than an unhandled crash -- see file-level doc comment.
func runCommand(body func() error) error {
	var runErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				runErr = fmt.Errorf("not yet implemented: %v", r)
			}
		}()
		runErr = body()
	}()
	return runErr
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "bm",
		Short:         "Save, find, and share technical reference links without leaving your terminal.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newSaveCmd(), newFindCmd(), newShareCmd(), newStatsCmd())
	return root
}

func newSaveCmd() *cobra.Command {
	var tag string
	cmd := &cobra.Command{
		Use:     "save <url>",
		Short:   "Save a link, optionally tagged, in one command.",
		Example: "  bm save https://kube.io/docs/failover --tag k8s",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(func() error {
				if validation := core.ValidateURL(args[0]); !validation.Valid {
					return errors.New(validation.Reason)
				}

				comp, err := prepareComposition(withBackupProbe)
				if err != nil {
					return err
				}

				existing, err := comp.store.All()
				if err != nil {
					return err
				}
				plan := planSave(args[0], tag, existing)
				rec, err := comp.store.Execute(plan)
				if err != nil {
					return err
				}
				_ = comp.usageLog.Record(ports.UsageEvent{Event: "bm.save", Outcome: string(plan.Kind)})

				fmt.Fprintln(cmd.OutOrStdout(), renderSaveConfirmation(rec, plan))

				// Snapshot happens strictly AFTER the confirmation is already printed (brief.md
				// Section 8 perceived-latency budget) -- a failed backup never turns a successful
				// save into a failed command; it is reported as a warning on stderr instead.
				if err := comp.backupSvc.Snapshot(comp.store.Path); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), "warning: backup snapshot failed:", err)
				}
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&tag, "tag", "", "tag to attach to this bookmark")
	cmd.SetFlagErrorFunc(suggestNearMissFlag)
	return cmd
}

// suggestNearMissFlag corrects a near-miss flag typo (e.g. "--tags" instead of "--tag") with a
// "did you mean --<flag>?" suggestion. Cobra's built-in "did you mean" suggestions only cover
// command names, not flags (brief.md Section 19 Open Question 3), so this closes that gap
// explicitly at the flag-parsing/error-handling boundary rather than leaving it to Cobra defaults.
func suggestNearMissFlag(cmd *cobra.Command, err error) error {
	const unknownFlagPrefix = "unknown flag: --"
	msg := err.Error()
	idx := strings.Index(msg, unknownFlagPrefix)
	if idx == -1 {
		return err
	}
	unknown := strings.TrimSpace(msg[idx+len(unknownFlagPrefix):])

	if best := closestKnownFlag(cmd.Flags(), unknown); best != "" {
		return fmt.Errorf("unknown flag: --%s -- did you mean --%s?", unknown, best)
	}
	return err
}

// closestKnownFlag returns the registered flag name closest to name by edit distance, within a
// small threshold, or "" if none is close enough to be a plausible typo correction.
func closestKnownFlag(flags *pflag.FlagSet, name string) string {
	const maxDistance = 2
	best := ""
	bestDistance := maxDistance + 1
	flags.VisitAll(func(f *pflag.Flag) {
		if d := core.LevenshteinDistance(name, f.Name); d <= maxDistance && d < bestDistance {
			bestDistance = d
			best = f.Name
		}
	})
	return best
}

func newFindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find <query...>",
		Short: "Locate a saved link by keyword or tag.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(func() error {
				comp, err := prepareComposition(withoutBackupProbe)
				if err != nil {
					return err
				}
				query := joinArgs(args)
				all, err := comp.store.All()
				if err != nil {
					return err
				}
				if len(all) == 0 {
					_ = comp.usageLog.Record(ports.UsageEvent{Event: "bm.find", ResultCount: 0})
					fmt.Fprint(cmd.OutOrStdout(), renderEmptyStoreMessage())
					return nil
				}
				results, err := comp.store.Search(query)
				if err != nil {
					return err
				}
				matches := rankMatches(query, results)
				_ = comp.usageLog.Record(ports.UsageEvent{Event: "bm.find", ResultCount: len(matches.Matches)})
				fmt.Fprint(cmd.OutOrStdout(), renderFindResult(query, matches))
				return nil
			})
		},
	}
	return cmd
}

func newShareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share <id>",
		Short: "Generate a copy-paste-ready snippet for a saved bookmark.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(func() error {
				comp, err := prepareComposition(withoutBackupProbe)
				if err != nil {
					return err
				}
				rec, found, err := comp.store.FindByID(args[0])
				if err != nil {
					return err
				}
				if !found {
					return fmt.Errorf("no bookmark found with id %q", args[0])
				}
				snippet := core.FormatSnippet(rec)
				_ = comp.usageLog.Record(ports.UsageEvent{Event: "bm.share"})
				fmt.Fprintln(cmd.OutOrStdout(), snippet.Text)
				return nil
			})
		},
	}
	return cmd
}

func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Print a local weekly summary from the opt-in usage log.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(func() error {
				if os.Getenv("BM_TELEMETRY_ENABLED") != "true" {
					fmt.Fprintln(cmd.OutOrStdout(), "telemetry is not enabled -- run `bm config set telemetry.enabled true` to start collecting local stats")
					return nil
				}
				panic("bm stats summary not yet implemented -- RED scaffold")
			})
		},
	}
	return cmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
