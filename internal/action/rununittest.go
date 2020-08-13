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
	// exec  util.Executor
	// xcode xcode.BuildService
}

func NewActionRun(api api.API) api.Action {
	return actionRunTest{api}
}

func (a actionRunTest) Run(ctx context.Context, config api.Config) error {
	// Creating a temp folder to contains the test results
	path, err := util.TempFileName("dothething", ".xcresult")
	if err != nil {
		return err
	}

	xce := xcode.ParseXCodeBuildError(a.runXCodebuildTest(ctx, path, config))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce

}
func (a actionRunTest) runXCodebuildTest(ctx context.Context,
	path string,
	config api.Config) error {
	fmt.Println(color.BlueString("Running test on %v (%v)", config.Destination.Name, config.Destination.Id))

	return RunCmd(a.API.Exec().CommandContext(ctx,
		xcode.Cmd,
		a.API.XCodeBuildService().GetArg(),
		a.API.XCodeBuildService().GetProjectPath(),
		xcode.ActionTest,
		xcode.FlagScheme, config.Scheme,
		xcode.FlagDestination, fmt.Sprintf("id=%s", config.Destination.Id),
		xcode.FlagResultBundlePath, path,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO"))
}
