---
name: deployment-low-cost
description: Low-cost deployment of the Alumni Connect Go API and PostgreSQL, using Docker, a single small VM or PaaS, managed Postgres, CI, backups, and environment configuration, plus Flutter release builds. Use when asked about hosting, Docker, CI/CD, environments, backups, monitoring, scaling, or Play Store builds.
---

# Low-Cost Deployment

Principle: infrastructure stays affordable while usage is small; scale only on measured need.

## Backend
- Multi-stage `Dockerfile` (build with `golang`, run on `distroless`/`alpine`), non-root user, `/healthz` and `/readyz` endpoints.
- Start with **one small instance + one managed Postgres** (or Postgres on the same VM with daily backups) and an S3-compatible bucket for photos. No Kubernetes, no queues, no Redis.
- Run migrations as a release step (`goose up`) before starting the new version.
- Config via env vars: `DATABASE_URL`, `JWT_SECRET`, `LINKEDIN_CLIENT_ID/SECRET`, `LINKEDIN_REDIRECT_URI`, `S3_*`, `ENV`.
- Backups: automated daily Postgres dump with retention; test a restore once before launch. Posts and verification decisions must not silently disappear.
- Availability: browsing should tolerate brief restarts; use graceful shutdown and health checks. Add redundancy later.

## CI (GitHub Actions)
- Backend: `go vet`, `golangci-lint`, `go test -race`, build image.
- Frontend: `flutter analyze`, `flutter test`, build APK on tags.

## Flutter release
- Android first: flavors `dev` / `prod`, `--dart-define=API_BASE_URL`, app signing keys stored outside the repo, `flutter build appbundle --release`.

## Observability
JSON logs with request IDs, basic uptime check, error tracking (e.g. Sentry free tier). Add metrics only when needed.
