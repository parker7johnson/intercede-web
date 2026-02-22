package mocks

// MockAuthVerifier implements middleware.TokenVerifier for use in tests.
type MockAuthVerifier struct {
	VerifyTokenFn func(token string) error
}

func (m *MockAuthVerifier) VerifyToken(token string) error {
	return m.VerifyTokenFn(token)
}
