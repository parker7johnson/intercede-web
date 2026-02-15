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


func (h *ApiHandlers) CreatePrayerRequest(w http.ResponseWriter, r *http.Request) {
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

	//check the code supplied is valid

	err := r.ParseForm()
	if err != nil {
		h.log.LogBadRequest(r, code, "Failed to parse form", err)
		http.Error(w, "Bad form supplied", http.StatusBadRequest)
		return
	}

	req := &models.Submission{
		Title:       r.FormValue("title"),
		Body:        r.FormValue("body"),
		ContactInfo: r.FormValue("contactInfo"),
		ChurchCode:  code,
	}

	err = h.sh.CreatePrayerRequest(req)
	if err != nil {
		h.log.Error("Failed to create prayer request - Request: %s %s from %s, Church-Code: %s, Title: %s, Error: %v",
			r.Method, r.URL.Path, r.RemoteAddr, code, req.Title, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Success("Prayer request created - Church-Code: %s, Title: %s, Body: %s, ContactInfo: %s",
		code, req.Title, req.Body, req.ContactInfo)

	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("Prayer Reqeust successfully created"))

}


func (h *ApiHandlers) CreatePraiseReport(w http.ResponseWriter, r *http.Request) {
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


	err := r.ParseForm()
	if err != nil {
		h.log.LogBadRequest(r, code, "Failed to parse form", err)
		http.Error(w, "Bad form supplied", http.StatusBadRequest)
		return
	}

	req := &models.Submission{
		Title:       r.FormValue("title"),
		Body:        r.FormValue("body"),
		ContactInfo: r.FormValue("contactInfo"),
		ChurchCode:  code,
	}

	err = h.sh.CreatePraiseReport(req)
	if err != nil {
		h.log.Error("Failed to create praise report - report: %s %s from %s, Church-Code: %s, Title: %s, Error: %v",
			r.Method, r.URL.Path, r.RemoteAddr, code, req.Title, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Success("Praise report created - Church-Code: %s, Title: %s, Body: %s, ContactInfo: %s",
		code, req.Title, req.Body, req.ContactInfo)

	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("Prayer Reqeust successfully created"))

}
