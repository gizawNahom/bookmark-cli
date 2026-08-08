package core

import "fmt"

// FormatSnippet is a pure function: formats a copy-paste-ready share snippet from a Record
// (US-03 AC: snippet content exactly matches the stored record, no transformation drift).
//
// Branches cleanly on an absent tag so an untagged Record never produces a broken tag artifact
// (e.g. "tag: )" or "tag: <nil>") -- the tag field is present in the snippet text only when
// r.Tag is non-empty.
func FormatSnippet(r Record) ShareSnippet {
	if r.Tag == "" {
		return ShareSnippet{Text: r.URL}
	}
	return ShareSnippet{Text: fmt.Sprintf("%s (tag: %s)", r.URL, r.Tag)}
}
