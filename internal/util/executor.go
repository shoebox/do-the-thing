package util

import (
	"context"
	"os/exec"
)

// Executor command helper base interface
type Executor interface {
	// CommandContext allow to execute a command with Context
	CommandContext(ctx context.Context, cmd string, args ...string) Cmd
}

type executor struct{}

// NewExecutor create a new instance of the implemented Cmd interface
func NewExecutor() Executor {
	return &executor{}
}

// CommandContext run a command with context
func (executor *executor) CommandContext(ctx context.Context, cmd string, args ...string) Cmd {
	return (*cmdWrapper)(exec.CommandContext(ctx, cmd, args...))
}
