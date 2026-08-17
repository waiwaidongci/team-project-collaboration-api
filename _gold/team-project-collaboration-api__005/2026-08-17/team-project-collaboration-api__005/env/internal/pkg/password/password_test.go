package password

import "testing"

func TestHashAndCheck(t *testing.T) {
	hash, err := Hash("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if hash == "" || hash == "correct-horse-battery-staple" {
		t.Fatalf("Hash() returned unexpected value: %q", hash)
	}
	if !Check(hash, "correct-horse-battery-staple") {
		t.Fatal("Check() returned false for correct password")
	}
	if Check(hash, "wrong-password") {
		t.Fatal("Check() returned true for wrong password")
	}
}
