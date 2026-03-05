package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad_ValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `
version: 1
settings:
  title: "Test"
  theme: "dark"
groups:
  - name: "Test Group"
    services:
      - name: "Test Service"
        url: "http://localhost:8080"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("Version = %d, want 1", cfg.Version)
	}
	if cfg.Settings.Title != "Test" {
		t.Errorf("Title = %q, want %q", cfg.Settings.Title, "Test")
	}
	if len(cfg.Groups) != 1 {
		t.Errorf("Groups = %d, want 1", len(cfg.Groups))
	}
}

func TestLoad_InvalidVersion(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `
version: 0
groups:
  - name: "Test"
    services:
      - name: "Test"
        url: "http://localhost"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("Load() expected error for invalid version")
	}
}

func TestLoad_MissingGroups(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `
version: 1
settings:
  title: "Test"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("Load() expected error for missing groups")
	}
}

func TestLoad_EnvExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Set test env var
	t.Setenv("TEST_TITLE", "Expanded Title")

	content := `
version: 1
settings:
  title: "${TEST_TITLE}"
groups:
  - name: "Test"
    services:
      - name: "Test"
        url: "http://localhost"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Settings.Title != "Expanded Title" {
		t.Errorf("Title = %q, want %q", cfg.Settings.Title, "Expanded Title")
	}
}

func TestLoad_Defaults(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `
version: 1
groups:
  - name: "Test"
    services:
      - name: "Test"
        url: "http://localhost"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Settings.Theme != "dark" {
		t.Errorf("Theme = %q, want %q", cfg.Settings.Theme, "dark")
	}
	if cfg.Settings.Title != "homectl" {
		t.Errorf("Title = %q, want %q", cfg.Settings.Title, "homectl")
	}
	if cfg.Settings.Cache.DefaultTTL != 30 {
		t.Errorf("Cache.DefaultTTL = %d, want 30", cfg.Settings.Cache.DefaultTTL)
	}
}

func TestValidate_ServiceMissingURL(t *testing.T) {
	cfg := &Config{
		Version: 1,
		Groups: []Group{
			{
				Name: "Test",
				Services: []Service{
					{Name: "Test", URL: ""},
				},
			},
		},
	}

	err := validate(cfg)
	if err == nil {
		t.Error("validate() expected error for missing URL")
	}
}

func TestGetTimeout(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"10s", 10 * time.Second},
		{"1m", 1 * time.Minute},
		{"", 10 * time.Second},
		{"invalid", 10 * time.Second},
	}

	for _, tt := range tests {
		s := &Settings{RequestTimeout: tt.input}
		result := s.GetTimeout()
		if result != tt.expected {
			t.Errorf("GetTimeout(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
