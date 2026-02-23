package mocks

import (
	"context"
	"database/sql"
)

// MockChurchDb implements the churchDatabase interface used by ChurchService.
type MockChurchDb struct {
	NamedExecFn  func(query string, arg interface{}) (sql.Result, error)
	GetContextFn func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (m *MockChurchDb) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return m.NamedExecFn(query, arg)
}

func (m *MockChurchDb) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.GetContextFn(ctx, dest, query, args...)
}
