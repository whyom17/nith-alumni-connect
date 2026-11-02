---
name: admin-approval-workflow
description: Manual admin approval and verification of alumni profiles by the Alumni Association, including pending/approved/rejected states and the verified badge. Use when working on onboarding, approval pending screens, admin endpoints, verification records, roll number handling, or the verified badge logic.
---

# Admin Approval Workflow (Manual, Phase 1)

## States
`pending` (default after LinkedIn signup) then `approved` | `rejected` | `suspended`.

## Flow
1. User completes profile (name, batch, branch, city, company/role, bio, photo) and submits verification info (e.g. roll number, optional evidence).
2. Backend stores submission in `verification_records`; user sees "Approval pending".
3. Admin lists pending users (`GET /v1/admin/verifications?status=pending`), reviews, then `POST /v1/admin/verifications/{id}/decision` with `{ decision, reason }`.
4. Decision runs in **one transaction**: update `users.status`, write `decided_by/decided_at/reason`. On failure, return an error; never silently drop it.
5. A `verification_approved` or `verification_rejected` notification is written in the same transaction (`notifications-inbox`).
6. Approved users get access to directory, announcements, connections and the verified badge.

## Rules
- Admin endpoints require role `admin`; seed admins via migration or CLI, never via public API.
- Roll numbers and evidence are **never** returned outside admin endpoints and are never logged.
- Rejected users can see the reason and resubmit.
- Audit trail: keep every decision (append rows, don't overwrite).
- Admin UI for Phase 1 can be minimal (API + simple internal tool). Don't build a full dashboard in the mobile app unless asked.
