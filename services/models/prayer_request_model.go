package models

type Submission struct {
	Title       string `db:"title"`
	Body        string `db:"body"`
	ContactInfo string `db:"contact_info"`
	ChurchCode  string `db:"church_code"`
}
