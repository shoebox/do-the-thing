package action

import (
	"context"
	"dothething/internal/util"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

type ActionArchive interface {
	Run(ctx context.Context, config xcode.Config) error
}

func NewArchive(xcode xcode.BuildService, exec util.Executor) ActionArchive {
	return actionArchive{exec, xcode}
}

type actionArchive struct {
	exec  util.Executor
	xcode xcode.BuildService
}

func (a actionArchive) Run(ctx context.Context, config xcode.Config) error {
	xce := xcode.ParseXCodeBuildError(a.build(ctx, config))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionArchive) build(ctx context.Context, config xcode.Config) error {
	log.Info().Msg("Archiving")
	return RunCmd(a.exec.CommandContext(ctx,
		xcode.Cmd,
		a.xcode.GetArg(),
		a.xcode.GetProjectPath(),
		xcode.ActionArchive,
		xcode.FlagScheme, config.Scheme,
		"CODE_SIGNING_ALLOWED=NO"))
}
