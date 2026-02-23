package tests

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/parkerjohnson/intercede/handlers/api"
	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/tests/mocks"
	stripe "github.com/stripe/stripe-go/v82"
	gotruetypes "github.com/supabase-community/gotrue-go/types"
)

const (
	testWebhookSecret = "whsec_testsecret"
	testPriceID       = "price_test123"
	testBaseURL       = "http://localhost:3000"
	testStripeURL     = "https://checkout.stripe.com/pay/cs_test"
	testStripeSession = "cs_test_session123"
)

// stripeSignature computes a valid Stripe-Signature header value for testing.
func stripeSignature(secret, payload string, ts int64) string {
	signed := fmt.Sprintf("%d.%s", ts, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signed))
	sig := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("t=%d,v1=%s", ts, sig)
}

func newMockAuth(err error) *mocks.MockChurchAuthClient {
	return &mocks.MockChurchAuthClient{
		SignupFn: func(req gotruetypes.SignupRequest) (*gotruetypes.SignupResponse, error) {
			if err != nil {
				return nil, err
			}
			return &gotruetypes.SignupResponse{Session: gotruetypes.Session{AccessToken: testToken}}, nil
		},
	}
}

func newCheckoutHandler(svc *mocks.MockChurchService, stripeCheckout *mocks.MockStripeCheckout) *api.ChurchAPIHandlers {
	return api.NewChurchAPI(svc, stripeCheckout, testWebhookSecret, testPriceID, testBaseURL, newMockAuth(nil))
}

// newWebhookHandler creates a handler with a no-op stripe creator (webhook tests don't use it).
func newWebhookHandler(svc *mocks.MockChurchService) *api.ChurchAPIHandlers {
	return api.NewChurchAPI(svc, nil, testWebhookSecret, testPriceID, testBaseURL, newMockAuth(nil))
}

func newCheckoutRequest(method string, name, email string) *http.Request {
	form := url.Values{}
	form.Set("churchname", name)
	form.Set("adminemail", email)
	form.Set("password", testPassword)
	r := httptest.NewRequest(method, "/api/church/checkout", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func TestStartCheckout(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		churchName   string
		adminEmail   string
		authErr      error
		stripeErr    error
		dbErr        error
		wantStatus   int
		wantRedirect string
	}{
		{
			name:         "success",
			method:       http.MethodPost,
			churchName:   testChurchName,
			adminEmail:   testChurchEmail,
			wantStatus:   http.StatusSeeOther,
			wantRedirect: testStripeURL,
		},
		{
			name:       "bad method",
			method:     http.MethodGet,
			churchName: testChurchName,
			adminEmail: testChurchEmail,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "missing church name",
			method:     http.MethodPost,
			churchName: "",
			adminEmail: testChurchEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing admin email",
			method:     http.MethodPost,
			churchName: testChurchName,
			adminEmail: "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "auth error",
			method:     http.MethodPost,
			churchName: testChurchName,
			adminEmail: testChurchEmail,
			authErr:    errors.New("supabase unavailable"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "stripe error",
			method:     http.MethodPost,
			churchName: testChurchName,
			adminEmail: testChurchEmail,
			stripeErr:  errors.New("stripe unavailable"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "db error",
			method:     http.MethodPost,
			churchName: testChurchName,
			adminEmail: testChurchEmail,
			dbErr:      errors.New("db error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stripeCheckout := &mocks.MockStripeCheckout{
				NewFn: func(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
					return &stripe.CheckoutSession{
						ID:  testStripeSession,
						URL: testStripeURL,
					}, tc.stripeErr
				},
			}
			churchSvc := &mocks.MockChurchService{
				CreatePendingFn: func(ctx context.Context, church *models.Church) error {
					return tc.dbErr
				},
			}
			mockAuth := newMockAuth(tc.authErr)

			handler := api.NewChurchAPI(churchSvc, stripeCheckout, testWebhookSecret, testPriceID, testBaseURL, mockAuth)
			w := httptest.NewRecorder()
			handler.StartCheckout(w, newCheckoutRequest(tc.method, tc.churchName, tc.adminEmail))

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.wantRedirect != "" {
				loc := w.Header().Get("Location")
				if loc != tc.wantRedirect {
					t.Errorf("expected redirect to %q, got %q", tc.wantRedirect, loc)
				}
			}
		})
	}
}

func TestHandleStripeWebhook(t *testing.T) {
	validPayload := `{"id":"evt_test","type":"checkout.session.completed","data":{"object":{"id":"cs_test_session123","object":"checkout.session"}}}`
	ts := time.Now().Unix()
	validSig := stripeSignature(testWebhookSecret, validPayload, ts)

	cases := []struct {
		name        string
		payload     string
		sig         string
		activateErr error
		wantStatus  int
	}{
		{
			name:       "valid event activates church",
			payload:    validPayload,
			sig:        validSig,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid signature",
			payload:    validPayload,
			sig:        "t=1234,v1=badsignature",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "activate error returns 500",
			payload:     validPayload,
			sig:         validSig,
			activateErr: errors.New("db error"),
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name: "unhandled event type returns 200",
			payload: `{"id":"evt_other","type":"customer.created","data":{"object":{}}}`,
			sig: func() string {
				p := `{"id":"evt_other","type":"customer.created","data":{"object":{}}}`
				return stripeSignature(testWebhookSecret, p, ts)
			}(),
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var activateCalled bool
			churchSvc := &mocks.MockChurchService{
				ActivateFn: func(ctx context.Context, sessionID, customerID, subscriptionID string) error {
					activateCalled = true
					return tc.activateErr
				},
			}

			r := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", strings.NewReader(tc.payload))
			r.Header.Set("Stripe-Signature", tc.sig)
			w := httptest.NewRecorder()

			newWebhookHandler(churchSvc).HandleStripeWebhook(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.name == "valid event activates church" && !activateCalled {
				t.Error("expected ActivateChurch to be called")
			}
		})
	}
}
