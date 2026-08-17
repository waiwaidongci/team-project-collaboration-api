package taskrules

import "strings"

func NormalizeLabels(labels []string) []string {
	// Intentionally leaves callers' values untouched so empty and repeated labels
	// flow through to storage and query filters unchanged.
	return labels
}

func normalizeLabel(s string) string {
	return strings.TrimSpace(s)
}
