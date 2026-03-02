package tests

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/parkerjohnson/intercede/handlers/api"
	"github.com/parkerjohnson/intercede/handlers/webhandlers"
	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/tests/mocks"
)

func newWebHandler(dbResult *[]models.SubmissionDB, dbError error) *webhandlers.WebHandlers {
	mockService := &mocks.MockSubmissionService{
		GetSubmissionsByChurchCodeFn: func(churchCode string) (*[]models.SubmissionDB, error) {
			return dbResult, dbError
		},
		CreatePrayerRequestFn: func(sub *models.Submission) error {
			return nil
		},
		CreatePraiseReportFn: func(sub *models.Submission) error {
			return nil
		},
		GetChurchCodesFn: func() ([]string, error) {
			return []string{}, nil
		},
	}
	return webhandlers.NewWeb(nil, mockService)
}

func newHandler(dbResult sql.Result, dbError error) *api.ApiHandlers {
	sh := services.New(&mocks.MockDb{
		MockedFn: func(s string, i interface{}) (sql.Result, error) {
			return dbResult, dbError
		},
		MockedSelectContext: func(ctx context.Context, i1 interface{}, s string, i2 ...interface{}) error {
			slicePtr, ok := i1.(*[]string)
			if !ok {
				return errors.New("bad type casting")
			}
			*slicePtr = []string{testChurchCode}
			return nil
		},
	})
	return api.NewApi(sh, nil)
}

func newQueryRequest(method, submissionType string) *http.Request {
	params := url.Values{}
	params.Add("type", submissionType)

	r := httptest.NewRequest(method, "/admin/dashboard/submission?"+params.Encode(), nil)
	r.Header.Set(api.CHURCH_CODE, testChurchCode)

	return r
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

func TestGetSumbissionsHttp(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		submissionType string
		churchCode     string
		wantStatus     int
		dbResult       *[]models.SubmissionDB
		dbError        error
	}{
		{"success", http.MethodGet, "prayer", testChurchCode, http.StatusOK, nil, nil},
		{"bad method", http.MethodPost, "prayer", testChurchCode, http.StatusMethodNotAllowed, nil, nil},
		{"no submission type", http.MethodGet, "", testChurchCode, http.StatusBadRequest, nil, nil},
		{"db error", http.MethodGet, testChurchCode, testChurchCode, http.StatusInternalServerError, nil, errors.New("mocked db result")},
		{"no church code", http.MethodGet, testChurchCode, "", http.StatusBadRequest, nil, errors.New("mocked db result")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := newQueryRequest(tc.method, tc.submissionType)
			if tc.churchCode == "" {
				r.Header.Set(api.CHURCH_CODE, "")
			}
			newWebHandler(tc.dbResult, tc.dbError).DashboardSubmissions(w, r)

			if w.Code != tc.wantStatus {
				t.Fatalf("expected status %d got %d", tc.wantStatus, w.Code)
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
