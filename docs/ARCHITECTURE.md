# Architecture

## Overview

homectl is a self-hosted homepage/dashboard for homelab administrators. It provides a fast, secure, and configurable interface to organize and monitor services.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              homectl Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐     ┌─────────────────────────────────────────────────┐   │
│  │   Browser    │────▶│              Frontend (React + Vite)            │   │
│  │              │◀────│  - Service tiles, widgets, search, theme toggle │   │
│  └──────────────┘     └────────────────────┬────────────────────────────┘   │
│                                            │ HTTP/JSON                        │
│                                            ▼                                  │
│  ┌──────────────┐     ┌─────────────────────────────────────────────────┐   │
│  │  Config File │────▶│              Backend (Go + Fiber)               │   │
│  │  (YAML)      │     │  - Config validation & hot reload               │   │
│  └──────────────┘     │  - Service status checks (HTTP, TCP, ping)      │   │
│                       │  - Widget execution engine                      │   │
│                       │  - SSRF protection (allowlist + CIDR)           │   │
│                       │  - Rate limiting                                │   │
│                       └────────────────────┬────────────────────────────┘   │
│                                            │                                  │
│              ┌─────────────────────────────┼─────────────────────────────┐   │
│              │                             │                             │   │
│              ▼                             ▼                             ▼   │
│     ┌─────────────────┐         ┌─────────────────┐         ┌──────────────┐│
│     │  SQLite (opt.)  │         │  In-Memory Cache│         │  External    ││
│     │  - Incidents    │         │  - LRU + TTL    │         │  Services    ││
│     │  - Preferences  │         │  - Per-widget   │         │  (LAN/Cloud) ││
│     │  - Sessions     │         │  - Backoff      │         │              ││
│     └─────────────────┘         └─────────────────┘         └──────────────┘│
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology | Rationale |
|-------|------------|-----------|
| **Frontend** | React 18 + Vite + TypeScript | Fast HMR, small bundle, excellent DX |
| **UI** | Tailwind CSS + shadcn/ui | Zero-runtime CSS, accessible, customizable |
| **Backend** | Go 1.21 + Fiber | Minimal memory, fast startup, strong typing |
| **Config** | YAML (gopkg.in/yaml.v3) | Human-readable, widely adopted |
| **Validation** | JSON Schema (sanidhana/go-jsonschema) | Clear error messages, versioning |
| **Cache** | In-memory LRU (hashicorp/golang-lru) | Sub-millisecond access, TTL support |
| **Storage** | SQLite (modernc.org/sqlite) | Pure Go, no CGO, optional |
| **Container** | Alpine Linux (multi-stage) | ~15MB final image |

## Directory Structure

```
homectl/
├── apps/
│   ├── web/                    # React frontend
│   │   ├── src/
│   │   │   ├── components/     # UI components
│   │   │   ├── hooks/          # React hooks
│   │   │   ├── lib/            # Utilities
│   │   │   ├── stores/         # State management
│   │   │   └── types/          # TypeScript types
│   │   ├── public/             # Static assets
│   │   ├── index.html
│   │   ├── vite.config.ts
│   │   └── package.json
│   │
│   └── server/                 # Go backend
│       ├── cmd/
│       │   └── main.go         # Entry point
│       ├── internal/
│       │   ├── config/         # Config loading & validation
│       │   ├── cache/          # Cache layer
│       │   ├── widgets/        # Widget implementations
│       │   ├── network/        # SSRF-safe HTTP client
│       │   ├── handlers/       # HTTP handlers
│       │   └── models/         # Data models
│       ├── scripts/            # Import/migration scripts
│       ├── go.mod
│       └── Dockerfile
│
├── packages/
│   └── shared/                 # Shared types & schemas
│       ├── schema/             # JSON Schema definitions
│       └── types/              # TypeScript types
│
├── data/
│   ├── icons/                  # Local icon storage
│   └── db/                     # SQLite database
│
├── docs/                       # Documentation
├── docker-compose.yml
├── config.yaml.example
└── README.md
```

## Data Flow

### 1. Configuration Loading

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│ config.yaml │───▶│   Parser     │───▶│  Validator  │───▶│  In-Memory   │
│             │    │  (yaml.v3)   │    │ (JSON Schema)│   │   Config     │
└─────────────┘    └──────────────┘    └─────────────┘    └──────────────┘
                          │                                      │
                          ▼                                      ▼
                   ┌──────────────┐                     ┌──────────────┐
                   │   Migration  │                     │  File Watch  │
                   │  (v1 → v2)   │                     │  (fsnotify)  │
                   └──────────────┘                     └──────────────┘
```

### 2. Service Status Check Flow

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│   Request   │───▶│  Cache Check │───▶│   Execute   │───▶│  Update Cache│
│  /api/services│   │  (LRU + TTL) │    │   Checker   │    │  + Response  │
└─────────────┘    └──────────────┘    └─────────────┘    └──────────────┘
                          │                  │
                    [HIT] │            [MISS]│
                          ▼                  ▼
                   ┌──────────────┐   ┌──────────────┐
                   │ Return Cached│   │ SSRF Check   │
                   │   Response   │   │ (Allowlist)  │
                   └──────────────┘   └──────────────┘
```

### 3. Widget Execution

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│   Widget    │───▶│  Scheduler   │───▶│  Executor   │───▶│   Result     │
│  Definition │    │ (per-interval)│   │ (HTTP/TCP)  │    │   Cache      │
└─────────────┘    └──────────────┘    └─────────────┘    └──────────────┘
                                                                  │
                                                                  ▼
                                                           ┌──────────────┐
                                                           │   Frontend   │
                                                           │   (Polling)  │
                                                           └──────────────┘
```

## Caching Strategy

### In-Memory Cache (LRU + TTL)

| Cache Type | TTL Default | Max Entries | Eviction |
|------------|-------------|-------------|----------|
| Service Status | 30s | 500 | LRU |
| Widget Results | Per-widget (default 60s) | 1000 | LRU + TTL |
| Config | Until file change | 1 | N/A |
| Icons | 24h | 200 | LRU |

### Cache Keys

```
service:{serviceId}:status          # Service health status
widget:{serviceId}:{widgetId}       # Widget result
config:current                      # Current configuration
icon:{iconName}                     # Resolved icon
```

### Retry & Backoff

```go
// Exponential backoff for failing checks
// Initial: 1s, Max: 60s, Multiplier: 2
// Reset on success
```

## Security Model

### SSRF Protection

```
┌─────────────────────────────────────────────────────────────────┐
│                    SSRF Protection Layers                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Allowlist Check                                              │
│     - Only hosts in config.settings.allowHosts permitted        │
│     - CIDR notation supported (192.168.0.0/16)                  │
│                                                                  │
│  2. Blocklist (Always Denied)                                    │
│     - Cloud metadata: 169.254.169.254/32                        │
│     - Link-local: 169.254.0.0/16                                │
│     - Loopback: 127.0.0.0/8 (unless explicitly allowed)         │
│     - Private ranges (if blockPrivateMetaIPs: true)             │
│                                                                  │
│  3. DNS Rebinding Protection                                     │
│     - Resolve hostname, verify IP still allowed                 │
│     - Block if resolution changes to blocked IP                 │
│                                                                  │
│  4. Connection Enforcement                                       │
│     - Custom Dialer with IP verification                        │
│     - Timeout enforcement (default 10s)                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Authentication (Optional)

```yaml
# Disabled by default
settings:
  auth:
    enabled: false
    provider: "local"  # local, github, google
    session:
      maxAge: 24h
```

### Secrets Handling

- **Environment variables**: `HOMECTL_SECRET_KEY`, `HOMECTL_OAUTH_CLIENT_SECRET`
- **Mounted files**: `/run/secrets/oauth_client_secret`
- **Never** stored in config file or returned by API

## Plugin/Widget System

### Widget Interface (Go)

```go
type Widget interface {
    Type() string
    Execute(ctx context.Context, cfg WidgetConfig) (*WidgetResult, error)
    CacheTTL() time.Duration
}
```

### Built-in Widgets

| Widget | Type | Description |
|--------|------|-------------|
| `httpStatus` | HTTP | Check HTTP status code and latency |
| `tcpPort` | TCP | Verify TCP port is open |
| `httpJson` | HTTP | Fetch JSON, extract field via JSONPath |
| `httpHtml` | HTTP | Scrape HTML, extract via CSS selector |
| `ping` | ICMP | Ping host (requires capabilities) |

### Widget Registration

```go
// internal/widgets/registry.go
var registry = map[string]Widget{
    "httpStatus": &HTTPStatusWidget{},
    "tcpPort":    &TCPPortWidget{},
    "httpJson":   &HTTPJSONWidget{},
}
```

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/health` | No | Health check |
| GET | `/api/config` | No* | Sanitized config (no secrets) |
| GET | `/api/services` | No* | Services with status |
| GET | `/api/widgets/:id` | No* | Specific widget result |
| POST | `/api/auth/login` | N/A | Login (if auth enabled) |
| POST | `/api/auth/logout` | Yes | Logout |

*Auth required if `settings.auth.enabled: true`

## Versioning & Migration

### Config Versioning

```yaml
# Current version
version: 1

# Future versions will have migration paths
# v1 → v2: automatic, additive changes only
```

### Migration Strategy

1. **Additive changes**: New fields are optional, backward compatible
2. **Breaking changes**: Require version bump, provide migration script
3. **Auto-migration**: Server migrates config in-memory, logs warnings

### Deprecation Policy

- Deprecated fields: Supported for 2 minor versions
- Removal: Announced in release notes, migration guide provided

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Initial render (LAN) | < 1s | Time to interactive |
| API response (cached) | < 50ms | p95 latency |
| API response (fresh) | < 500ms | p95 latency |
| Memory usage | < 50MB | RSS at steady state |
| Binary size | < 20MB | Compressed |
| Container size | < 25MB | Total image |

## Failure Modes

| Failure | Behavior | Recovery |
|---------|----------|----------|
| Config invalid | Use last valid config, log error | Fix config, hot reload |
| Service unreachable | Show "down" status, retry | Auto-retry with backoff |
| Cache corruption | Clear cache, recompute | Automatic |
| SQLite locked | Queue writes, retry | Automatic |
| Docker socket unavailable | Skip auto-discovery | Log warning |
