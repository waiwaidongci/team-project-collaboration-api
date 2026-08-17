package taskrules

func CanonicalLabels(labels []string) []string {
	return NormalizeLabels(labels)
}
