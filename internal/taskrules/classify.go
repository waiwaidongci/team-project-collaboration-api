package taskrules

import "strings"

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "timeout")
}
