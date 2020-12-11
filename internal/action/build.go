package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

func NewBuild(api api.API, cfg *api.Config) api.Action {
	return actionBuild{api, cfg}
}

type actionBuild struct {
	api.API
	*api.Config
}

func (a actionBuild) Run(ctx context.Context) error {
	xce := xcode.ParseXCodeBuildError(a.build(ctx))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionBuild) build(ctx context.Context) error {
	log.Info().Msg("Building")
	return RunCmd(a.API.Exec().CommandContext(ctx,
		xcode.Cmd,
		a.API.XCodeBuildService().GetArg(),
		a.Config.Path,
		xcode.ActionBuild,
		xcode.FlagScheme, a.Config.Scheme,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO"))
}
