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
	api.API
	*api.Config
}

func NewActionRun(api api.API, cfg *api.Config) api.Action {
	return actionRunTest{api, cfg}
}

func (a actionRunTest) Run(ctx context.Context) error {
	log.Info().Msg("Running unit tests")

	if err := a.API.SignatureService.Run(ctx); err != nil {
		return err
	}

	// defer deletion of the keychain
	defer a.API.KeyChain.Delete(ctx)

	// Creating a temp folder to contains the test results
	outputPath := a.API.PathService.XCResult()

	xce := xcode.ParseXCodeBuildError(a.runXCodebuildTest(ctx, outputPath))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce

}
func (a actionRunTest) runXCodebuildTest(ctx context.Context, path string) error {
	fmt.Println(color.BlueString("Running test on %v (%v)", a.Config.Destination.Name, a.Config.Destination.ID))

	return RunCmd(a.API.Exec.CommandContext(ctx,
		xcode.Cmd,
		a.API.BuildService.GetArg(),
		a.API.Config.Path,
		xcode.ActionTest,
		xcode.FlagScheme, a.Config.Scheme,
		xcode.FlagDestination, fmt.Sprintf("id=%s", a.Config.Destination.ID),
		xcode.FlagResultBundlePath, path,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO"))
}
