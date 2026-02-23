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
	const testUserID = "user-id-abc123"

	sentinel := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cases := []struct {
		name         string
		cookie       *http.Cookie
		userID       string
		verifyErr    error
		wantStatus   int
		wantLocation string
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
			userID:     testUserID,
			verifyErr:  nil,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			verifier := &mocks.MockAuthVerifier{
				VerifyTokenFn: func(token string) (string, error) {
					return tc.userID, tc.verifyErr
				},
			}

			var capturedUserID string
			sentinelCapture := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUserID = middleware.GetUserID(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.RequireAuth(verifier)(sentinelCapture)

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
			if tc.wantStatus == http.StatusOK && tc.userID != "" {
				if capturedUserID != tc.userID {
					t.Errorf("context user ID: got %q, want %q", capturedUserID, tc.userID)
				}
			}
		})
	}

	// Suppress unused variable warning for the unused sentinel
	_ = sentinel
}
