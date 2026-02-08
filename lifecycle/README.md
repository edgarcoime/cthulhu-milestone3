# Lifecycle

Microservice that manages bucket lifecycles: stores expiry times for buckets and periodically purges expired ones by calling the filemanager to delete buckets and their files. Part of [Cthulhu](../README.md); see the root README for full stack setup, prerequisites (Go 1.26, mise, Docker), and running everything with Docker Compose.

## What it does

- **gRPC API**: Create, get, and delete lifecycle records (bucket slug + expiry time).
- **Cleanup daemon**: Runs every X (usually 15minutes but can be decreased for demonstration purposes), finds expired lifecycles, calls the filemanager to delete those buckets, then removes the lifecycle records.

## Prerequisites

- **Go 1.26** (or version in root `mise.toml`)
- Proto/grpc code: from repo root run `mise install` then `mise run proto`

## Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the service binary to `bin/service` |
| `make dev` | Build and run the service (loads `.env` if present) |
| `make test` | Run tests |
| `make lint` | Run golangci-lint |
| `make clean` | Remove `bin/` and `tmp/` |

Proto/grpc code is generated from the **repo root**: `mise run proto`.

## Setup

1. Use the root `.env` or copy `lifecycle/.env.example` to `.env` in this directory. Set `FILEMANAGER_GRPC_URL` to your filemanager gRPC address (e.g. `localhost:48051` when running locally).
2. Run `make dev` (or `make build` then `./bin/service`).

## Run with Docker Compose

This service is included in the root `docker compose` stack. From the repo root: `docker compose up --build`. See [root README](../README.md#docker-compose). In Compose, `FILEMANAGER_GRPC_URL` is set to `filemanager:48051`.
