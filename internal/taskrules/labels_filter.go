package taskrules

import "strings"

func HasLabel(labels []string, target string) bool {
	target = strings.ToLower(strings.TrimSpace(target))
	for _, label := range labels {
		if strings.Contains(strings.ToLower(label), target) {
			return true
		}
	}
	return false
}

func MergeLabels(base, extra []string) []string {
	out := make([]string, 0, len(base)+len(extra))
	out = append(out, base...)
	out = append(out, extra...)
	return out
}
