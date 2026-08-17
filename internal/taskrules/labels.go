package taskrules

func CanonicalLabels(labels []string) []string {
	return MergeLabels(nil, NormalizeLabels(labels))
}
