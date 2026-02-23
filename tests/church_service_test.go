package tests

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/tests/mocks"
)

const (
	testSessionID      = "cs_test_session123"
	testChurchName     = "Grace Community Church"
	testChurchEmail    = "admin@grace.com"
	testChurchCodeFmt  = "GRAC" // expected 4-letter prefix
)

func newChurchService(namedExecFn func(string, interface{}) (sql.Result, error), getContextFn func(context.Context, interface{}, string, ...interface{}) error) *services.ChurchService {
	return services.NewChurchService(&mocks.MockChurchDb{
		NamedExecFn:  namedExecFn,
		GetContextFn: getContextFn,
	})
}

func TestCreatePendingChurch(t *testing.T) {
	cases := []struct {
		name    string
		dbErr   error
		wantErr bool
	}{
		{"success", nil, false},
		{"db error", errors.New("db error"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := newChurchService(
				func(q string, arg interface{}) (sql.Result, error) { return nil, tc.dbErr },
				nil,
			)
			church := &models.Church{
				Name:            testChurchName,
				AdminEmail:      testChurchEmail,
				StripeSessionID: testSessionID,
			}
			err := svc.CreatePendingChurch(context.Background(), church)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error=%v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetChurchBySessionID(t *testing.T) {
	cases := []struct {
		name       string
		getCtxErr  error
		wantErr    bool
		wantChurch *models.Church
	}{
		{
			name:      "success",
			getCtxErr: nil,
			wantErr:   false,
			wantChurch: &models.Church{
				Name:            testChurchName,
				AdminEmail:      testChurchEmail,
				StripeSessionID: testSessionID,
				Status:          "pending",
			},
		},
		{
			name:      "not found",
			getCtxErr: sql.ErrNoRows,
			wantErr:   true,
		},
		{
			name:      "db error",
			getCtxErr: errors.New("connection refused"),
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := newChurchService(
				nil,
				func(ctx context.Context, dest interface{}, q string, args ...interface{}) error {
					if tc.getCtxErr != nil {
						return tc.getCtxErr
					}
					c := dest.(*models.Church)
					c.Name = testChurchName
					c.AdminEmail = testChurchEmail
					c.StripeSessionID = testSessionID
					c.Status = "pending"
					return nil
				},
			)

			church, err := svc.GetChurchBySessionID(context.Background(), testSessionID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error=%v, got %v", tc.wantErr, err)
			}
			if !tc.wantErr {
				if church.Name != tc.wantChurch.Name {
					t.Errorf("expected name %q, got %q", tc.wantChurch.Name, church.Name)
				}
				if church.StripeSessionID != tc.wantChurch.StripeSessionID {
					t.Errorf("expected session ID %q, got %q", tc.wantChurch.StripeSessionID, church.StripeSessionID)
				}
			}
		})
	}
}

func TestActivateChurch(t *testing.T) {
	cases := []struct {
		name         string
		getCtxErr    error // error on first GetContext (GetChurchBySessionID)
		namedExecErr error
		wantErr      bool
	}{
		{"success", nil, nil, false},
		{"church not found", sql.ErrNoRows, nil, true},
		{"update error", nil, errors.New("db error"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			callCount := 0
			var capturedCode string

			svc := newChurchService(
				func(q string, arg interface{}) (sql.Result, error) {
					if m, ok := arg.(map[string]interface{}); ok {
						if code, ok := m["church_code"].(string); ok {
							capturedCode = code
						}
					}
					return nil, tc.namedExecErr
				},
				func(ctx context.Context, dest interface{}, q string, args ...interface{}) error {
					callCount++
					if callCount == 1 {
						// GetChurchBySessionID
						if tc.getCtxErr != nil {
							return tc.getCtxErr
						}
						c := dest.(*models.Church)
						c.Name = testChurchName
						c.AdminEmail = testChurchEmail
						c.StripeSessionID = testSessionID
						c.Status = "pending"
						return nil
					}
					// uniqueness check — return ErrNoRows to indicate the code is unique
					return sql.ErrNoRows
				},
			)

			err := svc.ActivateChurch(context.Background(), testSessionID, "cus_123", "sub_123")
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error=%v, got: %v", tc.wantErr, err)
			}

			if !tc.wantErr {
				// Verify church code format: "XXXX-YYYY" (4 letters, dash, 4 alphanumeric)
				parts := strings.Split(capturedCode, "-")
				if len(parts) != 2 {
					t.Fatalf("expected church code with one dash, got %q", capturedCode)
				}
				if len(parts[0]) != 4 {
					t.Errorf("expected 4-char prefix, got %q", parts[0])
				}
				if len(parts[1]) != 4 {
					t.Errorf("expected 4-char suffix, got %q", parts[1])
				}
				if parts[0] != testChurchCodeFmt {
					t.Errorf("expected prefix %q, got %q", testChurchCodeFmt, parts[0])
				}
			}
		})
	}
}
