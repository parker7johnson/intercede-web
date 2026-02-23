package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"github.com/parkerjohnson/intercede/services/models"
	"github.com/parkerjohnson/intercede/utils"
)

const alphanumChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type churchDatabase interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type ChurchService struct {
	log *utils.Logger
	db  churchDatabase
}

func NewChurchService(db churchDatabase) *ChurchService {
	return &ChurchService{
		log: utils.NewLogger("ChurchService"),
		db:  db,
	}
}

// CreatePendingChurch inserts a church with status='pending'. StripeSessionID must be set.
func (cs *ChurchService) CreatePendingChurch(ctx context.Context, church *models.Church) error {
	query := `INSERT INTO churches (name, admin_email, stripe_session_id, user_id, status)
	VALUES (:name, :admin_email, :stripe_session_id, :user_id, 'pending')`
	_, err := cs.db.NamedExec(query, church)
	return err
}

// ActivateChurch generates a church code and sets status='active' for the given session.
func (cs *ChurchService) ActivateChurch(ctx context.Context, sessionID, customerID, subscriptionID string) error {
	code, err := cs.generateUniqueChurchCode(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("generating church code: %w", err)
	}

	query := `UPDATE churches
	          SET status = 'active',
	              church_code = :church_code,
	              stripe_customer_id = :stripe_customer_id,
	              stripe_subscription_id = :stripe_subscription_id
	          WHERE stripe_session_id = :stripe_session_id`
	_, err = cs.db.NamedExec(query, map[string]interface{}{
		"church_code":            code,
		"stripe_customer_id":     customerID,
		"stripe_subscription_id": subscriptionID,
		"stripe_session_id":      sessionID,
	})
	return err
}

// GetChurchBySessionID fetches a church by its Stripe checkout session ID.
func (cs *ChurchService) GetChurchBySessionID(ctx context.Context, sessionID string) (*models.Church, error) {
	var church models.Church
	query := `SELECT id, name, admin_email,
	                 COALESCE(church_code, '') AS church_code,
	                 stripe_session_id,
	                 COALESCE(stripe_customer_id, '') AS stripe_customer_id,
	                 COALESCE(stripe_subscription_id, '') AS stripe_subscription_id,
	                 status
	          FROM churches WHERE stripe_session_id = $1`
	if err := cs.db.GetContext(ctx, &church, query, sessionID); err != nil {
		return nil, err
	}
	return &church, nil
}

// GetChurchByUserID fetches a church by its owner's user ID.
func (cs *ChurchService) GetChurchByUserID(ctx context.Context, userID string) (*models.Church, error) {
	var church models.Church
	query := `SELECT id, name, admin_email,
	                 COALESCE(church_code, '') AS church_code,
	                 stripe_session_id,
	                 COALESCE(stripe_customer_id, '') AS stripe_customer_id,
	                 COALESCE(stripe_subscription_id, '') AS stripe_subscription_id,
	                 status
	          FROM churches WHERE user_id = $1`
	if err := cs.db.GetContext(ctx, &church, query, userID); err != nil {
		return nil, err
	}
	return &church, nil
}

func (cs *ChurchService) generateUniqueChurchCode(ctx context.Context, sessionID string) (string, error) {
	church, err := cs.GetChurchBySessionID(ctx, sessionID)
	if err != nil {
		return "", err
	}

	prefix := namePrefix(church.Name)

	for range 10 {
		suffix, err := randomAlphanumeric(4)
		if err != nil {
			return "", err
		}
		code := prefix + "-" + suffix

		var existingID string
		err = cs.db.GetContext(ctx, &existingID, `SELECT id FROM churches WHERE church_code = $1`, code)
		if err == sql.ErrNoRows {
			return code, nil
		}
		if err != nil {
			return "", err
		}
		// collision — retry
	}
	return "", fmt.Errorf("could not generate unique church code after 10 attempts")
}

func namePrefix(name string) string {
	var letters []rune
	for _, r := range strings.ToUpper(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			letters = append(letters, r)
			if len(letters) == 4 {
				break
			}
		}
	}
	for len(letters) < 4 {
		letters = append(letters, 'X')
	}
	return string(letters)
}

func randomAlphanumeric(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	result := make([]byte, n)
	for i, v := range b {
		result[i] = alphanumChars[int(v)%len(alphanumChars)]
	}
	return string(result), nil
}
