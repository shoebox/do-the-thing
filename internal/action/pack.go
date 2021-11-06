package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"

	"github.com/rs/zerolog/log"
)

type actionPackage struct {
	*api.API
}

func NewActionPackage(api *api.API) api.Action {
	return actionPackage{API: api}
}

func (a actionPackage) Run(ctx context.Context) error {
	err := xcode.ParseXCodeBuildError(a.pack(ctx))
	if err != nil {
		log.Err(err)
	}

	return err
}

func (a actionPackage) pack(ctx context.Context) error {
	// defer deletion of the keychain
	defer a.API.KeyChain.Delete(ctx)

	// Resolving signature
	if err := a.API.SignatureService.Run(ctx); err != nil {
		return err
	}

	// Compute export options plist
	if err := a.API.ExportOptionService.Compute(); err != nil {
		return err
	}

	// The arguments
	args := []string{
		xcode.ActionPackage,
		xcode.FlagArchivePath, a.API.PathService.Archive(),
		xcode.FlagExportPath, a.API.PathService.Package(),
		xcode.FlagExportOptionsPlist, a.API.PathService.ExportPList(),
		a.API.PathService.ObjRoot(),
		// a.API.PathService.SymRoot(),
	}

	cmd, err := a.API.Exec.XCodeCommandContext(ctx, args...)
	if err != nil {
		return err
	}

	return RunCmd(*cmd)
}
