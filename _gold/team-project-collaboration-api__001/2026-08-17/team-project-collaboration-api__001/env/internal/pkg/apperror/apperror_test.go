package apperror

import (
	"errors"
	"testing"
)

func TestErrorUnwrap(t *testing.T) {
	inner := errors.New("inner")
	err := Internal(inner)
	if !errors.Is(err, inner) {
		t.Fatal("errors.Is() did not find wrapped error")
	}
	var appErr *Error
	if !As(err, &appErr) {
		t.Fatal("As() did not find apperror.Error")
	}
	if appErr.Status == 0 || appErr.Code == "" || appErr.Message == "" {
		t.Fatalf("unexpected app error: %+v", appErr)
	}
}
