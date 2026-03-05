package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIError represents an API error response
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// ErrorHandler handles Fiber errors
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(APIError{
		Error:   httpStatusText(code),
		Message: err.Error(),
		Code:    code,
	})
}

func httpStatusText(code int) string {
	switch code {
	case fiber.StatusBadRequest:
		return "Bad Request"
	case fiber.StatusUnauthorized:
		return "Unauthorized"
	case fiber.StatusForbidden:
		return "Forbidden"
	case fiber.StatusNotFound:
		return "Not Found"
	case fiber.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Error"
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  int64  `json:"uptime"`
}

var startTime = time.Now()

// HealthHandler returns the health status
func HealthHandler(version string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(HealthResponse{
			Status:  "healthy",
			Version: version,
			Uptime:  int64(time.Since(startTime).Seconds()),
		})
	}
}
