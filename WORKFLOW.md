# WORKFLOW.md: Development Plan for the Coding Agent

This file tells you (the coding agent) **what to build, in what order, and how to work**. Read it at the start of every session, together with `README.md`, `PRD.md`, and the `AGENTS.md` in the folder you are working in (`frontend/` or `backend/`).

Last updated: 3 Nov 2026. Targets: **Phase 1 mid December 2026**, **Phase 2 mid February 2027**.

## 0. About the developer
The developer knows web development well but is **new to Flutter and mobile**. So:
- Explain Flutter/Dart concepts briefly the first time they appear, using web analogies (widget = component, `ref.watch` = `useState`-style subscription).
- Prefer simple, readable code over clever abstractions.
- Build **one task at a time**, run it, and show how to verify it. Don't generate a whole milestone in one go.

## 1. Rules of engagement

1. **Pick the next unchecked task** in section 5 and work only on it. Do not start the next milestone until the current one is accepted.
2. **Stay in phase.** Before adding anything, check `phase-scope-guard`. If something belongs to a later phase, stop and say so.
3. **Contract first.** For any API change: update `backend/api/openapi.yaml`, then backend, then the Flutter DTO and repository. Frontend may use a fake repository while the backend isn't ready.
4. **Use the skills.** Each task lists the skills to read before starting (`.agents/skills/<name>/SKILL.md`). Always read `wireframe-screen-specs` before building a screen and `privacy-and-security` before touching personal data.
5. **Follow the wireframes** (`docs/wireframes/`). Layout is binding; colours come from the theme tokens.
6. **Small commits** using Conventional Commits (`git-workflow` skill), one logical change each. Never commit secrets or `.env`.
7. **Quality gates before saying "done":**
   - Backend: `make lint test` passes; migration reversible; OpenAPI updated.
   - Frontend: `flutter analyze` clean and `flutter test` passing; loading/error/empty states handled; checked in light and dark mode.
8. **Update the progress tracker** (section 4) in this file when a task or milestone is finished.
9. **Stop and ask** instead of guessing when: a requirement is ambiguous, a "blocking question" (section 3) isn't answered, a new dependency or infrastructure is needed, or a change touches auth or privacy rules.
10. **End every milestone** with a short report: what was built, how to run/verify it, what was deferred, and any new questions. Wait for approval.

## 2. Architecture reminders
- Backend: Go **modular monolith**, one Postgres, no queues/caches/microservices. Modules: auth, users, profiles, verification, directory, connections, announcements, notifications, admin.
- Frontend: Flutter, Riverpod, go_router, feature-first folders.
- Privacy is enforced **server-side**; hidden data is absent from responses.
- Roll numbers and verification records never leave admin endpoints.

## 3. Blocking questions (ask the developer before the milestone that needs them)

| Before | Question |
|---|---|
| M1 | LinkedIn developer app created? Redirect/app-link scheme and package name decided? |
| M2 | What proves alumni status (roll number only, college email, document)? Are current students allowed (the wireframe shows a 2028 batch)? Who are the first admins? |
| M4 | Connection model: mutual acceptance (assumed)? |
| M5 | Who may post announcements in Phase 1? (assumed admin-only) |
| M8 | Android only at launch? Hosting provider chosen? Privacy policy and terms ready? |
| P2-M1 | Final naming: "Groups" vs "Communities"? Who may create groups? |

If the developer hasn't answered, proceed with the assumption stated in `PRD.md` and flag it.

## 4. Progress tracker
Mark `[x]` when accepted by the developer.

**Phase 1**
- [ ] M0 Foundation
- [ ] M1 Auth and session
- [ ] M2 Profile, onboarding, verification
- [ ] M3 Search and profile view
- [ ] M4 Connections and contact privacy
- [ ] M5 Announcements feed
- [ ] M6 Inbox and notifications
- [ ] M7 Settings, theme, polish
- [ ] M8 Hardening and release

**Phase 2**
- [ ] P2-M1 Groups
- [ ] P2-M2 Posts, comments, likes
- [ ] P2-M3 Moderation, blocking, reporting
- [ ] P2-M4 Release

**Phase 3** (aspirational, each item needs explicit approval to start)
- [ ] Recommendations and Suggestions row
- [ ] Advanced search
- [ ] Institutional services
- [ ] SOS

## 5. Phase 1: MVP (target mid December 2026)

### M0 Foundation (target: first week of Nov)
**Goal:** both projects run, are linted, tested, and wired to a database. A "hello" request travels from phone to backend.

Backend
- [x] `go mod init`, `cmd/api/main.go`, chi router, config from env, `slog` JSON logging
- [x] `GET /healthz` and `GET /readyz` (DB ping)
- [x] `docker-compose.yml` with Postgres; `Makefile` (`run`, `test`, `lint`, `sqlc`, `migrate-*`)
- [x] goose + sqlc set up; first migration (`users`); `.env.example`
- [x] Error envelope and request-ID/recover/logging middleware
- [x] `api/openapi.yaml` skeleton; CI workflow (vet, lint, test)

Frontend
- [ ] `flutter create` in `frontend/`; add dependencies from `frontend/AGENTS.md`
- [ ] Folder structure, `AppColors`, light + dark `ThemeData`, `themeModeProvider`
- [ ] go_router with the 4-tab shell (Home, Groups, Search, Profile) and drawer; placeholder screens
- [ ] dio client reading `API_BASE_URL`; call `/healthz` and show the result on a debug screen
- [ ] CI workflow (`flutter analyze`, `flutter test`)

Skills: `go-backend-architecture`, `postgres-migrations-sqlc`, `flutter-architecture`, `flutter-theme-spiced-chai`, `deployment-low-cost`, `git-workflow`
**Done when:** `make run` + `flutter run` on a real phone shows a successful health check, themes switch, CI is green.

### M1 Auth and session (target: early to mid Nov)
**Goal:** a user can sign up or log in with LinkedIn and stay logged in.

Backend
- [ ] Migrations: `refresh_tokens`; users get `linkedin_sub`, `status`, `role`
- [ ] `POST /v1/auth/linkedin` (code exchange, ID-token verification, upsert user, `is_new_user`), `/auth/refresh` (rotating), `/auth/logout`
- [ ] JWT middleware; `RequireApproved`, `RequireAdmin`; rate limit on `/auth/*`
- [ ] Tests: token rotation and reuse detection, invalid ID token

Frontend
- [ ] Welcome, Create Account, Login screens (`wireframe-screen-specs`)
- [ ] LinkedIn sign-in via app link; secure token storage; dio interceptor with refresh
- [ ] Router redirects by auth state; Log Out clears everything

Skills: `linkedin-auth-flow`, `flutter-api-integration`, `privacy-and-security`, `api-contract-openapi`
**Done when:** fresh install, sign up, kill the app, reopen and still signed in; logout works; refresh works after access token expiry.

### M2 Profile, onboarding, verification (target: mid to late Nov)
**Goal:** new users complete a profile, submit verification, and an admin approves them.

Backend
- [ ] Migrations: `profiles` (incl. pronouns), `privacy_settings`, `contact_details`, `verification_records`
- [ ] `GET /v1/me`, `PATCH /v1/me/profile`, `POST /v1/me/photo-upload-url`, `POST /v1/me/verification`
- [ ] Admin: `GET /v1/admin/verifications`, `POST .../{id}/decision` (single transaction, audit rows)
- [ ] City normalisation; photo limits (type, size, EXIF stripped); seed first admin via CLI
- [ ] Tests: status transitions, pending users blocked from approved-only routes

Frontend
- [ ] Profile completion form (validation, photo picker and upload); Approval pending screen; rejected state with reason and resubmit
- [ ] Edit profile screen; own Profile tab (view)
- [ ] Router: incomplete profile to `/onboarding`, pending to `/approval-pending`

Admin tooling: minimal (API plus a simple CLI or script). Don't build an admin UI in the app unless asked.
Skills: `admin-approval-workflow`, `privacy-and-security`, `postgres-migrations-sqlc`, `flutter-api-integration`
**Done when:** a test account goes from signup to approved, sees the verified badge, and the verification record is unreachable from non-admin endpoints.

### M3 Search and profile view (target: late Nov)
**Goal:** approved users can find alumni and open their profiles.

Backend
- [ ] `pg_trgm` and B-tree indexes; `GET /v1/directory` (q, batch, branch, city, industry, keyset pagination, cap 50)
- [ ] `GET /v1/directory/filters`; `GET /v1/alumni/{id}` (no `contact` yet)
- [ ] Tests incl. only-approved and requester-excluded; check `EXPLAIN ANALYZE` on 10k seeded profiles (add a seed script)

Frontend
- [ ] Search screen: field, Batch/Branch/City chips with bottom-sheet pickers, rows with Connect button (static for now), debounce, pagination, states
- [ ] Other-alumni Profile screen (`/alumni/:id`)

Skills: `alumni-directory-search`, `wireframe-screen-specs`, `flutter-architecture`
**Done when:** searching with 3 filters returns quickly on seeded data, the screen handles empty/error/slow network, and profiles open from rows.

### M4 Connections and contact privacy (target: late Nov to early Dec)
**Goal:** alumni can connect, and contact buttons respect privacy.

Backend
- [ ] `connections` migration (unique pair, no self); endpoints for request/accept/decline/remove/list
- [ ] `GET/PATCH /v1/me/privacy`; `contact` object on `GET /v1/alumni/{id}` per the exposure rules
- [ ] `connection_status` in directory and profile responses; daily request limit
- [ ] Tests: every row of the contact exposure matrix (connected/not, chat on/off, contact on/off)

Frontend
- [ ] Connect/Pending/Accept/Connected states in Search rows and Profile (optimistic updates)
- [ ] WhatsApp (`wa.me`) and Email buttons shown only when `contact` data is present; Connections list from Profile; Privacy toggles (can live in the Settings screen shell)

Skills: `connections-and-privacy`, `privacy-and-security`, `flutter-testing`
**Done when:** the full matrix of visibility cases is verified by tests; no phone/email in any list or log.

### M5 Announcements feed (target: early Dec)
**Goal:** Home shows the announcements feed; users can share and save.

Backend
- [ ] Migrations: `announcements`, `saved_announcements`
- [ ] List (category filter, pinned first, hide expired, `is_saved`), detail, admin create/update/delete, save/unsave, saved list
- [ ] Tests incl. admin-only writes, idempotent save

Frontend
- [ ] Home post cards (Share, Save only), category chips, pull-to-refresh, pagination, states
- [ ] Saved Posts screen; share via `share_plus`; external links via `url_launcher` (https only)

Skills: `announcements-board`, `wireframe-screen-specs`
**Done when:** admin creates posts via API/CLI, they appear on Home, can be saved and found in Saved Posts. No Like/Comment UI exists.

### M6 Inbox and notifications (target: early Dec)
**Goal:** users see connection requests and key notifications.

Backend
- [ ] `notifications` migration; write notifications in the same transaction as connection and verification events
- [ ] Inbox list, unread count, mark read, read-all; 90-day cleanup
- [ ] Tests: ownership, transactional creation

Frontend
- [ ] Inbox screen (Requests section plus notifications); unread badge on Home app bar; navigation from items

Skills: `notifications-inbox`, `connections-and-privacy`
**Done when:** sending a request creates a badge and an inbox item for the recipient; accepting notifies the sender.

### M7 Settings, theme, polish (target: first half of Dec)
**Goal:** the app feels complete and consistent with the wireframes.

- [ ] Settings & Privacy screen (Privacy, Account, About; red Log Out); drawer complete (user header, theme toggle, Saved Posts, Settings, Log Out)
- [ ] Groups tab "Coming soon" state
- [ ] Empty/error/loading polish everywhere; accessibility pass (contrast, tap targets, text scale 1.3); app icon and splash
- [ ] Golden tests for key screens in light and dark
- [ ] Compare every screen with `docs/wireframes/` and list any intentional differences

Skills: `flutter-theme-spiced-chai`, `wireframe-screen-specs`, `flutter-testing`
**Done when:** a screen-by-screen review against the wireframes has no unexplained differences.

### M8 Hardening and release (target: mid Dec)
**Goal:** a safe pilot launch to a small group, then Phase 1 release.

- [ ] Security review with the `privacy-and-security` checklist; rate limits verified; dependency audit
- [ ] Dockerfile, deployment to the chosen host, HTTPS, env secrets, daily Postgres backups and a **tested restore**
- [ ] Seed data removed from production; first admins created; LinkedIn production app and redirect URIs set
- [ ] Android release: flavors (`dev`/`prod`), signing key kept outside the repo, `flutter build appbundle --release`; performance check in release mode on a real device
- [ ] Pilot with 10 to 20 alumni; fix issues; privacy policy and terms linked in About
- [ ] Basic error tracking and uptime check
- [ ] Release notes; tag `v1.0.0`

Skills: `deployment-low-cost`, `privacy-and-security`, `git-workflow`
**Done when:** pilot users complete signup, approval, search, connect and read announcements without blocking issues.

## 6. Phase 2 (target mid February 2027)
Start only after Phase 1 is released and the developer says so. Create the Phase 2 modules then, not before. Add new skills for communities and moderation when this phase starts.

### P2-M1 Groups
- Backend: `communities`, `community_members`; create (verified alumni), join/leave, discover/search, joined list; moderator role (verified plus admin-approved); types: chapter, interest, special
- Frontend: Groups tab replaces the placeholder: "Joined Groups", "Discover new groups", group detail, join/leave; drawer "Communities" entry
- Done when: a user can discover, join, and leave a group; creation is limited to verified alumni

### P2-M2 Posts, comments, likes
- Backend: posts in communities, comments, likes; pagination; edit/delete own
- Frontend: Like and Comment actions on post cards; community feed; notifications for replies (extend inbox)
- Done when: members can post and comment; non-members cannot

### P2-M3 Moderation, blocking, reporting
- Backend: `reports`, `blocks`; report content/user; moderator and admin review queue; hide/delete; block exclusions added to directory, connections, feeds
- Frontend: report and block actions; blocked-users list in Settings
- Done when: blocked users can't see or contact each other; reports reach moderators

### P2-M4 Release
- Regression pass of Phase 1; load check at expected scale; update PRD/README; tag `v2.0.0`

## 7. Phase 3 (aspirational; do not start without explicit approval)
- Recommendations (alumni and communities) and the profile **Suggestions** row
- Advanced search
- Institutional services (marksheets, guest house booking, campus entry)
- SOS
- Candidates to revisit: in-app messaging, push notifications

Each item starts with a short design note (scope, data needed, privacy impact, cost), not code.

## 8. Standard task loop

1. Read this file, the relevant `AGENTS.md`, and the skills listed for the task.
2. Restate the task and any assumptions in two or three lines.
3. Update the OpenAPI spec if the API changes.
4. Implement the smallest slice that works; add tests.
5. Run the quality gates (rule 7).
6. Show how to run and verify it (exact commands, expected result).
7. Commit; tick the tracker; propose the next task.

## 9. Session hand-off
At the end of a session, append a short entry to `docs/CHANGELOG-DEV.md` (create if missing): date, tasks completed, decisions made, open questions, next task. This keeps context for the next session.
