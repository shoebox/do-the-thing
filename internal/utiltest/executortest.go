package utiltest

import (
	"context"
	"dothething/internal/util"

	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) MockCommandContext(cmd string, args []string, res string, err error) {
	c := new(MockExecutorCmd)
	c.
		On("Output").
		Return(res, err)

	m.On("CommandContext",
		mock.Anything,
		cmd, args).
		Return(c)
}

// ContextExec mock execute the program with the provided arguments and context
func (m *MockExecutor) CommandContext(ctx context.Context, cmd string, args ...string) util.Cmd {
	r := m.Called(ctx, cmd, args)
	return r.Get(0).(util.Cmd)
}
