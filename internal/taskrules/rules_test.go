package taskrules

import (
	"reflect"
	"testing"
)

func TestLabelRules(t *testing.T) {
	raw := []string{"  urgent ", "", "urgent", "high"}
	want := []string{"urgent", "high"}
	if got := NormalizeLabels(raw); !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeLabels(%q) = %#v, want %#v", raw, got, want)
	}
	if got := DeduplicateLabels(raw); !reflect.DeepEqual(got, want) {
		t.Fatalf("DeduplicateLabels(%q) = %#v, want %#v", raw, got, want)
	}
	if HasLabel([]string{"high-priority"}, "high") {
		t.Fatal("HasLabel should require a full label match")
	}
	if !HasLabel([]string{"high"}, "HIGH") {
		t.Fatal("HasLabel should ignore case")
	}
	if got := CanonicalLabels([]string{" high ", "high"}); !reflect.DeepEqual(got, []string{"high"}) {
		t.Fatalf("CanonicalLabels() = %#v", got)
	}
}
