package api

import (
	"net/http"
)

func (ah *ApiHandlers) Login(w http.ResponseWriter, r *http.Request) {
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
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
}
