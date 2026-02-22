package mocks

import "database/sql"

type MockDb struct {
	MockedFn func(string, interface{}) (sql.Result, error)
}

func (m *MockDb) NamedExec(query string, args interface{}) (sql.Result, error) {
	return m.MockedFn(query, args)
}
