# Game Engine

A Multi-Tenant Real-Time Gaming & Engagement Backend Engine. The goal of the project is to dive deep into various Redis features, zero frontend and zero database. It is driven entirely via a command-line interface.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) with Docker Compose
- [Go](https://go.dev/dl/) 1.24+

## Redis Setup

This project uses **Redis Stack** (Redis + modules for JSON, Search, TimeSeries, etc.) via Docker.

Start the Redis container:

```bash
docker compose up -d
```

This starts:
| Service | Port | Description |
|---|---|---|
| Redis | `6379` | Primary connection used by the app |
| RedisInsight | `8001` | Browser UI for inspecting Redis data — open `http://localhost:8001` |

Data is persisted in a Docker volume (`redis-data`) and snapshotted every 60 seconds.

To stop:

```bash
docker compose down
```

To wipe all data:

```bash
docker compose down -v
```

## Running the App

Install dependencies and build:

```bash
go mod download
go build -o game-engine ./cmd
```

Test the Redis connection:

```bash
./game-engine ping
```

Expected output when Redis is reachable:

```
✅ Redis connected successfully
# Server
redis_version:x.x.x
...
PONG
```

## CLI Commands

```
./game-engine ping     # Verify Redis connection
./game-engine flush    # Flush all data (dev only — prompts for confirmation)
```

## Connection Details

The app connects to Redis using these defaults (see [internal/client/client.go](internal/client/client.go)):

| Setting      | Default          |
| ------------ | ---------------- |
| Address      | `localhost:6379` |
| Password     | none             |
| Database     | `0`              |
| Dial timeout | 5s               |
| Pool size    | 10 connections   |

The client startup will `PING` Redis and print server info. If the connection fails, the app exits with an error — ensure the Docker container is running before starting the app.
