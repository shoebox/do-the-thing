package util

import (
	"context"
	"dothething/internal/api"
	"fmt"
	"os/exec"
)

type executor struct{}

// NewExecutor create a new instance of the implemented Cmd interface
func NewExecutor() api.Executor {
	return &executor{}
}

// CommandContext run a command with context
func (executor *executor) CommandContext(ctx context.Context, cmd string, args ...string) api.Cmd {
	fmt.Println(cmd, args)
	return (*cmdWrapper)(exec.CommandContext(ctx, cmd, args...))
}
