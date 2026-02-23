package mocks

// MockAuthVerifier implements middleware.TokenVerifier for use in tests.
type MockAuthVerifier struct {
	VerifyTokenFn func(token string) (string, error)
}

func (m *MockAuthVerifier) VerifyToken(token string) (string, error) {
	return m.VerifyTokenFn(token)
}
