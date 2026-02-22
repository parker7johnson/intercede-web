package api

import (
	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/utils"
	"github.com/supabase-community/gotrue-go/types"
)

type submissionService interface {
	CreatePrayerRequest(*models.Submission) error
	CreatePraiseReport(*models.Submission) error
}

type authService interface {
	Login(string, string) (types.Session, error)
}

type ApiHandlers struct {
	submissionHandler submissionService
	authHandler       authService
	log               *utils.Logger
}

func NewApi(sh submissionService, ah authService) *ApiHandlers {
	return &ApiHandlers{
		submissionHandler: sh,
		log:               utils.NewLogger("IntercedeAPILogger"),
		authHandler:       ah,
	}
}
