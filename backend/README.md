# Backend — DevOps Command Center

Go API built with **Clean Architecture** under `internal/`.

## Layout

```
cmd/server/          Application entrypoint
config/              Viper configuration
internal/
  auth/              JWT + password helpers
  controllers/       HTTP handlers
  database/          GORM connect, migrate, seed
  dto/               Request/response contracts
  middleware/        JWT, RBAC, CORS, rate limit, security
  models/            GORM entities
  repositories/      Persistence interfaces + impl
  routes/            Gin route registration
  scheduler/         Cron collectors
  services/          Business + integrations
  websocket/         Live event hub
migrations/          SQL reference schema
pkg/                 Shared logger, redis, response helpers
tests/               Unit tests
docs/                Swagger metadata
```

The folders requested at the repo root (`controllers/`, `services/`, …) live under `internal/` per Go best practices (non-importable from outside the module root packages).

## Run

```bash
go mod tidy
go run ./cmd/server
```

Config: `config/config.yaml` (override with `DCC_*` env vars).
