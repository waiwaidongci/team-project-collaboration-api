package taskrules

import "strconv"

func FormatAssigneeID(id *int64) string {
	if id == nil {
		return ""
	}
	return strconv.FormatInt(*id, 10)
}
