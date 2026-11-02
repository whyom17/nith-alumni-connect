---
name: flutter-architecture
description: Feature-first Flutter architecture for the Alumni Connect app using Riverpod, go_router, freezed, clean data/domain/presentation layers, and the app shell (bottom navigation + drawer). Use when creating or restructuring any Flutter screen, feature, provider, route, repository, model, or navigation flow, or when asked where code should live in the frontend.
---

# Flutter Architecture

## Layout
`lib/features/<feature>/{data,domain,presentation}`; shared code in `lib/core`; app wiring and shell in `lib/app`.

## Layers
- **domain**: pure Dart entities and abstract repositories. No Flutter or dio imports.
- **data**: DTOs (`freezed` + `json_serializable`), API calls, repository implementations, DTO to entity mapping.
- **presentation**: screens, widgets, Riverpod notifiers. Widgets read providers; no business logic.

## State (Riverpod)
- `@riverpod` codegen. `AsyncNotifier` for screen state; `AsyncValue.when` renders loading/error/data; design an explicit empty state.
- Paginated lists: notifier holds items + `nextCursor` + `isLoadingMore`.
- Cross-feature access goes through domain interfaces.
- `themeModeProvider` (system/light/dark) persisted via `shared_preferences`.

## Routing (go_router)
Central `routerProvider` with a `redirect` driven by auth/approval state:
- signed out then `/welcome` (also `/signup`, `/login`)
- signed in, profile incomplete then `/onboarding`
- status `pending` then `/approval-pending`
- approved then the shell

Shell = `StatefulShellRoute.indexedStack` with a `NavigationBar` of 4 branches in this order: `/home`, `/groups`, `/search`, `/profile`. Each branch keeps its own navigation stack.
- The `Scaffold` for Home owns the `Drawer` (hamburger) and the Inbox action.
- Non-tab routes pushed above the shell: `/inbox`, `/saved`, `/settings`, `/alumni/:id`, `/profile/edit`.
- `/groups` is a `ComingSoon` placeholder in Phase 1.

## Checklist for a new feature
1. Entity + repository interface (domain)
2. DTO + API + repo impl (data)
3. Notifier + screen + widgets (presentation)
4. Route registration
5. Tests (`flutter-testing`)
6. Check against `wireframe-screen-specs` and `phase-scope-guard`

## Don'ts
No `setState` for shared state, no `BuildContext` across async gaps, no globals as singletons, no feature-to-feature `presentation` imports.
