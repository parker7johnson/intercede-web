package mocks

import (
	"github.com/supabase-community/gotrue-go/types"
)

type MockSupabaseClient struct {
	SignInFn      func(string, string) (types.Session, error)
	AutoRefreshFn func(types.Session)
}

func (msc *MockSupabaseClient) SignInWithEmailPassword(email, password string) (types.Session, error) {
	return msc.SignInFn(email, password)
}
func (msc *MockSupabaseClient) EnableTokenAutoRefresh(session types.Session) {
	return
}
