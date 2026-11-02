# AGENTS.md: Frontend (Flutter)

Mobile app for **NIT Hamirpur Alumni Connect** (in-app name: **NITH Alumni Connect**). Read the root `README.md` for product context and phase scope.

## Source of truth for UI
Wireframes live in `docs/wireframes/` (01 auth, 02 home + sidebar + settings, 03 profile + search + groups). The `wireframe-screen-specs` skill turns them into screen-by-screen rules. **Follow the layout, not the wireframe colours** (blue/green/yellow buttons and the dark canvas are placeholders). Colours always come from the theme.

## Stack
- Flutter (stable), Dart 3, Material 3
- State: Riverpod (`flutter_riverpod`, `riverpod_annotation` codegen)
- Routing: `go_router`
- HTTP: `dio` with an auth interceptor
- Models: `freezed` + `json_serializable`
- Secure storage: `flutter_secure_storage` (tokens only); `shared_preferences` (theme mode only)
- Misc: `url_launcher` (WhatsApp, email, external links), `share_plus` (share), `cached_network_image`
- Tests: `flutter_test`, `mocktail`, golden tests (light and dark)

## Project structure (feature-first)
```
lib/
├── main.dart
├── app/                 # App widget, router, shell (bottom nav + drawer), bootstrap
├── core/
│   ├── theme/           # Spiced Chai tokens, light + dark ThemeData, theme mode provider
│   ├── network/         # dio client, interceptors, error mapping
│   ├── storage/
│   ├── widgets/         # Avatar, VerifiedBadge, ConnectButton, EmptyState, ComingSoon
│   └── utils/
└── features/
    ├── auth/            # Welcome, Create Account, Login (LinkedIn), session
    ├── onboarding/      # Profile completion + verification submit, approval-pending
    ├── announcements/   # HOME feed, saved announcements
    ├── inbox/           # Notifications + incoming connection requests
    ├── directory/       # Search alumni + filters
    ├── profile/         # Own profile, other alumni profile, edit profile
    ├── connections/     # Connect / accept / decline / list
    ├── settings/        # Settings & Privacy
    └── groups/          # Phase 2. Phase 1 = "Coming soon" placeholder only
    # Each feature: data/ (api, dto, repo impl), domain/ (entities, repo interface), presentation/ (screens, widgets, providers)
```
Features never import from other features' `presentation/`. Share through `core/` or domain interfaces.

## Navigation
- **Bottom nav (4 tabs, in this order):** Home, Groups, Search, Profile.
- **Drawer (hamburger on Home app bar):** user header (avatar, name, `batch | branch`), dark/light toggle, Saved Posts, Settings & Privacy, Log Out. The "Communities" entry is added in Phase 2.
- **Home app bar:** hamburger (left), title "Alumni Connect", Inbox icon with unread badge (right).
- Routes: `/welcome`, `/signup`, `/login`, `/onboarding`, `/approval-pending`, `/home`, `/groups`, `/search`, `/profile`, `/alumni/:id`, `/profile/edit`, `/inbox`, `/saved`, `/settings`.

## Design system: Spiced Chai (light + dark)
Never hardcode colors. Use `AppColors` / `Theme.of(context).colorScheme`. The app supports **light and dark**, toggled from the drawer (default: follow system; persist choice with `shared_preferences`).

Brand tokens (from the palette):

| Token | Hex |
|---|---|
| `cream` | `#FDFBD4` |
| `cinnamon` | `#D47E30` |
| `brown` | `#8D5A2B` |
| `nutmeg` | `#825E34` |

Derived tokens:

| Token | Hex | Use |
|---|---|---|
| `ink` | `#3B2712` | Light: body text; Dark: text on cinnamon |
| `surface` | `#FFFDEA` | Light: cards, sheets, inputs |
| `espresso` | `#1C140D` | Dark: scaffold background |
| `espressoSurface` | `#2A1F15` | Dark: cards, sheets, inputs |
| `error` | `#B3261E` | Light: errors, Log Out |
| `errorDark` | `#FF8A80` | Dark: errors, Log Out |

| Role | Light | Dark |
|---|---|---|
| scaffold | cream | espresso |
| surface / card | surface | espressoSurface |
| text | ink | cream |
| primary (filled buttons, links) | brown (text cream) | cinnamon (text ink) |
| accent (FAB, badges, selected chips) | cinnamon | cinnamon |
| borders, muted icons | nutmeg | nutmeg |

Rules:
- Small text must have 4.5:1 contrast. Do **not** put white/cream small text on cinnamon; use ink.
- `brown` and `nutmeg` are close in hue; never use them together to convey different meanings.
- Verified badge (diamond icon next to the name) is cinnamon with a thin brown outline in light mode, and always has `Semantics(label: 'Verified alumnus')`.
- Log Out is shown in the error colour, at the bottom of the drawer and the Settings page.
- Min tap target 48x48; support text scale to 1.3.

## Screens
See `wireframe-screen-specs`. Phase 1 ships: Welcome, Create Account, Login, Profile completion, Approval pending, Home (announcements feed), Inbox, Drawer, Saved Posts, Search, Profile (own and others), Edit profile, Settings & Privacy, and a Groups placeholder.

## Conventions
- Null-safe, no `dynamic`, no `print` (use `logger`).
- Widgets small and `const` where possible; split any `build` over ~80 lines.
- Business logic lives in providers/notifiers, never in widgets.
- All API errors map to a typed `AppFailure`; screens render loading / error / empty / data states explicitly.
- Lists are paginated (cursor) and use `ListView.builder`.
- Images: `cached_network_image`; always show an initials fallback avatar.
- Keep user-facing strings in one place to ease later localization.
- Never log tokens, phone numbers or emails. Don't persist contact details locally.
- Privacy is decided by the server: a hidden field is absent from the response. The UI shows WhatsApp/Email buttons only when the data is present.

## Commands
```bash
flutter pub get
dart run build_runner build --delete-conflicting-outputs
flutter analyze        # must be clean
dart format .
flutter test
flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080
```

## Definition of done
- `flutter analyze` and `flutter test` pass
- Loading, error, empty states handled
- Checked in light **and** dark, at 1x and 1.3x text scale on a 360dp-wide screen
- Matches the wireframe layout; no feature outside the current phase (see `phase-scope-guard`)

## Skills to use
`wireframe-screen-specs`, `flutter-architecture`, `flutter-theme-spiced-chai`, `flutter-api-integration`, `flutter-testing`, `linkedin-auth-flow`, `alumni-directory-search`, `announcements-board`, `notifications-inbox`, `connections-and-privacy`, `privacy-and-security`, `phase-scope-guard`, `git-workflow`.
