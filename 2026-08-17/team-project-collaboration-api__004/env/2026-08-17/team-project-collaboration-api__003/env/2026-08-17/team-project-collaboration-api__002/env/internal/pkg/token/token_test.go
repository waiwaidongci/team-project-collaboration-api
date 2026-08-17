package token

import (
	"testing"
	"time"
)

func TestGenerateAndParse(t *testing.T) {
	raw, err := Generate(42, "alice@example.com", "Alice", "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	claims, err := Parse(raw, "test-secret")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.UserID != 42 || claims.Email != "alice@example.com" || claims.Name != "Alice" {
		t.Fatalf("Parse() claims = %+v", claims)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	raw, err := Generate(1, "alice@example.com", "Alice", "secret", time.Hour)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := Parse(raw, "other-secret"); err == nil {
		t.Fatal("Parse() expected error for wrong secret")
	}
}
