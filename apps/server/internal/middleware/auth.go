package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/palta-dev/homectl/apps/server/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// Auth middleware checks for a password if one is set in the configuration
func Auth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// If no password is set, skip auth
		if cfg.Settings.Password == "" {
			return c.Next()
		}

		// Check X-HOMECTL-AUTH header, password query param, or homectl_auth cookie
		authHeader := c.Get("X-HOMECTL-AUTH")
		authQuery := c.Query("password")
		authCookie := c.Cookies("homectl_auth")

		passwordProvided := authHeader
		if passwordProvided == "" {
			passwordProvided = authQuery
		}
		if passwordProvided == "" {
			passwordProvided = authCookie
		}

		// Use bcrypt to compare the provided password with the stored hash
		err := bcrypt.CompareHashAndPassword([]byte(cfg.Settings.Password), []byte(passwordProvided))
		if err == nil {
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid password",
		})
	}
}
