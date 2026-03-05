# Widget System

Widgets provide real-time data from your services or the host system.

## System Widgets

Displays host system metrics in the top bar.

```yaml
widgets:
  - type: system
    label: CPU
    options:
      metric: cpu
  - type: system
    label: RAM
    options:
      metric: mem
```

### Supported Metrics
- `cpu`: CPU utilization percentage.
- `mem`: RAM utilization percentage.
- `disk`: Disk usage (default path `/`).
- `temp`: CPU temperature.
- `uptime`: System uptime.

## HTTP JSON Widget

Fetches a JSON endpoint and extracts a value using GJSON syntax.

```yaml
widgets:
  - type: httpJson
    url: "https://api.example.com/stats"
    jsonPath: "data.active_users"
    label: "Users"
```

## HTTP Status Widget

Simply monitors the status code and latency of a URL.

```yaml
widgets:
  - type: httpStatus
    url: "https://google.com"
    label: "Google"
```
