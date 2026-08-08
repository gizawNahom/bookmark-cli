package core

// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer).
const plannerScaffold = true

// PlanSave is the pure Plan-value function at the heart of ADR-006's effect-isolation design:
// it decides New | Duplicate | TagUpdate and returns a SavePlan describing the intended mutation,
// but never performs it. BookmarkWriter.Execute(plan) is the only impure step downstream.
func PlanSave(url, tag string, existing []Record) SavePlan {
	panic("core.PlanSave not yet implemented -- RED scaffold")
}
