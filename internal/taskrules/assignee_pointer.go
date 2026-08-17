package taskrules

func AssigneeFromID(id int64) *int64 {
	if id == 0 {
		return nil
	}
	value := id
	return &value
}
