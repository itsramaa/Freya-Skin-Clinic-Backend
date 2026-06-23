package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

var ErrUserNotFound = errors.New("user tidak ditemukan")

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateSessionID(ctx context.Context, userID, sessionID string) error
	GetSessionID(ctx context.Context, userID string) (*string, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, is_default_password, session_id, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash,
		&user.IsDefaultPassword, &user.SessionID,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, is_default_password = false, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *userRepository) UpdateSessionID(ctx context.Context, userID, sessionID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET session_id = $1, updated_at = NOW() WHERE id = $2`,
		sessionID, userID,
	)
	return err
}

func (r *userRepository) GetSessionID(ctx context.Context, userID string) (*string, error) {
	var sessionID *string
	err := r.db.QueryRow(ctx, `SELECT session_id FROM users WHERE id = $1`, userID).Scan(&sessionID)
	if err != nil {
		return nil, err
	}
	return sessionID, nil
}
