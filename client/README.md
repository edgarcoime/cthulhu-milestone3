# Cthulhu Web Client

Next.js frontend for Cthulhu. Part of [Cthulhu](../README.md); see the root README for full stack setup, prerequisites (Node.js, Docker), and running everything with Docker Compose.

## Prerequisites

- **Node.js** (LTS, e.g. 20+)
- API URL: set `NEXT_PUBLIC_API_URL` (e.g. `http://localhost:7777`) for the gateway; this is build-time for the Next.js app.

## Development

From this directory:

```bash
make dev
```

Or with npm:

```bash
npm install
npm run dev
```

The app runs at [http://localhost:3000](http://localhost:3000) and uses the gateway at `NEXT_PUBLIC_API_URL` for auth and file APIs.

## Run with Docker Compose

This service is included in the root `docker compose` stack. From the repo root: `docker compose up --build`. The image is built with `NEXT_PUBLIC_API_URL` (default `http://localhost:7777`). See [root README](../README.md#docker-compose).
