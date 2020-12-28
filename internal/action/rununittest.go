package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"fmt"

	"github.com/fatih/color"
)

type actionRunTest struct {
	api.API
	*api.Config
}

func NewActionRun(api api.API, cfg *api.Config) api.Action {
	return actionRunTest{api, cfg}
}

func (a actionRunTest) Run(ctx context.Context) error {
	// Creating a temp folder to contains the test results
	path, err := util.TempFileName("dothething", ".xcresult")
	if err != nil {
		return err
	}

	xce := xcode.ParseXCodeBuildError(a.runXCodebuildTest(ctx, path))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce

}
func (a actionRunTest) runXCodebuildTest(ctx context.Context, path string) error {
	fmt.Println(color.BlueString("Running test on %v (%v)", a.Config.Destination.Name, a.Config.Destination.ID))

	return RunCmd(a.API.Exec().CommandContext(ctx,
		xcode.Cmd,
		a.API.XCodeBuildService().GetArg(),
		a.API.XCodeBuildService().GetProjectPath(),
		xcode.ActionTest,
		xcode.FlagScheme, a.Config.Scheme,
		xcode.FlagDestination, fmt.Sprintf("id=%s", a.Config.Destination.ID),
		xcode.FlagResultBundlePath, path,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO"))
}
