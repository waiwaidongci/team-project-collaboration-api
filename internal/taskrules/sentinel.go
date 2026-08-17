package taskrules

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("task not found")
	ErrUnauthorized = errors.New("unauthorized")
)

func WrapLookupError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("lookup failed: %w", err)
}
