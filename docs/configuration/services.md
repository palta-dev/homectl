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
