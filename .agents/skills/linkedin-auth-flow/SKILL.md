---
name: linkedin-auth-flow
description: Sign in with LinkedIn (OpenID Connect) for Alumni Connect, plus app JWT access and refresh tokens, on both Flutter and Go sides. Use when touching login, signup, session, token refresh, logout, LinkedIn OAuth, deep links for auth callback, or route guards based on auth state.
---

# LinkedIn Auth Flow

Signup/login is **LinkedIn only** in Phase 1. No passwords.

## Flow
1. App opens LinkedIn authorization in a system browser/custom tab (`flutter_appauth` or `url_launcher` + app link) with `state`, `nonce`, PKCE (`code_challenge`).
2. LinkedIn redirects to the app link with `code`.
3. App sends `{ code, code_verifier, redirect_uri }` to `POST /v1/auth/linkedin`.
4. Backend exchanges the code with LinkedIn, **verifies the ID token** (signature via JWKS, `iss`, `aud`, `exp`, `nonce`), reads `sub`, name, email, picture.
5. Backend upserts `users` by `linkedin_sub` (new users start `status=pending`), returns `{ access_token, refresh_token, user: { status, profile_complete } }`.
6. App stores tokens in secure storage and the router redirects by status (see `flutter-architecture`).

## Screens (wireframe 01)
Welcome (`/welcome`, "Get Started") to Create Account (`/signup`, "Sign up with LinkedIn") or Login (`/login`, "Log in with LinkedIn"). Both buttons run the identical flow above; the backend upserts, so Login with an unknown account simply creates a pending account. The response includes `is_new_user`, which the app may use to show a one-time welcome. See `wireframe-screen-specs`.

## Backend rules
- Access token: JWT, ~15 min, claims `sub`, `role`, `status`. Refresh token: opaque random, stored **hashed**, rotated on every use, revoked on logout and reuse detection.
- `POST /v1/auth/refresh`, `POST /v1/auth/logout`.
- Always load status from the DB on sensitive actions; do not trust a stale claim for approval.
- LinkedIn client secret in env only. Rate-limit `/auth/*`.
- LinkedIn identity alone does **not** make someone a verified alumnus. Approval is a separate step (`admin-approval-workflow`).

## Frontend rules
- Never embed the client secret in the app. Only the client ID and redirect URI.
- Clear secure storage and providers on logout.
