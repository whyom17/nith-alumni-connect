---
name: flutter-api-integration
description: How the Flutter app talks to the Go backend: dio client, JWT auth interceptor with refresh, error mapping, cursor pagination, DTOs, and image upload. Use when adding or changing any API call, repository, DTO, token handling, pagination, file or photo upload, or error handling in the frontend.
---

# Flutter API Integration

## Client
- Single `Dio` provider in `core/network`, base URL from `--dart-define=API_BASE_URL`.
- Interceptor attaches `Authorization: Bearer <access>`; on 401 it calls refresh once (queue concurrent requests), retries, else signs out.
- Tokens in `flutter_secure_storage` only.

## Contract
- Backend defines truth in `backend/api/openapi.yaml`. Mirror field names (`snake_case`) via `@JsonKey` or `fieldRename: FieldRename.snake`.
- Errors arrive as `{ "error": { code, message, fields } }`. Map to `AppFailure` subclasses: `Unauthorized`, `Forbidden` (e.g. account not approved), `Validation(fields)`, `NotFound`, `Network`, `Unknown`.
- Pagination: `?limit=20&cursor=` returns `{ data, next_cursor }`.

## Rules
- Repositories return `Result`/throw `AppFailure`; never leak `DioException` to the UI.
- Search inputs are debounced (300ms) and previous requests cancelled via `CancelToken`.
- Photo upload: request a pre-signed URL from the API, upload directly, then PATCH the profile with the resulting key. Resize and compress on device (max ~1080px).
- A hidden field (e.g. private phone) is simply absent. DTO fields for such data must be nullable.
- Never log request bodies containing tokens or phone numbers.
