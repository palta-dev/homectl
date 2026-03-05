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

Fetches a JSON endpoint and extracts a value using [GJSON syntax](https://github.com/tidwall/gjson).

```yaml
widgets:
  - type: httpJson
    url: "https://api.example.com/stats"
    jsonPath: "data.active_users"
    label: "Users"
    format: "raw"
```

### Formatting Options
- `raw`: Default. Displays the value as-is.
- `bytes`: Converts numbers to human-readable sizes (e.g., `1.2 GB`).
- `duration`: Converts milliseconds to time (e.g., `5.2s`).
- `percent`: Appends `%` to the number.

## HTTP HTML Widget

Scrapes text from an HTML page using CSS selectors.

```yaml
widgets:
  - type: httpHtml
    url: "https://status.myapp.com"
    selector: ".status-text"
    label: "Status"
```

### Scraper Options
- `selector`: A CSS selector (e.g., `#main > .status`).
- `attribute`: (Optional) Extract an attribute value instead of inner text (e.g., `title` or `src`).
