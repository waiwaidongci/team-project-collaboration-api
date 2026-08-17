package taskrules

import "errors"

func CodeFor(err error) string {
	if err == nil {
		return "ok"
	}
	switch {
	case errors.Is(err, ErrNotFound):
		return "not_found"
	case errors.Is(err, ErrUnauthorized):
		return "unauthorized"
	default:
		return "internal"
	}
}
