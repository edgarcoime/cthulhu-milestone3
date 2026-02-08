# syntax=docker/dockerfile:1

# --- Build Stage ---
FROM golang:1.25.6-bookworm AS builder

WORKDIR /src

# Copy local dependencies (paths relative to repo root build context)
COPY common/ ./common/
COPY proto/ ./proto/
COPY auth/go.mod auth/go.sum ./auth/
COPY auth/ ./auth/
COPY filemanager/go.mod filemanager/go.sum ./filemanager/
COPY filemanager/ ./filemanager/
COPY lifecycle/go.mod lifecycle/go.sum ./lifecycle/
COPY lifecycle/ ./lifecycle/

# Copy gateway module files
COPY gateway/go.mod gateway/go.sum ./gateway/

# Download dependencies
WORKDIR /src/gateway
RUN go mod download

# Copy the rest of the gateway source code
COPY gateway/ ./

# Build the binary (no CGO required)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/gateway-service ./cmd/service

# --- Final Stage ---
FROM debian:bookworm-slim

# Install CA certificates for HTTPS/gRPC
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/gateway-service .

# Default environment variables (can be overridden by docker-compose)
ENV APP_HOST=0.0.0.0 \
    APP_PORT=7777 \
    CORS_ORIGIN=* \
    AUTH_GRPC_URL= \
    FILEMANAGER_GRPC_URL= \
    LIFECYCLE_GRPC_URL=

# Run as non-root user
USER 1000:1000

EXPOSE 7777

CMD ["./gateway-service"]
