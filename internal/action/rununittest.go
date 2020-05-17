package action

import (
	"context"
	"crypto/rand"
	"dothething/internal/destination"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/output"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

const (
	xCRun      string = "xcrun"
	resultTool string = "xcresulttool"
	actionGet  string = "get"
	flagFormat string = "--format"
	formatJson string = "json"
	flagPath   string = "--path"
)

type ActionRunTest interface {
	Run(ctx context.Context, dest destination.Destination, config xcode.Config) error
}

type actionRunTest struct {
	exec  util.Executor
	xcode xcode.XCodeBuildService
}

func NewActionRun(service xcode.XCodeBuildService, exec util.Executor) ActionRunTest {
	return actionRunTest{xcode: service, exec: exec}
}

func (a actionRunTest) Run(ctx context.Context, d destination.Destination, config xcode.Config) error {
	// Creating a temp folder to contains the test results
	path, err := tempFileName("dothething", ".xcresult")
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
	config xcode.Config,
	dest destination.Destination) error {
	fmt.Println(color.BlueString("Running test on %v (%v)", dest.Name, dest.Id))

	cmd := a.exec.CommandContext(ctx,
		xcode.XCodeBuild,
		a.xcode.GetArg(),
		a.xcode.GetProjectPath(),
		xcode.ActionClean,
		xcode.ActionTest,
		xcode.FlagScheme, config.Scheme,
		xcode.FlagDestination, fmt.Sprintf("id=%s", dest.Id),
		xcode.FlagResultBundlePath, path,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO")

	pout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	perr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		f := output.NewFormatter(output.SimpleReporter{})
		f.Parse(pout)
		f.Parse(perr)
	}()

	//
	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (a actionRunTest) decodeXCResultFile(ctx context.Context, path string) ([]byte, error) {
	return a.exec.CommandContext(ctx,
		xCRun, resultTool,
		actionGet,
		flagFormat, formatJson,
		flagPath, path,
	).Output()
}

// TempFileName generates a temporary filename for use in testing or whatever
func tempFileName(prefix, suffix string) (string, error) {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix), nil
}
