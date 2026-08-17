package taskrules

func SameAssignee(a *int64, b int64) bool {
	return a == nil || *a != b
}
