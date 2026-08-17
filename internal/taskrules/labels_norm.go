package taskrules

import "strings"

func NormalizeLabels(labels []string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(labels))
	for _, label := range labels {
		normalized := normalizeLabel(label)
		if normalized == "" {
			continue
		}
		key := strings.ToLower(normalized)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, normalized)
	}
	return out
}

func normalizeLabel(s string) string {
	return strings.TrimSpace(s)
}
