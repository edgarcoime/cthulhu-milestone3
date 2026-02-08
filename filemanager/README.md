# Filemanager

Microservice that manages file buckets and S3-compatible storage: presigned upload/download, bucket metadata, optional password protection, and auth integration. Part of [Cthulhu](../README.md); see the root README for full stack setup, prerequisites (Go 1.26, mise, Docker), and running everything with Docker Compose.

## What it does

- **Uploads**: Two-phase presigned URL flow — PrepareUpload returns presigned PUT URLs; client uploads to S3; ConfirmUpload persists file metadata in SQLite.
- **Downloads**: PrepareDownload returns a presigned GET URL; for password-protected buckets, a bucket access token is required.
- **Buckets**: Create buckets (with optional password), list files, get bucket admins (via auth service), check if protected, authenticate (password or user) to get a bucket access token.
- **Storage**: S3-compatible backend (e.g. AWS S3 or LocalStack); talks to the auth service for user/admin resolution.

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
2. Set `AUTH_GRPC_URL` to your auth gRPC address (e.g. `localhost:49051` when running locally).
3. Configure S3: `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, `S3_ENDPOINT` (e.g. `http://localhost:4566` for LocalStack), `S3_PRESIGNED_ENDPOINT`, `S3_REGION`, `S3_BUCKET_NAME`. Use `S3_FORCE_PATH_STYLE=true` for LocalStack.
4. Run `make dev`.

## Run with Docker Compose

This service is included in the root `docker compose` stack. From the repo root: `docker compose up --build`. See [root README](../README.md#docker-compose). For S3, configure `S3_ENDPOINT` / `S3_PRESIGNED_ENDPOINT` in the root `.env` (e.g. LocalStack on host: `http://host.docker.internal:4566`).
