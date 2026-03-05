# Installation

There are several ways to install and run **homectl**, depending on your setup.

## 🚀 Quick Install (Recommended)

The fastest way to get a production-ready instance of **homectl** running is using our official installation script.

```bash
curl -sSL https://homectl.xyz/install.sh | sh
```

This script will:
1.  Verify you have **Docker** and **Docker Compose** installed.
2.  Create a `homectl` directory.
3.  Download a production-ready `docker-compose.yml`.
4.  Pull the latest image and start the container.

---

## 🐋 Docker Compose (Manual)

If you prefer to manage your containers manually, create a `docker-compose.yml` file:

```yaml
services:
  homectl:
    image: ghcr.io/palta-dev/homectl:latest
    container_name: homectl
    network_mode: host # Recommended for Tailscale detection
    volumes:
      - ./data:/app/data
      - /var/run/docker.sock:/var/run/docker.sock:ro
    restart: unless-stopped
```

### Run
```bash
docker compose up -d
```

---

## 🛠 Advanced: Build from Source

If you want to contribute to development or modify the source code:

### 1. Clone the repository
```bash
git clone https://github.com/palta-dev/homectl.git
cd homectl
```

### 2. Build and run
```bash
docker compose up --build
```

### 3. Build and Run Manually (Local Dev)

#### Prerequisites:
- **Go 1.24+**
- **Node.js 20+**

**Frontend:**
```bash
cd apps/web
npm install
npm run dev
```

**Backend:**
```bash
cd apps/server
go mod download
go run cmd/main.go
```
