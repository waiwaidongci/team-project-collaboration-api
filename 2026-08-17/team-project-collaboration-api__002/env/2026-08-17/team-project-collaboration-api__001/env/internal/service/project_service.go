package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"team-project-task-api/internal/model"
	"team-project-task-api/internal/pkg/apperror"
)

type ProjectService struct {
	projects ProjectRepository
	users    UserRepository
}

func NewProjectService(projects ProjectRepository, users UserRepository) *ProjectService {
	return &ProjectService{projects: projects, users: users}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID int64, name, description string) (model.Project, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return model.Project{}, apperror.BadRequest("project name is required", nil)
	}
	project, err := s.projects.CreateProject(ctx, name, description, userID)
	if err != nil {
		return model.Project{}, apperror.Internal(err)
	}
	return project, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, userID int64) ([]model.Project, error) {
	projects, err := s.projects.ListProjectsByUser(ctx, userID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return projects, nil
}

func (s *ProjectService) GetProject(ctx context.Context, userID, projectID int64) (model.Project, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return model.Project{}, err
	}
	project, err := s.projects.GetProjectByID(ctx, projectID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Project{}, apperror.NotFound("project not found")
	}
	if err != nil {
		return model.Project{}, apperror.Internal(err)
	}
	return project, nil
}

func (s *ProjectService) InviteMember(ctx context.Context, actorID, projectID int64, email, role string) (model.ProjectMember, error) {
	if err := s.ensureMember(ctx, projectID, actorID); err != nil {
		return model.ProjectMember{}, err
	}
	project, err := s.projects.GetProjectByID(ctx, projectID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ProjectMember{}, apperror.NotFound("project not found")
	}
	if err != nil {
		return model.ProjectMember{}, apperror.Internal(err)
	}
	if project.OwnerID != actorID {
		return model.ProjectMember{}, apperror.Forbidden("only the project owner can invite members")
	}

	role = strings.TrimSpace(role)
	if role == "" {
		role = model.ProjectRoleMember
	}
	if !model.ValidProjectRoles[role] {
		return model.ProjectMember{}, apperror.BadRequest("invalid project role", nil)
	}

	user, _, err := s.users.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ProjectMember{}, apperror.NotFound("user not found")
	}
	if err != nil {
		return model.ProjectMember{}, apperror.Internal(err)
	}
	if err := s.projects.AddMember(ctx, projectID, user.ID, role); err != nil {
		return model.ProjectMember{}, apperror.Internal(err)
	}
	return model.ProjectMember{
		ProjectID: projectID,
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      role,
	}, nil
}

func (s *ProjectService) ListMembers(ctx context.Context, userID, projectID int64) ([]model.ProjectMember, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	members, err := s.projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return members, nil
}

func (s *ProjectService) ListActivities(ctx context.Context, userID, projectID int64) ([]model.Activity, error) {
	if err := s.ensureMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	activities, err := s.projects.ListActivities(ctx, projectID, 50)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return activities, nil
}

func (s *ProjectService) IsMember(ctx context.Context, userID, projectID int64) (bool, error) {
	isMember, err := s.projects.IsMember(ctx, projectID, userID)
	if err != nil {
		return false, apperror.Internal(err)
	}
	return isMember, nil
}

func (s *ProjectService) ensureMember(ctx context.Context, projectID, userID int64) error {
	isMember, err := s.projects.IsMember(ctx, projectID, userID)
	if err != nil {
		return apperror.Internal(err)
	}
	if !isMember {
		return apperror.Forbidden("project member access required")
	}
	return nil
}
