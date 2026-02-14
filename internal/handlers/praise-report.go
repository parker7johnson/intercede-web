package handlers

import (
	"net/http"

	"github.com/parkerjohnson/intercede/web/templates/pages"
)

// PraiseReport handles the praise report page request
func (h *Handlers) PraiseReport(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if this is an HTMX request
	isHTMX := r.Header.Get("HX-Request") == "true"

	var err error
	if isHTMX {
		// For HTMX requests, return just the content fragment
		err = pages.PraiseReportContent().Render(r.Context(), w)
	} else {
		// For regular requests, return the full page with layout
		err = pages.PraiseReport().Render(r.Context(), w)
	}

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
