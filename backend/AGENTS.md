# AGENTS.md: Backend (Go)

REST API for **NIT Hamirpur Alumni Connect**. Read the root `README.md` for product context and phase scope.

## Guiding principle
**Modular monolith.** One Go binary, one PostgreSQL database. No microservices, message queues, or caches unless a measured problem demands it. Cost must stay low while usage is small (hundreds to tens of thousands of alumni).

## Stack
- Go (latest stable), standard `net/http` with `chi` router
- PostgreSQL 16 via `pgx/v5`; queries with `sqlc`; migrations with `goose`
- Auth: Sign in with LinkedIn (OIDC), then app-issued short-lived JWT access token + rotating refresh token
- Config: environment variables, `.env` for local only
- Logging: `log/slog` (JSON); validation: `go-playground/validator`
- Profile photos: S3-compatible bucket with pre-signed uploads
- Tests: `testing`, `testify`, `testcontainers-go` for DB integration

## Project structure
```
backend/
├── cmd/api/main.go            # Wiring only
├── internal/
│   ├── config/
│   ├── platform/              # db, http server, middleware, logger, storage
│   ├── auth/                  # LinkedIn OIDC, JWT, sessions
│   ├── users/                 # accounts, roles
│   ├── profiles/              # alumni profile + privacy settings
│   ├── verification/          # admin approval workflow
│   ├── directory/             # search + filters
│   ├── connections/
│   ├── announcements/         # feed + saved announcements (bookmarks)
│   ├── notifications/         # inbox items + unread count
│   └── admin/
│   # Phase 2 (do not create yet): communities/, posts/, moderation/, reports/
├── db/
│   ├── migrations/
│   └── queries/               # sqlc .sql files
├── api/openapi.yaml           # Source of truth for the contract
├── Makefile
├── Dockerfile
└── docker-compose.yml         # Local Postgres
```
Each module has `handler.go` (HTTP), `service.go` (business rules), `repository.go` (DB via sqlc), `types.go`. Handlers never touch the DB. Modules talk through service interfaces, never by importing another module's repository. This keeps community features changeable without breaking authentication.

## Domain model (Phase 1)
- `users`: id, linkedin_sub, email, role (`alumnus`, `admin`), status (`pending`, `approved`, `rejected`, `suspended`)
- `profiles`: user_id, name, pronouns (optional, max 20 chars), batch (year), branch, current_city, country, company, role_title, industry, bio, photo_url
- `privacy_settings`: user_id, chat_enabled (default true; WhatsApp allowed), contact_visible (default false; covers phone **and** email)
- `contact_details`: user_id, phone (private by default)
- `verification_records`: user_id, submitted_info, decision, decided_by, decided_at, reason (**never** exposed publicly; roll numbers live only here)
- `connections`: requester_id, addressee_id, status (`pending`, `accepted`, `declined`), timestamps
- `announcements`: id, author_id, category (`initiative`, `job`, `event`, `reunion`, `donation`), title, body, created_at, expires_at, pinned
- `saved_announcements`: user_id, announcement_id, created_at (PK on both ids)
- `notifications`: id, user_id, type, payload (jsonb, ids and names only), read_at, created_at
- Phase 2 role: `community_moderator` (must be verified and admin-approved).

## API conventions
- Base path `/v1`, JSON, `snake_case` fields, RFC 3339 UTC timestamps
- Cursor pagination: `?limit=20&cursor=...` returns `{ "data": [...], "next_cursor": "..." }`
- Errors: `{ "error": { "code": "string", "message": "string", "fields": {} } }` with correct status codes
- Update `api/openapi.yaml` in the same change as any endpoint change
- Only **approved** users may access directory, announcements, and connections. Pending users can only access their own profile and status.

## Phase 1 endpoints (summary)
- Auth: `POST /v1/auth/linkedin` (returns `is_new_user`, `status`, `profile_complete`), `/auth/refresh`, `/auth/logout`
- Me: `GET /v1/me`, `PATCH /v1/me/profile`, `POST /v1/me/photo-upload-url`, `POST /v1/me/verification`, `GET|PATCH /v1/me/privacy`, `GET /v1/me/saved-announcements`
- Alumni: `GET /v1/alumni/{id}`, `GET /v1/directory`, `GET /v1/directory/filters`
- Connections: `POST /v1/connections`, `POST /v1/connections/{id}/accept|decline`, `DELETE /v1/connections/{id}`, `GET /v1/connections`
- Announcements: `GET /v1/announcements`, `GET /{id}`, `POST|PATCH|DELETE` (admin), `PUT|DELETE /{id}/save`
- Inbox: `GET /v1/inbox`, `GET /v1/inbox/unread-count`, `POST /v1/inbox/{id}/read`, `POST /v1/inbox/read-all`
- Admin: `GET /v1/admin/verifications`, `POST /v1/admin/verifications/{id}/decision`

## Contact exposure rules
- `GET /v1/alumni/{id}` includes a `contact` object only if the viewer is an **accepted connection** of the target and the target allows it:
  - `contact.email` needs `contact_visible = true`
  - `contact.phone` needs `chat_enabled = true` AND `contact_visible = true`
- Otherwise omit `contact` entirely. Directory/list endpoints never return phone or email.
- Batch years are not restricted to the past until the product owner decides whether current students may join (wireframes show a 2028 batch).

## Security and privacy (non-negotiable)
- Authorization is enforced in the service layer, not only in middleware
- Privacy flags are enforced **server-side**: when `contact_visible=false`, the phone is omitted from the response entirely; same for chat
- Roll numbers and verification records are never returned from public endpoints or logged
- Parameterized queries only (sqlc); no string-built SQL
- Rate-limit auth and search endpoints; cap page size at 50
- Verify LinkedIn `state`, `nonce`, ID token signature, `aud` and `iss`
- Secrets from env only; never commit `.env`
- Verification decisions and announcements are written in transactions and must not silently fail

## Performance
- Index directory filter columns (`batch`, `branch`, `current_city`, `industry`) plus a trigram/full-text index for name and company
- Avoid N+1; return list DTOs, not full profiles
- Add a cache only after measuring a slow endpoint

## Commands
```bash
make run              # run API
make test             # go test ./... -race
make lint             # golangci-lint run
make sqlc             # regenerate DB code
make migrate-up / migrate-down / migrate-new name=add_x
```

## Definition of done
- `go vet`, `golangci-lint`, and `go test ./... -race` pass
- Migration is reversible; `openapi.yaml` updated
- Tests cover happy path, authorization failure, and privacy-flag behaviour
- No Phase 2/3 feature added (see `phase-scope-guard` skill)

## Skills to use
`go-backend-architecture`, `postgres-migrations-sqlc`, `linkedin-auth-flow`, `admin-approval-workflow`, `alumni-directory-search`, `connections-and-privacy`, `announcements-board`, `notifications-inbox`, `api-contract-openapi`, `privacy-and-security`, `deployment-low-cost`, `phase-scope-guard`, `git-workflow`.
