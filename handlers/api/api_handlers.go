package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/utils"
)

type ApiHandlers struct {
	db *sqlx.DB
	sh *services.SubmissionHandler 
	log *utils.Logger
}

func NewApi(db *sqlx.DB) *ApiHandlers {
	return &ApiHandlers{
			db: db,
			sh: services.New(db),
			log: utils.NewLogger("IntercedeAPILogger"),
	}
}
