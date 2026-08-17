package taskrules

import "testing"

func TestAssigneeRules(t *testing.T) {
	if ValidAssigneeID(nil) {
		t.Fatal("nil assignee should be invalid")
	}
	zero := int64(0)
	if ValidAssigneeID(&zero) {
		t.Fatal("zero assignee should be invalid")
	}
	valid := int64(7)
	if !ValidAssigneeID(&valid) {
		t.Fatal("positive assignee should be valid")
	}

	if got := FormatAssigneeID(nil); got != "" {
		t.Fatalf("FormatAssigneeID(nil) = %q, want empty string", got)
	}
}

func TestAssigneeHelpers(t *testing.T) {
	seven := int64(7)
	if got := FormatAssigneeID(&seven); got != "7" {
		t.Fatalf("FormatAssigneeID() = %q", got)
	}
	if !SameAssignee(&seven, 7) {
		t.Fatal("SameAssignee should return true for equal values")
	}
	if SameAssignee(nil, 7) {
		t.Fatal("SameAssignee should return false when current assignee is nil")
	}
	if AssigneeFromID(0) != nil {
		t.Fatal("AssigneeFromID(0) should return nil")
	}
}
