# SSRF Protection

Server-Side Request Forgery (SSRF) is a vulnerability where an attacker uses your dashboard to make requests to internal services that should not be publicly accessible.

**homectl** includes a hardened network layer to mitigate this risk.

## How it works

When a widget or health check requests a URL, **homectl**:
1. Resolves the hostname to its IP address.
2. Checks the IP against a blocklist of private and sensitive ranges.
3. Only allows the connection if the IP is public or explicitly permitted.

## Default Blocked Ranges

- **Cloud Metadata:** `169.254.169.254` (Used by AWS, GCP, etc. to store secrets).
- **Loopback:** `127.0.0.0/8`, `::1` (The dashboard host itself).
- **Private Networks:** `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`.

## Overriding Restrictions

If you need to access internal services, add them to the `allowHosts` setting in `config.yaml`:

```yaml
settings:
  allowHosts:
    - "192.168.1.10"
    - "nas.local"
```
