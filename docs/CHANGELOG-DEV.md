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
- Resolved: Flutter SDK installed and available on PATH (`Flutter 3.47.9`, `Dart 3.13.7`).

### Next Task
- Frontend M0 Foundation complete. Ready for Phase 1 M1: Auth and session.

## 2026-11-04: Phase 1 M0 Foundation (Frontend)

### Tasks Completed
- Added required dependencies to `frontend/pubspec.yaml` (Riverpod, go_router, dio, shared_preferences, flutter_secure_storage, freezed annotations, mocktail) and ran `flutter pub get`.
- Built Spiced Chai design system tokens in `core/theme/app_colors.dart` (Cream `#FDFBD4`, Cinnamon `#D47E30`, Brown `#8D5A2B`, Nutmeg `#825E34`, Ink `#3B2712`, Espresso `#1C140D`, etc.).
- Implemented `lightTheme` and `darkTheme` in `core/theme/app_theme.dart` with Material 3 component themes (FilledButton, Card, Chips, NavigationBar, Drawer, AppBar, Dividers).
- Created `themeModeProvider` in `core/theme/theme_mode_provider.dart` using Riverpod Notifier, persisted via `shared_preferences`, with light/dark toggle.
- Created `Dio` client provider in `core/network/api_client.dart` reading `API_BASE_URL` with typed domain error mapping to `AppFailure` (`NetworkFailure`, `UnauthorizedFailure`, `ForbiddenFailure`, `ValidationFailure`, `ServerFailure`).
- Created `HealthRepository` and `healthCheckProvider` in `core/network/health_repository.dart` to probe `/healthz` and `/readyz`.
- Built 4-tab application shell (`AppShell`) with `go_router` `StatefulShellRoute.indexedStack` (Home, Groups, Search, Profile) and drawer.
- Built `HomeScreen` with hamburger, Inbox action, announcements placeholder, and live `/healthz` diagnostic card.
- Built Phase 1 placeholders: `GroupsScreen` (with `ComingSoon` component), `SearchScreen`, `ProfileScreen`, `InboxScreen`, `SavedPostsScreen`, and `SettingsScreen`.
- Built `AppDrawer` with user header, light/dark mode icon toggle, Saved Posts, Settings & Privacy, and Log Out in error token.
- Created GitHub Actions CI workflow in `.github/workflows/frontend.yml`.
- Wrote unit and widget tests: `dart analyze` has 0 issues, and all 10 `flutter test` cases pass.

### Decisions Made
- `HomeScreen` scaffold directly owns the `AppDrawer` and hamburger trigger as specified in `flutter-architecture`.
- Default `API_BASE_URL` set to `http://10.0.2.2:8080` on mobile and `http://localhost:8080` on web/desktop via `const String.fromEnvironment('API_BASE_URL')`.

### Next Task
- Milestone 0 is complete. Propose Phase 1 M1: Auth and session.
