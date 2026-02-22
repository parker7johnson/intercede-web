package tests

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/parkerjohnson/intercede/handlers/api"
	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/tests/mocks"
)

func newHandler(dbResult sql.Result, dbError error) *api.ApiHandlers {
	sh := services.New(&mocks.MockDb{
		MockedFn: func(s string, i interface{}) (sql.Result, error) {
			return dbResult, dbError
		},
	})
	return api.NewApi(sh, nil)
}

func newFormRequest(method string, churchCode string) *http.Request {
	form := url.Values{}
	form.Add("title", title)
	form.Add("body", body)
	form.Add("contactInfo", contact)

	r := httptest.NewRequest(method, "/createPrayer", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if churchCode != "" {
		r.Header.Set(api.CHURCH_CODE, churchCode)
	}
	return r
}

func TestCreatePrayerRequest(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		churchCode string
		wantStatus int
		dbResult   sql.Result
		dbError    error
	}{
		{"success", http.MethodPost, testChurchCode, http.StatusOK, nil, nil},
		{"bad method", http.MethodGet, testChurchCode, http.StatusMethodNotAllowed, nil, nil},
		{"no church code", http.MethodPost, "", http.StatusBadRequest, nil, nil},
		{"db error", http.MethodPost, testChurchCode, http.StatusInternalServerError, nil, errors.New("mocked db result")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			newHandler(tc.dbResult, tc.dbError).CreatePrayerRequest(w, newFormRequest(tc.method, tc.churchCode))

			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}


func TestCreatePraiseReport(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		churchCode string
		wantStatus int
		dbResult   sql.Result
		dbError    error
	}{
		{"success", http.MethodPost, testChurchCode, http.StatusOK, nil, nil},
		{"bad method", http.MethodGet, testChurchCode, http.StatusMethodNotAllowed, nil, nil},
		{"no church code", http.MethodPost, "", http.StatusBadRequest, nil, nil},
		{"db error", http.MethodPost, testChurchCode, http.StatusInternalServerError, nil, errors.New("mocked db result")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			newHandler(tc.dbResult, tc.dbError).CreatePraiseReport(w, newFormRequest(tc.method, tc.churchCode))

			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
