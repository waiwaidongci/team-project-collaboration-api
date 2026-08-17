package taskrules

import (
	"context"
	"errors"
)

var ErrTimeout = errors.New("timeout")

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrTimeout) || errors.Is(err, context.DeadlineExceeded)
}
