package webhandlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/parkerjohnson/intercede/handlers/api"
	"github.com/parkerjohnson/intercede/services/models"
)

func (h *WebHandlers) DashboardSubmissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	churchCode := r.Header.Get(api.CHURCH_CODE)
	subType := r.URL.Query().Get("type")

	if churchCode == "" {
		h.log.LogBadRequest(r, churchCode, "No churchCode supplied", errors.New("No churchCode supplied"))
		http.Error(w, "Error fetching submission", http.StatusBadRequest)
		return
	}

	if subType == "" {
		h.log.LogBadRequest(r, churchCode, "No submission type supplied", errors.New("No submission type supplied"))
		http.Error(w, "Error fetching submission", http.StatusBadRequest)
		return
	}

	submissions, err := h.submissionService.GetSubmissionsByChurchCode(churchCode)
	if err != nil {
		h.log.LogBadRequest(r, churchCode, "Failed getting submissions", err)
		http.Error(w, "Error fetching submissions", http.StatusInternalServerError)
		return
	}

	// Filter by type
	var all []models.SubmissionDB
	if submissions != nil {
		all = *submissions
	}
	var filtered []*models.SubmissionDB
	for _, s := range all {
		if s.SubmissionType == subType {
			filtered = append(filtered, &s)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if len(filtered) == 0 {
		fmt.Fprintf(w, `<p class="prayer-subheader text-center py-8">No %s yet.</p>`, subType)
		return
	}

	for _, s := range filtered {
		fmt.Fprintf(w, `
			<div class="submission border-b py-2">
				<h3 class="font-semibold">%s</h3>
				<p class="text-sm text-gray-700">%s</p>
				<small class="text-gray-500">Submitted by %s</small>
			</div>
		`, s.Title, s.Body, derefString(s.SubmitterName))
	}
}

// Helper
func derefString(s *string) string {
	if s == nil {
		return "Anonymous"
	}
	return *s
}
