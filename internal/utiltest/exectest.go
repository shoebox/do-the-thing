package utiltest

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockExec struct {
	mock.Mock
}

var mockExecutor *MockExec
var exec *MockExec

func SetupMockExec() {
	exec = new(MockExec)
}

// ContextExec mock execute the program with the provided arguments and context
func (m *MockExec) ContextExec(ctx context.Context,
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
func (m *MockExec) Exec(path *string, name string, extra ...string) ([]byte, error) {
	args := m.Called(name, extra)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return []byte(args.String(0)), nil
}
