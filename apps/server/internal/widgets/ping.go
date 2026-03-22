package widgets

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
)

// PingWidget checks host reachability via ICMP ping
type PingWidget struct{}

func (w *PingWidget) Type() string {
	return "ping"
}

func (w *PingWidget) CacheTTL() time.Duration {
	return 30 * time.Second
}

func (w *PingWidget) Execute(ctx context.Context, cfg config.Widget, client *network.Client) (*Result, error) {
	host := cfg.Host
	if host == "" {
		return &Result{Error: "host is required", State: "error"}, nil
	}

	// Verify host is allowed via SSRF protection
	if err := client.CheckHost(ctx, host); err != nil {
		return &Result{Error: fmt.Sprintf("SSRF check failed: %v", err), State: "error"}, nil
	}

	start := time.Now()
	reachable, err := ping(host)
	latency := time.Since(start)

	if err != nil {
		return &Result{
			Value:     "Down",
			State:     "error",
			Error:     err.Error(),
			LastUpdate: time.Now(),
		}, nil
	}

	state := "good"
	if !reachable {
		state = "error"
	} else if latency > 500*time.Millisecond {
		state = "warning"
	}

	value := "Down"
	if reachable {
		value = "Up"
	}

	return &Result{
		Label:      "Ping",
		Value:      value,
		Formatted:  fmt.Sprintf("%s (%dms)", value, latency.Milliseconds()),
		State:      state,
		LastUpdate: time.Now(),
	}, nil
}

func ping(host string) (bool, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", host)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", host)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, nil // Consider unreachable instead of erroring out
	}

	// Basic check for success in output
	outStr := strings.ToLower(string(output))
	if strings.Contains(outStr, "ttl=") || strings.Contains(outStr, "received = 1") {
		return true, nil
	}

	return false, nil
}
