// Package util provide low level util function to interact with the system
package util

import (
	"context"
	"dothething/internal/api"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type executor struct {
	*api.API
}

// NewExecutor create a new instance of the implemented Cmd interface
func NewExecutor(api *api.API) api.Executor {
	return &executor{api}
}

// CommandContext run a command with context
func (e *executor) CommandContext(ctx context.Context, cmd string, args ...string) api.Cmd {
	log.Info().
		Str("Cmd", cmd).
		Strs("Args", args).
		Msg("Running command with context")

	return (*cmdWrapper)(exec.CommandContext(ctx, cmd, args...))
}

// CommandContext run a command with context
func (e *executor) XCodeCommandContext(ctx context.Context, args ...string) (*api.Cmd, error) {
	log.Info().
		Strs("Args", args).
		Msg("Running XCode command with context")

	i, err := e.API.XcodeSelectService.Find(ctx)
	if err != nil {
		return nil, err
	}

	// executing the command
	cmd := e.CommandContext(ctx, "xcodebuild", args...)
	cmd.SetEnv([]string{"DEVELOPER_DIR", i.DevPath})

	return &cmd, nil
}
