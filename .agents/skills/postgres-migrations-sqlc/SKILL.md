---
name: postgres-migrations-sqlc
description: PostgreSQL schema design, goose migrations, indexes, and sqlc query authoring for the Alumni Connect backend. Use when creating or altering tables, writing SQL queries, adding indexes, seeding data, or generating repository code.
---

# Postgres, Migrations and sqlc

## Migrations (goose)
- `make migrate-new name=add_x` creates `db/migrations/NNNN_add_x.sql` with `-- +goose Up` / `Down`.
- Every migration is reversible and idempotent in intent. Never edit an applied migration; add a new one.
- Use `uuid` primary keys (`gen_random_uuid()`), `timestamptz` with `default now()`, explicit `NOT NULL`, `CHECK` constraints for enums (or Postgres enums with care).
- Foreign keys with sensible `ON DELETE` rules.

## Core tables (Phase 1)
`users`, `profiles`, `privacy_settings`, `contact_details`, `verification_records`, `connections`, `announcements`, `saved_announcements`, `notifications`, `refresh_tokens`.
- `connections`: unique pair constraint (`LEAST/GREATEST` unique index) to prevent duplicates or reversed duplicates; `CHECK (requester_id <> addressee_id)`.
- `verification_records` holds roll number and evidence; it is a separate table so it can never be joined into public queries by accident.

## Indexes
`profiles(batch)`, `(branch)`, `(current_city)`, `(industry)`; `pg_trgm` GIN on `lower(name)` and `lower(company)`; `announcements(category, created_at DESC)`; `connections(addressee_id, status)`; `notifications(user_id, read_at, created_at DESC)`; `saved_announcements(user_id, created_at DESC)`.

## sqlc
- Queries in `db/queries/<module>.sql`, named `-- name: ListAlumni :many`.
- Select explicit columns (no `SELECT *`) and use separate list vs detail queries.
- Keyset pagination on `(created_at, id)` or `(name, id)`; never OFFSET for large lists.
- Run `make sqlc` and commit generated code.

## Don'ts
No dynamic SQL string concatenation. If filters are optional, use `sqlc.narg()` with `IS NULL OR` patterns or a vetted query builder.
