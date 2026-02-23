package mocks

import (
	"context"

	"github.com/parkerjohnson/intercede/services/models"
	stripe "github.com/stripe/stripe-go/v82"
)

// MockChurchService implements api.churchService and webhandlers.churchLookup.
type MockChurchService struct {
	CreatePendingFn  func(ctx context.Context, church *models.Church) error
	ActivateFn       func(ctx context.Context, sessionID, customerID, subscriptionID string) error
	GetBySessionIDFn func(ctx context.Context, sessionID string) (*models.Church, error)
	GetByUserIDFn    func(ctx context.Context, userID string) (*models.Church, error)
}

func (m *MockChurchService) CreatePendingChurch(ctx context.Context, church *models.Church) error {
	return m.CreatePendingFn(ctx, church)
}

func (m *MockChurchService) ActivateChurch(ctx context.Context, sessionID, customerID, subscriptionID string) error {
	return m.ActivateFn(ctx, sessionID, customerID, subscriptionID)
}

func (m *MockChurchService) GetChurchBySessionID(ctx context.Context, sessionID string) (*models.Church, error) {
	return m.GetBySessionIDFn(ctx, sessionID)
}

func (m *MockChurchService) GetChurchByUserID(ctx context.Context, userID string) (*models.Church, error) {
	return m.GetByUserIDFn(ctx, userID)
}

// MockStripeCheckout implements api.stripeSessionCreator.
type MockStripeCheckout struct {
	NewFn func(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
}

func (m *MockStripeCheckout) New(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	return m.NewFn(params)
}
