package webhandlers

import (
	"context"

	"github.com/parkerjohnson/intercede/services/models"
)

type churchLookup interface {
	GetChurchBySessionID(ctx context.Context, sessionID string) (*models.Church, error)
	GetChurchByUserID(ctx context.Context, userID string) (*models.Church, error)
}

// WebHandlers holds dependencies for HTTP handlers
type WebHandlers struct {
	churchService churchLookup
}

// NewWeb creates a new WebHandlers instance with dependencies
func NewWeb(cs churchLookup) *WebHandlers {
	return &WebHandlers{
		churchService: cs,
	}
}
