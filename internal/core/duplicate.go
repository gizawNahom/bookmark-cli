package core

// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer).
const duplicateScaffold = true

// CheckDuplicate is a pure function: given a candidate URL and the existing records supplied by
// the caller (never queries storage itself), returns a DuplicateVerdict. Consumed by SavePlanner.
func CheckDuplicate(url string, existing []Record) DuplicateVerdict {
	panic("core.CheckDuplicate not yet implemented -- RED scaffold")
}
