package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"
	"fmt"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

type actionRunTest struct {
	*api.API
}

func NewActionRunTest(api *api.API) api.Action {
	return actionRunTest{api}
}

func (a actionRunTest) Run(ctx context.Context) error {
	log.Info().Msg("Running unit tests")

	if err := a.API.SignatureService.Run(ctx); err != nil {
		return err
	}

	// defer deletion of the keychain
	defer a.API.KeyChain.Delete(ctx)

	xce := xcode.ParseXCodeBuildError(a.runXCodebuildTest(ctx))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce

}
func (a actionRunTest) runXCodebuildTest(ctx context.Context) error {
	fmt.Println("run tests")
	// listing possible destinations
	dd, err := a.API.DestinationService.List(ctx, a.Config.Scheme)
	if err != nil {
		return err
	}

	// using the last destination for now as a test
	d := dd[len(dd)-1]

	// Booting the destination
	if err := a.API.DestinationService.Boot(ctx, d); err != nil {
		return err
	}

	defer a.API.DestinationService.ShutDown(ctx, d)

	cmd, err := a.API.Exec.XCodeCommandContext(ctx,
		xcode.ActionTest,
		a.API.BuildService.GetArg(),
		a.API.Config.Path,
		xcode.FlagResultBundlePath, a.API.PathService.XCResult(),
		xcode.FlagDerivedData, a.API.PathService.DerivedData(),
		xcode.FlagScheme, a.API.Config.Scheme,
		xcode.FlagDestination, fmt.Sprintf("id=%s", d.ID),
		xcode.FlagCodeCoverage, "YES",
	)
	fmt.Println(cmd, err)

	if err != nil {
		return err
	}

	return RunCmd(*cmd)
}
