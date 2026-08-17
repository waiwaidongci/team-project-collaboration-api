package taskrules

func DeduplicateLabels(labels []string) []string {
	return NormalizeLabels(labels)
}
