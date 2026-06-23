package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"freya-skin-clinic-backend/internal/config"
	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/pkg/hash"
	"freya-skin-clinic-backend/internal/pkg/jwt"
	"freya-skin-clinic-backend/internal/repository"
)

var ErrInvalidCredentials = errors.New("Kredensial tidak valid")

type AuthService interface {
	Login(ctx context.Context, username, password string) (*model.LoginResponse, error)
	ChangePassword(ctx context.Context, userID, username, newPassword string) (*model.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Login(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := hash.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate session_id baru — invalidate session lama di device lain
	sessionID := uuid.New().String()
	if err := s.userRepo.UpdateSessionID(ctx, user.ID, sessionID); err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, sessionID, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User: model.UserPayload{
			ID:       user.ID,
			Username: user.Username,
		},
		IsDefaultPassword: user.IsDefaultPassword,
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID, username, newPassword string) (*model.LoginResponse, error) {
	passwordHash, err := hash.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.UpdatePassword(ctx, userID, passwordHash); err != nil {
		return nil, err
	}

	// Generate session_id baru setelah ganti password
	sessionID := uuid.New().String()
	if err := s.userRepo.UpdateSessionID(ctx, userID, sessionID); err != nil {
		return nil, err
	}

	// Return token baru agar frontend tidak perlu re-login
	token, err := jwt.GenerateToken(userID, username, sessionID, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User: model.UserPayload{
			ID:       userID,
			Username: username,
		},
		IsDefaultPassword: false,
	}, nil
}
