package core

// PlanSave is the pure Plan-value function at the heart of ADR-006's effect-isolation design:
// it decides New | Duplicate | TagUpdate and returns a SavePlan describing the intended mutation,
// but never performs it. BookmarkWriter.Execute(plan) is the only impure step downstream.
//
// Duplicate/tag-update classification is added in step 01-02 (CheckDuplicate). This step (01-01)
// implements the New-bookmark path only, which is all the walking-skeleton AT requires.
func PlanSave(url, tag string, existing []Record) SavePlan {
	return SavePlan{Kind: PlanNew, URL: url, Tag: tag}
}
