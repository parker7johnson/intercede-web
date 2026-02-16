package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/utils"
)

type SubmissionHandler struct {
	log *utils.Logger
	db  *sqlx.DB
}

func New(db *sqlx.DB) *SubmissionHandler {
	return &SubmissionHandler{
		log: utils.NewLogger("PrayerRequestService"),
		db:  db,
	}
}

func (sh *SubmissionHandler) CreatePrayerRequest(prayerReq *models.Submission) error {
	query := `insert into prayer_requests (title, body, contact_info, church_code)
	values (:title, :body, :contact_info, :church_code)`

	_, err := sh.db.NamedExec(query, prayerReq)

	if err != nil {
		return err
	}
	return nil
}

func (sh *SubmissionHandler) CreatePraiseReport(prayerReq *models.Submission) error {
	query := `insert into praise_reports (title, body, contact_info, church_code)
values (:title, :body, :contact_info, :church_code)`

	_, err := sh.db.NamedExec(query, prayerReq)

	if err != nil {
		return err
	}
	return nil
}
