package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the root configuration structure
type Config struct {
	Version int        `yaml:"version" json:"version"`
	Settings Settings   `yaml:"settings,omitempty" json:"settings,omitempty"`
	Groups   []Group    `yaml:"groups" json:"groups"`
	Icons    *IconsConfig `yaml:"icons,omitempty" json:"icons,omitempty"`
}

// Settings contains global settings
type Settings struct {
	Title               string            `yaml:"title,omitempty" json:"title,omitempty"`
	Theme               string            `yaml:"theme,omitempty" json:"theme,omitempty"`
	Background          string            `yaml:"background,omitempty" json:"background,omitempty"`
	BackgroundOpacity   *float64          `yaml:"backgroundOpacity,omitempty" json:"backgroundOpacity,omitempty"`
	AllowHosts          []string          `yaml:"allowHosts,omitempty" json:"allowHosts,omitempty"`
	BlockPrivateMetaIPs bool              `yaml:"blockPrivateMetaIPs,omitempty" json:"blockPrivateMetaIPs,omitempty"`
	RequestTimeout      string            `yaml:"requestTimeout,omitempty" json:"requestTimeout,omitempty"`
	Cache               *CacheConfig      `yaml:"cache,omitempty" json:"cache,omitempty"`
	Auth                *AuthConfig       `yaml:"auth,omitempty" json:"auth,omitempty"`
	Docker              *DockerConfig     `yaml:"docker,omitempty" json:"docker,omitempty"`
	Password            string            `yaml:"password,omitempty" json:"-"` // Hidden from JSON
}

// CacheConfig contains cache settings
type CacheConfig struct {
	DefaultTTL int            `yaml:"defaultTTL,omitempty" json:"defaultTTL,omitempty"`
	MaxEntries int            `yaml:"maxEntries,omitempty" json:"maxEntries,omitempty"`
	WidgetTTL  map[string]int `yaml:"widgetTTL,omitempty" json:"widgetTTL,omitempty"`
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	Enabled  bool          `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Provider string        `yaml:"provider,omitempty" json:"provider,omitempty"`
	Session  *SessionConfig `yaml:"session,omitempty" json:"session,omitempty"`
	GitHub   *GitHubConfig  `yaml:"github,omitempty" json:"github,omitempty"`
}

// SessionConfig contains session settings
type SessionConfig struct {
	MaxAge string `yaml:"maxAge,omitempty" json:"maxAge,omitempty"`
}

// GitHubConfig contains GitHub OAuth settings
type GitHubConfig struct {
	ClientID     string   `yaml:"clientId,omitempty" json:"clientId,omitempty"`
	ClientSecret string   `yaml:"clientSecret,omitempty" json:"clientSecret,omitempty"`
	AllowedUsers []string `yaml:"allowedUsers,omitempty" json:"allowedUsers,omitempty"`
}

// DockerConfig contains Docker auto-discovery settings
type DockerConfig struct {
	Enabled     bool         `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Socket      string       `yaml:"socket,omitempty" json:"socket,omitempty"`
	LabelPrefix string       `yaml:"labelPrefix,omitempty" json:"labelPrefix,omitempty"`
	Hosts       []HostConfig `yaml:"hosts,omitempty" json:"hosts,omitempty"`
	Subnet      string       `yaml:"subnet,omitempty" json:"subnet,omitempty"`
	Ignore      []string     `yaml:"ignore,omitempty" json:"ignore,omitempty"`
}

// HostConfig represents a host to scan for services
type HostConfig struct {
	Name    string   `yaml:"name,omitempty" json:"name,omitempty"`
	Address string   `yaml:"address" json:"address"`
	Ports   []int    `yaml:"ports,omitempty" json:"ports,omitempty"`
	Tags    []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// Group represents a service group
type Group struct {
	Name      string    `yaml:"name" json:"name"`
	Layout    string    `yaml:"layout,omitempty" json:"layout,omitempty"`
	Collapsed bool      `yaml:"collapsed,omitempty" json:"collapsed,omitempty"`
	Services  []Service `yaml:"services" json:"services"`
}

// Service represents a single service
type Service struct {
	Name        string   `yaml:"name" json:"name"`
	URL         string   `yaml:"url" json:"url"`
	Icon        string   `yaml:"icon,omitempty" json:"icon,omitempty"`
	Favicon     string   `yaml:"favicon,omitempty" json:"favicon,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	NewTab      bool     `yaml:"newTab,omitempty" json:"newTab,omitempty"`
	PingEnabled bool     `yaml:"pingEnabled,omitempty" json:"pingEnabled,omitempty"`
	Checks      []Check  `yaml:"checks,omitempty" json:"checks,omitempty"`
	Widgets     []Widget `yaml:"widgets,omitempty" json:"widgets,omitempty"`
}

// Check represents a health check
type Check struct {
	Type              string            `yaml:"type" json:"type"`
	URL               string            `yaml:"url,omitempty" json:"url,omitempty"`
	Host              string            `yaml:"host,omitempty" json:"host,omitempty"`
	Port              int               `yaml:"port,omitempty" json:"port,omitempty"`
	Method            string            `yaml:"method,omitempty" json:"method,omitempty"`
	ExpectStatus      int               `yaml:"expectStatus,omitempty" json:"expectStatus,omitempty"`
	ExpectBodyContains string           `yaml:"expectBodyContains,omitempty" json:"expectBodyContains,omitempty"`
	Headers           map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Timeout           string            `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	IntervalSeconds   int               `yaml:"intervalSeconds,omitempty" json:"intervalSeconds,omitempty"`
	Retries           int               `yaml:"retries,omitempty" json:"retries,omitempty"`
	Count             int               `yaml:"count,omitempty" json:"count,omitempty"`
}

// Widget represents a widget configuration
type Widget struct {
	Type      string            `yaml:"type" json:"type"`
	URL       string            `yaml:"url,omitempty" json:"url,omitempty"`
	Host      string            `yaml:"host,omitempty" json:"host,omitempty"`
	Port      int               `yaml:"port,omitempty" json:"port,omitempty"`
	JSONPath  string            `yaml:"jsonPath,omitempty" json:"jsonPath,omitempty"`
	Selector  string            `yaml:"selector,omitempty" json:"selector,omitempty"`
	Attribute string            `yaml:"attribute,omitempty" json:"attribute,omitempty"`
	Label     string            `yaml:"label,omitempty" json:"label,omitempty"`
	Format    string            `yaml:"format,omitempty" json:"format,omitempty"`
	CacheTTL  int               `yaml:"cacheTTL,omitempty" json:"cacheTTL,omitempty"`
	Options   map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
}

// IconsConfig contains icon source configuration
type IconsConfig struct {
	Sources []IconSource `yaml:"sources" json:"sources"`
}

// IconSource represents an icon source
type IconSource struct {
	Type        string `yaml:"type" json:"type"`
	Path        string `yaml:"path,omitempty" json:"path,omitempty"`
	BaseURL     string `yaml:"baseUrl,omitempty" json:"baseUrl,omitempty"`
	PathTemplate string `yaml:"pathTemplate,omitempty" json:"pathTemplate,omitempty"`
	Cache       *bool  `yaml:"cache,omitempty" json:"cache,omitempty"`
	CacheTTL    *int   `yaml:"cacheTTL,omitempty" json:"cacheTTL,omitempty"`
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	// Expand environment variables
	data = expandEnv(data)

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	// Validate
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	// Apply defaults
	applyDefaults(&cfg)

	return &cfg, nil
}

// expandEnv replaces ${VAR} patterns with environment variable values
var envPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

func expandEnv(data []byte) []byte {
	return envPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		varName := string(match[2 : len(match)-1])
		if val := os.Getenv(varName); val != "" {
			return []byte(val)
		}
		// Try reading from secret file
		secretPath := fmt.Sprintf("/run/secrets/%s", varName)
		if data, err := os.ReadFile(secretPath); err == nil {
			return []byte(strings.TrimSpace(string(data)))
		}
		return match // Keep original if not found
	})
}

// validate checks configuration validity
func validate(cfg *Config) error {
	if cfg.Version < 1 {
		return fmt.Errorf("invalid version: %d", cfg.Version)
	}
	if len(cfg.Groups) == 0 {
		return fmt.Errorf("at least one group is required")
	}
	for i, g := range cfg.Groups {
		if g.Name == "" {
			return fmt.Errorf("group[%d]: name is required", i)
		}
		if len(g.Services) == 0 {
			return fmt.Errorf("group[%d]: at least one service is required", i)
		}
		for j, s := range g.Services {
			if s.Name == "" {
				return fmt.Errorf("group[%d].service[%d]: name is required", i, j)
			}
			if s.URL == "" {
				return fmt.Errorf("group[%d].service[%d]: url is required", i, j)
			}
		}
	}
	return nil
}

// applyDefaults sets default values for optional fields
func applyDefaults(cfg *Config) {
	if cfg.Settings.Theme == "" {
		cfg.Settings.Theme = "dark"
	}
	if cfg.Settings.Title == "" {
		cfg.Settings.Title = "homectl"
	}
	if cfg.Settings.RequestTimeout == "" {
		cfg.Settings.RequestTimeout = "10s"
	}
	if cfg.Settings.Cache == nil {
		cfg.Settings.Cache = &CacheConfig{}
	}
	if cfg.Settings.Cache.DefaultTTL == 0 {
		cfg.Settings.Cache.DefaultTTL = 30
	}
	if cfg.Settings.Cache.MaxEntries == 0 {
		cfg.Settings.Cache.MaxEntries = 500
	}
	for i := range cfg.Groups {
		if cfg.Groups[i].Layout == "" {
			cfg.Groups[i].Layout = "grid"
		}
		for j := range cfg.Groups[i].Services {
			if cfg.Groups[i].Services[j].NewTab == false {
				cfg.Groups[i].Services[j].NewTab = true
			}
		}
	}
}

// GetTimeout parses the request timeout string
func (s *Settings) GetTimeout() time.Duration {
	if s.RequestTimeout == "" {
		return 10 * time.Second
	}
	d, err := time.ParseDuration(s.RequestTimeout)
	if err != nil {
		return 10 * time.Second
	}
	return d
}

const defaultYaml = `version: 1
settings:
  title: "homectl"
  theme: "dark"
  requestTimeout: "10s"
groups:
  - name: "Getting Started"
    services:
      - name: "homectl GitHub"
        url: "https://github.com/palta-dev/homectl"
        description: "Documentation and source code"
`

// EnsureExists checks if config exists, creates default if not
func EnsureExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create parent directories if they don't exist
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}
		return os.WriteFile(path, []byte(defaultYaml), 0644)
	}
	return nil
}

	// Save writes the configuration to the specified path in YAML format
func (cfg *Config) Save(path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshalling YAML: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
