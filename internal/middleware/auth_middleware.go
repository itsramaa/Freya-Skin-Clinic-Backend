package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"freya-skin-clinic-backend/internal/config"
	"freya-skin-clinic-backend/internal/pkg/jwt"
	"freya-skin-clinic-backend/internal/pkg/response"
	"freya-skin-clinic-backend/internal/repository"
)

func JWTMiddleware(cfg *config.Config, userRepo repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, http.StatusUnauthorized, "Authorization token required", nil)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Error(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
		}

		claims, err := jwt.ValidateToken(parts[1], cfg.JWTSecret)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				return response.Error(c, http.StatusUnauthorized, "Token expired", nil)
			}
			return response.Error(c, http.StatusUnauthorized, "Invalid token", nil)
		}

		// Validasi session_id dari DB — cegah multi-device session
		dbSessionID, err := userRepo.GetSessionID(c.Context(), claims.UserID)
		if err != nil || dbSessionID == nil || *dbSessionID != claims.SessionID {
			return response.Error(c, http.StatusUnauthorized, "Sesi telah berakhir. Silakan login kembali.", nil)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("session_id", claims.SessionID)

		return c.Next()
	}
}
