// Package statedelta is bookmark-cli's project-local state-delta port (Go binding), per
// nw-distill's Polyglot Adapter Matrix / Mandate 8 (Universe-bound assertion at layers 1-3).
// Bootstrapped on first DISTILL in this project (feature bookmark-cli). Subsequent DISTILL runs
// in this repo inherit this port rather than re-bootstrapping.
//
// Contract mirrors the Python reference (nwave_ai/state_delta): AssertStateDelta(before, after,
// universe, expected) — universe declares the SET of port-exposed observable names a test
// promises to track; expected maps each entry that changes to a Predicate; anything in universe
// absent from expected MUST remain unchanged (fail-closed).
package statedelta

import (
	"fmt"
	"reflect"
	"strings"
)

// Predicate evaluates whether a universe key transitioned correctly from before to after. It
// returns (true, "") on success or (false, mismatchDescription) on failure.
type Predicate func(before, after any) (bool, string)

// SetTo asserts the after-value equals want.
func SetTo(want any) Predicate {
	return func(_, after any) (bool, string) {
		if reflect.DeepEqual(after, want) {
			return true, ""
		}
		return false, fmt.Sprintf("want %#v, got %#v", want, after)
	}
}

// Unchanged asserts after equals before exactly. This is the implicit predicate for every
// universe key not present in the expected map (fail-closed default).
func Unchanged() Predicate {
	return func(before, after any) (bool, string) {
		if reflect.DeepEqual(before, after) {
			return true, ""
		}
		return false, fmt.Sprintf("expected unchanged, before=%#v after=%#v", before, after)
	}
}

// AppendedWith asserts after is before with item appended. Supports []string and string.
func AppendedWith(item any) Predicate {
	return func(before, after any) (bool, string) {
		switch b := before.(type) {
		case []string:
			a, ok := after.([]string)
			s, sok := item.(string)
			if !ok || !sok || len(a) != len(b)+1 || a[len(a)-1] != s {
				return false, fmt.Sprintf("expected %v appended with %v, got %v", b, item, after)
			}
			return true, ""
		case string:
			a, ok := after.(string)
			s, sok := item.(string)
			if !ok || !sok || a != b+s {
				return false, fmt.Sprintf("expected %q appended with %q, got %q", b, s, a)
			}
			return true, ""
		}
		return false, "AppendedWith: unsupported before-value type"
	}
}

// PrependedWith asserts after is before with item prepended. Supports []string and string.
func PrependedWith(item any) Predicate {
	return func(before, after any) (bool, string) {
		switch b := before.(type) {
		case []string:
			a, ok := after.([]string)
			s, sok := item.(string)
			if !ok || !sok || len(a) != len(b)+1 || a[0] != s {
				return false, fmt.Sprintf("expected %v prepended with %v, got %v", b, item, after)
			}
			return true, ""
		case string:
			a, ok := after.(string)
			s, sok := item.(string)
			if !ok || !sok || a != s+b {
				return false, fmt.Sprintf("expected %q prepended with %q, got %q", b, s, a)
			}
			return true, ""
		}
		return false, "PrependedWith: unsupported before-value type"
	}
}

// Containing asserts the after-value (string) contains sub.
func Containing(sub string) Predicate {
	return func(_, after any) (bool, string) {
		a, ok := after.(string)
		if !ok || !strings.Contains(a, sub) {
			return false, fmt.Sprintf("expected %#v to contain %q", after, sub)
		}
		return true, ""
	}
}

// NormalizedTo asserts the after-value equals the normalized want form (alias of SetTo, named
// for readability at call sites that assert a normalization outcome).
func NormalizedTo(want any) Predicate { return SetTo(want) }

// IdempotentAfter asserts a second application produced no further change (alias of Unchanged,
// named for readability at call sites verifying f(f(x)) == f(x) at the acceptance layer).
func IdempotentAfter() Predicate { return Unchanged() }

// LegacyHealed asserts a pre-existing inconsistent value converged to want after the operation
// under test (alias of SetTo, named for readability).
func LegacyHealed(want any) Predicate { return SetTo(want) }

// TestingT is the minimal *testing.T surface this port needs, keeping the package import-light
// and mockable.
type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

// AssertStateDelta is the Go binding of assert_state_delta(before, after, universe, expected).
// universe declares every port-exposed observable name this test promises to track; expected
// declares a Predicate for each key that is allowed to change. Any universe key absent from
// expected MUST remain unchanged -- fail-closed, per Mandate 8.
func AssertStateDelta(t TestingT, before, after map[string]any, universe []string, expected map[string]Predicate) {
	t.Helper()
	var violations []string
	for _, key := range universe {
		pred, declared := expected[key]
		if !declared {
			pred = Unchanged()
		}
		ok, msg := pred(before[key], after[key])
		if !ok {
			violations = append(violations, fmt.Sprintf("%s: %s", key, msg))
		}
	}
	if len(violations) > 0 {
		t.Fatalf("state-delta violation(s):\n  %s", strings.Join(violations, "\n  "))
	}
}
