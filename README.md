# Intercede

A web application for church congregations to submit prayer requests and praise reports. Members submit through a simple form, church administrators log in to manage submissions.

## How It Works

- Members visit the site and submit a prayer request or praise report via an HTMX-powered form
- Submissions are tagged with a church code passed via the `X-Church-Code` request header
- Administrators log in at `/adminlogin` using email and password, authenticated through Supabase
- Submissions are stored in PostgreSQL

## Tech Stack

- **Go** - HTTP server using the standard library (`net/http`)
- **PostgreSQL** - Submission storage via `sqlx` and `pgx`
- **Supabase** - Admin authentication (email/password, session cookies)
- **HTMX** - Form submissions and page interactions without custom JavaScript
- **Templ** - Type-safe, compiled HTML templates
- **Tailwind CSS** - Utility-first styling
- **Fly.io** - Deployment target

## Prerequisites

- Go 1.24 or higher
- Node.js 18 or higher (for Tailwind)
- Make
- A PostgreSQL database
- A Supabase project (for admin auth)

## Environment Variables

Create a `.env` file in the project root:

```
PORT=3000
DB_URL=postgres://...
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
```

## Development

Install dependencies:

```bash
make setup
```

Start the development server with hot reload:

```bash
make dev
```

The Go server runs on `http://localhost:3000`. Access via `http://localhost:7331` for Templ hot reload.

## Production Build

```bash
make build
./bin/server
```

Static assets (CSS, JS) are embedded in the binary at build time.

## Project Structure

```
intercede/
├── cmd/web/
│   └── main.go                  # Entry point, routing
├── handlers/
│   ├── api/                     # API handlers (submissions, auth)
│   └── webhandlers/             # Page handlers
├── internal/
│   └── middleware/              # Logger, security headers, auth middleware
├── services/
│   ├── models/                  # Data models (Submission)
│   ├── submissions_service.go   # Prayer request and praise report DB logic
│   └── auth_service.go          # Supabase authentication
├── utils/                       # Logger
├── web/
│   ├── templates/               # Templ templates (layouts, pages, components)
│   └── static/                  # CSS, JS assets
├── .github/workflows/           # CI (build, vet, test)
├── Makefile
└── fly.toml                     # Fly.io deployment config
```

## Available Commands

- `make setup` - Install all dependencies
- `make dev` - Start development server with hot reload
- `make build` - Build production binary
- `make clean` - Remove generated files

## CI

GitHub Actions runs on every push and on pull requests to `main`. The pipeline generates Templ files, builds, vets, and runs all tests.
