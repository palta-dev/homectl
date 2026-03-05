package widgets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
	"github.com/tidwall/gjson"
)

// HTTPJSONWidget fetches JSON and extracts a field
type HTTPJSONWidget struct{}

func (w *HTTPJSONWidget) Type() string {
	return "httpJson"
}

func (w *HTTPJSONWidget) CacheTTL() time.Duration {
	return 60 * time.Second
}

func (w *HTTPJSONWidget) Execute(ctx context.Context, cfg config.Widget, client *network.Client) (*Result, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", cfg.URL, nil)
	if err != nil {
		return &Result{Error: err.Error(), State: "error"}, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return &Result{Error: err.Error(), State: "error"}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Result{Error: err.Error(), State: "error"}, nil
	}

	// Extract JSON field
	result := gjson.GetBytes(body, cfg.JSONPath)
	if !result.Exists() {
		return &Result{Error: "JSONPath not found: " + cfg.JSONPath, State: "error"}, nil
	}

	value := result.Value()
	formatted := formatValue(value, cfg.Format)
	state := determineState(value, cfg.Format)

	return &Result{
		Label:      cfg.Label,
		Value:      value,
		Formatted:  formatted,
		State:      state,
		LastUpdate: time.Now(),
	}, nil
}

func formatValue(value interface{}, format string) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		switch format {
		case "bytes":
			return formatBytes(int64(v))
		case "duration":
			return formatDuration(int(v))
		case "percent":
			return fmt.Sprintf("%.1f%%", v)
		default:
			if v == float64(int64(v)) {
				return strconv.FormatInt(int64(v), 10)
			}
			return fmt.Sprintf("%.2f", v)
		}
	default:
		data, _ := json.Marshal(v)
		return string(data)
	}
}

func determineState(value interface{}, format string) string {
	switch v := value.(type) {
	case string:
		lower := strings.ToLower(v)
		if lower == "ok" || lower == "healthy" || lower == "up" || lower == "true" {
			return "good"
		}
		if lower == "error" || lower == "down" || lower == "false" {
			return "error"
		}
		return "warning"
	case bool:
		if v {
			return "good"
		}
		return "error"
	case float64:
		if v == 0 {
			return "warning"
		}
		return "good"
	default:
		return "good"
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(ms int) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.1fs", float64(ms)/1000)
}
