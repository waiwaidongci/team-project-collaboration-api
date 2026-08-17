package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"team-project-task-api/internal/model"
)

func (r *Repo) CreateTask(ctx context.Context, projectID, createdBy int64, input model.CreateTaskInput) (model.Task, error) {
	var task model.Task
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO tasks (project_id, title, description, priority, status, assignee_id, due_date, labels, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, project_id, title, description, priority, status, assignee_id, due_date, labels, created_by, created_at, updated_at`,
		projectID,
		input.Title,
		input.Description,
		input.Priority,
		input.Status,
		input.AssigneeID,
		input.DueDate,
		pq.Array(input.Labels),
		createdBy,
	).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.AssigneeID,
		&task.DueDate,
		pq.Array(&task.Labels),
		&task.CreatedBy,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}

func (r *Repo) GetTask(ctx context.Context, projectID, taskID int64) (model.Task, error) {
	var task model.Task
	err := r.db.QueryRowContext(ctx, `
		SELECT id, project_id, title, description, priority, status, assignee_id, due_date, labels, created_by, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND project_id = $2`,
		taskID, projectID,
	).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.AssigneeID,
		&task.DueDate,
		pq.Array(&task.Labels),
		&task.CreatedBy,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}

func (r *Repo) ListTasks(ctx context.Context, projectID int64, filter model.TaskFilter) ([]model.Task, int64, error) {
	pagination := NewPagination(filter.Page, filter.PageSize)

	where := "project_id = $1"
	args := []any{projectID}
	if filter.Status != "" {
		args = append(args, filter.Status)
		where += fmt.Sprintf(" AND status = $%d", len(args))
	}
	if filter.AssigneeID > 0 {
		args = append(args, filter.AssigneeID)
		where += fmt.Sprintf(" AND assignee_id = $%d", len(args))
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tasks WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, project_id, title, description, priority, status, assignee_id, due_date, labels, created_by, created_at, updated_at
		FROM tasks
		WHERE ` + where + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)
	args = append(args, pagination.PageSize, pagination.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.AssigneeID,
			&task.DueDate,
			pq.Array(&task.Labels),
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *Repo) UpdateTask(ctx context.Context, projectID, taskID int64, input model.UpdateTaskInput) (model.Task, error) {
	setParts := make([]string, 0)
	args := []any{}
	if input.Title != nil {
		args = append(args, *input.Title)
		setParts = append(setParts, fmt.Sprintf("title = $%d", len(args)))
	}
	if input.Description != nil {
		args = append(args, *input.Description)
		setParts = append(setParts, fmt.Sprintf("description = $%d", len(args)))
	}
	if input.Priority != nil {
		args = append(args, *input.Priority)
		setParts = append(setParts, fmt.Sprintf("priority = $%d", len(args)))
	}
	if input.Status != nil {
		args = append(args, *input.Status)
		setParts = append(setParts, fmt.Sprintf("status = $%d", len(args)))
	}
	if input.AssigneeID != nil {
		args = append(args, *input.AssigneeID)
		setParts = append(setParts, fmt.Sprintf("assignee_id = $%d", len(args)))
	}
	if input.DueDate != nil {
		args = append(args, *input.DueDate)
		setParts = append(setParts, fmt.Sprintf("due_date = $%d", len(args)))
	}
	if input.Labels != nil {
		args = append(args, pq.Array(*input.Labels))
		setParts = append(setParts, fmt.Sprintf("labels = $%d", len(args)))
	}
	if len(setParts) == 0 {
		return r.GetTask(ctx, projectID, taskID)
	}

	setParts = append(setParts, "updated_at = now()")
	args = append(args, taskID, projectID)
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d AND project_id = $%d", strings.Join(setParts, ", "), len(args)-1, len(args))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return model.Task{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}
	if affected == 0 {
		return model.Task{}, sql.ErrNoRows
	}
	return r.GetTask(ctx, projectID, taskID)
}

func (r *Repo) DeleteTask(ctx context.Context, projectID, taskID int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1 AND project_id = $2", taskID, projectID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repo) UpdateTaskStatus(ctx context.Context, projectID, taskID int64, status string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE tasks SET status = $1, updated_at = now()
		WHERE id = $2 AND project_id = $3`,
		status, taskID, projectID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repo) UpdateTaskAssignee(ctx context.Context, projectID, taskID int64, assigneeID int64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE tasks SET assignee_id = $1, updated_at = now()
		WHERE id = $2 AND project_id = $3`,
		assigneeID, taskID, projectID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repo) CreateComment(ctx context.Context, taskID, userID int64, content string) (model.Comment, error) {
	var comment model.Comment
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO comments (task_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, task_id, user_id, content, created_at`,
		taskID, userID, content,
	).Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		return model.Comment{}, err
	}
	err = r.db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = $1", userID).Scan(&comment.AuthorName)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (r *Repo) ListComments(ctx context.Context, taskID int64) ([]model.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.task_id, c.user_id, u.name, c.content, c.created_at
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.task_id = $1
		ORDER BY c.created_at ASC, c.id ASC`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.AuthorName, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (r *Repo) ListOverdueTasks(ctx context.Context, userID int64) ([]model.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.project_id, p.name, t.title, t.description, t.priority, t.status, t.assignee_id, t.due_date, t.labels, t.created_by, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN projects p ON p.id = t.project_id
		WHERE t.assignee_id = $1
		  AND t.due_date < CURRENT_DATE
		  AND t.status <> 'done'
		ORDER BY t.due_date ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.ProjectName,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.AssigneeID,
			&task.DueDate,
			pq.Array(&task.Labels),
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
