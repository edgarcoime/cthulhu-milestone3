# syntax=docker/dockerfile:1

# --- Build Stage ---
FROM golang:1.25.6-bookworm AS builder

# Install build dependencies for CGO (required by go-sqlite3)
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

# Copy local dependencies first to leverage Docker cache
# These paths are relative to the repo root (build context)
COPY common/ ./common/
COPY proto/ ./proto/
COPY auth/go.mod auth/go.sum ./auth/
COPY auth/ ./auth/

# Copy filemanager module files
COPY filemanager/go.mod filemanager/go.sum ./filemanager/

# Download dependencies
WORKDIR /src/filemanager
RUN go mod download

# Copy the rest of the filemanager source code
COPY filemanager/ ./

# Build the binary
# CGO_ENABLED=1 is mandatory for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/filemanager-service ./cmd/service

# --- Final Stage ---
FROM debian:bookworm-slim

# Install CA certificates for S3/HTTPS and sqlite3 for DB management
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/filemanager-service .

# Create a directory for the SQLite database
RUN mkdir -p /data && chown -R 1000:1000 /data

# Default environment variables (can be overridden by docker-compose)
ENV APP_HOST=0.0.0.0 \
    APP_PORT=48051 \
    SQLITE_DB_FILE=/data/filemanager.db \
    AUTH_GRPC_URL= \
    BUCKET_TOKEN_SECRET_KEY= \
    S3_ACCESS_KEY_ID= \
    S3_SECRET_ACCESS_KEY= \
    S3_ENDPOINT= \
    S3_REGION= \
    S3_BUCKET_NAME= \
    S3_FORCE_PATH_STYLE=true

# Run as a non-root user for security
USER 1000:1000

EXPOSE 48051

CMD ["./filemanager-service"]
