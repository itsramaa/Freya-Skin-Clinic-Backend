package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/pkg/response"
	"freya-skin-clinic-backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid", nil)
	}

	if req.Username == "" || req.Password == "" {
		return response.Error(c, http.StatusBadRequest, "Username dan password wajib diisi", nil)
	}

	loginResp, err := h.authService.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return response.Error(c, http.StatusUnauthorized, "Kredensial tidak valid", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "Terjadi kesalahan server", nil)
	}

	return response.Success(c, http.StatusOK, "Login berhasil", loginResp)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req model.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid", nil)
	}

	if len(req.PasswordBaru) < 8 {
		return response.Error(c, http.StatusBadRequest, "Password minimal 8 karakter", nil)
	}

	userID, _ := c.Locals("user_id").(string)
	username, _ := c.Locals("username").(string)
	if userID == "" {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	loginResp, err := h.authService.ChangePassword(c.Context(), userID, username, req.PasswordBaru)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengubah password", nil)
	}

	return response.Success(c, http.StatusOK, "Password berhasil diperbarui", loginResp)
}
