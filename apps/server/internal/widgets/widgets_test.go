package widgets

import (
	"context"
	"testing"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
)

func TestHTTPJSONWidget_Execute(t *testing.T) {
	widget := &HTTPJSONWidget{}
	client, _ := network.NewClient(network.Config{Timeout: 1 * time.Second})

	// Test with mock server would go here
	// For now, test error handling
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cfg := config.Widget{
		Type:     "httpJson",
		URL:      "http://invalid-host-that-does-not-exist/api",
		JSONPath: "$.value",
		Label:    "Test",
	}

	result, err := widget.Execute(ctx, cfg, client)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.State != "error" {
		t.Errorf("State = %q, want %q", result.State, "error")
	}
	if result.Error == "" {
		t.Error("Expected error message")
	}
}

func TestHTTPJSONWidget_FormatValue(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		format string
		want   string
	}{
		{"string", "hello", "raw", "hello"},
		{"bool true", true, "raw", "true"},
		{"bool false", false, "raw", "false"},
		{"number int", 42.0, "raw", "42"},
		{"number float", 3.14, "raw", "3.14"},
		{"bytes", 1024.0, "bytes", "1.0 KB"},
		{"percent", 95.5, "percent", "95.5%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatValue(tt.value, tt.format)
			if got != tt.want {
				t.Errorf("formatValue(%v, %q) = %q, want %q", tt.value, tt.format, got, tt.want)
			}
		})
	}
}

func TestHTTPJSONWidget_DetermineState(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"string ok", "ok", "good"},
		{"string healthy", "healthy", "good"},
		{"string up", "up", "good"},
		{"string error", "error", "error"},
		{"string down", "down", "error"},
		{"bool true", true, "good"},
		{"bool false", false, "error"},
		{"number zero", 0.0, "warning"},
		{"number positive", 1.0, "good"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineState(tt.value, "")
			if got != tt.want {
				t.Errorf("determineState(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestTCPPortWidget_Execute_InvalidConfig(t *testing.T) {
	widget := &TCPPortWidget{}
	client, _ := network.NewClient(network.Config{Timeout: 1 * time.Second})
	ctx := context.Background()

	tests := []struct {
		name string
		cfg  config.Widget
	}{
		{"missing host", config.Widget{Type: "tcpPort", Port: 80}},
		{"invalid port zero", config.Widget{Type: "tcpPort", Host: "localhost", Port: 0}},
		{"invalid port negative", config.Widget{Type: "tcpPort", Host: "localhost", Port: -1}},
		{"invalid port too high", config.Widget{Type: "tcpPort", Host: "localhost", Port: 70000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := widget.Execute(ctx, tt.cfg, client)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if result.State != "error" {
				t.Errorf("State = %q, want %q", result.State, "error")
			}
		})
	}
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	client, _ := network.NewClient(network.Config{Timeout: 1 * time.Second})
	registry := NewRegistry(client)

	// Register built-ins
	RegisterBuiltins(registry)

	tests := []string{"httpJson", "httpHtml", "httpStatus", "tcpPort"}
	for _, widgetType := range tests {
		t.Run(widgetType, func(t *testing.T) {
			widget, ok := registry.Get(widgetType)
			if !ok {
				t.Fatalf("Get(%q) returned not found", widgetType)
			}
			if widget.Type() != widgetType {
				t.Errorf("Type() = %q, want %q", widget.Type(), widgetType)
			}
		})
	}
}

func TestRegistry_Execute_UnknownWidget(t *testing.T) {
	client, _ := network.NewClient(network.Config{Timeout: 1 * time.Second})
	registry := NewRegistry(client)
	ctx := context.Background()

	cfg := config.Widget{
		Type: "unknownWidget",
	}

	result, err := registry.Execute(ctx, cfg)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Error == "" {
		t.Error("Expected error for unknown widget type")
	}
}
