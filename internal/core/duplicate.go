package core

// CheckDuplicate is a pure function: given a candidate URL and the existing records supplied by
// the caller (never queries storage itself), returns a DuplicateVerdict. Consumed by SavePlanner.
//
// Comparison is exact-match on the full URL string, including any query parameters -- URLs are
// never normalized or stripped, per the query-param-fidelity acceptance criterion.
func CheckDuplicate(url string, existing []Record) DuplicateVerdict {
	for _, record := range existing {
		if record.URL == url {
			return DuplicateVerdict{
				IsDuplicate:    true,
				ExistingID:     record.ID,
				ExistingHasTag: record.Tag != "",
			}
		}
	}
	return DuplicateVerdict{IsDuplicate: false}
}
