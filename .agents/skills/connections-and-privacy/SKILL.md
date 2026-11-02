---
name: connections-and-privacy
description: Alumni connections (request, accept, decline, remove), profile contact buttons for WhatsApp and Email, and the privacy toggles Chat ON/OFF and Contact visibility ON/OFF. Use when working on Connect buttons, connection requests, other-alumni profile screens, contact details, WhatsApp or email links, pronouns, or the Settings & Privacy screen.
---

# Connections, Contact Buttons and Privacy

## Connections
- States: `none`, `pending_outgoing`, `pending_incoming`, `accepted`, `declined`.
- Endpoints: `POST /v1/connections` (`{ addressee_id }`), `POST /v1/connections/{id}/accept|decline`, `DELETE /v1/connections/{id}`, `GET /v1/connections?status=`.
- Mutual acceptance (assumption). Unique pair constraint, no self-connections, per-user daily request limit. Both users must be `approved`.
- Creating/accepting a connection writes an inbox notification in the same transaction (`notifications-inbox`).
- The accepted list is reached from the Profile page (a "Connections" count/link), because the bottom nav has no Connections tab.

## Profile page (wireframe 03)
Fields: photo, name, verified badge, optional pronouns, batch, branch, "Role at Company", bio. Action area: WhatsApp icon, Email icon, full-width Connect/Edit button. No Suggestions (Phase 3).

`GET /v1/alumni/{id}` returns the profile, `connection_status`, and a `contact` object **only when allowed**:

| Contact item | Returned only if |
|---|---|
| `contact.email` | target `contact_visible = true` AND viewer is an **accepted connection** |
| `contact.phone` (E.164) | target `chat_enabled = true` AND `contact_visible = true` AND viewer is an accepted connection AND a phone exists |

If neither applies, omit `contact` entirely. The app shows the WhatsApp/Email buttons only when the corresponding field is present.

## Privacy toggles (Settings & Privacy)
| Setting | Default | Meaning |
|---|---|---|
| `chat_enabled` | ON | "Let connections message me on WhatsApp" |
| `contact_visible` | OFF | "Show my phone number and email to my connections" (labelled "Contact visibility") |

`GET/PATCH /v1/me/privacy`. A WhatsApp link needs the number, so the WhatsApp button requires **both** toggles. Enforce on the server; never send and then hide.

## Opening the links (Flutter)
- WhatsApp: `https://wa.me/<digits only>?text=<prefilled>` via `url_launcher` (external app mode).
- Email: `mailto:<email>`.
- Don't cache or log numbers/emails; handle "can't launch" with a snackbar.

## Phase boundary
Blocking and reporting arrive in Phase 2; keep connection and directory queries easy to extend with a block exclusion.
