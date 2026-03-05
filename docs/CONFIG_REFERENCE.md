# Configuration Reference

This document describes every field in the homectl configuration file (`config.yaml`).

## Quick Example

```yaml
version: 1
settings:
  title: "Home Lab"
  theme: "dark"
  allowHosts:
    - "192.168.0.0/16"
    - "10.0.0.0/8"
  blockPrivateMetaIPs: true

groups:
  - name: "Infrastructure"
    layout: "grid"
    services:
      - name: "Grafana"
        url: "http://grafana:3000"
        icon: "grafana"
        checks:
          - type: "http"
            url: "http://grafana:3000/api/health"
            expectStatus: 200
            intervalSeconds: 30
        widgets:
          - type: "httpJson"
            url: "http://grafana:3000/api/health"
            jsonPath: "$.database"
            label: "DB"

icons:
  sources:
    - type: "local"
      path: "/data/icons"
    - type: "simpleicons"
      cache: true
```

---

## Root Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `version` | integer | Yes | - | Config schema version. Current: `1` |
| `settings` | object | No | `{}` | Global settings |
| `groups` | array | Yes | - | Service groups |
| `icons` | object | No | `{}` | Icon configuration |

---

## Settings

```yaml
settings:
  title: "Home Lab"
  theme: "dark"
  allowHosts:
    - "192.168.0.0/16"
  blockPrivateMetaIPs: true
  auth:
    enabled: false
  cache:
    defaultTTL: 30
```

### Settings Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `title` | string | No | `"homectl"` | Dashboard title shown in browser tab |
| `theme` | string | No | `"dark"` | Default theme: `"dark"` or `"light"` |
| `allowHosts` | array | No | `[]` | Allowed hostnames/CIDRs for outbound requests |
| `blockPrivateMetaIPs` | boolean | No | `true` | Block cloud metadata IPs (169.254.169.254) |
| `auth` | object | No | `{}` | Authentication settings |
| `cache` | object | No | `{}` | Cache configuration |
| `requestTimeout` | string | No | `"10s"` | Default timeout for HTTP/TCP requests |

### Settings: allowHosts

List of allowed destinations for outbound requests. Supports:
- Hostnames: `"grafana"`, `"api.example.com"`
- CIDR ranges: `"192.168.0.0/16"`, `"10.0.0.0/8"`
- Single IPs: `"192.168.1.100"`

**Security Note**: If empty, only localhost is allowed (for testing).

### Settings: auth

```yaml
auth:
  enabled: false
  provider: "local"  # local, github, google
  session:
    maxAge: "24h"
  # For OAuth providers:
  github:
    clientId: "${GITHUB_CLIENT_ID}"
    clientSecret: "${GITHUB_CLIENT_SECRET}"
    allowedUsers:
      - "your-username"
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `enabled` | boolean | No | `false` | Enable authentication |
| `provider` | string | No | `"local"` | Auth provider: `local`, `github`, `google` |
| `session.maxAge` | string | No | `"24h"` | Session duration |
| `github.clientId` | string | Conditional | - | GitHub OAuth client ID |
| `github.clientSecret` | string | Conditional | - | GitHub OAuth client secret |
| `github.allowedUsers` | array | Conditional | - | Allowed GitHub usernames |

### Settings: cache

```yaml
cache:
  defaultTTL: 30
  maxEntries: 500
  widgetTTL:
    httpStatus: 30
    tcpPort: 60
    httpJson: 120
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `defaultTTL` | integer | No | `30` | Default TTL in seconds |
| `maxEntries` | integer | No | `500` | Maximum cache entries |
| `widgetTTL` | object | No | `{}` | Per-widget TTL overrides |

---

## Groups

```yaml
groups:
  - name: "Infrastructure"
    layout: "grid"
    collapsed: false
    services:
      - name: "Grafana"
        url: "http://grafana:3000"
        icon: "grafana"
```

### Group Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | - | Group display name |
| `layout` | string | No | `"grid"` | Layout: `"grid"`, `"list"`, `"compact"` |
| `collapsed` | boolean | No | `false` | Start collapsed |
| `services` | array | Yes | - | Services in this group |

### Layout Options

| Layout | Description |
|--------|-------------|
| `grid` | Card grid, 3-4 columns on desktop |
| `list` | Vertical list with details |
| `compact` | Minimal icons + names |

---

## Services

```yaml
services:
  - name: "Grafana"
    url: "http://grafana:3000"
    icon: "grafana"
    description: "Metrics and dashboards"
    tags:
      - "monitoring"
      - "critical"
    checks:
      - type: "http"
        url: "http://grafana:3000/api/health"
        expectStatus: 200
        intervalSeconds: 30
    widgets:
      - type: "httpJson"
        url: "http://grafana:3000/api/health"
        jsonPath: "$.database"
        label: "DB"
```

### Service Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | - | Service display name |
| `url` | string | Yes | - | Primary URL (clickable) |
| `icon` | string | No | `"globe"` | Icon name or URL |
| `description` | string | No | `""` | Short description |
| `tags` | array | No | `[]` | Tags for filtering/search |
| `checks` | array | No | `[]` | Health checks |
| `widgets` | array | No | `[]` | Widgets to display |
| `newTab` | boolean | No | `true` | Open URL in new tab |
| `pingEnabled` | boolean | No | `false` | Show ping latency |

### Icon Resolution

Icons are resolved in order:

1. **Named icons**: `"grafana"` → look up in icon sources
2. **URLs**: `"https://example.com/icon.svg"` → fetch directly
3. **Local paths**: `"/data/icons/custom.svg"` → read from filesystem
4. **Fallback**: Default globe icon

---

## Checks

Health checks run periodically to determine service status.

### Check Types

#### HTTP Check

```yaml
checks:
  - type: "http"
    url: "http://grafana:3000/api/health"
    method: "GET"
    expectStatus: 200
    expectBodyContains: "ok"
    headers:
      Authorization: "Bearer ${API_TOKEN}"
    timeout: "5s"
    intervalSeconds: 30
    retries: 2
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"http"` |
| `url` | string | Yes | - | URL to check |
| `method` | string | No | `"GET"` | HTTP method |
| `expectStatus` | integer | No | `200` | Expected status code |
| `expectBodyContains` | string | No | - | String that must be in body |
| `headers` | object | No | `{}` | Request headers |
| `timeout` | string | No | `"10s"` | Request timeout |
| `intervalSeconds` | integer | No | `60` | Check interval |
| `retries` | integer | No | `1` | Retries before marking down |

#### TCP Port Check

```yaml
checks:
  - type: "tcp"
    host: "postgres"
    port: 5432
    timeout: "3s"
    intervalSeconds: 30
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"tcp"` |
| `host` | string | Yes | - | Hostname or IP |
| `port` | integer | Yes | - | Port number |
| `timeout` | string | No | `"5s"` | Connection timeout |
| `intervalSeconds` | integer | No | `60` | Check interval |

#### Ping Check

```yaml
checks:
  - type: "ping"
    host: "192.168.1.1"
    count: 3
    intervalSeconds: 60
```

**Note**: Ping requires elevated capabilities. May not work in all container setups.

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"ping"` |
| `host` | string | Yes | - | Hostname or IP |
| `count` | integer | No | `3` | Number of pings |
| `intervalSeconds` | integer | No | `60` | Check interval |

---

## Widgets

Widgets display additional information on service cards.

### Widget Types

#### HTTP JSON Widget

Fetches JSON and extracts a field using JSONPath.

```yaml
widgets:
  - type: "httpJson"
    url: "http://grafana:3000/api/health"
    jsonPath: "$.database"
    label: "DB"
    format: "status"  # status, bytes, duration, raw
    cacheTTL: 120
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"httpJson"` |
| `url` | string | Yes | - | URL to fetch |
| `jsonPath` | string | Yes | - | JSONPath expression |
| `label` | string | No | - | Display label |
| `format` | string | No | `"raw"` | Value formatting |
| `cacheTTL` | integer | No | `60` | Cache TTL in seconds |

**Format Options**:
- `raw`: Display as-is
- `status`: Green/red based on value
- `bytes`: Format as KB/MB/GB
- `duration`: Format as ms/s/m
- `percent`: Add % suffix

#### HTTP HTML Widget

Scrapes HTML and extracts text using CSS selector.

```yaml
widgets:
  - type: "httpHtml"
    url: "http://status.example.com"
    selector: ".status-indicator"
    attribute: "data-status"  # optional, defaults to textContent
    label: "Status"
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"httpHtml"` |
| `url` | string | Yes | - | URL to scrape |
| `selector` | string | Yes | - | CSS selector |
| `attribute` | string | No | - | Attribute to extract |
| `label` | string | No | - | Display label |

#### TCP Port Widget

Shows TCP port status with optional response time.

```yaml
widgets:
  - type: "tcpPort"
    host: "postgres"
    port: 5432
    label: "Postgres"
    showLatency: true
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"tcpPort"` |
| `host` | string | Yes | - | Hostname or IP |
| `port` | integer | Yes | - | Port number |
| `label` | string | No | - | Display label |
| `showLatency` | boolean | No | `false` | Show connection time |

---

## Icons Configuration

```yaml
icons:
  sources:
    - type: "local"
      path: "/data/icons"
    - type: "simpleicons"
      cache: true
      cacheTTL: 86400
```

### Icon Sources

#### Local Source

```yaml
- type: "local"
  path: "/data/icons"
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"local"` |
| `path` | string | Yes | - | Directory path |

#### Simple Icons

[Simple Icons](https://simpleicons.org/) - 3000+ brand icons.

```yaml
- type: "simpleicons"
  cache: true
  cacheTTL: 86400
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"simpleicons"` |
| `cache` | boolean | No | `true` | Cache fetched icons |
| `cacheTTL` | integer | No | `86400` | Cache TTL in seconds |

#### Custom URL Source

```yaml
- type: "url"
  baseUrl: "https://icons.example.com"
  pathTemplate: "/{name}.svg"
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | - | `"url"` |
| `baseUrl` | string | Yes | - | Base URL |
| `pathTemplate` | string | No | `"/{name}.svg"` | URL template |

---

## Environment Variable Substitution

Use `${VAR_NAME}` syntax in config values:

```yaml
settings:
  auth:
    github:
      clientId: "${GITHUB_CLIENT_ID}"
      clientSecret: "${GITHUB_CLIENT_SECRET}"

services:
  - name: "API"
    url: "http://api:8080"
    checks:
      - type: "http"
        url: "http://api:8080/health"
        headers:
          Authorization: "Bearer ${API_TOKEN}"
```

**Resolution Order**:
1. Environment variables
2. Mounted secret files (`/run/secrets/secret_name`)
3. Literal string (if no match)

---

## Config Versioning

### Version 1 (Current)

This is the initial schema version.

### Future Versions

When breaking changes are needed:

```yaml
# Old config (v1)
version: 1
groups:
  - services: [...]

# New config (v2) - example
version: 2
groups:
  - items: [...]  # renamed field
```

Migration is automatic. The server will:
1. Load v1 config
2. Apply migration rules
3. Use v2 schema in memory
4. Log migration event

---

## Validation Errors

Example validation errors:

```
ERROR config validation failed:
  - groups[0].services[2].checks[0].url: required field missing
  - settings.allowHosts[1]: invalid CIDR notation
  - services[3].icon: icon "unknown-icon" not found in any source
```

Fix errors and the config will hot-reload (if enabled).

---

## Example: Complete Config

```yaml
version: 1

settings:
  title: "Home Lab Dashboard"
  theme: "dark"
  allowHosts:
    - "192.168.0.0/16"
    - "10.0.0.0/8"
    - "grafana"
    - "prometheus"
  blockPrivateMetaIPs: true
  requestTimeout: "10s"
  cache:
    defaultTTL: 30
    maxEntries: 500

groups:
  - name: "Monitoring"
    layout: "grid"
    services:
      - name: "Grafana"
        url: "http://grafana:3000"
        icon: "grafana"
        description: "Metrics and dashboards"
        tags:
          - "monitoring"
          - "critical"
        checks:
          - type: "http"
            url: "http://grafana:3000/api/health"
            expectStatus: 200
            intervalSeconds: 30
        widgets:
          - type: "httpJson"
            url: "http://grafana:3000/api/health"
            jsonPath: "$.database"
            label: "DB"
            format: "status"

      - name: "Prometheus"
        url: "http://prometheus:9090"
        icon: "prometheus"
        checks:
          - type: "http"
            url: "http://prometheus:9090/-/healthy"
            expectStatus: 200

  - name: "Services"
    layout: "grid"
    services:
      - name: "Portainer"
        url: "http://portainer:9000"
        icon: "docker"
        checks:
          - type: "tcp"
            host: "portainer"
            port: 9000
            intervalSeconds: 30

      - name: "Home Assistant"
        url: "http://hass:8123"
        icon: "home-assistant"
        checks:
          - type: "http"
            url: "http://hass:8123/api/config"
            headers:
              Authorization: "Bearer ${HASS_TOKEN}"

  - name: "Network"
    layout: "list"
    services:
      - name: "Router"
        url: "http://192.168.1.1"
        icon: "router"
        pingEnabled: true
        checks:
          - type: "ping"
            host: "192.168.1.1"
            count: 3

icons:
  sources:
    - type: "local"
      path: "/data/icons"
    - type: "simpleicons"
      cache: true
      cacheTTL: 86400
```
