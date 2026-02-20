package webhandlers 

import (
	"net/http"

	"github.com/parkerjohnson/intercede/web/templates/pages"
)

func (h *WebHandlers) AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"

	var err error
	if isHTMX {
		err = pages.AdminLoginContent().Render(r.Context(), w)
	} else {
		err = pages.AdminLogin().Render(r.Context(), w)
	}

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
