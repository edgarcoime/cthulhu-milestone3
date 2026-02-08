# Cthulhu

## What is CTHULHU?

CTHULHU is an anonymous, file sharing platform that lets anyone upload up to 1 GB of files without an account and share a secure URL that expires after 48 hours. Authorized users can extend retention up to 14 days, manage, and delete their uploads on demand, and optionally password-protect shared content. Built on a microservices architecture with a focus on scalability, security, and user privacy.

## <a name="toc">Table of Contents</a>

- [What is CTHULHU](#what-is-cthulhu)
- [Table of Contents](#toc)
- [Prerequisites](#prerequisites)
- [Development Environment](#development-environment)
  - [1. Environment and config](#1-environment-and-config)
  - [2. Setup each service](#2-setup-each-service)
  - [3. Start services (order and commands)](#3-start-services)
  - [4. Test the application](#4-test-the-application)
- [Run everything with Docker Compose](#docker-compose)
  - [Prerequisites and env](#docker-compose-prereqs)
  - [Build and run](#docker-compose-run)

---

## Prerequisites

Install the following before setting up the project:

| Tool      | Version   | Purpose                          | Links |
|-----------|-----------|----------------------------------|-------|
| **Go**    | 1.26      | Backend services (auth, filemanager, gateway, lifecycle) | [download](https://go.dev/doc/install) |
| **mise**  | latest    | Tool versions and tasks (proto, sqlc) | [install](https://mise.jdx.dev/getting-started.html) |
| **Terraform** | recent | Infrastructure (optional; see `terraform/`) | [download](https://developer.hashicorp.com/terraform/install) |
| **Docker**| recent    | Running dependencies and full stack via Compose | [download](https://docs.docker.com/engine/install/) |
| **Node.js** | LTS (e.g. 20+) | Client frontend (Next.js) | [download](https://nodejs.org/en/download) |
| **Make**  | —         | Running `make dev` in each service | Usually preinstalled on UNIX |

**Using mise (recommended):** From the repo root, run `mise install` to install the Go version and tools defined in `mise.toml` (e.g. `protoc`, `sqlc`, protobuf/grpc codegen). This keeps Go and codegen in sync with the project.

---

## Development Environment

Work from the **repository root** unless a step says otherwise.

### 1. Environment and config

- Copy the root env example and fill in secrets (AWS, OAuth, JWT, etc.):

  ```bash
  cp .env.example .env
  # Edit .env with your values
  ```

- Each Go service that needs its own config has an `.env.example` in its folder (e.g. `auth/`, `filemanager/`, `gateway/`, `lifecycle/`). Copy and edit as needed:

  ```bash
  cp auth/.env.example auth/.env
  cp filemanager/.env.example filemanager/.env
  # etc.
  ```

- Ensure **gRPC URLs** in `.env` point to where services run (e.g. `localhost:49051` for auth when running everything on the host).

### 2. Setup each service

| Service      | Path         | Setup steps |
|-------------|--------------|-------------|
| **Auth**    | `./auth`     | `cp .env.example .env`, fill OAuth/JWT. Uses SQLite; DB path in `.env`. |
| **Filemanager** | `./filemanager` | `cp .env.example .env`. Set S3/LocalStack vars, `AUTH_GRPC_URL`, and optionally RabbitMQ URL if used. |
| **Gateway** | `./gateway`  | Uses root `.env` or own `.env`; set `AUTH_GRPC_URL`, `FILEMANAGER_GRPC_URL`, `LIFECYCLE_GRPC_URL`, `CORS_ORIGIN`. |
| **Lifecycle** | `./lifecycle` | Uses root `.env` or own `.env`; set `FILEMANAGER_GRPC_URL`. |
| **Client**  | `./client`   | Next.js app. Set `NEXT_PUBLIC_API_URL` (e.g. `http://localhost:7777`) for API base URL. |

**Code generation (mise):**

- From repo root, install tools and generate proto/sqlc if needed:

  ```bash
  mise install
  mise run proto   # Protobuf/ gRPC Go code
  mise run sqlc    # SQLc for auth and filemanager
  ```

### <a name="3-start-services">3. Start services</a>

Start services **in this order** (each in its own terminal, from the given directory):

1. **Auth** — `./auth`  
   ```bash
   cd auth && make dev
   ```
2. **Filemanager** — `./filemanager`  
   ```bash
   cd filemanager && make dev
   ```
3. **Lifecycle** — `./lifecycle`  
   ```bash
   cd lifecycle && make dev
   ```
4. **Gateway** — `./gateway`  
   ```bash
   cd gateway && make dev
   ```
5. **Client** — `./client`  
   ```bash
   cd client && make dev
   ```

Generic pattern:

```bash
cd <service_folder>
make dev
```

### 4. Test the application

- **Client:** [http://localhost:3000](http://localhost:3000)
- **Gateway (API):** [http://localhost:7777](http://localhost:7777)
- **Auth, Filemanager, Lifecycle:** Not exposed directly; all traffic goes through the Gateway.

---

## <a name="docker-compose">Run everything with Docker Compose</a>

You can run the full stack with Docker Compose from the repo root.

### <a name="docker-compose-prereqs">Prerequisites and env</a>

- **Docker** (and Docker Compose v2) installed.
- Root **`.env`** file with the same variables as in `.env.example` (JWT, GitHub OAuth, S3/LocalStack, bucket token secret, etc.). Compose reads this file by default.

For Compose, gRPC URLs should use **service names** so containers can reach each other (e.g. `AUTH_GRPC_URL=auth:49051`). The `docker-compose.yml` defaults already use these; you can omit them or set explicitly.

### <a name="docker-compose-run">Build and run</a>

From the **repository root**:

```bash
# Ensure .env is present and populated
cp .env.example .env
# Edit .env as needed

# Build and start all services
docker compose up --build
```

To run in the background:

```bash
docker compose up --build -d
```

**Services and ports:**

| Service      | Port(s)     | Notes |
|-------------|-------------|--------|
| **auth**    | 49051       | gRPC |
| **filemanager** | 48051   | gRPC; may need LocalStack (e.g. on host) for S3 |
| **lifecycle**   | 50051   | gRPC |
| **gateway**     | 7777    | HTTP API; set `CORS_ORIGIN` (e.g. `http://localhost:3000`) |
| **client**      | 3000    | Next.js; `NEXT_PUBLIC_API_URL` is set at build time (e.g. `http://localhost:7777`) |

**Useful commands:**

```bash
docker compose down          # Stop and remove containers
docker compose up -d --build  # Rebuild and run in background
docker compose logs -f        # Follow logs
```

For S3-compatible storage (e.g. LocalStack), run it separately and set `S3_ENDPOINT` / `S3_PRESIGNED_ENDPOINT` in `.env` as in `.env.example` (e.g. `http://localhost:4566` for local, or `http://host.docker.internal:4566` for containers talking to host).
