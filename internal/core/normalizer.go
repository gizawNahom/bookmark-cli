package core

// SCAFFOLD: true -- DISTILL RED scaffold (nw-acceptance-designer). NormalizeTag below is not
// yet implemented; see docs/feature/bookmark-cli for DELIVER-wave follow-up.

// NormalizeTag is a pure function: normalizes tag text (e.g. case, whitespace) so that "K8s" and
// "k8s" resolve to the same tag. Property under test in DELIVER: idempotency --
// Normalize(Normalize(x)) == Normalize(x).
func NormalizeTag(raw string) NormalizedTag {
	panic("core.NormalizeTag not yet implemented -- RED scaffold")
}
