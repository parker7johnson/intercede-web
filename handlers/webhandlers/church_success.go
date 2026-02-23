package webhandlers

import (
	"net/http"

	"github.com/parkerjohnson/intercede/web/templates/pages"
)

func (h *WebHandlers) ChurchSuccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Missing session_id", http.StatusBadRequest)
		return
	}

	church, err := h.churchService.GetChurchBySessionID(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "Church not found", http.StatusNotFound)
		return
	}

	if err := pages.ChurchSuccess(church.ChurchCode, church.Name, sessionID).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// ChurchCodeStatus is polled by HTMX until the webhook has fired and the church code is set.
func (h *WebHandlers) ChurchCodeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Missing session_id", http.StatusBadRequest)
		return
	}

	church, err := h.churchService.GetChurchBySessionID(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "Church not found", http.StatusNotFound)
		return
	}

	var renderErr error
	if church.ChurchCode == "" {
		renderErr = pages.ChurchCodePending(sessionID).Render(r.Context(), w)
	} else {
		renderErr = pages.ChurchCodeBlock(church.ChurchCode).Render(r.Context(), w)
	}
	if renderErr != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *WebHandlers) ChurchCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := pages.ChurchCancel().Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
