package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"team-project-task-api/internal/model"
	"team-project-task-api/internal/pkg/apperror"
	"team-project-task-api/internal/pkg/password"
	"team-project-task-api/internal/pkg/token"
)

type AuthService struct {
	users       UserRepository
	tokenSecret string
	tokenTTL    time.Duration
}

type AuthResult struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func NewAuthService(users UserRepository, tokenSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{users: users, tokenSecret: tokenSecret, tokenTTL: tokenTTL}
}

func (s *AuthService) Register(ctx context.Context, email, name, rawPassword string) (AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)
	if email == "" || name == "" {
		return AuthResult{}, apperror.BadRequest("email and name are required", nil)
	}
	if len(rawPassword) < 8 {
		return AuthResult{}, apperror.BadRequest("password must be at least 8 characters", nil)
	}

	exists, err := s.users.IsUserExists(ctx, email)
	if err != nil {
		return AuthResult{}, apperror.Internal(err)
	}
	if exists {
		return AuthResult{}, apperror.Conflict("email already registered")
	}

	passwordHash, err := password.Hash(rawPassword)
	if err != nil {
		return AuthResult{}, apperror.Internal(err)
	}
	user, err := s.users.CreateUser(ctx, email, name, passwordHash)
	if err != nil {
		return AuthResult{}, apperror.Internal(err)
	}
	return s.authResult(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, rawPassword string) (AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, passwordHash, err := s.users.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return AuthResult{}, apperror.Unauthorized("invalid email or password")
	}
	if err != nil {
		return AuthResult{}, apperror.Internal(err)
	}
	if !password.Check(passwordHash, rawPassword) {
		return AuthResult{}, apperror.Unauthorized("invalid email or password")
	}
	return s.authResult(ctx, user)
}

func (s *AuthService) GetUser(ctx context.Context, userID int64) (model.User, error) {
	user, err := s.users.GetUserByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, apperror.NotFound("user not found")
	}
	if err != nil {
		return model.User{}, apperror.Internal(err)
	}
	return user, nil
}

func (s *AuthService) authResult(ctx context.Context, user model.User) (AuthResult, error) {
	jwtToken, err := token.Generate(user.ID, user.Email, user.Name, s.tokenSecret, s.tokenTTL)
	if err != nil {
		return AuthResult{}, apperror.Internal(err)
	}
	return AuthResult{Token: jwtToken, User: user}, nil
}
