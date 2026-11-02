---
name: notifications-inbox
description: The Inbox (top-right icon on Home) for Alumni Connect: in-app notifications and incoming connection requests, unread badge, and the notifications table and endpoints. Use when working on the inbox screen, notification items, unread counts, connection request alerts, verification result notices, or any "notify users about activity" requirement.
---

# Notifications and Inbox (Phase 1)

The Inbox is **not** a chat. One-to-one conversation happens on WhatsApp (link-out). Push notifications (FCM) are out of scope for Phase 1; the inbox is pull-based with an unread badge.

## Content
1. **Requests** section on top: incoming pending connections with Accept / Decline.
2. **Notifications** list, newest first.

## Types (Phase 1)
`connection_request`, `connection_accepted`, `verification_approved`, `verification_rejected` (with reason).
Add `announcement_pinned` only if requested. Community/post types arrive in Phase 2.

## Backend
- Table `notifications(id, user_id, type, payload jsonb, read_at, created_at)`; index `(user_id, read_at, created_at DESC)`.
- Create a notification **in the same transaction** as the triggering event (connection created/accepted, verification decision), so none silently disappear.
- Endpoints: `GET /v1/inbox?limit=&cursor=`, `GET /v1/inbox/unread-count`, `POST /v1/inbox/{id}/read`, `POST /v1/inbox/read-all`.
- Payload holds ids and display names only (never phone, email, roll number).
- Retention: prune read notifications older than 90 days (scheduled task or on read).
- A user can only read their own notifications.

## Flutter
- Home app bar Inbox icon with count badge (99+ cap), refreshed on app resume and after relevant actions; no polling faster than ~60s.
- Tapping a notification marks it read and navigates (to `/alumni/:id` or `/profile`).
- States: loading, empty ("You're all caught up"), error with retry.
