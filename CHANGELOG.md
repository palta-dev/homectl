# Changelog

All notable changes to homectl will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release with core functionality

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
