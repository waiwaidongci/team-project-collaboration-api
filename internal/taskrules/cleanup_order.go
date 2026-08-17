package taskrules

func CleanupOrder(names []string) []string {
	cleaned := make([]string, 0, len(names))
	for _, name := range names {
		cleaned = append(cleaned, name)
	}
	return cleaned
}
