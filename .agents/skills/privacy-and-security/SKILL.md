---
name: privacy-and-security
description: Security and privacy checklist for Alumni Connect covering authorization, impersonation prevention, sensitive data (roll numbers, phone numbers, verification records), secure storage, rate limiting, input validation, and logging. Use when reviewing or writing auth, profile, directory, admin, upload, or any code handling personal data, and before finishing any backend feature.
---

# Privacy and Security Checklist

## Requirements this protects
- Users must not access private data or impersonate verified alumni.
- Roll numbers and sensitive verification records must never be exposed publicly.

## Checklist
1. **AuthZ in services**, not just routes. Check ownership (`user_id == actor`) and status (`approved`).
2. **Field-level privacy** enforced server-side (phone, chat availability); hidden means absent.
3. **Verification data isolated** in `verification_records`; only admin endpoints read it; excluded from logs, errors, analytics.
4. **Impersonation:** verified badge derives only from `status=approved` set by an admin; never from client input or LinkedIn data alone. Name/batch/branch edits after approval may require re-review (flag `needs_review`).
5. **Input validation:** length limits, enums, URL scheme allowlist (`https` only) for external links; escape text.
6. **SQL:** sqlc/parameterized only.
7. **Uploads:** pre-signed URLs with content-type and size limits (e.g. image/jpeg|png|webp, max 5 MB); strip EXIF GPS; random object keys.
8. **Tokens:** short-lived access JWT, hashed rotating refresh tokens; Flutter secure storage only.
9. **Rate limiting** on auth, search, connection requests.
10. **Logging:** no tokens, phone numbers, roll numbers, or full emails. Include request ID.
11. **Secrets:** env vars only; `.env` is gitignored; provide `.env.example` with placeholders.
12. **Transport:** HTTPS only in production; certificate pinning not required for MVP.
13. **Contact exposure:** email and phone appear only in `GET /v1/alumni/{id}` for accepted connections who pass the privacy toggles (`connections-and-privacy`); never in lists, notifications, shares, or logs. Pronouns are optional free text, length-limited and escaped.
14. **Account deletion/data export:** keep in mind for compliance; design `ON DELETE` rules accordingly.

## Review output
When reviewing, list violations as `file:line: issue, fix`, ordered by severity.
