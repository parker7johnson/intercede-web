package mocks

import "github.com/supabase-community/gotrue-go/types"

// MockChurchAuthClient implements api.churchAuthClient for use in tests.
type MockChurchAuthClient struct {
	SignupFn func(req types.SignupRequest) (*types.SignupResponse, error)
}

func (m *MockChurchAuthClient) Signup(req types.SignupRequest) (*types.SignupResponse, error) {
	return m.SignupFn(req)
}
