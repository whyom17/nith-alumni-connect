---
name: go-backend-architecture
description: Modular-monolith Go architecture for the Alumni Connect API: chi router, handler/service/repository layering, dependency wiring, middleware, config, logging, and error handling. Use when creating or changing any Go package, module, HTTP handler, middleware, service, or when deciding where backend code belongs.
---

# Go Backend Architecture

## Shape
One binary (`cmd/api`), one Postgres. Feature modules under `internal/<module>`: `handler.go`, `service.go`, `repository.go`, `types.go`.

## Rules
- **Handlers:** decode and validate, call service, encode. No SQL, no business rules.
- **Services:** business and authorization rules; depend on repository and other modules' **service interfaces** (define the interface in the consumer package).
- **Repositories:** wrap `sqlc` generated code; translate `pgx` errors (not found, unique violation) to domain errors.
- Constructor injection; wire everything in `cmd/api/main.go`. No globals, no `init()` side effects.
- Pass `context.Context` first everywhere; set request timeouts.

## Middleware order
request ID, recover, structured logging (`slog`), CORS (if needed), rate limit, auth (JWT), `RequireApproved` / `RequireAdmin`.

## Errors
Domain errors (`ErrNotFound`, `ErrForbidden`, `ErrConflict`, `ErrValidation`) mapped in one place to the JSON error envelope from `AGENTS.md`. Never return raw DB or internal errors to clients. Log with request ID.

## Extensibility without over-building
Phase 2 modules (communities, posts, moderation) will be added as new packages depending only on `users` and `profiles` service interfaces. Do not create them early.

## Tests
Service unit tests with fake repositories; repository tests against a real Postgres (`testcontainers-go`); handler tests via `httptest`.
