---
name: api-contract-openapi
description: Maintains the REST API contract in backend/api/openapi.yaml, with naming, pagination, error envelope, versioning, and status code conventions shared by Go and Flutter. Use when adding, changing, or documenting any endpoint, request/response shape, or error, or when the frontend and backend disagree about a field.
---

# API Contract (OpenAPI)

## Source of truth
`backend/api/openapi.yaml` (OpenAPI 3.1). Change the spec in the same commit as the handler and the Flutter DTO.

## Conventions
- Prefix `/v1`; nouns, plural; `snake_case` JSON; ids are UUID strings; times RFC 3339 UTC.
- Pagination: `limit` (default 20, max 50), `cursor`; response `{ data: [], next_cursor }`.
- Errors: `{ error: { code, message, fields? } }`.
  - 400 validation, 401 unauthenticated, 403 forbidden (including `account_not_approved`), 404, 409 conflict, 422 semantic, 429 rate limited, 500.
- `PATCH` for partial updates; `POST` for actions (`/accept`).
- Never expose: roll numbers, verification evidence, private phone numbers, refresh tokens.
- Define reusable schemas under `components` (`AlumniSummary` (search row), `AlumniProfile` (with optional `Contact`), `Announcement` (with `is_saved`), `Connection`, `Notification`, `PrivacySettings`, `Error`).
- Breaking changes require a new version path, not silent edits.

## Workflow
1. Edit spec, 2. validate (`make openapi-lint`), 3. implement handler, 4. update Flutter DTO and repository, 5. add tests on both sides.
