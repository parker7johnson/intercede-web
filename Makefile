.PHONY: setup dev build clean help

# Default target
.DEFAULT_GOAL := help

# Templ binary path (use from PATH or fallback to GOPATH/bin)
TEMPL := $(shell command -v templ 2>/dev/null || echo "$$(go env GOPATH)/bin/templ")

## setup: Install dependencies and download HTMX
setup:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing Node dependencies..."
	npm install
	@echo "Installing Templ CLI..."
	go install github.com/a-h/templ/cmd/templ@latest
	@echo "Downloading HTMX..."
	curl -o web/static/js/htmx.min.js https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js
	@echo "Setup complete!"

## dev: Start development server with hot reload
dev:
	@echo "Starting development server..."
	@echo "Server: http://localhost:3000"
	@echo "Proxy (with hot reload): http://localhost:7331"
	@echo ""
	@trap 'kill 0' EXIT; \
	npm run dev:css & \
	$(TEMPL) generate --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/web"

## build: Build production binary
build:
	@echo "Generating Tailwind CSS (minified)..."
	npm run build:css
	@echo "Generating Templ files..."
	$(TEMPL) generate
	@echo "Building Go binary..."
	go build -o bin/server ./cmd/web
	@echo "Build complete! Binary: bin/server"

## clean: Remove generated files and build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f web/static/css/tailwind.css
	find . -name "*_templ.go" -type f -delete
	@echo "Clean complete!"

## help: Show this help message
help:
	@echo "Available commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
