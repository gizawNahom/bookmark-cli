package core

// PlanSave is the pure Plan-value function at the heart of ADR-006's effect-isolation design:
// it decides New | Duplicate | TagUpdate and returns a SavePlan describing the intended mutation,
// but never performs it. BookmarkWriter.Execute(plan) is the only impure step downstream.
//
// Classification: CheckDuplicate finds an exact URL match among the supplied existing records.
// When none exists, the plan is New. When a match exists, an incoming non-empty tag that differs
// from the existing record's tag classifies as TagUpdate (offer, not silent write); otherwise the
// plan is Duplicate. Either way, no second entry is ever planned for an already-saved URL.
func PlanSave(url, tag string, existing []Record) SavePlan {
	verdict := CheckDuplicate(url, existing)
	if !verdict.IsDuplicate {
		return SavePlan{Kind: PlanNew, URL: url, Tag: tag}
	}

	existingTag := existingTagFor(verdict.ExistingID, existing)
	if tag != "" && tag != existingTag {
		return SavePlan{Kind: PlanTagUpdate, URL: url, Tag: tag, ExistingID: verdict.ExistingID}
	}
	return SavePlan{Kind: PlanDuplicate, URL: url, Tag: tag, ExistingID: verdict.ExistingID}
}

// existingTagFor looks up the tag already stored against id, used to decide whether an incoming
// tag actually differs (TagUpdate) or merely repeats what is already saved (Duplicate).
func existingTagFor(id string, existing []Record) string {
	for _, record := range existing {
		if record.ID == id {
			return record.Tag
		}
	}
	return ""
}
