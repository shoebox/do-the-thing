package utiltest

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockExec struct {
	mock.Mock
}

var MockExecutor *MockExec
var Exec *MockExec

func SetupMockExec() {
	Exec = new(MockExec)
}

func TearDownMockExec() {
	Exec = nil
}

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

func (m *MockExec) Exec(path *string, name string, extra ...string) ([]byte, error) {
	args := m.Called(name, extra)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return []byte(args.String(0)), nil
}
