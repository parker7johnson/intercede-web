# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Development (hot reload via templ proxy at :7331, server at :3000)
make dev

# Run all tests
go test ./...

# Run a single test
go test ./tests/ -run TestFunctionName

# Build production binary (generates CSS + templ first)
make build

# Generate templ files only (required after editing .templ files)
$(go env GOPATH)/bin/templ generate

# Install all dependencies (Go modules, Node, templ CLI, HTMX)
make setup
```

> `templ` is not in PATH by default — always use `$(go env GOPATH)/bin/templ`.

## Architecture

**Entry point:** `cmd/web/main.go` — wires all dependencies, registers routes on a plain `http.ServeMux`, and applies middleware as a chain (`Logger → SecurityHeaders → handler`). `RequireAuth` is applied per-route as a wrapper, not globally.

**Handler packages:**
- `handlers/api/` — `ApiHandlers` (prayer/praise submissions + admin login) and `ChurchAPIHandlers` (Stripe checkout + webhook)
- `handlers/webhandlers/` — `WebHandlers` (page renders for all user-facing routes)

**Services:**
- `services/submissions_service.go` — prayer request / praise report DB logic
- `services/auth_service.go` — Supabase sign-in, wraps `gotrue` client
- `services/church_service.go` — pending/activate church, unique code generation

**Models:** `services/models/` — structs with `db:` tags used by `sqlx`.

**Templates:** `web/templates/pages/*.templ` — compiled to `*_templ.go` by the templ CLI. Edit only `.templ` files and regenerate; never edit `*_templ.go` directly.

**Middleware:** `internal/middleware/` — `Logger`, `SecurityHeaders`, `RequireAuth`. Auth reads a `session` cookie, calls `TokenVerifier.VerifyToken`, and stores the user ID in context under `middleware.UserIDKey`. Retrieve it with `middleware.GetUserID(ctx)`.

**Static assets** (CSS, JS) are embedded in the binary via `web/static/` and `go:embed`.

## Key Patterns

**Service interfaces are declared in the consuming package**, not the service package. Each handler struct defines the minimal interface it needs (e.g., `churchLookup` in `webhandlers`, `churchService` in `api`).

**DB access:** `NamedExec` for INSERT/UPDATE with named struct/map params; `GetContext` for single-row SELECT.

**HTMX-aware handlers** check the `HX-Request` header and render only the content fragment vs. the full page layout. Forms that redirect off-site (e.g., to Stripe) use `hx-boost="false"`.

**Tests** live in `tests/` package with mocks in `tests/mocks/`. Mocks implement the same minimal interfaces declared in handler packages.

## Church Registration Flow

1. `GET /church-create` → form (name, admin email, password)
2. `POST /api/church/checkout` → Supabase signup → Stripe checkout session → insert pending church → redirect to Stripe
3. Stripe webhook `POST /webhooks/stripe` → `checkout.session.completed` → activate church, generate code (`PREFIX-XXXX` format)
4. `GET /church/success?session_id=xxx` → show code if ready, otherwise HTMX polling card
5. `GET /church/code-status?session_id=xxx` → polled by HTMX; returns code block when ready

**Local Stripe testing:**
```bash
stripe listen --forward-to localhost:3000/webhooks/stripe --events checkout.session.completed
# Copy the whsec_... secret into .env as STRIPE_WEBHOOK_SECRET, then restart server
# Test card: 4242 4242 4242 4242
```

## Environment Variables

```
PORT                    # defaults to 3000
DB_URL                  # postgres connection string
SUPABASE_URL
SUPABASE_KEY
STRIPE_SECRET_KEY
STRIPE_WEBHOOK_SECRET
STRIPE_PRICE_ID
BASE_URL                # defaults to http://localhost:<PORT>
```

## Database Schema

The `churches` table must be created manually:

```sql
CREATE TABLE churches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    admin_email TEXT NOT NULL,
    church_code TEXT UNIQUE,
    stripe_session_id TEXT UNIQUE NOT NULL,
    stripe_customer_id TEXT,
    stripe_subscription_id TEXT,
    user_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT now()
);
```
