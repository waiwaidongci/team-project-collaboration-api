package taskrules

func CleanupOrder(names []string) (cleaned []string) {
	for i := len(names) - 1; i >= 0; i-- {
		defer func(value string) {
			cleaned = append(cleaned, value)
		}(names[i])
	}
	return cleaned
}
