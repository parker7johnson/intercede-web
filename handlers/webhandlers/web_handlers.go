package webhandlers

import (
	"context"

	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/utils"
)

type churchLookup interface {
	GetChurchBySessionID(ctx context.Context, sessionID string) (*models.Church, error)
	GetChurchByUserID(ctx context.Context, userID string) (*models.Church, error)
}

type submissionLookup interface {
	GetSubmissionsByChurchCode(string) (*[]models.SubmissionDB, error)
	CreatePrayerRequest(*models.Submission) error
	CreatePraiseReport(*models.Submission) error
	GetChurchCodes() ([]string, error)
}

// WebHandlers holds dependencies for HTTP handlers
type WebHandlers struct {
	churchService churchLookup
	submissionService submissionLookup
	log utils.Logger
}

// NewWeb creates a new WebHandlers instance with dependencies
func NewWeb(cs churchLookup, ss submissionLookup) *WebHandlers {
	return &WebHandlers{
		churchService:     cs,
		submissionService: ss,
		log: *utils.NewLogger("WebApiLogger"),
	}
}
