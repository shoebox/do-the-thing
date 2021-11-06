package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"

	"github.com/rs/zerolog/log"
)

func NewArchive(api *api.API) api.Action {
	return ActionArchive{api}
}

type ActionArchive struct {
	*api.API
}

func (a ActionArchive) Run(ctx context.Context) error {
	err := xcode.ParseXCodeBuildError(a.archive(ctx))
	if err != nil {
		log.Err(err)
	}

	return err
}

func (a ActionArchive) archive(ctx context.Context) error {
	log.Info().Msg("Archiving")
	// defer deletion of the keychain
	defer a.API.KeyChain.Delete(ctx)

	// Resolving signature configuration
	if err := a.API.SignatureService.Run(ctx); err != nil {
		return err
	}

	// The archiving arguments
	args := []string{
		a.API.BuildService.GetArg(),
		a.API.Config.Path,
		xcode.ActionArchive,
		xcode.FlagScheme, a.API.Config.Scheme,
		xcode.FlagConfiguration, a.API.Config.Configuration,
		xcode.FlagArchivePath, a.API.PathService.Archive(),
		a.API.PathService.ObjRoot(),
		// a.API.PathService.SymRoot(),
	}

	//
	cmd, err := a.API.Exec.XCodeCommandContext(ctx, args...)
	if err != nil {
		return err
	}

	return RunCmd(*cmd)
}
