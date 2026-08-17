package service

import (
	"testing"

	"team-project-task-api/internal/model"
)

func TestValidateTaskInput(t *testing.T) {
	if err := validateTaskInput("Title", model.TaskPriorityMedium, model.TaskStatusTodo); err != nil {
		t.Fatalf("valid input returned error: %v", err)
	}
	if err := validateTaskInput("", model.TaskPriorityMedium, model.TaskStatusTodo); err == nil {
		t.Fatal("empty title should fail")
	}
	if err := validateTaskInput("Title", "unknown", model.TaskStatusTodo); err == nil {
		t.Fatal("unknown priority should fail")
	}
	if err := validateTaskInput("Title", model.TaskPriorityMedium, "unknown"); err == nil {
		t.Fatal("unknown status should fail")
	}
}
