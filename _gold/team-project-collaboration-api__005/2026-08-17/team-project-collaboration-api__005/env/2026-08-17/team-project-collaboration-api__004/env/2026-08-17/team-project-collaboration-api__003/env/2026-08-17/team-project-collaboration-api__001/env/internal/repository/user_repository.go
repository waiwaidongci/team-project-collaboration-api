package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"team-project-task-api/internal/model"
)

func (r *Repo) CreateUser(ctx context.Context, email, name, passwordHash string) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, created_at`,
		strings.ToLower(email), name, passwordHash,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	return user, err
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (model.User, string, error) {
	var user model.User
	var passwordHash string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, created_at
		FROM users
		WHERE email = $1`,
		strings.ToLower(email),
	).Scan(&user.ID, &user.Email, &user.Name, &passwordHash, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, "", sql.ErrNoRows
	}
	return user, passwordHash, err
}

func (r *Repo) GetUserByID(ctx context.Context, userID int64) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, created_at
		FROM users
		WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	return user, err
}

func (r *Repo) IsUserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		strings.ToLower(email),
	).Scan(&exists)
	return exists, err
}
