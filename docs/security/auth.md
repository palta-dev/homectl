# Authentication & Security

**homectl** includes professional-grade security features to ensure your dashboard can be safely exposed to the internet.

## Session-Based Authentication

**homectl** features a dedicated login system to protect your configuration. When a password is set, the Settings page is inaccessible until you authenticate.

### How it works
1. **Initial Setup:** When you first deploy, the dashboard is open.
2. **Setting a Password:** Go to **Settings > Security**, enter a new password, and save.
3. **Protection:** The server hashes your password using **bcrypt** and saves it to `data/config.yaml`.
4. **Login:** Subsequent visits to the Settings page will redirect you to a secure login screen. Your session is maintained in your browser so you only have to log in once.

### Recovery
If you forget your password, you can manually clear the `password` field in `data/config.yaml` and restart the server to restore access.

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
