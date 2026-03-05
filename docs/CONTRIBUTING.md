# Contributing to homectl

Thank you for considering contributing to homectl! This document provides guidelines and instructions for contributing.

## Prerequisites

To contribute to homectl, you'll need the following installed:

- **Node.js**: `v18.0.0` or newer
- **Go**: `v1.24.0` or newer
- **npm**: `v10.0.0` or newer (for workspace support)

## Getting Started

1. **Fork the repository** on GitHub.
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/homectl.git
   cd homectl
   ```
3. **Install dependencies**:
   ```bash
   npm install
   ```
   *This will install dependencies for the root and all workspaces (web, server, shared).*
4. **Create a branch**:
   ```bash
   git checkout -b feature/my-feature
   ```

## Development

homectl uses a monorepo structure with npm workspaces.

```bash
# Run both frontend and backend development servers concurrently
npm run dev

# Run tests (Web and Server)
npm run test

# Lint all codebases
npm run lint

# Type check (Frontend)
npm run typecheck
```

### Hot Reloading

- **Frontend**: The React app uses Vite and will automatically reload on changes.
- **Backend**: The Go server includes a configuration watcher that reloads settings on `config.yaml` changes, but code changes currently require a restart (standard `go run` behavior).

### Configuration

On the first run, the server will automatically generate a default `config.yaml` in the project root if it doesn't exist. You can modify this file to test different dashboard configurations.

## Project Structure

```
homectl/
├── apps/
│   ├── web/           # React frontend (Vite, TypeScript, Tailwind)
│   └── server/        # Go backend (Fiber)
├── packages/
│   └── shared/        # Shared types & JSON schemas
├── data/              # Default storage for SQLite and icons
└── docs/              # Documentation source (MkDocs)
```

## Making Changes

### Frontend (React/TypeScript)

- Follow existing code style and naming conventions.
- Use TypeScript strictly; avoid using `any`.
- Write tests for new components using Vitest.
- Ensure all UI changes are responsive and follow the monochrome technical aesthetic.

### Backend (Go)

- Follow Go idioms and conventions.
- Add unit tests for new functionality in `_test.go` files.
- Document public APIs and internal functions where appropriate.
- Handle errors explicitly and with context.

### Documentation

- Update relevant files in `docs/` for any new features or configuration changes.
- Ensure `CONFIG_REFERENCE.md` is updated if new configuration options are added.

## Pull Request Process

1. **Before submitting:**
   - Run `npm run lint && npm run test` to ensure everything is correct.
   - Update documentation if your changes introduce new features or settings.
   - Rebase your branch on the latest `main`.

2. **PR Title:**
   - Use [Conventional Commits](https://www.conventionalcommits.org/) format.
   - Examples: `feat: add TCP check widget`, `fix: config validation error`, `docs: update installation guide`.

3. **PR Description:**
   - Clearly describe what changed and why.
   - Link any related issues (e.g., `Closes #123`).
   - Include screenshots or GIFs for UI-related changes.

4. **Review:**
   - Be responsive to feedback and make requested changes.
   - Re-request a review once you've addressed all comments.

## Issue Guidelines

### Bug Reports

- Provide clear steps to reproduce the issue.
- Include your sanitized configuration file.
- Attach relevant error messages or logs from the browser console or server output.

### Feature Requests

- Describe the use case and why the feature would be valuable.
- Provide examples or mockups if possible.

## Release Process

homectl follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## Questions?

- Open an issue for discussion or clarification.
- Check existing issues and the documentation site first.

Thank you for building the best homelab dashboard with us! 🎉
