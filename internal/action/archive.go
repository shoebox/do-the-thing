package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

func NewArchive(api *api.API) api.Action {
	return ActionArchive{api}
}

type ActionArchive struct {
	*api.API
}

func (a ActionArchive) Run(ctx context.Context) error {
	xce := xcode.ParseXCodeBuildError(a.build(ctx))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a ActionArchive) build(ctx context.Context) error {
	log.Info().Msg("Archiving")

	if err := a.API.SignatureService.Run(ctx); err != nil {
		return err
	}

	// defer deletion of the keychain
	defer a.API.KeyChain.Delete(ctx)

	// The archiving arguments
	args := []string{
		a.API.BuildService.GetArg(),
		a.API.Config.Path,
		xcode.ActionArchive,
		xcode.FlagScheme, a.API.Config.Scheme,
		xcode.FlagConfiguration, a.API.Config.Configuration,
		xcode.FlagArchivePath, a.API.PathService.Archive(),
		a.API.PathService.ObjRoot(),
		a.API.PathService.SymRoot(),
	}

	return RunCmd(a.API.Exec.CommandContext(ctx, xcode.Cmd, args...))
}
