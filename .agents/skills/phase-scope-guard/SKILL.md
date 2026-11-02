---
name: phase-scope-guard
description: Keeps work inside the current roadmap phase of NIT Hamirpur Alumni Connect. Use before starting ANY new feature, screen, endpoint, table, or dependency, and whenever a request mentions communities, groups, likes, comments, posts, moderation, reporting, blocking, suggestions, recommendations, SOS, marksheets, guest house, campus entry, in-app chat, or push notifications. Also use when deciding whether to add infrastructure such as caches, queues, or extra services.
---

# Phase Scope Guard

## Current phase: Phase 1 (MVP, mid December)
In scope:
- LinkedIn-only signup/login (Welcome, Create Account, Login screens)
- Profile: name, pronouns (optional), batch, branch, current city, company/role, bio, photo
- Manual admin approval of profiles; verified badge
- Search alumni (name/batch/branch/city, filters Batch, Branch, City)
- Connections (request/accept, LinkedIn-style)
- WhatsApp and Email contact buttons (privacy-gated); no in-app chat
- **Home = announcements feed** shown as post cards, with Share and Save
- Saved Posts (saved announcements)
- Inbox: connection requests + in-app notifications (no push)
- Drawer, dark/light theme toggle, Settings & Privacy (Chat ON/OFF, Contact visibility ON/OFF), Log Out
- **Groups tab present as a "Coming soon" placeholder** so navigation doesn't change between releases

## Deferred
- **Phase 2 (mid Feb):** communities/groups (joined, discover, create, join), posts, **Like and Comment**, drawer "Communities" entry, basic moderation, blocking, reporting
- **Phase 3 (aspirational):** profile **Suggestions** row, recommendations, advanced search, institutional services, SOS
- Not planned: in-app chat, push notifications (revisit after Phase 1)

## Procedure
1. Map the request to a Phase 1 bullet. If it maps, proceed.
2. If it belongs to Phase 2/3, **stop**, say so, and propose the smallest Phase 1 alternative or a design note.
3. Do not add tables, endpoints, routes, or packages for future phases. Leave extension points (clean module boundaries) instead.
4. Reject new infrastructure (Redis, queues, extra services, search engines) unless a measured problem is documented. Postgres first.
5. Prefer the simplest design that meets today's requirements; note tradeoffs in the PR description.

## Decision test
"Would the MVP launch be incomplete without this?" If no, defer it.
