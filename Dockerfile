# syntax=docker/dockerfile:1

# Build arguments for version
ARG VERSION=0.2.0
ARG BUILD_DATE

# ============================================
# Frontend Build Stage
# ============================================
FROM node:20-alpine AS frontend-build
WORKDIR /app

# Copy root package.json for workspaces
COPY package*.json ./
COPY apps/web/package*.json ./apps/web/
COPY packages/shared/package*.json ./packages/shared/

# Install dependencies using workspaces
RUN npm install --include=dev --workspace=@homectl/web --workspace=@homectl/shared

# Copy source code
COPY packages/shared/ ./packages/shared/
COPY apps/web/ ./apps/web/

# Build frontend
RUN npm run build --workspace=@homectl/web

# ============================================
# Backend Build Stage
# ============================================
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS backend-build
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY apps/server/go.mod apps/server/go.sum ./apps/server/

# Download dependencies
WORKDIR /app/apps/server
RUN go mod download

# Copy source code
COPY apps/server/ .

# Build backend with optimizations
ARG VERSION
ARG BUILD_DATE
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
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
    tzdata \
    iputils

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
EXPOSE 7777

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=10s \
    CMD wget -q --spider http://localhost:7777/api/health || exit 1

# Entrypoint
ENTRYPOINT ["./homectl"]
CMD ["--config", "/app/data/config.yaml"]
