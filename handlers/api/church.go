package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/utils"
	stripe "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
	"github.com/supabase-community/gotrue-go/types"
)

type churchService interface {
	CreatePendingChurch(ctx context.Context, church *models.Church) error
	ActivateChurch(ctx context.Context, sessionID, customerID, subscriptionID string) error
	GetChurchBySessionID(ctx context.Context, sessionID string) (*models.Church, error)
}

// stripeSessionCreator abstracts Stripe checkout session creation for testability.
type stripeSessionCreator interface {
	New(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
}

// churchAuthClient abstracts Supabase auth signup for testability.
type churchAuthClient interface {
	Signup(req types.SignupRequest) (*types.SignupResponse, error)
}

type realStripeSession struct{}

func (r *realStripeSession) New(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	return session.New(params)
}

type ChurchAPIHandlers struct {
	churchService churchService
	stripe        stripeSessionCreator
	webhookSecret string
	priceID       string
	baseURL       string
	auth          churchAuthClient
	log           *utils.Logger
}

// NewChurchAPI constructs the handler. Pass nil for sc to use the real Stripe client.
func NewChurchAPI(cs churchService, sc stripeSessionCreator, webhookSecret, priceID, baseURL string, auth churchAuthClient) *ChurchAPIHandlers {
	if sc == nil {
		sc = &realStripeSession{}
	}
	return &ChurchAPIHandlers{
		churchService: cs,
		stripe:        sc,
		webhookSecret: webhookSecret,
		priceID:       priceID,
		baseURL:       baseURL,
		auth:          auth,
		log:           utils.NewLogger("ChurchAPIHandlers"),
	}
}

func (ch *ChurchAPIHandlers) StartCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("churchname")
	email := r.FormValue("adminemail")
	password := r.FormValue("password")
	if name == "" || email == "" || password == "" {
		http.Error(w, "Church name and admin email are required", http.StatusBadRequest)
		return
	}

	res, err := ch.auth.Signup(types.SignupRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		http.Error(w, "Error creating user profile", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    res.AccessToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(ch.priceID),
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: stripe.String(email),
		SuccessURL:    stripe.String(ch.baseURL + "/church/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:     stripe.String(ch.baseURL + "/church/cancel"),
	}

	stripeSession, err := ch.stripe.New(params)
	if err != nil {
		ch.log.Error("creating stripe session: %v", err)
		http.Error(w, "Failed to create checkout session", http.StatusInternalServerError)
		return
	}

	church := &models.Church{
		Name:            name,
		AdminEmail:      email,
		StripeSessionID: stripeSession.ID,
		UserId:          res.User.ID.String(),
	}
	if err := ch.churchService.CreatePendingChurch(r.Context(), church); err != nil {
		ch.log.Error("saving pending church: %v", err)
		http.Error(w, "Failed to save church record", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, stripeSession.URL, http.StatusSeeOther)
}

func (ch *ChurchAPIHandlers) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	sig := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEventWithOptions(body, sig, ch.webhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		ch.log.Error("webhook signature verification failed: %v", err)
		http.Error(w, "Invalid webhook signature", http.StatusBadRequest)
		return
	}

	if event.Type == stripe.EventTypeCheckoutSessionCompleted {
		var checkoutSession stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
			ch.log.Error("parsing checkout session: %v", err)
			http.Error(w, "Failed to parse event data", http.StatusInternalServerError)
			return
		}

		customerID := ""
		if checkoutSession.Customer != nil {
			customerID = checkoutSession.Customer.ID
		}
		subscriptionID := ""
		if checkoutSession.Subscription != nil {
			subscriptionID = checkoutSession.Subscription.ID
		}

		if err := ch.churchService.ActivateChurch(r.Context(), checkoutSession.ID, customerID, subscriptionID); err != nil {
			ch.log.Error("activating church: %v", err)
			http.Error(w, "Failed to activate church", http.StatusInternalServerError)
			return
		}

		ch.log.Success("church activated for session %s", checkoutSession.ID)
	}

	w.WriteHeader(http.StatusOK)
}
