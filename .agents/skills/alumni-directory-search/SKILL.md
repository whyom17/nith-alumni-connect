---
name: alumni-directory-search
description: Search Alumni screen and API for Alumni Connect, with a text search plus Batch, Branch, and City filters, Connect buttons in result rows, Postgres indexing, and keyset pagination. Use when building or modifying search, filters, alumni list rows, or anything about finding alumni by name, batch, branch, city, company, or industry.
---

# Search Alumni

## UI (wireframe 03)
- Title "Find Alumni". Search field placeholder "Search by name, batch, branch, ...".
- Filters row: chips **Batch**, **Branch**, **City**. Tapping a chip opens a bottom sheet picker (multi-select); selected chips turn cinnamon and show the count.
- Result row: avatar (initials fallback), **name**, `2016 | ECE`, **city**, trailing **Connect** button, divider. Tapping the row opens `/alumni/:id`.
- Connect button state comes from `connection_status`: `none` Connect, `pending_outgoing` Pending (disabled), `pending_incoming` Accept, `accepted` Connected. Update optimistically and roll back on failure.
- Debounce typing 300ms and cancel stale requests. States: loading skeleton, empty ("No alumni match your filters"), error with retry.
- Industry is supported by the API but has no chip in the Phase 1 wireframe; don't add one without asking.

## API
`GET /v1/directory?q=&batch=&branch=&city=&industry=&limit=20&cursor=`
- Repeated params for multi-value filters.
- Row DTO: `id, name, batch, branch, city, photo_url, is_verified, connection_status` (company/role may be included but is not shown in the row). **No** phone, email, roll number.
- `GET /v1/directory/filters` returns the distinct `batches`, `branches`, `cities` of approved profiles, so pickers are data-driven and never hardcoded.
- Only approved requesters; only approved profiles are returned. Exclude the requester.

## Backend
- Name/company text search with `pg_trgm` on lowercased columns plus GIN indexes; also match batch/branch tokens in `q`. No Elasticsearch.
- B-tree indexes on batch, branch, city, industry.
- Keyset pagination on `(name, id)`; cap `limit` at 50.
- Normalise `current_city` on write (trim, collapse spaces, title-case) so "Delhi" and "delhi" group together. Consider a city autocomplete later.
- Target p95 under about 300ms at 10k profiles; measure with `EXPLAIN ANALYZE` before caching.
- Design the query so a Phase 2 block-exclusion join is easy to add.

## Out of scope
Profile "Suggestions", recommendations, advanced/semantic search (Phase 3).
