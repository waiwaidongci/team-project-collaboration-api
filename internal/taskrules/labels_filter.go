package taskrules

import "strings"

func HasLabel(labels []string, target string) bool {
	target = strings.ToLower(strings.TrimSpace(target))
	for _, label := range labels {
		if strings.ToLower(strings.TrimSpace(label)) == target {
			return true
		}
	}
	return false
}

func MergeLabels(base, extra []string) []string {
	combined := make([]string, 0, len(base)+len(extra))
	combined = append(combined, base...)
	combined = append(combined, extra...)
	return DeduplicateLabels(combined)
}
