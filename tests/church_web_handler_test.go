package tests

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/parkerjohnson/intercede/handlers/webhandlers"
	"github.com/parkerjohnson/intercede/internal/middleware"
	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/tests/mocks"
)

func newWebHandlers(svc *mocks.MockChurchService) *webhandlers.WebHandlers {
	return webhandlers.NewWeb(svc, nil)
}

func TestChurchCreate(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		htmx       bool
		wantStatus int
	}{
		{"GET full page", http.MethodGet, false, http.StatusOK},
		{"GET HTMX partial", http.MethodGet, true, http.StatusOK},
		{"POST not allowed", http.MethodPost, false, http.StatusMethodNotAllowed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// ChurchCreate does not use the church service.
			wh := newWebHandlers(&mocks.MockChurchService{})

			r := httptest.NewRequest(tc.method, "/church-create", nil)
			if tc.htmx {
				r.Header.Set("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			wh.ChurchCreate(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestChurchSuccess(t *testing.T) {
	testChurch := &models.Church{
		Name:            testChurchName,
		AdminEmail:      testChurchEmail,
		StripeSessionID: testSessionID,
		ChurchCode:      "GRAC-7X3K",
		Status:          "active",
	}
	pendingChurch := &models.Church{
		Name:            testChurchName,
		AdminEmail:      testChurchEmail,
		StripeSessionID: testSessionID,
		ChurchCode:      "",
		Status:          "pending",
	}

	cases := []struct {
		name        string
		method      string
		sessionID   string
		svcResult   *models.Church
		svcErr      error
		wantStatus  int
		wantInBody  string
		wantMissing string
	}{
		{
			name:       "success with code",
			method:     http.MethodGet,
			sessionID:  testSessionID,
			svcResult:  testChurch,
			wantStatus: http.StatusOK,
			wantInBody: "GRAC-7X3K",
		},
		{
			name:        "code pending shows polling state",
			method:      http.MethodGet,
			sessionID:   testSessionID,
			svcResult:   pendingChurch,
			wantStatus:  http.StatusOK,
			wantInBody:  "hx-get",
			wantMissing: "GRAC-7X3K",
		},
		{
			name:       "missing session_id",
			method:     http.MethodGet,
			sessionID:  "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "church not found",
			method:     http.MethodGet,
			sessionID:  testSessionID,
			svcErr:     sql.ErrNoRows,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "service error",
			method:     http.MethodGet,
			sessionID:  testSessionID,
			svcErr:     errors.New("db error"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "bad method",
			method:     http.MethodPost,
			sessionID:  testSessionID,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wh := newWebHandlers(&mocks.MockChurchService{
				GetBySessionIDFn: func(ctx context.Context, sessionID string) (*models.Church, error) {
					return tc.svcResult, tc.svcErr
				},
				GetByUserIDFn: func(ctx context.Context, userID string) (*models.Church, error) {
					return nil, errors.New("not used")
				},
			})

			target := "/church/success"
			if tc.sessionID != "" {
				target += "?session_id=" + tc.sessionID
			}
			r := httptest.NewRequest(tc.method, target, nil)
			w := httptest.NewRecorder()

			wh.ChurchSuccess(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}

			if tc.wantInBody != "" {
				body := w.Body.String()
				if !strings.Contains(body, tc.wantInBody) {
					t.Errorf("expected %q in response body", tc.wantInBody)
				}
			}
			if tc.wantMissing != "" {
				body := w.Body.String()
				if strings.Contains(body, tc.wantMissing) {
					t.Errorf("did not expect %q in response body", tc.wantMissing)
				}
			}
		})
	}
}

func TestChurchCodeStatus(t *testing.T) {
	readyChurch := &models.Church{
		Name:            testChurchName,
		StripeSessionID: testSessionID,
		ChurchCode:      "GRAC-7X3K",
		Status:          "active",
	}
	pendingChurch := &models.Church{
		Name:            testChurchName,
		StripeSessionID: testSessionID,
		ChurchCode:      "",
		Status:          "pending",
	}

	cases := []struct {
		name        string
		sessionID   string
		svcResult   *models.Church
		svcErr      error
		wantStatus  int
		wantInBody  string
		wantMissing string
	}{
		{
			name:       "code ready returns code block",
			sessionID:  testSessionID,
			svcResult:  readyChurch,
			wantStatus: http.StatusOK,
			wantInBody: "GRAC-7X3K",
		},
		{
			name:        "code pending keeps polling",
			sessionID:   testSessionID,
			svcResult:   pendingChurch,
			wantStatus:  http.StatusOK,
			wantInBody:  "hx-get",
			wantMissing: "GRAC-7X3K",
		},
		{
			name:       "missing session_id returns 400",
			sessionID:  "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "church not found returns 404",
			sessionID:  testSessionID,
			svcErr:     sql.ErrNoRows,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wh := newWebHandlers(&mocks.MockChurchService{
				GetBySessionIDFn: func(ctx context.Context, sessionID string) (*models.Church, error) {
					return tc.svcResult, tc.svcErr
				},
				GetByUserIDFn: func(ctx context.Context, userID string) (*models.Church, error) {
					return nil, errors.New("not used")
				},
			})

			target := "/church/code-status"
			if tc.sessionID != "" {
				target += "?session_id=" + tc.sessionID
			}
			r := httptest.NewRequest(http.MethodGet, target, nil)
			w := httptest.NewRecorder()

			wh.ChurchCodeStatus(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.wantInBody != "" {
				body := w.Body.String()
				if !strings.Contains(body, tc.wantInBody) {
					t.Errorf("expected %q in response body, got:\n%s", tc.wantInBody, body)
				}
			}
			if tc.wantMissing != "" {
				body := w.Body.String()
				if strings.Contains(body, tc.wantMissing) {
					t.Errorf("did not expect %q in response body", tc.wantMissing)
				}
			}
		})
	}
}

func TestChurchCancel(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"GET", http.MethodGet, http.StatusOK},
		{"POST not allowed", http.MethodPost, http.StatusMethodNotAllowed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wh := newWebHandlers(&mocks.MockChurchService{})

			r := httptest.NewRequest(tc.method, "/church/cancel", nil)
			w := httptest.NewRecorder()

			wh.ChurchCancel(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAdminDashboard(t *testing.T) {
	const testUserID = "user-abc-123"
	testChurch := &models.Church{
		Name:       testChurchName,
		ChurchCode: "GRAC-7X3K",
	}

	cases := []struct {
		name       string
		method     string
		userID     string
		htmx       bool
		svcResult  *models.Church
		svcErr     error
		wantStatus int
		wantInBody string
	}{
		{
			name:       "full page with church code",
			method:     http.MethodGet,
			userID:     testUserID,
			svcResult:  testChurch,
			wantStatus: http.StatusOK,
			wantInBody: "GRAC-7X3K",
		},
		{
			name:       "HTMX partial with church code",
			method:     http.MethodGet,
			userID:     testUserID,
			htmx:       true,
			svcResult:  testChurch,
			wantStatus: http.StatusOK,
			wantInBody: "GRAC-7X3K",
		},
		{
			name:       "graceful degrade when church not found",
			method:     http.MethodGet,
			userID:     testUserID,
			svcErr:     sql.ErrNoRows,
			wantStatus: http.StatusOK,
			wantInBody: "Code loading...",
		},
		{
			name:       "no user ID in context",
			method:     http.MethodGet,
			userID:     "",
			wantStatus: http.StatusOK,
			wantInBody: "Code loading...",
		},
		{
			name:       "bad method",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wh := newWebHandlers(&mocks.MockChurchService{
				GetBySessionIDFn: func(ctx context.Context, sessionID string) (*models.Church, error) {
					return nil, errors.New("not used")
				},
				GetByUserIDFn: func(ctx context.Context, userID string) (*models.Church, error) {
					return tc.svcResult, tc.svcErr
				},
			})

			r := httptest.NewRequest(tc.method, "/admin/dashboard", nil)
			if tc.htmx {
				r.Header.Set("HX-Request", "true")
			}
			if tc.userID != "" {
				ctx := context.WithValue(r.Context(), middleware.UserIDKey, tc.userID)
				r = r.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			wh.AdminDashboard(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.wantInBody != "" {
				body := w.Body.String()
				if !strings.Contains(body, tc.wantInBody) {
					t.Errorf("expected %q in response body", tc.wantInBody)
				}
			}
		})
	}
}
