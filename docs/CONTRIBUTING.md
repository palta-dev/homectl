# Contributing to homectl

Thank you for considering contributing to homectl! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Welcome newcomers and help them learn

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/homectl.git`
3. Install dependencies: `npm install`
4. Create a branch: `git checkout -b feature/my-feature`

## Development

```bash
# Run development server
npm run dev

# Run tests
npm run test

# Lint code
npm run lint

# Type check
npm run typecheck
```

## Project Structure

```
homectl/
├── apps/
│   ├── web/           # React frontend
│   └── server/        # Go backend
├── packages/
│   └── shared/        # Shared types & schemas
└── docs/              # Documentation
```

## Making Changes

### Frontend (React/TypeScript)

- Follow existing code style
- Use TypeScript strictly (no `any`)
- Write tests for new components
- Ensure responsive design

### Backend (Go)

- Follow Go idioms and conventions
- Add tests for new functionality
- Document public APIs
- Handle errors appropriately

### Documentation

- Update docs for new features
- Keep CONFIG_REFERENCE.md up to date
- Add examples where helpful

## Pull Request Process

1. **Before submitting:**
   - Run `npm run lint && npm run test`
   - Ensure all checks pass
   - Update documentation if needed

2. **PR Title:**
   - Use conventional commits format
   - Examples: `feat: add TCP check widget`, `fix: config validation error`

3. **PR Description:**
   - Describe what changed and why
   - Link related issues
   - Include screenshots for UI changes

4. **Review:**
   - Respond to feedback promptly
   - Make requested changes
   - Re-request review when ready

## Issue Guidelines

### Bug Reports

- Include steps to reproduce
- Provide configuration (sanitized)
- Include error messages/logs
- Specify version and environment

### Feature Requests

- Describe the use case
- Explain why it's needed
- Provide examples if possible

## Release Process

Releases follow semantic versioning:

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## Questions?

- Open an issue for discussion
- Check existing issues and docs first
- Be patient - maintainers are volunteers

Thank you for contributing! 🎉
