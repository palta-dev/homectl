# Remote Access

By default, discovered services point to `localhost`. When you access your dashboard remotely (e.g., via a VPS IP or domain), clicking a link to `localhost` will fail because your browser tries to find the service on your laptop/phone, not your server.

**homectl** provides two main ways to solve this.

## 1. Tailscale (Recommended)

Tailscale allows you to access your server using a private IP address (e.g., `100.x.y.z`) that works from anywhere.

### Automatic Detection
The **homectl** backend automatically scans your server for Tailscale. If it detects a Tailscale IP, it will use that IP for all auto-discovered Docker containers instead of `localhost`.

### Setup
1. Install Tailscale on your server and your client device.
2. Log in to both.
3. Access your dashboard using the server's Tailscale IP. 

### Manual Override (Priority)
If automatic detection fails or you want to use a specific IP for all services:
1. Go to **Settings > Discovery**.
2. Enter your IP (e.g., `100.67.67.67`) into the **Host IP Override** field.
3. Click **Save & Rescan**.
4. Every discovered Docker container will now use this IP for its generated URL.

---

## 2. Reverse Proxy (Split-DNS)

If you have a public domain (e.g., `example.com`), you can assign subdomains to each service.

1. Set up **Nginx Proxy Manager** or **Traefik**.
2. Create a proxy host for each service (e.g., `plex.example.com` -> `localhost:32400`).
3. In your **homectl** settings, manually update the service URL to use the public domain instead of the discovered IP.

## 3. Manual Overrides

For any service, you can manually override the discovered URL in the **Settings** page. Simply click **Configure** on a service card and enter the URL you wish to use globally.
