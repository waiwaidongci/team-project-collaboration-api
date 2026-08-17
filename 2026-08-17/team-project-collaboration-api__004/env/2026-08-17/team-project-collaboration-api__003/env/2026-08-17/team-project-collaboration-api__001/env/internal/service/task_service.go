package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"team-project-task-api/internal/model"
	"team-project-task-api/internal/pkg/apperror"
)

type TaskService struct {
	projects ProjectRepository
	tasks    TaskRepository
}

func NewTaskService(projects ProjectRepository, tasks TaskRepository) *TaskService {
	return &TaskService{projects: projects, tasks: tasks}
}

func (s *TaskService) CreateTask(ctx context.Context, userID, projectID int64, input model.CreateTaskInput) (model.Task, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return model.Task{}, err
	}
	if input.Status == "" {
		input.Status = model.TaskStatusTodo
	}
	if input.Priority == "" {
		input.Priority = model.TaskPriorityMedium
	}
	if err := validateTaskInput(input.Title, input.Priority, input.Status); err != nil {
		return model.Task{}, err
	}
	if input.AssigneeID != nil {
		if err := s.ensureAssigneeMember(ctx, projectID, *input.AssigneeID); err != nil {
			return model.Task{}, err
		}
	}
	task, err := s.tasks.CreateTask(ctx, projectID, userID, input)
	if err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	_ = s.projects.CreateActivity(ctx, model.Activity{
		ProjectID: projectID,
		TaskID:    &task.ID,
		UserID:    userID,
		Action:    model.ActivityTaskCreated,
	})
	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, userID, projectID int64, filter model.TaskFilter) ([]model.Task, int64, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return nil, 0, err
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	tasks, total, err := s.tasks.ListTasks(ctx, projectID, filter)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return tasks, total, nil
}

func (s *TaskService) GetTask(ctx context.Context, userID, projectID, taskID int64) (model.Task, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return model.Task{}, err
	}
	task, err := s.tasks.GetTask(ctx, projectID, taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, apperror.NotFound("task not found")
	}
	if err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, userID, projectID, taskID int64, input model.UpdateTaskInput) (model.Task, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return model.Task{}, err
	}
	oldTask, err := s.tasks.GetTask(ctx, projectID, taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, apperror.NotFound("task not found")
	}
	if err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	if input.Status != nil {
		if !model.ValidTaskStatuses[*input.Status] {
			return model.Task{}, apperror.BadRequest("invalid task status", nil)
		}
	}
	if input.Priority != nil {
		if !model.ValidTaskPriorities[*input.Priority] {
			return model.Task{}, apperror.BadRequest("invalid task priority", nil)
		}
	}
	if input.AssigneeID != nil {
		if err := s.ensureAssigneeMember(ctx, projectID, *input.AssigneeID); err != nil {
			return model.Task{}, err
		}
	}

	updatedTask, err := s.tasks.UpdateTask(ctx, projectID, taskID, input)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, apperror.NotFound("task not found")
	}
	if err != nil {
		return model.Task{}, apperror.Internal(err)
	}

	if input.Status != nil && oldTask.Status != *input.Status {
		if err := s.projects.CreateActivity(ctx, activityForChange(projectID, taskID, userID, model.ActivityStatusChanged, oldTask.Status, *input.Status)); err != nil {
			return model.Task{}, apperror.Internal(err)
		}
	}
	if input.AssigneeID != nil && !sameInt64(oldTask.AssigneeID, *input.AssigneeID) {
		if err := s.projects.CreateActivity(ctx, activityForChange(projectID, taskID, userID, model.ActivityAssigneeChanged, assigneeValue(oldTask.AssigneeID), strconv.FormatInt(*input.AssigneeID, 10))); err != nil {
			return model.Task{}, apperror.Internal(err)
		}
	}
	return updatedTask, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, userID, projectID, taskID int64, status string) (model.Task, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return model.Task{}, err
	}
	if !model.ValidTaskStatuses[status] {
		return model.Task{}, apperror.BadRequest("invalid task status", nil)
	}
	oldTask, err := s.tasks.GetTask(ctx, projectID, taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, apperror.NotFound("task not found")
	}
	if err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	if oldTask.Status == status {
		return oldTask, nil
	}
	if err := s.tasks.UpdateTaskStatus(ctx, projectID, taskID, status); err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	if err := s.projects.CreateActivity(ctx, activityForChange(projectID, taskID, userID, model.ActivityStatusChanged, oldTask.Status, status)); err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	return s.tasks.GetTask(ctx, projectID, taskID)
}

func (s *TaskService) UpdateAssignee(ctx context.Context, userID, projectID, taskID, assigneeID int64) (model.Task, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return model.Task{}, err
	}
	if err := s.ensureAssigneeMember(ctx, projectID, assigneeID); err != nil {
		return model.Task{}, err
	}
	oldTask, err := s.tasks.GetTask(ctx, projectID, taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, apperror.NotFound("task not found")
	}
	if err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	if oldTask.AssigneeID != nil && *oldTask.AssigneeID == assigneeID {
		return oldTask, nil
	}
	if err := s.tasks.UpdateTaskAssignee(ctx, projectID, taskID, assigneeID); err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	if err := s.projects.CreateActivity(ctx, activityForChange(projectID, taskID, userID, model.ActivityAssigneeChanged, assigneeValue(oldTask.AssigneeID), strconv.FormatInt(assigneeID, 10))); err != nil {
		return model.Task{}, apperror.Internal(err)
	}
	return s.tasks.GetTask(ctx, projectID, taskID)
}

func (s *TaskService) AddComment(ctx context.Context, userID, projectID, taskID int64, content string) (model.Comment, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return model.Comment{}, err
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return model.Comment{}, apperror.BadRequest("comment content is required", nil)
	}
	if _, err := s.tasks.GetTask(ctx, projectID, taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Comment{}, apperror.NotFound("task not found")
		}
		return model.Comment{}, apperror.Internal(err)
	}
	comment, err := s.tasks.CreateComment(ctx, taskID, userID, content)
	if err != nil {
		return model.Comment{}, apperror.Internal(err)
	}
	_ = s.projects.CreateActivity(ctx, model.Activity{
		ProjectID: projectID,
		TaskID:    &taskID,
		UserID:    userID,
		Action:    model.ActivityCommentAdded,
	})
	return comment, nil
}

func (s *TaskService) ListComments(ctx context.Context, userID, projectID, taskID int64) ([]model.Comment, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	if _, err := s.tasks.GetTask(ctx, projectID, taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("task not found")
		}
		return nil, apperror.Internal(err)
	}
	comments, err := s.tasks.ListComments(ctx, taskID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return comments, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, userID, projectID, taskID int64) error {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return err
	}
	if _, err := s.tasks.GetTask(ctx, projectID, taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound("task not found")
		}
		return apperror.Internal(err)
	}
	if err := s.tasks.DeleteTask(ctx, projectID, taskID); err != nil {
		return apperror.Internal(err)
	}
	_ = s.projects.CreateActivity(ctx, model.Activity{
		ProjectID: projectID,
		TaskID:    &taskID,
		UserID:    userID,
		Action:    model.ActivityTaskDeleted,
	})
	return nil
}

func (s *TaskService) ListOverdueTasks(ctx context.Context, userID int64) ([]model.Task, error) {
	tasks, err := s.tasks.ListOverdueTasks(ctx, userID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return tasks, nil
}

func (s *TaskService) ensureMember(ctx context.Context, projectID, userID int64) error {
	isMember, err := s.projects.IsMember(ctx, projectID, userID)
	if err != nil {
		return apperror.Internal(err)
	}
	if !isMember {
		return apperror.Forbidden("project member access required")
	}
	return nil
}

func (s *TaskService) ensureAssigneeMember(ctx context.Context, projectID, assigneeID int64) error {
	isMember, err := s.projects.IsMember(ctx, projectID, assigneeID)
	if err != nil {
		return apperror.Internal(err)
	}
	if !isMember {
		return apperror.BadRequest("assignee must be a project member", nil)
	}
	return nil
}

func validateTaskInput(title, priority, status string) error {
	if strings.TrimSpace(title) == "" {
		return apperror.BadRequest("task title is required", nil)
	}
	if !model.ValidTaskStatuses[status] {
		return apperror.BadRequest("invalid task status", nil)
	}
	if !model.ValidTaskPriorities[priority] {
		return apperror.BadRequest("invalid task priority", nil)
	}
	return nil
}

func activityForChange(projectID, taskID, userID int64, action, oldValue, newValue string) model.Activity {
	return model.Activity{
		ProjectID: projectID,
		TaskID:    &taskID,
		UserID:    userID,
		Action:    action,
		OldValue:  &oldValue,
		NewValue:  &newValue,
	}
}

func assigneeValue(assigneeID *int64) string {
	if assigneeID == nil {
		return ""
	}
	return strconv.FormatInt(*assigneeID, 10)
}

func sameInt64(a *int64, b int64) bool {
	return a != nil && *a == b
}
