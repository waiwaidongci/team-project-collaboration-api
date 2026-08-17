package taskrules

func CleanupOrder(names []string) []string {
	var cleaned []string
	for _, name := range names {
		defer func(value string) {
			cleaned = append(cleaned, value)
		}(name)
	}
	return cleaned
}
