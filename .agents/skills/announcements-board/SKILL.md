---
name: announcements-board
description: The Announcements feed that serves as the Home page of Alumni Connect, rendered as post cards, covering initiatives, job postings, events, reunions, and donations, plus Share, Save, and the Saved Posts screen. Use when building the home screen, post/announcement cards, categories, pinned items, expiry, saving or sharing, or related endpoints.
---

# Announcements Feed (Home)

## Categories
`initiative`, `job`, `event`, `reunion`, `donation`.

## Card (wireframe 02)
- Header: avatar + author name (+ category tag in cinnamon).
- Body: title + text preview (tap to expand/open detail). Optional event date/location and external link button.
- Divider, then the action row. **Phase 1 actions: Share, Save.** Like and Comment are Phase 2 (posts); do not render disabled placeholders.
- **Share** uses `share_plus` with title, short text, and the external link if any (no deep links yet).
- **Save** toggles a bookmark (optimistic). The drawer's **Saved Posts** opens `/saved`, a paginated list of saved announcements, with an empty state.
- Above the feed: category chips (All + 5). Pull-to-refresh, loading skeleton, empty state ("Nothing here yet").

## API
- `GET /v1/announcements?category=&limit=&cursor=`: pinned first, then newest; hides expired; each item includes `is_saved`.
- `GET /v1/announcements/{id}`
- `PUT /v1/announcements/{id}/save`, `DELETE /v1/announcements/{id}/save` (idempotent)
- `GET /v1/me/saved-announcements?limit=&cursor=`
- `POST /v1/announcements`, `PATCH`, `DELETE`: Phase 1 default is **admin-only authoring** (the wireframe has no compose UI). The author name on the card is the poster's name (e.g. "Alumni Association"), so opening posting to alumni later needs no UI change. Authors/admins may edit or delete.
- Fields: title (max 120), body (max 5000, plain text), category, optional `event_date`, `location`, `external_link` (https only), `expires_at`, `pinned`.
- Table `saved_announcements(user_id, announcement_id, created_at)` with PK `(user_id, announcement_id)`; deleting an announcement cascades.

## Rules
- Never render raw HTML; escape text.
- Writes are transactional; surface failures to the user.
- Donations: external link only, no payments in the app.
- No push notifications in Phase 1.
