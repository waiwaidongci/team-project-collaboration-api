package repository

import (
	"context"
	"errors"

	"github.com/lib/pq"

	"team-project-task-api/internal/model"
)

func (r *Repo) CreateProject(ctx context.Context, name, description string, ownerID int64) (model.Project, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Project{}, err
	}
	defer tx.Rollback()

	var project model.Project
	err = tx.QueryRowContext(ctx, `
		INSERT INTO projects (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, owner_id, created_at, updated_at`,
		name, description, ownerID,
	).Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return model.Project{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)`,
		project.ID, ownerID, model.ProjectRoleOwner,
	)
	if err != nil {
		return model.Project{}, err
	}
	if err := tx.Commit(); err != nil {
		return model.Project{}, err
	}
	return project, nil
}

func (r *Repo) GetProjectByID(ctx context.Context, projectID int64) (model.Project, error) {
	var project model.Project
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, owner_id, created_at, updated_at
		FROM projects
		WHERE id = $1`,
		projectID,
	).Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
	return project, err
}

func (r *Repo) ListProjectsByUser(ctx context.Context, userID int64) ([]model.Project, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.name, p.description, p.owner_id, p.created_at, p.updated_at
		FROM projects p
		INNER JOIN project_members pm ON pm.project_id = p.id
		WHERE pm.user_id = $1
		ORDER BY p.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]model.Project, 0)
	for rows.Next() {
		var project model.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}

func (r *Repo) IsMember(ctx context.Context, projectID, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM project_members
			WHERE project_id = $1 AND user_id = $2
		)`,
		projectID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *Repo) AddMember(ctx context.Context, projectID, userID int64, role string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, user_id)
		DO UPDATE SET role = EXCLUDED.role`,
		projectID, userID, role,
	)
	return err
}

func (r *Repo) ListMembers(ctx context.Context, projectID int64) ([]model.ProjectMember, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT pm.id, pm.project_id, pm.user_id, u.name, u.email, pm.role, pm.created_at
		FROM project_members pm
		INNER JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.created_at ASC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]model.ProjectMember, 0)
	for rows.Next() {
		var member model.ProjectMember
		if err := rows.Scan(&member.ID, &member.ProjectID, &member.UserID, &member.Name, &member.Email, &member.Role, &member.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (r *Repo) CreateActivity(ctx context.Context, activity model.Activity) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO activities (project_id, task_id, user_id, action, old_value, new_value)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		activity.ProjectID, activity.TaskID, activity.UserID, activity.Action, activity.OldValue, activity.NewValue,
	)
	return err
}

func (r *Repo) ListActivities(ctx context.Context, projectID int64, limit int) ([]model.Activity, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.project_id, a.task_id, a.user_id, u.name, a.action, a.old_value, a.new_value, a.created_at
		FROM activities a
		INNER JOIN users u ON u.id = a.user_id
		WHERE a.project_id = $1
		ORDER BY a.created_at DESC, a.id DESC
		LIMIT $2`,
		projectID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := make([]model.Activity, 0)
	for rows.Next() {
		var activity model.Activity
		if err := rows.Scan(
			&activity.ID,
			&activity.ProjectID,
			&activity.TaskID,
			&activity.UserID,
			&activity.UserName,
			&activity.Action,
			&activity.OldValue,
			&activity.NewValue,
			&activity.CreatedAt,
		); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, rows.Err()
}

func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
