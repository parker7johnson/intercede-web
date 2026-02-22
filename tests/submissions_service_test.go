package tests

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/tests/mocks"
)

var testSubmission = &models.Submission{
	Title:       title,
	Body:        body,
	ContactInfo: contact,
	ChurchCode:  testChurchCode,
}

func newSubmissionService(dbResult sql.Result, dbErr error) *services.SubmissionHandler {
	return services.New(&mocks.MockDb{
		MockedFn: func(s string, i interface{}) (sql.Result, error) {
			return dbResult, dbErr
		},
	})
}

func TestServiceCreatePrayerRequest(t *testing.T) {
	cases := []struct {
		name     string
		dbResult sql.Result
		dbErr    error
		wantErr  bool
	}{
		{"success", nil, nil, false},
		{"db error", nil, errors.New("db error"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := newSubmissionService(tc.dbResult, tc.dbErr).CreatePrayerRequest(testSubmission)

			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestServiceCreatePraiseReport(t *testing.T) {
	cases := []struct {
		name     string
		dbResult sql.Result
		dbErr    error
		wantErr  bool
	}{
		{"success", nil, nil, false},
		{"db error", nil, errors.New("db error"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := newSubmissionService(tc.dbResult, tc.dbErr).CreatePraiseReport(testSubmission)

			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
