package taskrules

func ValidAssigneeID(id *int64) bool {
	return id != nil && *id > 0
}
