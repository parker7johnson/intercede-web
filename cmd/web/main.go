package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/parkerjohnson/intercede/internal/handlers"
	"github.com/parkerjohnson/intercede/internal/middleware"
	"github.com/parkerjohnson/intercede/web"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize handlers with dependencies
	h := handlers.New()

	// Create router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/prayerrequest", h.PrayerRequest)
	mux.HandleFunc("/praisereport", h.PraiseReport)

	// Serve static files
	// In production (with embedded files), this serves from embed.FS
	// In development, this serves from the filesystem
	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Build middleware chain: Logger -> SecurityHeaders -> Router
	handler := middleware.Logger(middleware.SecurityHeaders(mux))

	// Start server
	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
