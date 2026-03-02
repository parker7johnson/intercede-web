package api

import (
	"net/http"

)

func (ah *ApiHandlers) Login(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	
	session, err := ah.authHandler.Login(email, password)
	if err != nil {
		http.Error(w, "Invalid email or password supplied", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.AccessToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	w.Header().Set("HX-Redirect", "/admin/dashboard")
	w.WriteHeader(http.StatusOK)
}
