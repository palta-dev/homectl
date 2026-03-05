package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/palta-dev/homectl/apps/server/internal/cache"
	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/discovery"
	"github.com/palta-dev/homectl/apps/server/internal/widgets"
	"golang.org/x/crypto/bcrypt"
)

// ServicesResponse represents the services API response
type ServicesResponse struct {
	Groups []GroupWithStatus `json:"groups"`
}

// GroupWithStatus represents a group with service statuses
type GroupWithStatus struct {
	Name      string            `json:"name"`
	Layout    string            `json:"layout,omitempty"`
	Collapsed bool              `json:"collapsed,omitempty"`
	Services  []ServiceWithStatus `json:"services"`
}

// ServiceWithStatus represents a service with its status
type ServiceWithStatus struct {
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	Icon        string        `json:"icon,omitempty"`
	Favicon     string        `json:"favicon,omitempty"`
	Description string        `json:"description,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	NewTab      bool          `json:"newTab,omitempty"`
	PingEnabled bool          `json:"pingEnabled,omitempty"`
	Status      ServiceStatus `json:"status"`
	Widgets     []WidgetResult `json:"widgets,omitempty"`
	IsDiscovered bool          `json:"isDiscovered,omitempty"`
}

// ServiceStatus represents the health status of a service
type ServiceStatus struct {
	State     string  `json:"state"` // up, down, degraded, unknown
	Latency   *int64  `json:"latency,omitempty"`
	LastCheck *string `json:"lastCheck,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// WidgetResult represents a widget result
type WidgetResult struct {
	Label      string      `json:"label,omitempty"`
	Value      interface{} `json:"value"`
	Formatted  string      `json:"formatted,omitempty"`
	State      string      `json:"state,omitempty"`
	LastUpdate *string     `json:"lastUpdated,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// ServicesHandler returns services with their current status
func ServicesHandler(cfg *config.Config, cacheManager *cache.Manager, widgetRegistry *widgets.Registry) fiber.Handler {
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

		if !isAuthenticated {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Password required",
			})
		}

		// Try cache first
		cacheKey := "services:all"
		if cached, ok := cacheManager.Get(cacheKey); ok {
			return c.JSON(cached)
		}

		// Merge configured and discovered services
		var allGroups []GroupWithStatus
// Perform Auto-Discovery
var discoveredServices []config.Service
if cfg.Settings.Docker != nil && cfg.Settings.Docker.Enabled {
	log.Printf("[SERVICES] Discovery starting via %s...", cfg.Settings.Docker.Socket)
	dockerDiscoverer, err := discovery.NewDockerDiscoverer(cfg.Settings.Docker.Socket, cfg.Settings.Docker.LabelPrefix, cfg.Settings.BaseHost)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		ds, err := dockerDiscoverer.DiscoverServices(ctx)
		cancel()
		if err == nil {
			discoveredServices = ds
			log.Printf("[SERVICES] Discovery found %d services", len(ds))
		} else {
			log.Printf("[SERVICES ERROR] Discovery failed: %v", err)
		}
		dockerDiscoverer.Close()
	} else {
		log.Printf("[SERVICES ERROR] Failed to create discoverer: %v", err)
	}
}

	// Filter discovered services based on Ignore list
	var filteredDiscovered []config.Service
	ignoreMap := make(map[string]bool)
if cfg.Settings.Docker != nil {
	for _, name := range cfg.Settings.Docker.Ignore {
		ignoreMap[name] = true
	}
}

for _, svc := range discoveredServices {
	if !ignoreMap[svc.Name] {
		filteredDiscovered = append(filteredDiscovered, svc)
	} else {
		log.Printf("[SERVICES] Ignoring discovered service: %s", svc.Name)
	}
}

// Add discovered services as a separate group if any exist
if len(filteredDiscovered) > 0 {
	log.Printf("[SERVICES] Adding %d discovered services to response", len(filteredDiscovered))
	discoveredGroup := GroupWithStatus{
				Collapsed: false,
				Services:  make([]ServiceWithStatus, len(filteredDiscovered)),
			}
			for i, svc := range filteredDiscovered {
				discoveredGroup.Services[i] = ServiceWithStatus{
					Name:         svc.Name,
					URL:          svc.URL,
					Icon:         svc.Icon,
					Favicon:      svc.Favicon,
					Description:  svc.Description,
					Tags:         svc.Tags,
					NewTab:       svc.NewTab,
					Status:       getServiceStatus(svc, cacheManager),
					IsDiscovered: true,
				}
			}
			allGroups = append(allGroups, discoveredGroup)
		}

		// Add configured groups
		for _, group := range cfg.Groups {
			allGroups = append(allGroups, processGroup(group, cacheManager, widgetRegistry))
		}

		response := ServicesResponse{
			Groups: allGroups,
		}

		// Cache for 30 seconds
		cacheManager.SetWithTTL(cacheKey, response, 30*time.Second)

		return c.JSON(response)
	}
}

func processGroup(group config.Group, cacheManager *cache.Manager, widgetRegistry *widgets.Registry) GroupWithStatus {
	gResp := GroupWithStatus{
		Name:      group.Name,
		Layout:    group.Layout,
		Collapsed: group.Collapsed,
		Services:  make([]ServiceWithStatus, len(group.Services)),
	}

	for i, svc := range group.Services {
		gResp.Services[i] = processService(svc, cacheManager, widgetRegistry)
	}

	return gResp
}

func processService(svc config.Service, cacheManager *cache.Manager, widgetRegistry *widgets.Registry) ServiceWithStatus {
	sResp := ServiceWithStatus{
		Name:        svc.Name,
		URL:         svc.URL,
		Icon:        svc.Icon,
		Description: svc.Description,
		Tags:        svc.Tags,
		NewTab:      svc.NewTab,
		PingEnabled: svc.PingEnabled,
		Status:      getServiceStatus(svc, cacheManager),
	}

	// Process widgets
	if len(svc.Widgets) > 0 {
		sResp.Widgets = make([]WidgetResult, len(svc.Widgets))
		for i, w := range svc.Widgets {
			sResp.Widgets[i] = executeWidget(w, cacheManager, widgetRegistry)
		}
	}

	return sResp
}

func getServiceStatus(svc config.Service, cacheManager *cache.Manager) ServiceStatus {
	now := time.Now().UTC().Format(time.RFC3339)
	state := "unknown"

	// If discovered and has a state tag
	for _, tag := range svc.Tags {
		if tag == "exited" || tag == "stopped" {
			state = "down"
			break
		}
		if tag == "running" {
			state = "up"
			break
		}
	}
	
	return ServiceStatus{
		State:     state,
		LastCheck: &now,
	}
}

func executeWidget(w config.Widget, cacheManager *cache.Manager, registry *widgets.Registry) WidgetResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := registry.Execute(ctx, w)
	if err != nil {
		return WidgetResult{
			Label: w.Label,
			Error: err.Error(),
			State: "error",
		}
	}

	var lastUpdate *string
	if !result.LastUpdate.IsZero() {
		s := result.LastUpdate.UTC().Format(time.RFC3339)
		lastUpdate = &s
	}

	return WidgetResult{
		Label:      result.Label,
		Value:      result.Value,
		Formatted:  result.Formatted,
		State:      result.State,
		LastUpdate: lastUpdate,
		Error:      result.Error,
	}
}
