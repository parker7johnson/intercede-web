package models

type Church struct {
	ID                   string `db:"id"`
	Name                 string `db:"name"`
	AdminEmail           string `db:"admin_email"`
	ChurchCode           string `db:"church_code"`
	StripeSessionID      string `db:"stripe_session_id"`
	StripeCustomerID     string `db:"stripe_customer_id"`
	StripeSubscriptionID string `db:"stripe_subscription_id"`
	Status               string `db:"status"`
	UserId							 string `db:"user_id"`
}
