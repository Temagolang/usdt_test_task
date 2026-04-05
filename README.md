# USDT Rates Service

gRPC service for fetching USDT exchange rates (ask/bid) from Grinex exchange with PostgreSQL storage.

> **Note:** Module path `github.com/example/grinex-rates-service` is a placeholder.
> Replace with the actual module path before publishing.

## Quick Start

Acceptance flow — proves the app builds and starts:

```bash
make build
docker-compose up -d                       # starts postgres only
docker-compose run --rm app ./app          # starts the service, ctrl-C to stop
```

Run tests:

```bash
make test
```

For development with published ports (gRPC 50051, HTTP 8080):

```bash
make run
```

## Project Structure

```
main.go                              # Thin bootstrap
cmd/                                 # CLI commands (grpc, migrate)
internal/
  app/                               # Composition root, wiring, lifecycle
  config/                            # Section-based config (env + flags)
  transport/grpc/                    # gRPC transport (thin, no business logic)
  service/rates/                     # Business logic
  client/grinex/                     # Grinex HTTP client (resty)
  repo/rates/                        # Repository interface
  repo/rates/postgres/               # PostgreSQL implementation
  migration/postgres/                # SQL migrations
  observability/otel/                # OpenTelemetry foundation
api/proto/rates/v1/                  # Protobuf definitions
```

## Commands

```bash
make build          # Build binary
make test           # Run tests
make lint           # Run golangci-lint
make docker-build   # Build Docker image
make run            # Start postgres + app with published ports (dev)
make gen            # Generate proto + sqlc code
make migrate-up     # Apply pending migrations
make migrate-down   # Rollback last migration
```

## Configuration

Configuration via environment variables and CLI flags. **Flags take priority over env vars.**

| Section  | Parameter      | Env Var        | CLI Flag         | Default             |
|----------|----------------|----------------|------------------|---------------------|
| App      | Log level      | `LOG_LEVEL`    | `--log-level`    | `info`              |
| GRPC     | Port           | `GRPC_PORT`    | `--grpc-port`    | `50051`             |
| HTTP     | Port (healthz) | `HTTP_PORT`    | `--http-port`    | `8080`              |
| Postgres | DSN            | `DATABASE_DSN` | `--database-dsn` | —                   |
| Grinex   | API base URL   | `GRINEX_URL`   | `--grinex-url`   | `https://grinex.io` |

See `.env.example` for a complete example.

## Available Now (skeleton)

- **gRPC Health**: standard `grpc.health.v1.Health` service
- **gRPC Reflection**: enabled for development
- **HTTP**: `GET /healthz` on HTTP port
- **Graceful shutdown**: SIGINT/SIGTERM with 15s timeout
- **Log level**: configurable via `LOG_LEVEL` / `--log-level`

## Not Yet Implemented

- `RatesService.GetRates` — proto contract defined, handler not registered
- Grinex HTTP client
- PostgreSQL repository and migrations
- Business logic (topN, avgNM algorithms)
- Unit tests
- Full OTel tracing, Prometheus metrics

## Go Version

go.mod pins `go 1.26.1`. Dockerfile uses `golang:1.26-alpine` (latest 1.26.x patch).