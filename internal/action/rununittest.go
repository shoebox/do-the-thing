package action

import (
	"context"
	"dothething/internal/config"
	"dothething/internal/destination"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"fmt"

	"github.com/fatih/color"
)

type ActionRunTest interface {
	Run(ctx context.Context, dest destination.Destination, config config.Config) error
}

type actionRunTest struct {
	exec  util.Executor
	xcode xcode.BuildService
}

func NewActionRun(service xcode.BuildService, exec util.Executor) ActionRunTest {
	return actionRunTest{xcode: service, exec: exec}
}

func (a actionRunTest) Run(ctx context.Context, d destination.Destination, config config.Config) error {
	// Creating a temp folder to contains the test results
	path, err := util.TempFileName("dothething", ".xcresult")
	if err != nil {
		return err
	}

	xce := xcode.ParseXCodeBuildError(a.runXCodebuildTest(ctx, path, config, d))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce

}
func (a actionRunTest) runXCodebuildTest(ctx context.Context,
	path string,
	config config.Config,
	dest destination.Destination) error {
	fmt.Println(color.BlueString("Running test on %v (%v)", dest.Name, dest.Id))

	return RunCmd(a.exec.CommandContext(ctx,
		xcode.Cmd,
		a.xcode.GetArg(),
		a.xcode.GetProjectPath(),
		xcode.ActionTest,
		xcode.FlagScheme, config.Scheme,
		xcode.FlagDestination, fmt.Sprintf("id=%s", dest.Id),
		xcode.FlagResultBundlePath, path,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO"))
}
