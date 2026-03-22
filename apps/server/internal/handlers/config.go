package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/palta-dev/homectl/apps/server/internal/cache"
	"github.com/palta-dev/homectl/apps/server/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// ConfigResponse represents the sanitized config response
type ConfigResponse struct {
	Title             string               `json:"title"`
	Theme             string               `json:"theme"`
	Background        string               `json:"background,omitempty"`
	BackgroundOpacity *float64             `json:"backgroundOpacity,omitempty"`
	AllowHosts        []string             `json:"allowHosts,omitempty"`
	RequestTimeout    string               `json:"requestTimeout,omitempty"`
	Docker            *config.DockerConfig `json:"docker,omitempty"`
	Groups            []GroupResponse      `json:"groups"`
	Icons             *config.IconsConfig  `json:"icons,omitempty"`
	PasswordProtected bool                 `json:"passwordProtected"`
}

// GroupResponse represents a group in the API response
type GroupResponse struct {
	Name      string            `json:"name"`
	Layout    string            `json:"layout,omitempty"`
	Collapsed bool              `json:"collapsed,omitempty"`
	Services  []ServiceResponse `json:"services"`
}

// ServiceResponse represents a service in the API response
type ServiceResponse struct {
	Name        string          `json:"name"`
	URL         string          `json:"url"`
	Icon        string          `json:"icon,omitempty"`
	Description string          `json:"description,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	NewTab      bool            `json:"newTab,omitempty"`
	PingEnabled bool            `json:"pingEnabled,omitempty"`
	Widgets     []config.Widget `json:"widgets,omitempty"`
	// Note: checks are not returned to client
}

// ConfigHandler returns the sanitized configuration
func ConfigHandler(cfg *config.Config, cacheManager *cache.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Determine authentication status
		isAuthenticated := false
		if cfg.Settings.Password != "" {
			authHeader := c.Get("X-HOMECTL-AUTH")
			authCookie := c.Cookies("homectl_auth")
			
			passwordProvided := authHeader
			if passwordProvided == "" {
				passwordProvided = authCookie
			}

			if passwordProvided != "" {
				err := bcrypt.CompareHashAndPassword([]byte(cfg.Settings.Password), []byte(passwordProvided))
				if err == nil {
					isAuthenticated = true
				}
			}
		} else {
			// If no password set, everything is public
			isAuthenticated = true
		}

		// Try cache first (cache separate results for auth/unauth)
		cacheKey := "config:response:unauth"
		if isAuthenticated {
			cacheKey = "config:response:auth"
		}

		if cached, ok := cacheManager.Get(cacheKey); ok {
			return c.JSON(cached)
		}

		// Build response (strip sensitive data)
		response := ConfigResponse{
			Title:             cfg.Settings.Title,
			Theme:             cfg.Settings.Theme,
			Background:        cfg.Settings.Background,
			BackgroundOpacity: cfg.Settings.BackgroundOpacity,
			PasswordProtected: cfg.Settings.Password != "",
		}

		// Only include sensitive network/discovery/service data if authenticated
		if isAuthenticated {
			response.AllowHosts = cfg.Settings.AllowHosts
			response.RequestTimeout = cfg.Settings.RequestTimeout
			response.Docker = cfg.Settings.Docker
			response.Icons = cfg.Icons
			response.Groups = make([]GroupResponse, len(cfg.Groups))

			for i, group := range cfg.Groups {
				gResp := GroupResponse{
					Name:      group.Name,
					Layout:    group.Layout,
					Collapsed: group.Collapsed,
					Services:  make([]ServiceResponse, len(group.Services)),
				}

				for j, svc := range group.Services {
					gResp.Services[j] = ServiceResponse{
						Name:        svc.Name,
						URL:         svc.URL,
						Icon:        svc.Icon,
						Description: svc.Description,
						Tags:        svc.Tags,
						NewTab:      svc.NewTab,
						PingEnabled: svc.PingEnabled,
						Widgets:     svc.Widgets,
					}
				}
				response.Groups[i] = gResp
			}
		}

		// Cache for 5 seconds
		cacheManager.SetWithTTL(cacheKey, response, 5*time.Second)

		return c.JSON(response)
	}
}

// UpdateConfigHandler updates the configuration settings
func UpdateConfigHandler(cfg *config.Config, cacheManager *cache.Manager, configPath string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use a map to handle partial updates
		var updates map[string]interface{}
		if err := c.BodyParser(&updates); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON: " + err.Error(),
			})
		}

		// Update Settings
		if val, ok := updates["title"].(string); ok {
			cfg.Settings.Title = val
		}
		if val, ok := updates["theme"].(string); ok {
			cfg.Settings.Theme = val
		}
		if val, ok := updates["background"].(string); ok {
			cfg.Settings.Background = val
		}
		if val, ok := updates["backgroundOpacity"].(float64); ok {
			cfg.Settings.BackgroundOpacity = &val
		}
		if val, ok := updates["requestTimeout"].(string); ok {
			cfg.Settings.RequestTimeout = val
		}
		if val, ok := updates["password"].(string); ok {
			if val == "" {
				cfg.Settings.Password = ""
			} else {
				// Hash the password before saving
				hashed, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to hash password: " + err.Error(),
					})
				}
				cfg.Settings.Password = string(hashed)
			}
		}
		if val, ok := updates["allowHosts"].([]interface{}); ok {
			hosts := make([]string, len(val))
			for i, h := range val {
				if s, ok := h.(string); ok {
					hosts[i] = s
				}
			}
			cfg.Settings.AllowHosts = hosts
		}

		if dockerUpdate, ok := updates["docker"].(map[string]interface{}); ok {
			if cfg.Settings.Docker == nil {
				cfg.Settings.Docker = &config.DockerConfig{}
			}
			if val, ok := dockerUpdate["enabled"].(bool); ok {
				cfg.Settings.Docker.Enabled = val
			}
			if val, ok := dockerUpdate["subnet"].(string); ok {
				cfg.Settings.Docker.Subnet = val
			}
			if val, ok := dockerUpdate["ignore"].([]interface{}); ok {
				ignore := make([]string, len(val))
				for i, h := range val {
					if s, ok := h.(string); ok {
						ignore[i] = s
					}
				}
				cfg.Settings.Docker.Ignore = ignore
			}
		}

		// Surgical update for Groups to preserve backend-only fields like Checks
		if val, ok := updates["groups"].([]interface{}); ok {
			newGroups := make([]config.Group, len(val))
			for i, g := range val {
				gMap, _ := g.(map[string]interface{})
				name, _ := gMap["name"].(string)

				// Try to find existing group to preserve its fields
				var existingGroup *config.Group
				for j := range cfg.Groups {
					if cfg.Groups[j].Name == name {
						existingGroup = &cfg.Groups[j]
						break
					}
				}

				group := config.Group{
					Name:      name,
					Layout:    "grid",
					Collapsed: false,
				}
				if layout, ok := gMap["layout"].(string); ok {
					group.Layout = layout
				}
				if collapsed, ok := gMap["collapsed"].(bool); ok {
					group.Collapsed = collapsed
				}

				if services, ok := gMap["services"].([]interface{}); ok {
					group.Services = make([]config.Service, len(services))
					for k, s := range services {
						sMap, _ := s.(map[string]interface{})
						sName, _ := sMap["name"].(string)
						sURL, _ := sMap["url"].(string)

						// Try to find existing service to preserve Checks/Widgets
						var existingSvc *config.Service
						if existingGroup != nil {
							for l := range existingGroup.Services {
								if existingGroup.Services[l].Name == sName {
									existingSvc = &existingGroup.Services[l]
									break
								}
							}
						}

						svc := config.Service{
							Name: sName,
							URL:  sURL,
						}

						// Preserve fields if existing service found
						if existingSvc != nil {
							svc.Checks = existingSvc.Checks
							svc.Widgets = existingSvc.Widgets
							svc.Icon = existingSvc.Icon
							svc.Description = existingSvc.Description
							svc.Tags = existingSvc.Tags
							svc.NewTab = existingSvc.NewTab
							svc.PingEnabled = existingSvc.PingEnabled
						}

						// Update with any new metadata from frontend
						if icon, ok := sMap["icon"].(string); ok {
							svc.Icon = icon
						}
						if desc, ok := sMap["description"].(string); ok {
							svc.Description = desc
						}
						if newTab, ok := sMap["newTab"].(bool); ok {
							svc.NewTab = newTab
						}

						group.Services[k] = svc
					}
				}
				newGroups[i] = group
			}
			cfg.Groups = newGroups
		}

		// Save to file
		if err := cfg.Save(configPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save configuration: " + err.Error(),
			})
		}

		// Clear cache
		cacheManager.Clear("config:*")
		cacheManager.Clear("services:*")

		return c.JSON(fiber.Map{
			"message": "Configuration updated successfully",
			"status":  "success",
		})
	}
}
