# USDT Rates Service

gRPC service for fetching USDT exchange rates (ask/bid) from Grinex exchange with PostgreSQL storage and OpenTelemetry tracing.

> **Note:** Module path `github.com/example/grinex-rates-service` is a placeholder.
> Replace with the actual module path before publishing.

## Quick Start

```bash
make build
docker-compose up -d                                  # starts postgres
docker-compose run --rm app ./app                     # migrates + starts gRPC server
```

The Docker entrypoint automatically runs `migrate up` before starting the server.
To run migrations manually: `docker-compose run --rm app ./app migrate up`.

Run tests:

```bash
make test
```

## Development

For development with published ports (gRPC 50051, HTTP 8080):

```bash
make run                                              # postgres + app with ports
```

To verify the service is working (requires [grpcurl](https://github.com/fullstorydev/grpcurl)):

```bash
grpcurl -plaintext -d '{"top_n":{"n":1}}' localhost:50051 rates.v1.RatesService/GetRates
```

> `docker-compose run` does not publish ports by default.
> Use `make run` or `docker-compose run --rm --service-ports app ./app` to expose ports.

## Project Structure

```
main.go                              # Thin bootstrap
cmd/                                 # CLI commands (grpc, migrate up/down)
internal/
  app/                               # Composition root, wiring, lifecycle
  config/                            # Section-based config (env + flags)
  transport/grpc/                    # gRPC handler + server (thin)
  service/rates/                     # Business logic, algorithms, domain types
  client/grinex/                     # Grinex HTTP client (resty)
  repo/rates/                        # Repository interface
  repo/rates/postgres/               # PostgreSQL implementation (sqlc)
  migration/postgres/                # SQL migrations (golang-migrate)
  observability/otel/                # OpenTelemetry tracing
api/proto/rates/v1/                  # Protobuf definitions
gen/rates/v1/                        # Generated proto Go code
```

## Commands

```bash
make build          # Build binary
make test           # Run tests
make lint           # Run golangci-lint
make docker-build   # Build Docker image
make run            # Start postgres + app with published ports (dev)
make gen            # Generate proto + sqlc code
make gen-proto      # Generate proto code only
make gen-sqlc       # Generate sqlc code only
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

OpenTelemetry tracing is opt-in via standard env var `OTEL_EXPORTER_OTLP_ENDPOINT`.

See `.env.example` for a complete example.

## API

### gRPC: `RatesService.GetRates`

Fetches USDT ask and bid prices from Grinex order book using a specified algorithm:

- **TopN** — value at position N in the order book array (1-based)
- **AvgNM** — average of entries in range [N, M] (1-based, inclusive)

The result is persisted to PostgreSQL on each call.

Request (proto `oneof`):
```
{"top_n": {"n": 3}}
{"avg_nm": {"n": 1, "m": 5}}
```

Response:
```json
{"ask": "80.86", "bid": "80.73", "timestamp": "2026-04-05T18:50:00Z"}
```

### Health Checks

- **gRPC Health**: standard `grpc.health.v1.Health` service
- **HTTP**: `GET /healthz` on HTTP port

## Implementation Status

| Component | Status |
|-----------|--------|
| gRPC GetRates handler | Implemented |
| Grinex HTTP client (resty) | Implemented |
| PostgreSQL repository (pgx/v5, sqlc) | Implemented |
| SQL migrations (golang-migrate) | Implemented |
| TopN / AvgNM algorithms | Implemented |
| Service orchestration | Implemented |
| Unit tests (client, algorithms, service) | Implemented |
| OpenTelemetry tracing (gRPC + client) | Implemented |
| Graceful shutdown | Implemented |
| Config: env + flags | Implemented |
| Prometheus metrics (`/metrics`) | Implemented |

## Go Version

go.mod pins `go 1.26.1`. Dockerfile uses `golang:1.26-alpine`.