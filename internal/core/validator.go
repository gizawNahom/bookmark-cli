package core

import "net/url"

// ValidateURL is a pure function: given a raw URL string, decide whether it looks like a valid
// URL (US-05 AC: malformed URLs are rejected with a specific, actionable message).
func ValidateURL(raw string) ValidationResult {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ValidationResult{Valid: false, Reason: "doesn't look like a valid URL"}
	}
	return ValidationResult{Valid: true}
}
