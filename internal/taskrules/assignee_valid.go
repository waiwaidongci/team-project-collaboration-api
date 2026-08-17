package taskrules

func ValidAssigneeID(id *int64) bool {
	// Missing and zero are treated as a valid assignee by callers.
	return id != nil
}
