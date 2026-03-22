package widgets

import (
	"context"
	"fmt"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemWidget handles system monitoring metrics
type SystemWidget struct{}

func (w *SystemWidget) Type() string {
	return "system"
}

func (w *SystemWidget) CacheTTL() time.Duration {
	return 5 * time.Second
}

func (w *SystemWidget) Execute(ctx context.Context, cfg config.Widget, client *network.Client) (*Result, error) {
	metricType := cfg.Options["metric"]
	if metricType == "" {
		metricType = "cpu" // Default to CPU
	}

	result := &Result{
		Label:      cfg.Label,
		LastUpdate: time.Now(),
		State:      "good",
	}

	switch metricType {
	case "cpu":
		percent, err := cpu.PercentWithContext(ctx, 0, false)
		if err != nil {
			return nil, err
		}
		if len(percent) > 0 {
			result.Value = percent[0]
			result.Formatted = fmt.Sprintf("%.1f%%", percent[0])
			if percent[0] > 90 {
				result.State = "error"
			} else if percent[0] > 70 {
				result.State = "warning"
			}
		}
	case "mem":
		v, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			return nil, err
		}
		result.Value = v.UsedPercent
		result.Formatted = fmt.Sprintf("%.1f%%", v.UsedPercent)
		if v.UsedPercent > 90 {
			result.State = "error"
		} else if v.UsedPercent > 75 {
			result.State = "warning"
		}
	case "disk":
		path := cfg.Options["path"]
		if path == "" {
			path = "/"
		}
		usage, err := disk.UsageWithContext(ctx, path)
		if err != nil {
			return nil, err
		}
		result.Value = usage.UsedPercent
		result.Formatted = fmt.Sprintf("%.1f%%", usage.UsedPercent)
		if usage.UsedPercent > 95 {
			result.State = "error"
		} else if usage.UsedPercent > 85 {
			result.State = "warning"
		}
	case "load":
		l, err := load.AvgWithContext(ctx)
		if err != nil {
			return nil, err
		}
		result.Value = l.Load1
		result.Formatted = fmt.Sprintf("%.2f", l.Load1)
	case "temp":
		temps, err := host.SensorsTemperaturesWithContext(ctx)
		if err != nil {
			return nil, err
		}
		if len(temps) > 0 {
			// Find a suitable temperature sensor
			sensor := temps[0]
			for _, t := range temps {
				if t.SensorKey == "package id 0" || t.SensorKey == "coretemp_package_id_0" {
					sensor = t
					break
				}
			}
			result.Value = sensor.Temperature
			result.Formatted = fmt.Sprintf("%.1f°C", sensor.Temperature)
			if sensor.Temperature > 85 {
				result.State = "error"
			} else if sensor.Temperature > 70 {
				result.State = "warning"
			}
		} else {
			result.Value = 0
			result.Formatted = "N/A"
			result.State = "unknown"
		}
	case "uptime":
		u, err := host.UptimeWithContext(ctx)
		if err != nil {
			return nil, err
		}
		result.Value = u
		result.Formatted = formatUptime(u)
	default:
		return nil, fmt.Errorf("unsupported system metric: %s", metricType)
	}

	return result, nil
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
