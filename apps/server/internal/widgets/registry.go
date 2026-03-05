package widgets

import (
	"context"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
)

// Result represents a widget execution result
type Result struct {
	Label      string      `json:"label,omitempty"`
	Value      interface{} `json:"value"`
	Formatted  string      `json:"formatted,omitempty"`
	State      string      `json:"state,omitempty"` // good, warning, error
	LastUpdate time.Time   `json:"lastUpdated,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// Widget defines the interface for all widgets
type Widget interface {
	Type() string
	Execute(ctx context.Context, cfg config.Widget, client *network.Client) (*Result, error)
	CacheTTL() time.Duration
}

// Registry holds all registered widgets
type Registry struct {
	widgets map[string]Widget
	client  *network.Client
}

// NewRegistry creates a new widget registry
func NewRegistry(client *network.Client) *Registry {
	return &Registry{
		widgets: make(map[string]Widget),
		client:  client,
	}
}

// Register adds a widget to the registry
func (r *Registry) Register(w Widget) {
	r.widgets[w.Type()] = w
}

// Get retrieves a widget by type
func (r *Registry) Get(widgetType string) (Widget, bool) {
	w, ok := r.widgets[widgetType]
	return w, ok
}

// Execute runs a widget and returns the result
func (r *Registry) Execute(ctx context.Context, widgetCfg config.Widget) (*Result, error) {
	w, ok := r.Get(widgetCfg.Type)
	if !ok {
		return &Result{
			Error: "unknown widget type: " + widgetCfg.Type,
			State: "error",
		}, nil
	}
	return w.Execute(ctx, widgetCfg, r.client)
}

// RegisterBuiltins registers all built-in widgets
func RegisterBuiltins(r *Registry) {
	r.Register(&HTTPStatusWidget{})
	r.Register(&HTTPJSONWidget{})
	r.Register(&HTTPHTMLWidget{})
	r.Register(&TCPPortWidget{})
	r.Register(&SystemWidget{})
}
