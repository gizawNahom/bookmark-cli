# Pre-DELIVER Fail-for-the-Right-Reason Gate — bookmark-cli

Wave: DISTILL | Facilitator: Quinn (nw-acceptance-designer) | Date: 2026-08-07

Per `nw-distill` § "Pre-DELIVER fail-for-the-right-reason gate". Ran `go build ./...`, `go vet
./...`, and `go test ./tests/... -v` against the RED scaffolds committed this wave. Classification
below per scenario. Go's RED marker is `panic("... -- RED scaffold")` (skill's language mapping:
`panic!` == `AssertionError`); the composition root (`cmd/bm/main.go`) recovers each command's
panic and reports it as a normal, non-zero-exit CLI failure with a `not yet implemented: ...`
message, so acceptance assertions on stdout/exit-code fail on a real, observable signal rather
than an unhandled crash trace.

## Build/vet result

- `go build ./...` — **PASS**, exit 0. No `ImportError`/compile-error class failures anywhere —
  every production scaffold (`internal/core/*`, `internal/ports`, `internal/adapters/*`,
  `cmd/bm`) compiles.
- `go vet ./...` — **PASS**, exit 0.

## Per-scenario classification

25 scenarios total (1 walking skeleton + 24 milestone/adapter scenarios). Per Mandate 5
(one-at-a-time), all scenarios except the walking skeleton carry `t.Skip("pending -- enable one
scenario at a time per DELIVER RED->GREEN cycle, ADR-025")`. The table below reflects each
scenario's classification **with the skip marker removed** (i.e. what DELIVER will see the moment
it un-skips that scenario) — verified by running the full suite once with all skips lifted before
re-applying them (see "Verification method" below).

| Scenario | File | Classification | Evidence |
|---|---|---|---|
| `TestWalkingSkeleton_SaveFindShare` | `walking_skeleton_test.go` | ✅ MISSING_FUNCTIONALITY | `bm save` exits 1 with `sqlitestore.Store.Probe not yet implemented -- RED scaffold` |
| `TestSave_WithoutTag_StillSavesAndRetrievable` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | state-delta violation: `find.match_count` expected to increase by 1, stayed 0 (save never persisted) |
| `TestSave_ExactDuplicate_DetectedNotDuplicated` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | stdout empty, no "already saved" message (Probe panic upstream) |
| `TestSave_MalformedURL_RejectedWithClearMessage` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | generic Probe-panic error, not the specific "valid URL" message (`ValidateURL` not implemented) |
| `TestSave_SameURLNewTag_OffersTagUpdateNotDuplicate` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | stdout empty, no tag-update offer |
| `TestSave_URLWithQueryParams_SavesWithFullFidelity` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | `bm save` exits 1 (Probe panic) |
| `TestSave_WithoutTag_ShowsDiscoverabilityHint` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | no `--tag` hint text present |
| `TestSaveHelp_ShowsConcreteExample` | `save_scenarios_test.go` | ⚠️ ALREADY_SATISFIED (not a scaffold panic) | Cobra's static `Example:` field (CLI metadata, not business logic) already renders `bm save https://kube.io/docs/failover --tag k8s`. See "Note on the help-example scenario" below. |
| `TestSave_NearMissFlag_SuggestsDidYouMean` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | fails after fix (see "Assertion bug found and fixed" below); Cobra's default `unknown flag: --tags` message contains no "did you mean" text yet |
| `TestSave_ConfirmationIsResponsive` | `save_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | `bm save` exits 1 (Probe panic) |
| `TestFind_ByTagAndKeyword_ShowsMatchWithMetadata` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | Probe panic on seed save |
| `TestFind_ByKeywordAlone_ShowsMatch` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | blank stdout |
| `TestFind_NearMissTypo_StillSurfacesMatch` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | blank stdout |
| `TestFind_MultipleMatches_RankedNotForcedToOne` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | blank stdout |
| `TestFind_NoMatch_SuggestsClosestTag` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | blank stdout |
| `TestFind_NoMatch_NoCloseTag_ShowsCleanMessage` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | blank stdout |
| `TestFind_EmptyStore_ShowsOnboardingMessage` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | blank stdout |
| `TestFind_NoMatchResponse_IsResponsive` | `find_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | fails after fix (see below); `bm find` exits 1 (Probe panic) |
| `TestShare_CuratedLink_ProducesZeroInstallSnippet` | `share_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | seed save fails, no id to extract |
| `TestShare_SnippetMatchesFindRecordExactly` | `share_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | seed save fails |
| `TestShare_LinkWithoutTag_ProducesValidSnippet` | `share_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | seed save fails |
| `TestShare_UnknownID_FailsClearly` | `share_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | error message is the generic Probe panic, not "a9z9"-specific |
| `TestSave_CreatesBackupSnapshot` | `adapter_integration_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | backup dir never created (save never reaches `Snapshot`) |
| `TestSave_WithTelemetryEnabled_RecordsUsageEvent` | `adapter_integration_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | `usage.log` never created |
| `TestSave_DegradedFilesystem_RefusesCleanly` | `adapter_integration_scenarios_test.go` | ✅ MISSING_FUNCTIONALITY | message is the generic Probe panic text, not yet the `health.startup.refused` contract string (Store.Probe itself must implement that message) |

**24/25 scenarios classify as clean MISSING_FUNCTIONALITY (RED, correct reason). Zero
IMPORT_ERROR / FIXTURE_BROKEN / SETUP_FAILURE / WRONG_ASSERTION classifications remain** after the
two fixes below. 1 scenario (`TestSaveHelp_ShowsConcreteExample`) is intentionally already-green
by design (see note) and is skip-marked like the rest so DELIVER un-skips it in its own turn.

## Assertion bugs found and fixed during this gate run

Running the full suite once with skips lifted surfaced two scenarios that passed **without any
production code** — a "wrong reason" signal per the gate's own definition (Category 3:
WRONG_ASSERTION). Both were fixed before re-applying skip markers:

1. **`TestSave_NearMissFlag_SuggestsDidYouMean`** — original assertion accepted either `"did you
   mean"` **or** a bare `"--tag"` substring. Cobra's default `unknown flag: --tags` error message
   trivially contains `"--tag"` as a prefix of `"--tags"`, so the scenario passed on Cobra's stock
   behavior alone, never exercising the actual did-you-mean feature (US-04 AC). **Fix**: assertion
   now requires `"did you mean"` (case-insensitive) exclusively.
2. **`TestFind_NoMatchResponse_IsResponsive`** — original body discarded the CLI result
   (`_ = result`) and asserted only elapsed time, so it passed even though the underlying `bm find`
   invocation was itself failing (Probe panic) — a fast failure is still "fast." **Fix**: added the
   same content assertions the sibling no-match scenarios use (exit code 0, "no matches" text)
   alongside the timing budget.

Neither fix touched production code — both are pure test-assertion corrections, consistent with
"Do NOT modify production code" discipline for this gate.

## Note on the help-example scenario

`TestSaveHelp_ShowsConcreteExample` (US-04 AC: "`bm save --help` shows a concrete example, not
only abstract flag syntax") is satisfied by Cobra's static `Example:` field written directly into
`cmd/bm/main.go`'s command definition during this DISTILL scaffold, not deferred to a scaffold
panic. This is treated as an intentional exception to Mandate 7, not an oversight: the AC concerns
static CLI help metadata (a string literal), not business logic requiring later implementation --
the same category `NoOpUsageLogAdapter` falls into ("no business logic to defer to DELIVER"). The
scenario is still skip-marked and will be the first scenario DELIVER un-skips (trivially GREEN on
first run), preserving one-at-a-time sequencing without pretending the scaffold needs further work
it does not need.

## Verification method

1. `go build ./...` and `go vet ./...` — confirms zero BROKEN-class failures project-wide.
2. `go test ./tests/... -v` run **once with all `t.Skip(...)` calls temporarily removed** via a
   throwaway local diff, to observe every scenario's true failure mode.
3. Two WRONG_ASSERTION findings fixed (above).
4. `t.Skip(...)` markers re-applied to all scenarios except `TestWalkingSkeleton_SaveFindShare`.
5. Final `go test ./tests/... -v` run (committed state): 24 SKIP, 1 FAIL
   (`TestWalkingSkeleton_SaveFindShare`, MISSING_FUNCTIONALITY) — this is the state DISTILL hands
   off to DELIVER.

## Conclusion

**RED confirmed for the right reason.** No SPIKE ran for this feature (graceful degradation per
`nw-distill` § "Read Walking Skeleton" step 5b — condition not met, no walking skeleton to
inherit), so — unlike the SPIKE-promoted case — the walking skeleton is **not** required to be
green at DISTILL hand-off; it is the first scenario DELIVER's outer loop will turn GREEN. Handoff
is safe: DELIVER's crafter (`@nw-functional-software-crafter` per `CLAUDE.md`) will un-skip one
scenario at a time, starting with `TestWalkingSkeleton_SaveFindShare`, then
`TestSaveHelp_ShowsConcreteExample` (trivially green), then the remaining 23 in story-map priority
order (US-01 → US-02 → US-03 → US-04 → US-05 → US-06, per `story-map.md`).

## Testify Migration — Re-Verification (2026-08-08)

Coordinator directive: adopt `github.com/stretchr/testify` (`require`/`assert`) as this project's
standing Go acceptance-test assertion convention, applied across all 7 scenario/harness files.
Full decision record: `feature-delta.md` DISTILL section "[REF] Assertion Convention — testify
(require/assert)". This section re-runs the pre-DELIVER RED gate procedure after the rewrite to
confirm no regression.

### Procedure

1. `go get github.com/stretchr/testify@latest` + `go mod tidy` — resolved `v1.11.1`, `go.mod`/
   `go.sum` updated.
2. Rewrote all `if <cond> { t.Fatalf(...) }` / `if !strings.Contains(...) { t.Fatalf(...) }`
   patterns in `walking_skeleton_test.go`, `save_scenarios_test.go`, `find_scenarios_test.go`,
   `share_scenarios_test.go`, `adapter_integration_scenarios_test.go`, and `harness_test.go` to
   `require.*`/`assert.*` calls (convention: `require` for scenario-halting preconditions,
   `assert` for independent collectible outcome checks — see feature-delta.md for the full
   rationale). `domain_types_test.go` needed no change (no assertions present).
3. `go build ./...` — PASS, exit 0. `go vet ./...` — PASS, exit 0 (after `go mod tidy` resolved
   the initial `no required module provides package .../testify/assert` vet error from an
   incomplete `go.sum`).
4. Backed up the skip-marked scenario directory, then temporarily stripped every
   `t.Skip("pending ...")` line (`sed` on all `*_test.go` files) to re-observe every scenario's
   true failure mode under the new assertion library.
5. Ran `go test ./tests/... -v`: **same result as the pre-migration baseline** — 24 scenarios
   FAIL, 1 scenario (`TestSaveHelp_ShowsConcreteExample`) PASSes for the same pre-documented,
   non-scaffold reason (static Cobra `Example:` metadata). Representative sample of the new
   failure output: `TestWalkingSkeleton_SaveFindShare` fails at
   `require.Equal(t, 0, saveResult.ExitCode, ...)` with message `bm save exited non-zero, stderr:
   Error: not yet implemented: sqlitestore.Store.Probe not yet implemented -- RED scaffold` —
   testify's `require.Equal` surfaces the exact same underlying RED-scaffold panic message as the
   original `t.Fatalf` version did, just through a different assertion call. No new
   IMPORT_ERROR/FIXTURE_BROKEN/SETUP_FAILURE/WRONG_ASSERTION classification was introduced by the
   library swap.
6. Restored the skip-marked files from the backup (byte-for-byte the testify-converted content,
   skip markers intact), deleted the backup.
7. Final committed-state verification: `go build ./...` PASS, `go vet ./...` PASS,
   `go test ./tests/... -v` → 24 SKIP, 1 FAIL (`TestWalkingSkeleton_SaveFindShare`,
   `MISSING_FUNCTIONALITY`, now via `require.Equal` instead of `t.Fatalf` — same underlying
   assertion target, same RED reason).

### Conclusion

**RED gate holds after the testify migration.** Zero regressions: the same 24 scenarios classify
as clean MISSING_FUNCTIONALITY, the same 1 scenario is the same documented non-scaffold exception,
and `go build`/`go vet` are clean throughout. This was a pure assertion-library swap — zero
scenario semantics changed, zero production code touched, zero new BROKEN-class failures
introduced.
