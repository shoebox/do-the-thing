package util

import (
	"context"
	"dothething/internal/api"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type executor struct{}

// NewExecutor create a new instance of the implemented Cmd interface
func NewExecutor() api.Executor {
	return &executor{}
}

// CommandContext run a command with context
func (executor *executor) CommandContext(ctx context.Context, cmd string, args ...string) api.Cmd {
	log.Debug().
		Str("Cmd", cmd).
		Strs("Args", args).
		Msg("Running command with context")
	return (*cmdWrapper)(exec.CommandContext(ctx, cmd, args...))
}
