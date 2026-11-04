# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Reminder for the owner

At the start of a session, remind the owner of the unchecked items under "Open tasks" at the end of `tasks/todo.md`. Remove this section once they are done.

## Repository layout

This is a monorepo with two independent services:

- **`/` (root)** — Nuxt 3 SSG frontend (Node/TypeScript)
- **`/backend`** — Go HTTP backend (module `ivmanto.com/backend`)

Each has its own `Dockerfile` and is deployed as a separate Cloud Run service via `cloudbuild.yaml`.

## Local development

**Backend** — requires a `backend/.env` file (copy from `backend/.env.example`):
```bash
cd backend && go run ./cmd/server   # listens on :8080
```
On startup the backend initialises all clients (GCS, Vertex AI, Google Calendar, SMTP). If any required env var is missing, it exits immediately — check the log for `"missing required environment variables"`.

For GCP APIs locally: `gcloud auth application-default login` then `gcloud auth application-default set-quota-project ivmanto-com-prod`.

**Frontend** — proxies `/api/*` to `localhost:8080` in dev mode:
```bash
npm install && npm run dev   # listens on :3000
```

## Common commands

```bash
# Frontend
npm run lint       # ESLint + auto-fix
npm run format     # Prettier
npm run generate   # SSG build (what Cloud Build runs)

# Backend
cd backend
go build ./...
go test ./...
go test ./internal/gcal/...   # run a single package's tests
```

## Architecture

### Request flow

```
Browser → Nuxt frontend (Cloud Run)
              └─ /api/* → Go backend (Cloud Run)
                              ├─ /api/booking/*    → Google Calendar API (via DWD impersonation)
                              ├─ /api/blog/*       → GCS bucket (markdown) → in-memory cache
                              ├─ /api/contact      → SMTP
                              ├─ /api/ideas        → Vertex AI (Gemini)
                              └─ /api/articles     → in-memory likes counter
```

### Backend packages (`backend/internal/`)

| Package | Role |
|---|---|
| `config` | Loads all env vars at startup; exits if any required var is missing |
| `gcal` | Google Calendar wrapper — fetches availability, books/cancels slots, creates Meet conferences |
| `booking` | HTTP handler for `/api/booking/*`; orchestrates gcal + email + analytics |
| `email` | SMTP service; HTML templates live in `email/*.html` |
| `blog` | GCS-backed blog pipeline with Pub/Sub-triggered in-memory cache refresh |
| `ideas` | Vertex AI GenAI handler; prompt template configurable via env/Secret Manager |
| `analytics` | GA4 Measurement Protocol client for server-side conversion tracking |
| `middleware` | CORS (allow-all) and structured JSON request logging |
| `contact` | Contact form → SMTP email |

### Google Calendar / booking

Available time slots are Google Calendar events whose `Summary` matches `GCAL_AVAILABLE_SLOT_SUMMARY` (env var, currently `AfB`). Booking updates the event in-place (atomic via ETag). Google Meet conferences are created via `ConferenceDataVersion(1)` on the update call.

**Auth:** The backend uses Domain-Wide Delegation (DWD) — `impersonate.CredentialsTokenSource` with `Subject` set to `GCAL_IMPERSONATE_USER`. This is required because plain service-account ADC cannot create Google Meet conferences. The DWD scope (`https://www.googleapis.com/auth/calendar`) must be authorised in the Workspace Admin Console for the service account's OAuth client ID.

### Frontend → backend in production

`NUXT_API_BASE_URL` build arg (set in `cloudbuild.yaml`) is `https://ivmanto.com`. At build time, SSG pages that need backend data (e.g. blog article list) fetch from this URL. Runtime `/api/*` calls from the browser hit the same domain and are load-balanced to the backend Cloud Run service.

## Deployment

Push to `main` triggers Cloud Build (`cloudbuild.yaml`), which:
1. Builds and pushes both Docker images to Artifact Registry (`europe-west3`)
2. Deploys frontend and backend to Cloud Run (`europe-west3`, project `ivmanto-com-prod`)
3. Backend secrets (`SMTP_PASS`, `GA_API_SECRET`, etc.) are injected from Secret Manager; non-secret config comes from substitution variables in the trigger

Working branches follow the pattern `dev-vX.Y.Z`. Merge to `main` via PR — branch protection requires checks to pass.
