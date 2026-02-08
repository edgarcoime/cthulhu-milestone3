# Gateway

HTTP API gateway that fronts the auth, filemanager, and lifecycle gRPC services. Single entry point for the client (e.g. Next.js) with CORS and optional-auth middleware. Part of [Cthulhu](../README.md); see the root README for full stack setup, prerequisites (Go 1.26, mise, Docker), and running everything with Docker Compose.

## What it does

- **Auth**: OAuth initiate/callback, token refresh, logout, validate.
- **Files**: Upload (prepare → confirm), bucket authenticate, get bucket, bucket admins, protected check, presigned download.
- **Lifecycle**: Get bucket lifecycle (expiry) by bucket ID.
- **Server**: Fiber app with CORS, request logging, and graceful shutdown; proxies requests to the backend microservices.

## Prerequisites

- **Go 1.26** (or version in root `mise.toml`)
- Proto/grpc code: from repo root run `mise install` then `mise run proto`

## Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run the gateway (`go run ./cmd/service`) |
| `make test` | Run tests |
| `make lint` | Run golangci-lint |
| `make clean` | Remove `tmp/` |

## Setup

1. Use the root `.env` or copy `gateway/.env.example` to `.env` in this directory.
2. Set `CORS_ORIGIN` to your frontend origin (e.g. `http://localhost:3000`).
3. Set gRPC URLs: `AUTH_GRPC_URL`, `FILEMANAGER_GRPC_URL`, `LIFECYCLE_GRPC_URL` (e.g. `localhost:49051`, `localhost:48051`, `localhost:50051` when all services run on host).
4. Run `make dev`. The gateway listens on port **7777** (or `APP_PORT` from env).

## Run with Docker Compose

This service is included in the root `docker compose` stack. From the repo root: `docker compose up --build`. See [root README](../README.md#docker-compose). In Compose, gRPC URLs use service names (e.g. `auth:49051`, `filemanager:48051`, `lifecycle:50051`).
