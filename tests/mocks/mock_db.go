package mocks

import "database/sql"
import ctx "context"

type MockDb struct {
	MockedFn            func(string, interface{}) (sql.Result, error)
	MockedSelectContext func(ctx.Context, interface{}, string, ...interface{}) error
}

func (m *MockDb) NamedExec(query string, args interface{}) (sql.Result, error) {
	return m.MockedFn(query, args)
}

func (m *MockDb) SelectContext(c ctx.Context, v interface{}, query string, args ...interface{}) error {
	return m.MockedSelectContext(c, v, query)
}
