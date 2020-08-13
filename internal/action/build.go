package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

func NewBuild(api api.API) api.Action {
	return actionBuild{api}
}

type actionBuild struct {
	api.API
	//exec  util.Executor
	//xcode xcode.BuildService
}

func (a actionBuild) Run(ctx context.Context, config api.Config) error {
	xce := xcode.ParseXCodeBuildError(a.build(ctx, config))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionBuild) build(ctx context.Context, config api.Config) error {
	log.Info().Msg("Building")
	return RunCmd(a.API.Exec().CommandContext(ctx,
		xcode.Cmd,
		a.API.XCodeBuildService().GetArg(),
		a.API.XCodeBuildService().GetProjectPath(),
		xcode.ActionBuild,
		xcode.FlagScheme, config.Scheme,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO"))
}
