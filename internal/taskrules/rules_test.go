package taskrules

import (
	"context"
	"reflect"
	"testing"
)

func TestCleanupOrder(t *testing.T) {
	got := CleanupOrder([]string{"db", "cache", "queue"})
	want := []string{"db", "cache", "queue"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CleanupOrder() = %#v, want %#v", got, want)
	}
}

func TestContextPropagation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	child, childCancel := WithTaskContext(parent)
	defer childCancel()
	cancel()
	select {
	case <-child.Done():
	default:
		t.Fatal("task context did not inherit parent cancellation")
	}
}

func TestCleanupResource(t *testing.T) {
	recorder := &CloseRecorder{}
	CleanupResource(recorder, "connection")
	if len(recorder.Closed) != 1 || recorder.Closed[0] != "connection" {
		t.Fatalf("resource should be closed after cleanup, got %#v", recorder.Closed)
	}
}

func TestRecoverPanic(t *testing.T) {
	if got := RecoverPanic(func() { panic("boom") }); got == nil {
		t.Fatal("RecoverPanic should recover panics")
	}
}
