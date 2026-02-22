package services

import (
	"github.com/parkerjohnson/intercede/utils"
	"github.com/supabase-community/gotrue-go/types"
)

type supabaseClient interface {
	SignInWithEmailPassword(email, password string) (types.Session, error)
	EnableTokenAutoRefresh(session types.Session)
}

type AuthService struct {
	log    *utils.Logger
	Client supabaseClient
}

func NewAuthService(client supabaseClient) *AuthService {
	return &AuthService{
		log:    utils.NewLogger("AuthApiLogger"),
		Client: client,
	}
}

func (as *AuthService) Login(email, password string) (types.Session, error) {
	session, err := as.Client.SignInWithEmailPassword(email, password)

	if err != nil {
		as.log.Error("Supabase auth error: %v - email: %s", err, email)
	}
  as.Client.EnableTokenAutoRefresh(session)
	return session, err

}
