package widgets

import (
	"context"
	"net/http"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
)

// HTTPStatusWidget checks HTTP status and latency
type HTTPStatusWidget struct{}

func (w *HTTPStatusWidget) Type() string {
	return "httpStatus"
}

func (w *HTTPStatusWidget) CacheTTL() time.Duration {
	return 30 * time.Second
}

func (w *HTTPStatusWidget) Execute(ctx context.Context, cfg config.Widget, client *network.Client) (*Result, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", cfg.URL, nil)
	if err != nil {
		return &Result{Error: err.Error(), State: "error"}, nil
	}

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return &Result{
			Label:      cfg.Label,
			Value:      "error",
			Formatted:  "down",
			State:      "error",
			LastUpdate: time.Now(),
			Error:      err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	var state string
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		state = "good"
	} else if resp.StatusCode < 500 {
		state = "warning"
	} else {
		state = "error"
	}

	return &Result{
		Label:      cfg.Label,
		Value:      resp.StatusCode,
		Formatted:  formatLatencyWidget(latency.Milliseconds()),
		State:      state,
		LastUpdate: time.Now(),
	}, nil
}

func formatLatencyWidget(ms int64) string {
	if ms < 100 {
		return "fast"
	}
	if ms < 500 {
		return "normal"
	}
	return "slow"
}
