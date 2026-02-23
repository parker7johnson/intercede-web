package webhandlers

import (
	"net/http"

	"github.com/parkerjohnson/intercede/internal/middleware"
	"github.com/parkerjohnson/intercede/web/templates/pages"
)

func (h *WebHandlers) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	churchCode := ""
	userID := middleware.GetUserID(r.Context())
	if userID != "" {
		if church, err := h.churchService.GetChurchByUserID(r.Context(), userID); err == nil {
			churchCode = church.ChurchCode
		}
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	var err error
	if isHTMX {
		err = pages.AdminDashboardContent(churchCode).Render(r.Context(), w)
	} else {
		err = pages.AdminDashboard(churchCode).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
