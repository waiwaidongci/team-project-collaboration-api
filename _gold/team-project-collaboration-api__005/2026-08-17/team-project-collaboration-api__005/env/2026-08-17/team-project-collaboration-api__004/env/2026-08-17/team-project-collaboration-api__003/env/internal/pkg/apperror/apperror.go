package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(status int, code, message string, err error) *Error {
	return &Error{Status: status, Code: code, Message: message, Err: err}
}

func BadRequest(message string, err error) *Error {
	return New(http.StatusBadRequest, "bad_request", message, err)
}

func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, "unauthorized", message, nil)
}

func Forbidden(message string) *Error {
	return New(http.StatusForbidden, "forbidden", message, nil)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, "not_found", message, nil)
}

func Conflict(message string) *Error {
	return New(http.StatusConflict, "conflict", message, nil)
}

func Internal(err error) *Error {
	return New(http.StatusInternalServerError, "internal_error", "internal server error", err)
}

func As(target error, want **Error) bool {
	return errors.As(target, want)
}
