package service

import (
	"context"

	"team-project-task-api/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, name, passwordHash string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, string, error)
	GetUserByID(ctx context.Context, userID int64) (model.User, error)
	IsUserExists(ctx context.Context, email string) (bool, error)
}

type ProjectRepository interface {
	CreateProject(ctx context.Context, name, description string, ownerID int64) (model.Project, error)
	GetProjectByID(ctx context.Context, projectID int64) (model.Project, error)
	ListProjectsByUser(ctx context.Context, userID int64) ([]model.Project, error)
	IsMember(ctx context.Context, projectID, userID int64) (bool, error)
	AddMember(ctx context.Context, projectID, userID int64, role string) error
	ListMembers(ctx context.Context, projectID int64) ([]model.ProjectMember, error)
	CreateActivity(ctx context.Context, activity model.Activity) error
	ListActivities(ctx context.Context, projectID int64, limit int) ([]model.Activity, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, projectID, createdBy int64, input model.CreateTaskInput) (model.Task, error)
	GetTask(ctx context.Context, projectID, taskID int64) (model.Task, error)
	ListTasks(ctx context.Context, projectID int64, filter model.TaskFilter) ([]model.Task, int64, error)
	UpdateTask(ctx context.Context, projectID, taskID int64, input model.UpdateTaskInput) (model.Task, error)
	DeleteTask(ctx context.Context, projectID, taskID int64) error
	UpdateTaskStatus(ctx context.Context, projectID, taskID int64, status string) error
	UpdateTaskAssignee(ctx context.Context, projectID, taskID int64, assigneeID int64) error
	CreateComment(ctx context.Context, taskID, userID int64, content string) (model.Comment, error)
	ListComments(ctx context.Context, taskID int64) ([]model.Comment, error)
	ListOverdueTasks(ctx context.Context, userID int64) ([]model.Task, error)
}
