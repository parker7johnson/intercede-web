package models

import (
	"time"
)

type Submission struct {
	Title       string `db:"title"`
	Body        string `db:"body"`
	ContactInfo string `db:"contact_info"`
	ChurchCode  string `db:"church_code"`
}

type SubmissionDB struct {
	ID             string     `db:"id"`
	CreatedAt      time.Time  `db:"created_at"`
	Title          string     `db:"title"`
	Body           string     `db:"body"`
	ContactInfo    *string    `db:"contact_info"`
	ChurchCode     string     `db:"church_code"`
	SubmitterName  *string    `db:"submitter_name"`
	Status         *string    `db:"status"`
	UpdatedAt      *time.Time `db:"updated_at"`
	SubmissionType string     `db:"submission_type"`
}
