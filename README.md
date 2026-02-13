# Intercede

A modern, production-ready boilerplate for building web applications with **Go**, **HTMX**, **Templ**, and **Tailwind CSS**.

## Features

- 🚀 **Fast Development** - Hot reload for Go, Templ, and Tailwind
- 🎯 **Type-Safe Templates** - Compile-time checking with Templ
- ⚡ **Dynamic Interactions** - Rich UI without complex JavaScript (HTMX)
- 🎨 **Utility-First Styling** - Tailwind CSS with automatic purging
- 📦 **Single Binary Deployment** - Embedded static assets
- 🏗️ **Clean Architecture** - Handlers, middleware, and dependency injection
- 🔒 **Security Headers** - Built-in security middleware

## Prerequisites

- **Go** 1.23 or higher
- **Node.js** 18 or higher
- **Make** (for build automation)

## Quick Start

### 1. Setup

Install dependencies and download HTMX:

```bash
make setup
```

This will:
- Install Go dependencies
- Install Node dependencies (Tailwind CSS)
- Install the Templ CLI
- Download HTMX from CDN

### 2. Development

Start the development server with hot reload:

```bash
make dev
```

This starts:
- **Go server** on `http://localhost:3000`
- **Templ proxy** on `http://localhost:7331` (use this for hot reload)
- **Tailwind** in watch mode

Visit `http://localhost:7331` to see your app with hot reload enabled.

### 3. Production Build

Build the production binary:

```bash
make build
```

Run the production server:

```bash
./bin/server
```

The binary includes all static assets (CSS, JS) embedded, so you can deploy it anywhere without additional files.

## Project Structure

```
intercede/
├── cmd/web/
│   └── main.go              # Server entry point, routing, middleware
├── internal/
│   ├── handlers/            # HTTP request handlers
│   │   ├── handlers.go      # Dependency injection container
│   │   └── home.go          # Home page handler
│   └── middleware/          # HTTP middleware
│       └── middleware.go    # Logger, security headers
├── web/
│   ├── embed.go             # Embed static files for production
│   ├── templates/           # Templ templates
│   │   ├── layouts/         # Base layouts
│   │   ├── pages/           # Page templates
│   │   └── components/      # Reusable components
│   └── static/              # Static assets
│       ├── css/             # Tailwind CSS
│       └── js/              # JavaScript (HTMX)
├── go.mod                   # Go dependencies
├── package.json             # Node dependencies
├── tailwind.config.js       # Tailwind configuration
├── Makefile                 # Build automation
└── README.md                # This file
```

## Available Commands

Run `make help` to see all available commands:

- `make setup` - Install all dependencies
- `make dev` - Start development server with hot reload
- `make build` - Build production binary
- `make clean` - Remove generated files

## Adding New Pages

### 1. Create a Templ template

Create a new file in `web/templates/pages/`:

```go
// web/templates/pages/about.templ
package pages

import "github.com/parkerjohnson/intercede/web/templates/layouts"

templ About() {
    @layouts.Base("About") {
        <h1>About Page</h1>
        <p>Your content here</p>
    }
}
```

### 2. Create a handler

Add a handler in `internal/handlers/`:

```go
// internal/handlers/about.go
package handlers

import (
    "net/http"
    "github.com/parkerjohnson/intercede/web/templates/pages"
)

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
    pages.About().Render(r.Context(), w)
}
```

### 3. Register the route

Add the route in `cmd/web/main.go`:

```go
mux.HandleFunc("/about", h.About)
```

## Working with HTMX

HTMX is included and ready to use. Example button that makes a POST request:

```html
<button
    hx-post="/api/endpoint"
    hx-target="#result"
    hx-swap="innerHTML"
>
    Click Me
</button>
<div id="result"></div>
```

Create a handler that returns HTML fragments:

```go
func (h *Handlers) MyEndpoint(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte("<p>Updated content!</p>"))
}
```

## Development Tips

### Hot Reload

The development server (`make dev`) automatically reloads when you change:
- `.templ` files (Templ templates)
- `.go` files (Go source code)
- Tailwind classes in templates

Always access the app via `http://localhost:7331` (proxy) for hot reload to work.

### Tailwind CSS

Tailwind is configured to scan all `.templ` files. Just use utility classes in your templates:

```go
<div class="bg-blue-500 text-white p-4 rounded-lg">
    Hello, Tailwind!
</div>
```

The CSS is automatically rebuilt in development and minified in production.

### Static Assets

In development, static files are served from `web/static/`.
In production, they're embedded in the binary via `//go:embed`.

## Tech Stack

- **[Go](https://go.dev/)** - Backend language and HTTP server
- **[HTMX](https://htmx.org/)** - Dynamic interactions without JavaScript
- **[Templ](https://templ.guide/)** - Type-safe HTML templates
- **[Tailwind CSS](https://tailwindcss.com/)** - Utility-first CSS framework

## License

MIT
