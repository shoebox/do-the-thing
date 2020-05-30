package action

import (
	"context"
	"dothething/internal/util"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

type ActionBuild interface {
	Run(ctx context.Context, config xcode.Config) error
}

func NewBuild(xcode xcode.BuildService, exec util.Executor) ActionBuild {
	return actionBuild{exec, xcode}
}

type actionBuild struct {
	exec  util.Executor
	xcode xcode.BuildService
}

func (a actionBuild) Run(ctx context.Context, config xcode.Config) error {
	xce := xcode.ParseXCodeBuildError(a.build(ctx, config))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionBuild) build(ctx context.Context, config xcode.Config) error {
	log.Info().Msg("Building")
	return RunCmd(a.exec.CommandContext(ctx,
		xcode.Cmd,
		a.xcode.GetArg(),
		a.xcode.GetProjectPath(),
		"build",
		xcode.FlagScheme, config.Scheme,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO"))
}
