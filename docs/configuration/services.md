# Services & Groups

**homectl** organizes your dashboard into Groups, which contain individual Services.

```yaml
groups:
  - name: "Development"
    layout: "grid"
    services:
      - name: "Gitea"
        url: "http://gitea.local"
        icon: "git"
```

## Group Options

| Key | Options | Description |
|-----|---------|-------------|
| `name` | `string` | The header for the group. |
| `layout` | `grid`, `list` | Visual arrangement of cards. |
| `collapsed` | `boolean` | Whether the group starts minimized. |

## Service Options

| Key | Description |
|-----|-------------|
| `name` | The name displayed on the card. |
| `url` | The destination when clicked. |
| `icon` | Name of the icon to display. |
| `description` | Subtitle text for the service. |
| `newTab` | Open link in a new window. |
| `pingEnabled` | Enable background ICMP/TCP ping checks. |

## Real-World Examples

### 1. Media Server (Plex)
```yaml
services:
  - name: "Plex"
    url: "http://192.168.1.100:32400"
    icon: "plex"
    description: "Media Library"
    widgets:
      - type: httpStatus
        url: "http://192.168.1.100:32400/identity"
        label: "API"
```

### 2. Network (Pi-hole)
```yaml
services:
  - name: "Pi-hole"
    url: "http://pi.hole/admin"
    icon: "pi-hole"
    widgets:
      - type: httpJson
        url: "http://pi.hole/admin/api.php?summary"
        jsonPath: "ads_blocked_today"
        label: "Blocked"
```

### 3. Smart Home (Home Assistant)
```yaml
services:
  - name: "Home Assistant"
    url: "https://hass.example.com"
    icon: "home-assistant"
    pingEnabled: true
```
