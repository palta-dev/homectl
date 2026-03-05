package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/palta-dev/homectl/apps/server/internal/cache"
	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/handlers"
	"github.com/palta-dev/homectl/apps/server/internal/middleware"
	"github.com/palta-dev/homectl/apps/server/internal/network"
	"github.com/palta-dev/homectl/apps/server/internal/widgets"
)

const (
	version     = "0.1.0"
	defaultPort = 8080
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	port := flag.Int("port", defaultPort, "Server port")
	flag.Parse()

	log.Printf("Starting homectl v%s", version)

	// Ensure config file exists (creates default if missing)
	if err := config.EnsureExists(*configPath); err != nil {
		log.Printf("Warning: Could not create default config: %v", err)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded configuration from %s", *configPath)

	// Initialize cache
	cacheManager := cache.New(cache.Config{
		DefaultTTL:   30 * time.Second,
		MaxEntries:   500,
		CleanupInterval: 1 * time.Minute,
	})
	log.Println("Cache initialized")

	// Initialize network client (SSRF protection)
	netClient, err := network.NewClient(network.Config{
		AllowHosts:          cfg.Settings.AllowHosts,
		BlockPrivateMetaIPs: cfg.Settings.BlockPrivateMetaIPs,
		Timeout:             cfg.Settings.GetTimeout(),
	})
	if err != nil {
		log.Fatalf("Failed to initialize network client: %v", err)
	}

	// Initialize widget registry
	widgetRegistry := widgets.NewRegistry(netClient)
	widgets.RegisterBuiltins(widgetRegistry)
	log.Println("Widget registry initialized")

	// Determine static directory path
	exePath, _ := os.Executable()
	baseDir := filepath.Dir(exePath)
	staticDir := filepath.Join(baseDir, "static")

	// Fallback if running with 'go run'
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = "./static"
		if _, err := os.Stat(staticDir); os.IsNotExist(err) {
			// Try one more common dev path
			staticDir = "apps/server/static"
		}
	}
	log.Printf("Serving static files from %s", staticDir)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "homectl/" + version,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorHandler: handlers.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,X-HOMECTL-AUTH",
	}))

	// API routes with rate limiting (MUST be before SPA fallback)
	api := app.Group("/api", middleware.DefaultRateLimiter())
	api.Get("/health", handlers.HealthHandler(version))
	api.Get("/config", handlers.ConfigHandler(cfg, cacheManager))
	api.Get("/config/auth-check", middleware.Auth(cfg), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"authenticated": true})
	})
	api.Post("/auth/login", handlers.LoginHandler(cfg))
	api.Post("/auth/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("homectl_auth")
		return c.JSON(fiber.Map{"message": "Logged out"})
	})

	// Strict limiter for config updates to prevent brute forcing
	configLimiter := limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many configuration update attempts. Please wait a minute.",
			})
		},
	})

	api.Put("/config", configLimiter, middleware.Auth(cfg), handlers.UpdateConfigHandler(cfg, cacheManager, *configPath))
	api.Get("/services", handlers.ServicesHandler(cfg, cacheManager, widgetRegistry))

	// Serve static files (embedded frontend)
	app.Static("/", staticDir, fiber.Static{
		Compress:      true,
		MaxAge:        3600,
		CacheDuration: time.Hour,
		Browse:        false,
		Index:         "index.html",
	})

	// SPA fallback - serve index.html for all non-API routes
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(staticDir, "index.html"))
	})

	// Start config watcher for hot reload
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.Watch(ctx, *configPath, func(newCfg *config.Config) {
		log.Println("Configuration reloaded")
		cfg = newCfg
		cacheManager.Clear("services:*")
		cacheManager.Clear("config:*")
	}); err != nil {
		log.Printf("Warning: Config watcher failed: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")
		cancel()
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
