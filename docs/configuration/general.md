# General Settings

The `settings` block in `data/config.yaml` controls the global behavior and appearance of your dashboard.

```yaml
settings:
  title: "My Lab"
  theme: "dark"
  background: "https://images.unsplash.com/photo-..."
  backgroundOpacity: 0.5
  requestTimeout: "10s"
```

## Available Options

| Key | Type | Description |
|-----|------|-------------|
| `title` | `string` | The text displayed in the top-left corner. |
| `theme` | `string` | `dark` or `light`. |
| `background` | `string` | URL to a background image. |
| `backgroundOpacity` | `float` | Transparency of the image tint (0.0 to 1.0). |
| `requestTimeout` | `string` | Duration for health checks (e.g., `5s`, `30s`). |
| `password` | `string` | **Hashed.** Used for session auth. Set via the UI. |
