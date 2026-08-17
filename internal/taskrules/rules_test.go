package taskrules

import (
	"errors"
	"testing"
)

func TestErrorClassification(t *testing.T) {
	lookupErr := WrapLookupError(ErrNotFound)
	if !errors.Is(lookupErr, ErrNotFound) {
		t.Fatal("wrapped error should preserve the not-found sentinel")
	}
	if CodeFor(lookupErr) != "not_found" {
		t.Fatalf("CodeFor() = %q, want not_found", CodeFor(lookupErr))
	}
	timeoutErr := errors.New("request timed out")
	if !IsRetryable(timeoutErr) {
		t.Fatal("timeout error should be retryable")
	}
	if NormalizeError(ErrNotFound) != ErrNotFound {
		t.Fatal("NormalizeError should keep the original sentinel")
	}
}
