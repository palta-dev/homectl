# Authentication & Security

**homectl** includes professional-grade security features to ensure your dashboard can be safely exposed to the internet.

## Session-Based Authentication

Settings are protected by a password system. When a password is set in `config.yaml`, the Settings page is hidden behind a login screen.

### Setting a Password
1. Navigate to **Settings > Security**.
2. Enter your desired password and click **Update Password**.
3. The server will hash your password using **bcrypt** and save it to your config.

## SSRF Protection

To prevent the dashboard from being used to attack your internal network, all widgets and health checks use a hardened HTTP client.

### Protected Ranges
By default, the following ranges are blocked:
- `127.0.0.0/8` (Loopback)
- `169.254.169.254` (Cloud Metadata)
- `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16` (Private Networks - if configured)

### Allowlist
You can explicitly allow certain hosts in your `config.yaml`:

```yaml
settings:
  allowHosts:
    - "192.168.1.50"
    - "my-internal-service.local"
```
