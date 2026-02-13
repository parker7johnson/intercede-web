package handlers

import (
	"net/http"

	"github.com/parkerjohnson/intercede/web/templates/pages"
)

// Home handles the home page request
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Render the home page component
	err := pages.Home().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
