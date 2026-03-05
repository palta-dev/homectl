# Security Guide

This document outlines security considerations and best practices for running homectl.

## SSRF Protection

homectl includes built-in protection against Server-Side Request Forgery (SSRF) attacks.

### How It Works

1. **Allowlist-based**: Only hosts/CIDRs explicitly configured in `allowHosts` are accessible
2. **Blocklist enforcement**: Certain IP ranges are always blocked regardless of allowlist
3. **DNS rebinding protection**: Hostnames are resolved and verified before each request

### Always Blocked

These IP ranges are blocked by default:

| Range | Reason |
|-------|--------|
| `169.254.169.254/32` | Cloud metadata endpoint |
| `169.254.0.0/16` | Link-local addresses |
| `127.0.0.0/8` | Loopback |
| `::1/128` | IPv6 loopback |

### Private Networks

When `blockPrivateMetaIPs: true` (default), these are also blocked:

| Range | Reason |
|-------|--------|
| `10.0.0.0/8` | Private Class A |
| `172.16.0.0/12` | Private Class B |
| `192.168.0.0/16` | Private Class C |
| `fc00::/7` | IPv6 private |

### Configuration Example

```yaml
settings:
  # Only allow specific networks
  allowHosts:
    - "192.168.1.0/24"      # Your LAN
    - "grafana"             # Specific hostname
    - "10.0.0.50"           # Specific IP
  
  # Block private networks (recommended for cloud deployments)
  blockPrivateMetaIPs: true
```

### Testing SSRF Protection

```bash
# This should fail (blocked IP)
curl http://localhost:8080/api/widgets/test?url=http://169.254.169.254/latest/meta-data

# This should succeed (allowed host)
curl http://localhost:8080/api/widgets/test?url=http://grafana:3000/api/health
```

## Docker Socket Security

⚠️ **WARNING**: Mounting the Docker socket grants container root-equivalent access to the host.

### Risk Assessment

| Risk Level | Configuration |
|------------|---------------|
| 🔴 HIGH | Docker socket mounted + exposed to network |
| 🟡 MEDIUM | Docker socket mounted, network isolated |
| 🟢 LOW | Docker socket not mounted |

### Recommendations

1. **Do not enable Docker auto-discovery** unless absolutely necessary
2. **If enabled**, ensure homectl is not exposed to untrusted networks
3. **Consider alternatives**: Use static configuration or Docker labels only

### Safe Configuration

```yaml
# Recommended: Disable Docker auto-discovery
settings:
  docker:
    enabled: false
```

```yaml
# If you must enable it:
settings:
  docker:
    enabled: true
    socket: "/var/run/docker.sock"
    labelPrefix: "homectl"
```

```yaml
# docker-compose.yml - restrict network exposure
services:
  homectl:
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro  # Read-only
    networks:
      - internal  # Internal network only
    # Do NOT expose ports to external network
```

## Authentication

### Enabling Auth

```yaml
settings:
  auth:
    enabled: true
    provider: "github"
    github:
      clientId: "${GITHUB_CLIENT_ID}"
      clientSecret: "${GITHUB_CLIENT_SECRET}"
      allowedUsers:
        - "your-username"
```

### Best Practices

1. **Always use environment variables** for secrets
2. **Restrict OAuth to specific users** via `allowedUsers`
3. **Use strong session timeouts** (default: 24h)
4. **Enable HTTPS** when exposing to network

## Secrets Management

### Environment Variables

```bash
# .env file (do not commit to git!)
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_secret
HOMECTL_SECRET_KEY=generate_with_openssl
```

### Mounted Secrets (Docker Swarm/Kubernetes)

```yaml
# docker-compose.yml
services:
  homectl:
    secrets:
      - github_client_secret
      - github_client_id

secrets:
  github_client_secret:
    external: true
  github_client_id:
    external: true
```

```yaml
# config.yaml
settings:
  auth:
    github:
      clientId: "${GITHUB_CLIENT_ID}"
      clientSecret: "${GITHUB_CLIENT_SECRET}"
```

## Rate Limiting

homectl includes built-in rate limiting to prevent abuse:

| Endpoint | Limit | Burst |
|----------|-------|-------|
| `/api/*` | 10 req/s | 20 |
| Auth endpoints | 5 req/s | 10 |

### Adjusting Limits

Modify in `cmd/main.go`:

```go
api := app.Group("/api", middleware.DefaultRateLimiter())
// Change to:
api := app.Group("/api", middleware.NewRateLimiter(middleware.RateLimiterConfig{
    RequestsPerSecond: 20,
    BurstSize: 40,
}).Middleware())
```

## Network Exposure

### Recommended Deployment

```
┌─────────────────┐     ┌─────────────┐     ┌─────────────┐
│   Internet      │────▶│   Reverse   │────▶│   homectl   │
│                 │     │   Proxy     │     │  (internal) │
│                 │     │  (Traefik)  │     │             │
└─────────────────┘     └─────────────┘     └─────────────┘
                              │
                              ▼
                        ┌─────────────┐
                        │   SSL/TLS   │
                        │  (Let's     │
                        │   Encrypt)  │
                        └─────────────┘
```

### Firewall Rules

```bash
# Allow only from trusted networks
ufw allow from 192.168.1.0/24 to any port 8080
ufw deny 8080
```

## Security Checklist

- [ ] `allowHosts` configured with minimum required access
- [ ] `blockPrivateMetaIPs: true` (unless on trusted LAN)
- [ ] Docker socket NOT mounted (or isolated)
- [ ] Authentication enabled if exposed to network
- [ ] Secrets in environment variables, not config file
- [ ] HTTPS enabled via reverse proxy
- [ ] Firewall rules restrict access
- [ ] Regular security updates applied
- [ ] Rate limiting configured appropriately

## Reporting Vulnerabilities

Please report security vulnerabilities to: security@homectl.xyz

**Do not** create public GitHub issues for security vulnerabilities.

## Security Updates

Subscribe to security advisories:
- GitHub Security Advisories
- Docker Hub security scanning
- Dependabot alerts
