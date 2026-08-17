package taskrules

import "strconv"

func FormatAssigneeID(id *int64) string {
	return strconv.FormatInt(*id, 10)
}
