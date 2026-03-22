# Changelog

All notable changes to homectl will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2024-03-21

### Added
- Background Monitoring Worker: Periodic health checks now run in the background.
- Historical Data: Added SQLite-backed latency history for the last 24 hours.
- Dynamic Sparklines: Service cards now render real latency graphs from historical data.
- ICMP Ping Widget: Native ping support for host reachability.
- Advanced HTTP JSON Widget: Added support for custom HTTP methods and headers (for API tokens).

### Changed
- Improved DNS Rebinding Protection by pinning resolved IPs during requests.
- Enhanced configuration validation with duration string parsing.
- Refactored server initialization to support persistent SQLite storage by default.
- Disabled Docker discovery by default to prioritize a security-first posture.

### Fixed
- Fixed bug in environment variable expansion returning literal placeholders when missing.
- Removed verbose debug logging in production builds.
- Corrected repository and organization metadata in dashboard footer.

## [0.1.0] - 2024-03-02

### Added

#### Core Features
- Config-driven dashboard with YAML configuration
- Service groups with grid, list, and compact layouts
- Health checks: HTTP, TCP, and Ping
- Widget system: httpJson, httpHtml, tcpPort, httpStatus
- Live config reload on file changes
- Dark/light theme toggle
- Search functionality across services
- Responsive design for mobile/desktop

#### Security
- SSRF protection with allowlist/blocklist
- Rate limiting middleware (10 req/s default)
- Environment variable expansion for secrets
- Cloud metadata IP blocking

#### Infrastructure
- Multi-stage Docker build (~25MB final image)
- Multi-arch support (linux/amd64, linux/arm64)
- Unit tests for backend and frontend

#### Documentation
- README with quickstart guide
- ARCHITECTURE.md with system design
- CONFIG_REFERENCE.md with full schema
- SECURITY.md with security best practices
- CONTRIBUTING.md with contribution guidelines

### Technical Details

#### Backend (Go + Fiber)
- Config loading with validation
- LRU cache with TTL per widget
- SQLite storage for incident history (optional)
- SSRF-safe HTTP client with CIDR matching
- Rate limiting middleware

#### Frontend (React + Vite)
- TypeScript for type safety
- Tailwind CSS for styling
- React Query for data fetching
- Zustand for state management
- Loading skeletons for better UX
- Collapsible service groups

### Known Issues
- Ping check requires elevated privileges (uses TCP fallback)
- Docker auto-discovery not yet implemented
- OAuth authentication not yet implemented

## Version History

| Version | Date | Notes |
|---------|------|-------|
| 0.1.0 | 2024-03-02 | Initial release |

[Unreleased]: https://github.com/palta-dev/homectl/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/palta-dev/homectl/releases/tag/v0.1.0
