package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

func NewBuild(api *api.API) api.Action {
	return actionBuild{api}
}

type actionBuild struct {
	*api.API
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
	args := []string{
		a.API.BuildService.GetArg(),
		a.API.Config.Path,
		xcode.ActionBuild,
		xcode.FlagScheme, a.API.Config.Scheme,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO",
	}

	cmd, err := a.API.Exec.XCodeCommandContext(ctx, args...)
	if err != nil {
		return err
	}

	return RunCmd(*cmd)
}
