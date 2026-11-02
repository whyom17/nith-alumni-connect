---
name: git-workflow
description: Git branching, commit message, and pull request conventions for the Alumni Connect monorepo (frontend and backend). Use when committing, creating branches, writing PR descriptions, or preparing a change for review.
---

# Git Workflow

## Branches
`main` is always releasable. Work in `feat/<short-name>`, `fix/<short-name>`, `chore/<short-name>`. Keep branches short-lived.

## Commits (Conventional Commits)
`feat(directory): add branch filter`, `fix(auth): rotate refresh token on reuse`, `chore(backend): bump pgx`. Scope = module or `frontend`/`backend`. One logical change per commit.

## Pull requests
Include: what and why, screenshots for UI (light theme, 360dp), API/spec changes, migration notes, test evidence, tradeoffs, and the phase it belongs to.

## Before pushing
- Backend: `make lint test`
- Frontend: `flutter analyze && flutter test`
- Spec updated if the API changed; no secrets or `.env` committed.

## Don'ts
No force-push to `main`; no unrelated refactors in feature PRs; no commented-out code.
