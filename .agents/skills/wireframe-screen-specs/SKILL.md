---
name: wireframe-screen-specs
description: Screen-by-screen UI specification for the Alumni Connect Flutter app, derived from the Excalidraw wireframes in docs/wireframes/. Use whenever building, changing, or reviewing ANY screen, layout, navigation, bottom bar, drawer, app bar, list item, button placement, or when asked what a screen should contain or which phase a UI element belongs to.
---

# Wireframe Screen Specs

Images: `docs/wireframes/01-auth-screens.png`, `02-home-sidebar-settings.png`, `03-profile-search-groups.png`. View them when unsure.

**Wireframe rules**
- Layout, labels, and hierarchy are binding. **Colours are not**: blue Connect button, green WhatsApp, yellow Email, and the dark canvas are placeholders. Use theme tokens.
- Everything not drawn (onboarding, approval pending, edit profile, inbox, saved) must follow the same visual language.

## Global shell
- Bottom nav, 4 items in order: **Home, Groups, Search, Profile**.
- Drawer from Home's hamburger. Home app bar: hamburger, title "Alumni Connect", Inbox icon with unread badge.

## Auth (image 01)
| Screen | Route | Contents |
|---|---|---|
| Welcome | `/welcome` | Logo circle (placeholder asset), "NITH ALUMNI CONNECT", **Get Started** button (to `/signup`), "Already have an account? Log In" (to `/login`) |
| Create Account | `/signup` | Heading "Create Account", **Sign up with LinkedIn** button, "Already have an account? Log in" |
| Login | `/login` | Heading "Log In", **Log in with LinkedIn** button (the wireframe says "Signup with LinkedIn", treated as a typo), "Don't have an account? Sign up" |

Both buttons call the same flow (`linkedin-auth-flow`); the screens differ only in copy and links.

## Home (image 02) `/home`
- Feed of announcement **post cards** (see `announcements-board`): header (avatar + author name), body, divider, action row.
- Action row in the wireframe: Like, Comment, Share, Save. **Phase 1 shows Share and Save only**; Like and Comment arrive with posts in Phase 2.
- Category filter chips above the feed (All, Initiatives, Jobs, Events, Reunions, Donations).

## Drawer (image 02)
- Header: avatar, name, `batch | branch` (e.g. "2028 | CSE"); top-right circle = **dark/light mode toggle**.
- Items: **Saved Posts**; **Communities** (Phase 2, hidden for now); **Settings & Privacy**.
- Bottom: **Log Out**.

## Settings & Privacy (image 02) `/settings`
- Back arrow + title "Settings & Privacy"; sections separated by dividers; **Log Out** in red at the bottom.
- Sections (the wireframe leaves them blank, so this is the proposal): **Privacy** (Chat on WhatsApp ON/OFF, Contact visibility ON/OFF), **Account** (Edit profile, LinkedIn email read-only), **About** (version, terms).

## Profile (image 03) `/profile` and `/alumni/:id`
- Avatar; name + verified diamond badge + optional `| pronouns`; "Batch of 20XX" (+ branch); "Role at Company"; Bio.
- Contact icon buttons (WhatsApp, Email): shown **only if the API returns the data** (`connections-and-privacy`).
- Full-width primary button: **Connect** on others (states: Connect / Pending / Accept / Connected); **Edit profile** on own profile.
- "Suggestions" row (4 tiles) is **Phase 3. Do not build it.**

## Search (image 03) `/search`
- Title "Find Alumni"; search field "Search by name, batch, branch, ..."; "Filters" row with chips **Batch, Branch, City**.
- Result rows: avatar, name, `2016 | ECE`, city, trailing Connect button; divider between rows; row tap opens the profile. Details in `alumni-directory-search`.

## Groups (image 03) `/groups`
- Title "Joined Groups"; "Discover new groups" entry; list of joined groups. **Phase 2.** In Phase 1 render the title and a "Coming soon" empty state (`core/widgets/ComingSoon`).

## Inbox `/inbox` (not drawn)
Incoming connection requests (accept/decline) above a list of notifications. See `notifications-inbox`.
