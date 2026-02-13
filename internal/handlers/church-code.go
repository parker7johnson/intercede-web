package handlers

import (
	"net/http"

	"github.com/parkerjohnson/intercede/web/templates/pages"
)

// ChurchCode handles the church code page request
func (h *Handlers) ChurchCode(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Render the church code page component
	err := pages.ChurchCode().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
