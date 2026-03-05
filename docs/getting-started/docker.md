# Docker Deployment

**homectl** is designed to run natively in Docker and can automatically discover other containers on your host.

## Volume Mounts

| Host Path | Container Path | Description |
|-----------|----------------|-------------|
| `./data` | `/app/data` | Persistent storage for `config.yaml` and icons. |
| `/var/run/docker.sock` | `/var/run/docker.sock` | **Optional.** Required for auto-discovery. |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TZ` | `UTC` | Timezone for system logs and widgets. |
| `HOMECTL_CONFIG` | `data/config.yaml` | Path to the configuration file inside the container. |

## Auto-Discovery

To enable auto-discovery, ensure you have mounted the Docker socket and set `docker.enabled: true` in your configuration.

```yaml
docker:
  enabled: true
  socket: /var/run/docker.sock
```
