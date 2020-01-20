package utiltest

import (
	"github.com/stretchr/testify/mock"
)

type MockExec struct {
	mock.Mock
}

var MockExecutor *MockExec

func (m *MockExec) Exec(name string, extra ...string) ([]byte, error) {
	args := m.Called(name, extra)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return []byte(args.String(0)), nil
}
