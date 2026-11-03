# NIT Hamirpur Alumni Connect

A mobile app that connects NIT Hamirpur alumni with each other: verified profiles, a searchable directory, and an announcements board. Built for the college community only.

## Why this exists

Alumni connections are scattered across WhatsApp groups, LinkedIn, and personal contacts. Three people this app should serve:

- An alumnus in Bengaluru looking for seniors in backend engineering.
- Someone in Germany looking for NIT Hamirpur graduates across Europe.
- A recent graduate looking for career guidance from an alumnus at Google or Amazon.

## Product pillars

1. **Identity and trust**: verified profiles (batch, branch, verification badge)
2. **Discovery**: search by name, batch, branch, company, skills, region
3. **Communities**: interest-based and regional groups (Phase 2)
4. **Communication**: connections, WhatsApp links, notifications

## Tech stack

| Layer | Choice |
|---|---|
| Mobile app | Flutter (Dart) |
| Backend | Go, modular monolith, REST API |
| Database | PostgreSQL |
| Auth | Sign in with LinkedIn (OIDC), app-issued JWT |
| Coding agent | Antigravity (`AGENTS.md` + `.agents/skills/`) |

Architecture rule: **one backend deployable, one database**. No microservices, queues, or caches until a measured need exists.

## Roadmap

### Phase 1: MVP (target: mid December)
- Signup/login via LinkedIn only (Welcome, Create Account, Login screens)
- Profile: name, pronouns (optional), batch, branch, current city, company/role, bio, profile photo
- Manual admin approval (Alumni Association) of new profiles; verified badge
- Search alumni (name, batch, branch, city; filters Batch, Branch, City)
- Connect with other alumni (LinkedIn-style); WhatsApp and Email buttons on profiles (privacy-gated)
- **Home = announcements feed** shown as post cards (initiatives, jobs, events, reunions, donations) with Share and Save
- Saved Posts, Inbox (connection requests and notifications), side drawer, dark/light mode toggle
- Settings & Privacy: Chat ON/OFF, Contact visibility ON/OFF, Log Out
- Groups tab exists as a "Coming soon" placeholder

### Phase 2 (target: mid February)
- Groups tab: joined groups, discover, create, join (chapter-based, interest-based, special groups)
- Likes; drawer "Communities" entry
- Posts and comments (Like and Comment actions on post cards)
- Basic moderation, blocking, reporting

### Phase 3 (aspirational)
- Alumni and community recommendations (profile "Suggestions" row), advanced search
- Institutional services (marksheets, guest house booking, campus entry)
- SOS

## Actors

Alumnus, Verified Alumnus, Community Moderator (verified + admin-approved), Community Member, Admin (Alumni Association).

## Screens and navigation

Wireframes are in `docs/wireframes/` (colours there are placeholders; the theme decides colours).

- **Bottom nav:** Home, Groups, Search, Profile
- **Drawer:** user header, dark/light toggle, Saved Posts, Settings & Privacy, Log Out
- **Home bar:** hamburger, "Alumni Connect", Inbox icon
- Per-screen rules: `.agents/skills/wireframe-screen-specs/SKILL.md`

## Repository layout

```
alumni-connect/
├── README.md
├── docs/
│   └── wireframes/        # Excalidraw exports (auth, home/sidebar/settings, profile/search/groups)
├── frontend/              # Flutter app
│   └── AGENTS.md          # Rules for the frontend agent
├── backend/               # Go API
│   └── AGENTS.md          # Rules for the backend agent
└── .agents/
    └── skills/            # Antigravity workspace skills (one folder per skill)
```

## Getting started

```bash
# Backend
cd backend
cp .env.example .env        # fill in LinkedIn + DB values
docker compose up -d db     # local Postgres
make migrate-up
make run                    # http://localhost:8080

# Frontend
cd frontend
flutter pub get
flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080
```

## Design

The app uses the **Spiced Chai** palette (cream, cinnamon, nutmeg) in **light and dark** modes, with a toggle in the drawer. See `frontend/AGENTS.md` and the `flutter-theme-spiced-chai` skill.

| Token | Hex |
|---|---|
| Cream | `#FDFBD4` |
| Cinnamon | `#D47E30` |
| Brown | `#8D5A2B` |
| Nutmeg | `#825E34` |

## Working with Antigravity

- Agents read the nearest `AGENTS.md`: work inside `frontend/` or `backend/` and the matching rules apply.
- Skills in `.agents/skills/` load on demand based on their `description`.
- Screens must match `docs/wireframes/` (see `wireframe-screen-specs`).
- Always check the phase scope (`phase-scope-guard` skill) before adding a feature.
