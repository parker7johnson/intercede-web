package services

import (
	"context"
	"database/sql"

	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/utils"
)

type database interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type SubmissionHandler struct {
	log *utils.Logger
	db  database
}

func New(db database) *SubmissionHandler {
	return &SubmissionHandler{
		log: utils.NewLogger("PrayerRequestService"),
		db:  db,
	}
}

func (sh *SubmissionHandler) GetSubmissionsByChurchCode(churchCode string) (*[]models.SubmissionDB, error) {
	query := `select * from submissions where church_code = $1`
	var submissions []models.SubmissionDB
	err := sh.db.SelectContext(context.Background(), &submissions, query, churchCode)
	if err != nil {
		return nil, err
	}
	return &submissions, nil
}

func (sh *SubmissionHandler) GetChurchCodes() ([]string, error) {
	query := `select distinct church_code from churches`
	var codes []string
	err := sh.db.SelectContext(context.Background(), &codes, query)
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func (sh *SubmissionHandler) CreatePrayerRequest(prayerReq *models.Submission) error {
	query := `insert into submissions (title, body, contact_info, church_code) 
	values (:title, :body, :contact_info, :church_code)`

	_, err := sh.db.NamedExec(query, prayerReq)

	if err != nil {
		return err
	}
	return nil
}

func (sh *SubmissionHandler) CreatePraiseReport(prayerReq *models.Submission) error {
	query := `insert into submissions (title, body, contact_info, church_code) 
values (:title, :body, :contact_info, :church_code)`

	_, err := sh.db.NamedExec(query, prayerReq)

	if err != nil {
		return err
	}
	return nil
}
