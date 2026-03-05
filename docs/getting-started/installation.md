# Installation

There are two primary ways to run **homectl**: using Docker (recommended) or building from source.

## Docker (Recommended)

The easiest way to get started is using Docker Compose.

1. Create a directory for homectl:
   ```bash
   mkdir homectl && cd homectl
   ```

2. Create a `docker-compose.yml` file:
   ```yaml
   services:
     homectl:
       image: ghcr.io/palta-dev/homectl:latest
       container_name: homectl
       ports:
         - "8080:8080"
       volumes:
         - ./data:/app/data
         - /var/run/docker.sock:/var/run/docker.sock:ro
       restart: unless-stopped
   ```

3. Start the dashboard:
   ```bash
   docker compose up -d
   ```

## Building from Source

If you prefer to build the binary yourself, you will need **Go 1.24+** and **Node.js 20+**.

1. Clone the repository:
   ```bash
   git clone https://github.com/palta-dev/homectl.git
   cd homectl
   ```

2. Build and run:
   ```bash
   docker compose up --build
   ```
