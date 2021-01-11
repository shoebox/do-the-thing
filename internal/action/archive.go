package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

func NewArchive(api api.API, cfg *api.Config) api.Action {
	return actionArchive{api, cfg}
}

type actionArchive struct {
	api.API
	*api.Config
}

func (a actionArchive) Run(ctx context.Context) error {

	xce := xcode.ParseXCodeBuildError(a.build(ctx))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionArchive) build(ctx context.Context) error {
	log.Info().Msg("Archiving")

	if err := a.API.SignatureService().Run(ctx); err != nil {
		return err
	}

	// defer deletion of the keychain
	defer a.API.KeyChainService().Delete(ctx)

	// The archiving arguments
	args := []string{
		a.API.XCodeBuildService().GetArg(),
		a.API.XCodeBuildService().GetProjectPath(),
		xcode.ActionArchive,
		xcode.FlagScheme, a.Config.Scheme,
		xcode.FlagConfiguration, a.Config.Configuration,
		xcode.FlagArchivePath, a.API.PathService().Archive(),
		a.API.PathService().ObjRoot(),
		a.API.PathService().SymRoot(),
	}

	return RunCmd(a.API.Exec().CommandContext(ctx, xcode.Cmd, args...))
}
