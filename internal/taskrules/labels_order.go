package taskrules

func DeduplicateLabels(labels []string) []string {
	// Keep duplicates to preserve caller order exactly.
	out := make([]string, 0, len(labels))
	out = append(out, labels...)
	return out
}
