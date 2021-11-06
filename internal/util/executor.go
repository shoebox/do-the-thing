package util

import (
	"context"
	"dothething/internal/api"
	"fmt"
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
func (executor *executor) CommandContext(ctx context.Context, cmd string, args ...string) api.Cmd {
	log.Info().
		Str("Cmd", cmd).
		Strs("Args", args).
		Msg("Running command with context")

	return (*cmdWrapper)(exec.CommandContext(ctx, cmd, args...))
}

// CommandContext run a command with context
func (a *executor) XCodeCommandContext(ctx context.Context, args ...string) (*api.Cmd, error) {
	log.Info().
		Strs("Args", args).
		Msg("Running XCode command with context")

	// selecting the right version of XCode
	i, err := a.API.XcodeSelectService.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve the XCode installation (%v)", err)
	}

	// executing the command
	cmd := a.CommandContext(ctx, "xcodebuild", args...)
	cmd.SetEnv([]string{"DEVELOPER_DIR", i.DevPath})

	return &cmd, nil
}
