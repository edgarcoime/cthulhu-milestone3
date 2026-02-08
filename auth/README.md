# Auth

Microservice that handles authentication: OAuth (e.g. GitHub), JWT access/refresh tokens, and session storage. Part of [Cthulhu](../README.md); see the root README for full stack setup, prerequisites (Go 1.26, mise, Docker), and running everything with Docker Compose.

## What it does

- **OAuth**: Initiate OAuth flow (PKCE) and handle callback; creates/updates users and returns access + refresh tokens.
- **Tokens**: Validate access tokens, refresh token rotation, logout (revoke refresh token).
- **Storage**: SQLite for users, refresh tokens, and OAuth session state.

## Prerequisites

- **Go 1.26** (or version in root `mise.toml`)
- For gRPC codegen: from repo root run `mise install` then `mise run proto`

## Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run the service (`go run ./cmd/service`) |
| `make sqlc` | Regenerate SQLc code (after changing `internal/repository/sqlc/*.sql`) |
| `make test` | Run tests |
| `make lint` | Run golangci-lint |
| `make clean` | Remove `tmp/` |

Proto/grpc code is generated from the **repo root**: `mise run proto`.

## Setup

1. Copy `.env.example` to `.env`.
2. Set `JWT_SECRET` (or `JWT_SECRET_KEY`) to a strong secret for signing JWTs.
3. For GitHub OAuth: set `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, and `GITHUB_REDIRECT_URI` from your GitHub OAuth app.
4. Run `make dev`.

## Run with Docker Compose

This service is included in the root `docker compose` stack. From the repo root: `docker compose up --build`. See [root README](../README.md#docker-compose).
