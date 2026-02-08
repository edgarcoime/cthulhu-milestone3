# Project Changes Summary

This document summarizes the major changes made to the Cthulhu AWS project.

---

## 1. Incorporating gRPC

Service-to-service communication was moved to **gRPC** with shared Protobuf definitions.

### Proto definitions (`proto/proto/`)

| File                | Purpose                                                                        |
| ------------------- | ------------------------------------------------------------------------------ |
| `auth.proto`        | Auth service: OAuth initiation/callback, token validation, refresh, logout     |
| `filemanager.proto` | Filemanager service: upload/download flows, bucket admins, auth, delete bucket |
| `lifecycle.proto`   | Lifecycle service: create/get/delete lifecycle records (bucket slug + expiry)  |

### Generated code

- Proto files are compiled to Go with `protoc` (via `mise run proto`).
- Generated packages live under `proto/` and are consumed by each service (e.g. `lifecycle/pkg/pb/`).

### Service usage

- **Auth** — Exposes `AuthService` over gRPC (port 49051).
- **Filemanager** — Exposes `FilemanagerService` over gRPC (port 48051); uses auth via gRPC client.
- **Lifecycle** — Exposes `LifecycleService` over gRPC (port 50051); calls filemanager via gRPC for purge.
- **Gateway** — HTTP API that talks to auth, filemanager, and lifecycle via gRPC clients (see `gateway/internal/connections/grpc.go`).

gRPC URLs are configurable via env (e.g. `AUTH_GRPC_URL`, `FILEMANAGER_GRPC_URL`, `LIFECYCLE_GRPC_URL`); in Docker Compose they use service names (e.g. `auth:49051`, `filemanager:48051`, `lifecycle:50051`).

---

## 2. Setup `mise.toml`

A root **`mise.toml`** was added to pin tool versions and define tasks for code generation.

### Tools

| Tool                                               | Version | Purpose                       |
| -------------------------------------------------- | ------- | ----------------------------- |
| `go`                                               | 1.25.6  | Backend services              |
| `protoc`                                           | latest  | Protobuf compiler             |
| `go:google.golang.org/protobuf/cmd/protoc-gen-go`  | latest  | Go protobuf codegen           |
| `go:google.golang.org/grpc/cmd/protoc-gen-go-grpc` | latest  | Go gRPC codegen               |
| `go:github.com/sqlc-dev/sqlc/cmd/sqlc`             | latest  | SQL → Go for auth/filemanager |

### Tasks

| Task             | Description                                                                                                       |
| ---------------- | ----------------------------------------------------------------------------------------------------------------- |
| `mise run proto` | Compile all `.proto` files in `proto/proto/` into Go under `proto/` (module `github.com/cthulhu-platform/proto`). |
| `mise run sqlc`  | Generate sqlc code for `auth` and `filemanager` from their `internal/repository/sqlc` definitions.                |

**Usage:** From repo root, run `mise install` then `mise run proto` (and `mise run sqlc` when SQL/sqlc config changes). READMEs for each service reference the root `mise.toml` for Go/codegen setup.

---

## 3. Create Lifecycle Microservice

A new **lifecycle** microservice was added to manage bucket expiry and periodic purge.

### Location

`lifecycle/`

### Responsibilities

- **gRPC API** (`LifecycleService`): Create, get, and delete lifecycle records (bucket slug + expiry time).
- **Cleanup daemon**: On a configurable interval (e.g. 15 minutes), finds expired lifecycles, calls the filemanager gRPC to delete those buckets and their files, then removes the lifecycle records.

### Structure (high level)

- `cmd/service/main.go` — Entrypoint.
- `internal/connections/grpc.go` — gRPC client to filemanager.
- `internal/daemon/cleanup.go` — Periodic purge logic.
- `internal/repository/` — Persistence (e.g. SQLite schema in `schemas/up.sql`).
- `internal/server/server.go` — gRPC server.
- `internal/service/service.go` — Business logic (purge expired buckets via filemanager client).
- `pkg/client/` — gRPC client used by the gateway.
- `pkg/pb/` — Generated lifecycle protobuf/grpc code (from `proto/proto/lifecycle.proto`).

### Configuration

- Env (e.g. `FILEMANAGER_GRPC_URL`) for filemanager gRPC address; in Docker Compose this is `filemanager:48051`.
- Optional `.env` from `lifecycle/.env.example`.

---

## 4. Create Dockerfiles for Each Service

Each service has its own **Dockerfile** for building a production-style image.

| Service         | Dockerfile               | Port  | Notes                                                                                      |
| --------------- | ------------------------ | ----- | ------------------------------------------------------------------------------------------ |
| **auth**        | `auth/Dockerfile`        | 49051 | Multi-stage; Go 1.25.6, SQLite (CGO); copies `common/`, `proto/`.                          |
| **filemanager** | `filemanager/Dockerfile` | 48051 | Multi-stage; SQLite; copies `common/`, `proto/`, `auth/` (for proto/client deps).          |
| **lifecycle**   | `lifecycle/Dockerfile`   | 50051 | Multi-stage; SQLite; copies `common/`, `proto/`, `auth/`, `filemanager/`.                  |
| **gateway**     | `gateway/Dockerfile`     | 7777  | Multi-stage; copies `common/`, `proto/`, and other services as needed for client packages. |
| **client**      | `client/Dockerfile`      | 3000  | Next.js app; build-time `NEXT_PUBLIC_API_URL` via build arg.                               |

### Common patterns

- **Build context**: Root of the repo (so Dockerfiles can `COPY common/`, `COPY proto/`, etc.).
- **Multi-stage**: Builder stage compiles Go (or Node for client); final stage is minimal (e.g. `debian:bookworm-slim` for Go services).
- **Non-root user**: Go services run as `USER 1000:1000`.
- **Data volumes**: Services that use SQLite create `/data` and expect it to be mounted (e.g. Compose volumes `auth-data`, `filemanager-data`, `lifecycle-data`).

---

## 5. Consolidated Dockerfiles Using Docker Compose

A single **`docker-compose.yml`** at the repo root builds and runs all services using the per-service Dockerfiles.

### Services in Compose

| Service         | Build                                              | Ports | Volumes          | Notes                                                                                              |
| --------------- | -------------------------------------------------- | ----- | ---------------- | -------------------------------------------------------------------------------------------------- |
| **auth**        | `context: .`, `dockerfile: auth/Dockerfile`        | 49051 | auth-data        | OAuth/JWT env vars.                                                                                |
| **filemanager** | `context: .`, `dockerfile: filemanager/Dockerfile` | 48051 | filemanager-data | S3/LocalStack, `AUTH_GRPC_URL=auth:49051`; `extra_hosts` for host.docker.internal.                 |
| **lifecycle**   | `context: .`, `dockerfile: lifecycle/Dockerfile`   | 50051 | lifecycle-data   | `FILEMANAGER_GRPC_URL=filemanager:48051`.                                                          |
| **gateway**     | `context: .`, `dockerfile: gateway/Dockerfile`     | 7777  | —                | `AUTH_GRPC_URL`, `FILEMANAGER_GRPC_URL`, `LIFECYCLE_GRPC_URL` set to service names; `CORS_ORIGIN`. |
| **client**      | `context: .`, `dockerfile: client/Dockerfile`      | 3000  | —                | Build arg `NEXT_PUBLIC_API_URL` (e.g. `http://localhost:7777`).                                    |

### Volumes

- `auth-data`, `filemanager-data`, `lifecycle-data` — Persist SQLite DBs and any service-specific data.

### Usage

- **Build and run:** From repo root, `docker compose up --build`.
- **Env:** Use root `.env` (see `.env.example`); gRPC URLs in Compose default to service names so containers can reach each other.

The root README documents prerequisites (Docker, env vars) and the Docker Compose workflow so that all services, including the lifecycle microservice, run together with a single command.
