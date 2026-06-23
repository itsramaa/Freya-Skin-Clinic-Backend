package model

import "time"

// User — domain model
type User struct {
	ID                string    `db:"id"`
	Username          string    `db:"username"`
	PasswordHash      string    `db:"password_hash"`
	IsDefaultPassword bool      `db:"is_default_password"`
	SessionID         *string   `db:"session_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// DTOs

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserPayload struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type LoginResponse struct {
	Token             string      `json:"token"`
	User              UserPayload `json:"user"`
	IsDefaultPassword bool        `json:"is_default_password"`
}

type ChangePasswordRequest struct {
	PasswordBaru string `json:"password_baru" validate:"required,min=8"`
}
