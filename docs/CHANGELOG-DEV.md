# Development Changelog

## 2026-11-03: Phase 1 M0 Foundation (Backend)

### Tasks Completed
- Initialized Go module (`github.com/nith-alumni-connect/backend`) with Go 1.26.
- Implemented environment configuration in `internal/config` (PORT, ENV, DATABASE_URL, LOG_LEVEL).
- Structured JSON logging using standard library `log/slog` in `internal/platform/logger`.
- Standard API error envelope `{ "error": { "code": "...", "message": "...", "fields": ... } }` and JSON response helpers in `internal/platform/httputil`.
- Platform middleware in `internal/platform/middleware`: `RequestID` (adds `X-Request-Id` to response and context), structured `RequestLogger`, and `Recoverer` (catches panics and formats 500 JSON error response).
- Health check endpoints in `internal/health`: `GET /healthz` (liveness) and `GET /readyz` (readiness with DB pinger check).
- HTTP server with graceful shutdown in `internal/platform/server` and wired in `cmd/api/main.go`.
- Configured local Postgres 16 in `backend/docker-compose.yml` with health check.
- Created `backend/Makefile` with targets (`run`, `test`, `lint`, `sqlc`, `migrate-*`, `docker-*`).
- Added initial database migration (`00001_create_users_table.sql`) with goose syntax and reversible down migration.
- Configured `backend/sqlc.yaml` and query definitions in `backend/db/queries/users.sql`; installed latest `sqlc` and `goose`, and generated type-safe database access layer in `internal/platform/db`.
- Created OpenAPI 3.1 contract skeleton in `backend/api/openapi.yaml` covering `/healthz` and `/readyz`.
- Created GitHub Actions CI workflow in `.github/workflows/backend.yml`.
- All backend tests (`go test -v -race ./...`) and `go vet ./...` passing cleanly.

### Decisions Made
- Module path chosen: `github.com/nith-alumni-connect/backend`.
- Database access layer: pgx/v5 with connection pooling in `internal/platform/db/pool.go` and sqlc generated code.
- Goose and sqlc installed as global Go binaries in `~/go/bin`.

### Open Questions / Notes
- Flutter SDK is not currently detected on PATH on this environment. Need Flutter SDK installed / PATH configured to proceed with Frontend M0 tasks.

### Next Task
- Frontend M0: `flutter create` in `frontend/`, configure dependencies, Spiced Chai theme system, and 4-tab shell navigation with `/healthz` debug screen.
