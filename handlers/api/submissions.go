package api

import (
	"github.com/parkerjohnson/intercede/services/models"
	"net/http"
)

const (
	CHURCH_CODE = "X-Church-Code"
	NO_CODE     = "No church code supplied"
	NOT_ALLOWED = "Method not Allowed"
	BAD_FORM    = "Bad form supplied"
)


func (h *ApiHandlers) handleSubmission(w http.ResponseWriter, r *http.Request, save func(*models.Submission) error) {
	code := r.Header.Get(CHURCH_CODE)

	if r.Method != http.MethodPost {
		h.log.LogBadRequest(r, code, "Method not allowed", nil)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if code == "" {
		h.log.LogBadRequest(r, "", "No church code supplied", nil)
		http.Error(w, NO_CODE, http.StatusBadRequest)
		return
	}

	// check the code supplied is valid

	err := r.ParseForm()
	if err != nil {
		h.log.LogBadRequest(r, code, "Failed to parse form", err)
		http.Error(w, BAD_FORM, http.StatusBadRequest)
		return
	}

	req := &models.Submission{
		Title:       r.FormValue("title"),
		Body:        r.FormValue("body"),
		ContactInfo: r.FormValue("contactInfo"),
		ChurchCode:  code,
	}

	if err = save(req); err != nil {
		h.log.Error("Failed to save submission - Request: %s %s from %s, Church-Code: %s, Title: %s, Error: %v",
			r.Method, r.URL.Path, r.RemoteAddr, code, req.Title, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Success("Submission created - Church-Code: %s, Title: %s, Body: %s, ContactInfo: %s",
		code, req.Title, req.Body, req.ContactInfo)

	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("Submission successfully created"))
}

func (h *ApiHandlers) CreatePrayerRequest(w http.ResponseWriter, r *http.Request) {
	h.handleSubmission(w, r, h.submissionHandler.CreatePrayerRequest)
}

func (h *ApiHandlers) CreatePraiseReport(w http.ResponseWriter, r *http.Request) {
	h.handleSubmission(w, r, h.submissionHandler.CreatePraiseReport)
}
