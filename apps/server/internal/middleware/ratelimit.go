package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

// RateLimiterConfig holds rate limiting configuration
type RateLimiterConfig struct {
	RequestsPerSecond float64
	BurstSize         int
	KeyFunc           func(*fiber.Ctx) string
}

// RateLimiter manages rate limiting for multiple clients
type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientLimiter
	config   RateLimiterConfig
	cleanup  time.Duration
	lastClean time.Time
}

type clientLimiter struct {
	limiter *rate.Limiter
	lastAccess time.Time
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	if cfg.RequestsPerSecond <= 0 {
		cfg.RequestsPerSecond = 10
	}
	if cfg.BurstSize <= 0 {
		cfg.BurstSize = 20
	}
	if cfg.KeyFunc == nil {
		// Default: use IP address
		cfg.KeyFunc = func(c *fiber.Ctx) string {
			return c.IP()
		}
	}

	return &RateLimiter{
		clients: make(map[string]*clientLimiter),
		config:  cfg,
		cleanup: time.Minute,
	}
}

// Middleware returns the Fiber middleware handler
func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := rl.config.KeyFunc(c)
		limiter := rl.getLimiter(key)

		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
				"retry":   limiter.Reserve().Delay().Seconds(),
			})
		}

		return c.Next()
	}
}

// getLimiter returns or creates a rate limiter for a client
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Cleanup old entries periodically
	if time.Since(rl.lastClean) > rl.cleanup {
		rl.cleanupOldClients()
	}

	client, exists := rl.clients[key]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize)
		client = &clientLimiter{
			limiter:    limiter,
			lastAccess: time.Now(),
		}
		rl.clients[key] = client
	}

	client.lastAccess = time.Now()
	return client.limiter
}

// cleanupOldClients removes clients that haven't been accessed recently
func (rl *RateLimiter) cleanupOldClients() {
	now := time.Now()
	for key, client := range rl.clients {
		if now.Sub(client.lastAccess) > 10*time.Minute {
			delete(rl.clients, key)
		}
	}
	rl.lastClean = now
}

// DefaultRateLimiter creates a standard rate limiter for API endpoints
func DefaultRateLimiter() fiber.Handler {
	rl := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
	})
	return rl.Middleware()
}

// StrictRateLimiter creates a stricter rate limiter for auth endpoints
func StrictRateLimiter() fiber.Handler {
	rl := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: 5,
		BurstSize:         10,
	})
	return rl.Middleware()
}
