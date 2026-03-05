# syntax=docker/dockerfile:1

# Build arguments for version
ARG VERSION=0.1.0
ARG BUILD_DATE

# ============================================
# Frontend Build Stage
# ============================================
FROM node:20-alpine AS frontend-build
WORKDIR /app

# Copy shared package first
COPY packages/shared/package*.json ./packages/shared/
# Copy web package
COPY apps/web/package*.json ./apps/web/

# Install dependencies (using --workspace for monorepos if needed, but here simple install)
# Since we are not using a root package.json for npm workspaces, we'll install separately
WORKDIR /app/apps/web
RUN npm install

# Copy source code
COPY packages/shared/ /app/packages/shared/
COPY apps/web/ /app/apps/web/

# Build frontend
RUN npm run build

# ============================================
# Backend Build Stage
# ============================================
FROM golang:1.24-alpine AS backend-build
WORKDIR /app/apps/server

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY apps/server/go.mod apps/server/go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY apps/server/ .

# Build backend with optimizations
ARG VERSION
ARG BUILD_DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE}" \
    -o /app/homectl \
    ./cmd

# ============================================
# Production Image
# ============================================
FROM alpine:3.19 AS production
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    wget \
    ca-certificates \
    tzdata

# Create non-root user
RUN addgroup -g 1000 homectl && \
    adduser -D -u 1000 -G homectl -h /app homectl

# Copy frontend build
COPY --from=frontend-build --chown=homectl:homectl \
    /app/apps/web/dist ./static

# Copy backend binary
COPY --from=backend-build --chown=homectl:homectl \
    /app/homectl ./

# Create data directories
RUN mkdir -p /app/data/icons /app/data/db && \
    chown -R homectl:homectl /app

# Switch to root user for system stats and docker socket access
USER root

# Environment variables
ENV HOMECTL_CONFIG=/app/data/config.yaml \
    TZ=UTC

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=10s \
    CMD wget -q --spider http://localhost:8080/api/health || exit 1

# Entrypoint
ENTRYPOINT ["./homectl"]
CMD ["--config", "/app/data/config.yaml"]
