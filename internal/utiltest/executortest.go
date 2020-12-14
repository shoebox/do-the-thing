package utiltest

import (
	"context"
	"dothething/internal/api"

	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
	mock.Mock
}

// CommandContext allow to execute a command with Context
func (m *MockExecutor) CommandContext(ctx context.Context, cmd string, args ...string) api.Cmd {
	c := m.Called(ctx, cmd, args)
	return c.Get(0).(api.Cmd)
}

func (m *MockExecutor) MockCommandContext(cmd string, args []string, res string, err error) {
	c := new(MockExecutorCmd)
	c.On("Output").Return(res, err)

	m.On("CommandContext", mock.Anything, cmd, args).Return(c)
}
