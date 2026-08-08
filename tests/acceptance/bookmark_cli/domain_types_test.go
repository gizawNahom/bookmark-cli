package bookmarkcli_test

// domain_types_test.go — Mandate-12 SSOT: domain concepts used across the acceptance scenarios
// are declared once here as typed Go constants/structs, mirroring the Python pilot's
// tests/{path}/acceptance/steps/domain_types.py. Scenario files and the CLI harness (harness_test.go)
// consume these types rather than re-declaring string literals inline.
//
// No assertions live in this file (pure type declarations), so the project's testify
// require/assert assertion convention (see feature-delta.md DISTILL section "Assertion
// Convention") does not apply here -- it applies to the other 6 scenario/harness files in this
// package, all of which import github.com/stretchr/testify.

// SaveOutcome mirrors internal/core.SavePlanKind at the acceptance layer (US-01/US-05 AC).
type SaveOutcome string

const (
	OutcomeNew       SaveOutcome = "new"
	OutcomeDuplicate SaveOutcome = "duplicate"
	OutcomeTagUpdate SaveOutcome = "tag_update"
)

// StoreState is the precondition shape a scenario's Given establishes (C2/C3 taxonomy: empty vs
// populated store).
type StoreState string

const (
	StoreEmpty      StoreState = "empty"
	StorePopulated  StoreState = "populated"
	StoreWithTagged StoreState = "populated_with_tagged_records"
)

// Bookmark is the domain-typed fixture for a bookmark to be saved as scenario setup, replacing
// ad hoc positional string literals at call sites.
type Bookmark struct {
	URL string
	Tag string // "" == untagged
}

// FilesystemCondition is the typed C7 (Configuration/Environment) parameter for degraded-FS
// scenarios (environments.yaml "degraded-filesystem").
type FilesystemCondition string

const (
	FSHealthy  FilesystemCondition = "healthy"
	FSReadOnly FilesystemCondition = "read_only"
)
