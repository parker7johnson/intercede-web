package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/parkerjohnson/intercede/internal/middleware"
	"github.com/parkerjohnson/intercede/tests/mocks"
)

func TestRequireAuth(t *testing.T) {
	sentinel := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cases := []struct {
		name           string
		cookie         *http.Cookie
		verifyErr      error
		wantStatus     int
		wantLocation   string
	}{
		{
			name:         "no session cookie redirects to adminlogin",
			cookie:       nil,
			wantStatus:   http.StatusSeeOther,
			wantLocation: adminLoginPath,
		},
		{
			name:         "invalid token redirects to adminlogin",
			cookie:       &http.Cookie{Name: sessionCookieName, Value: badToken},
			verifyErr:    errors.New("invalid token"),
			wantStatus:   http.StatusSeeOther,
			wantLocation: adminLoginPath,
		},
		{
			name:       "valid token passes through to next handler",
			cookie:     &http.Cookie{Name: sessionCookieName, Value: testToken},
			verifyErr:  nil,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			verifier := &mocks.MockAuthVerifier{
				VerifyTokenFn: func(token string) error {
					return tc.verifyErr
				},
			}

			handler := middleware.RequireAuth(verifier)(sentinel)

			req := httptest.NewRequest(http.MethodGet, adminPath, nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.wantLocation != "" {
				if loc := rec.Header().Get("Location"); loc != tc.wantLocation {
					t.Errorf("Location: got %q, want %q", loc, tc.wantLocation)
				}
			}
		})
	}
}
