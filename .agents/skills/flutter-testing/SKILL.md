---
name: flutter-testing
description: Testing strategy for the Flutter app: unit tests for notifiers and repositories, widget tests, golden tests, and mocking with mocktail and Riverpod overrides. Use when writing tests, fixing failing tests, adding a feature that needs coverage, or before declaring a frontend task done.
---

# Flutter Testing

## What to test
- **Notifiers/repositories:** unit tests with `mocktail` fakes; cover success, failure, empty, pagination.
- **Widgets:** key screens (sign-in, directory, announcements, settings) with `ProviderScope(overrides: [...])`. Assert loading, error, empty, and data states.
- **Golden tests:** Home card, Search row, Profile, Drawer; run in **light and dark** at 360x800 and text scale 1.3.
- **Router:** redirect logic for signed-out / incomplete profile / pending / approved.

## Conventions
- Files mirror `lib/` under `test/`; name `<thing>_test.dart`.
- Build test data with small factory helpers in `test/fixtures/`.
- No real network or secure storage in tests.
- Privacy: test that fields absent from the API (hidden phone) render no phone UI.

- Profile: WhatsApp/Email buttons render only when `contact` fields are present.
- Search row: Connect button states (Connect, Pending, Accept, Connected).
- Shell: bottom nav order Home, Groups, Search, Profile; Groups shows the Coming soon state.

## Gate
`flutter analyze && flutter test` must pass before finishing.
