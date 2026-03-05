package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/palta-dev/homectl/apps/server/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Password string `json:"password"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// If no password is set, we don't allow "logging in" - the UI should handle this
		if cfg.Settings.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No password configured"})
		}

		err := bcrypt.CompareHashAndPassword([]byte(cfg.Settings.Password), []byte(input.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
		}

		// In a full system, we'd use a JWT or Session ID. 
		// For this implementation, we'll return a simple success and the frontend will 
		// keep the password in memory (or we can set a cookie).
		// Let's set a simple cookie for "session" feel.
		c.Cookie(&fiber.Cookie{
			Name:     "homectl_auth",
			Value:    input.Password, // In production, use a token. For now, this facilitates the existing X-HOMECTL-AUTH header.
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: true,
			SameSite: "Lax",
		})

		return c.JSON(fiber.Map{"message": "Logged in successfully"})
	}
}
