// +build !test

package util

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockExec struct {
	mock.Mock
}

// ContextExec mock execute the program with the provided arguments and context
func (m *mockExec) ContextExec(ctx context.Context,
	name string,
	extra ...string) ([]byte, error) {

	l := []interface{}{ctx, name}
	ex := append(l, extra)

	args := m.Called(ex...)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return []byte(args.String(0)), nil
}

// Exec mock execute the program with the provided arguments
func (m *mockExec) Exec(path *string, name string, extra ...string) ([]byte, error) {
	args := m.Called(name, extra)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return []byte(args.String(0)), nil
}
