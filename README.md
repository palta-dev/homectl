# homectl

> A minimalist, high-performance dashboard for homelab administrators. Fast, sharp, and monochrome.

![License](https://img.shields.io/badge/license-MIT-black.svg)
![Go](https://img.shields.io/badge/go-1.24-black?logo=go)
![React](https://img.shields.io/badge/react-18-black?logo=react)
![Docker](https://img.shields.io/badge/docker-ready-black?logo=docker)

## What is homectl?

homectl is a technical, distraction-free dashboard designed for power users. It features a sharp, monochrome aesthetic, native Docker auto-discovery, and real-time system monitoring.

## 🚀 Advanced Setup

### Reverse Proxy (Recommended)

To access your dashboard securely over the internet, we recommend using a reverse proxy like **Nginx Proxy Manager** or **Traefik** with Let's Encrypt SSL.

1. Point your domain (e.g., `dash.example.com`) to your VPS.
2. Configure the proxy to point to `http://<vps-ip>:8080`.
3. Enable "Websockets Support" and "Force SSL".

### Docker Auto-Discovery

homectl can automatically find containers. Label your other containers to customize how they appear:

```yaml
services:
  my-app:
    image: my-app:latest
    labels:
      - "homectl.name=My Application"
      - "homectl.url=https://app.example.com"
      - "homectl.icon=rocket"
```

## 🛠 Tech Stack

- **Backend**: Go 1.24 (Fiber) - *Optimized for low latency and SSRF safety.*
- **Frontend**: React 18 (Vite, TypeScript) - *Fast, reactive, and type-safe.*
- **Security**: Bcrypt hashing, session-based auth, and hardened HTTP clients.

## 📖 Documentation

For full configuration options, widget details, and security guides, visit our documentation site:

👉 **[https://palta-dev.github.io/homectl](https://palta-dev.github.io/homectl)**

## Quick Start

### Docker Compose (Recommended)

```yaml
services:
  homectl:
    image: ghcr.io/palta-dev/homectl:latest
    container_name: homectl
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./data:/data
      - /var/run/docker.sock:/var/run/docker.sock:ro # Required for auto-discovery
    environment:
      - TZ=UTC
    restart: unless-stopped
```

### Installation from Source

```bash
# Clone the repository
git clone https://github.com/palta-dev/homectl.git
cd homectl

# Build and run with Docker
docker compose up --build
```

## Configuration

### System Monitoring Widgets

Add real-time stats to your header by including them in your `config.yaml`:

```yaml
groups:
  - name: System Health
    services:
      - name: Local Node
        widgets:
          - type: system
            label: CPU
            options: { metric: cpu }
          - type: system
            label: RAM
            options: { metric: mem }
          - type: system
            label: Temp
            options: { metric: temp }
```

### Background Customization

Set a modern background directly from the UI or config:

```yaml
settings:
  background: "https://images.unsplash.com/photo-1506744038136-46273834b3fb"
  backgroundOpacity: 0.6 # Adjust tint level (0.0 to 1.0)
```

## API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/config` | GET | Sanitized config |
| `/api/services` | GET | Services with status & widgets |

## Development

- **Backend**: Go 1.24 (Fiber)
- **Frontend**: React 18 (Vite, TypeScript)
- **Aesthetic**: Sharp edges (0px radius), monochrome, technical.

## License

MIT
