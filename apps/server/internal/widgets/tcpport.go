package widgets

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
)

// TCPPortWidget checks TCP port connectivity
type TCPPortWidget struct{}

func (w *TCPPortWidget) Type() string {
	return "tcpPort"
}

func (w *TCPPortWidget) CacheTTL() time.Duration {
	return 30 * time.Second
}

func (w *TCPPortWidget) Execute(ctx context.Context, cfg config.Widget, client *network.Client) (*Result, error) {
	host := cfg.Host
	if host == "" {
		return &Result{Error: "host is required", State: "error"}, nil
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return &Result{Error: "invalid port", State: "error"}, nil
	}

	// Verify host is allowed before connecting
	if err := client.CheckHost(ctx, host); err != nil {
		return &Result{
			Label:      cfg.Label,
			Value:      false,
			Formatted:  "blocked",
			State:      "error",
			LastUpdate: time.Now(),
			Error:      "SSRF check failed: " + err.Error(),
		}, nil
	}

	start := time.Now()

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", host, cfg.Port)
	conn, err := dialer.DialContext(ctx, "tcp", address)
	latency := time.Since(start)

	if err != nil {
		return &Result{
			Label:      cfg.Label,
			Value:      false,
			Formatted:  "down",
			State:      "error",
			LastUpdate: time.Now(),
			Error:      err.Error(),
		}, nil
	}
	defer conn.Close()

	return &Result{
		Label:      cfg.Label,
		Value:      true,
		Formatted:  fmt.Sprintf("up (%dms)", latency.Milliseconds()),
		State:      "good",
		LastUpdate: time.Now(),
	}, nil
}
