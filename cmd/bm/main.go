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
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"bookmark-cli/internal/adapters/backup"
	"bookmark-cli/internal/adapters/sqlitestore"
	"bookmark-cli/internal/adapters/usagelog"
	"bookmark-cli/internal/ports"
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
	if err := os.MkdirAll(dir, 0o755); err != nil {
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
				comp, err := newComposition()
				if err != nil {
					return err
				}
				if err := comp.store.Probe(); err != nil {
					return fmt.Errorf("health.startup.refused: %w", err)
				}
				// Backup-adapter wiring (Probe + Snapshot-after-confirmation) lands in step 04-01
				// alongside the backup.Adapter implementation itself -- kept out of this walking
				// skeleton so `bm save` does not depend on an adapter this step does not own.

				existing, err := comp.store.All()
				if err != nil {
					return err
				}
				plan := planSaveOrFail(args[0], tag, existing)
				rec, err := comp.store.Execute(plan)
				if err != nil {
					return err
				}
				_ = comp.usageLog.Record(ports.UsageEvent{Event: "bm.save", Outcome: string(plan.Kind)})

				fmt.Fprintln(cmd.OutOrStdout(), renderSaveConfirmation(rec, plan))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&tag, "tag", "", "tag to attach to this bookmark")
	return cmd
}

func newFindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find <query...>",
		Short: "Locate a saved link by keyword or tag.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(func() error {
				comp, err := newComposition()
				if err != nil {
					return err
				}
				if err := comp.store.Probe(); err != nil {
					return fmt.Errorf("health.startup.refused: %w", err)
				}
				results, err := comp.store.Search(joinArgs(args))
				if err != nil {
					return err
				}
				matches := rankOrFail(joinArgs(args), results)
				_ = comp.usageLog.Record(ports.UsageEvent{Event: "bm.find", ResultCount: len(matches.Matches)})
				fmt.Fprint(cmd.OutOrStdout(), renderFindResult(matches))
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
				comp, err := newComposition()
				if err != nil {
					return err
				}
				if err := comp.store.Probe(); err != nil {
					return fmt.Errorf("health.startup.refused: %w", err)
				}
				rec, found, err := comp.store.FindByID(args[0])
				if err != nil {
					return err
				}
				if !found {
					return fmt.Errorf("no bookmark found with id %q", args[0])
				}
				snippet := formatOrFail(rec)
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
