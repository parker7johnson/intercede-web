package tests

import (
	ctx "context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"slices"
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

func newSubmissionService(dbResult sql.Result, dbErr error, mockedDest interface{}) *services.SubmissionHandler {
	return services.New(&mocks.MockDb{
		MockedFn: func(s string, i interface{}) (sql.Result, error) {
			return dbResult, dbErr
		},
		MockedSelectContext: func(c ctx.Context, i interface{}, s string, args ...interface{}) error {

			v := reflect.ValueOf(i)
			if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
				return fmt.Errorf("dest must be pointer to slice")
			}

			// mockedDest must also be a pointer to a slice
			mv := reflect.ValueOf(mockedDest)
			if mv.Kind() != reflect.Ptr || mv.Elem().Kind() != reflect.Slice {
				return fmt.Errorf("mockedDest must be pointer to slice")
			}

			// Copy the slice value from mockedDest into i
			v.Elem().Set(mv.Elem())

			return dbErr
		},
	})
}

func TestGetChurchCodes(t *testing.T) {
	cases := []struct {
		name     string
		dbResult *[]string
		dbErr    error
	}{
		{"success", &[]string{testChurchCode}, nil},
		{"db error", &[]string{testChurchCode}, errors.New("mockedDbErr")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := newSubmissionService(nil, tc.dbErr, tc.dbResult).GetChurchCodes()

			if err != tc.dbErr {
				t.Fatalf("expected: %v got %v", tc.dbErr, err)
				return
			}

			if err == nil && !slices.EqualFunc(res, *tc.dbResult, func(a, b string) bool { return a == b}) {
				t.Fatalf("expected: %v got %v", *tc.dbResult, res)
			}
		})
	}
}

func TestServiceGetSubmissions(t *testing.T) {
	cases := []struct {
		name     string
		dbResult *[]models.SubmissionDB
		dbErr    error
	}{
		{"success", &[]models.SubmissionDB{}, nil},
		{"db error", &[]models.SubmissionDB{}, errors.New("mockedDbErr")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := newSubmissionService(nil, tc.dbErr, tc.dbResult).GetSubmissionsByChurchCode(testChurchCode)

			if err != tc.dbErr {
				t.Fatalf("expected: %v got %v", tc.dbErr, err)
			}

			if tc.dbErr == nil && !reflect.DeepEqual(res, tc.dbResult) {
				t.Fatalf("expected: %v got %v", tc.dbResult, res)
			}

		})
	}

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
			err := newSubmissionService(tc.dbResult, tc.dbErr, nil).CreatePrayerRequest(testSubmission)

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
			err := newSubmissionService(tc.dbResult, tc.dbErr, nil).CreatePraiseReport(testSubmission)

			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
