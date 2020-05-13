package unittest

import (
	"context"
	"crypto/rand"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/output"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
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
	Run(ctx context.Context, dest string) error
}

type actionRunTest struct {
	exec  util.Executor
	xcode xcode.XCodeBuildService
}

type failedTest struct {
	name   string
	msg    string
	target string
}

func NewActionRun(service xcode.XCodeBuildService, exec util.Executor) ActionRunTest {
	return actionRunTest{xcode: service, exec: exec}
}

func (a actionRunTest) Run(ctx context.Context, dest string) error {
	// Creating a temp folder to contains the test results
	path := tempFileName("dothething", ".xcresult")

	// Run test via xcodebuild
	_, err := a.runXCodebuildTest(ctx, path, dest)
	if err != nil {
		return xcode.ParseXCodeBuildError(err)
	}

	// output.Parse(strings.NewReader(res))
	/*
		// TODO: Handle xcodebuild result

		// Decode the content of the xcresult file
		b, err := a.decodeXCResultFile(ctx, path)
		if err != nil {
			return err
		}

		// Retrieving the test failure summaries
		issues, err := xcresult.ParseIssues(b)
		if err != nil {
			return err
		}
		for _, v := range issues {
			log.Error().
				Str("Case name", v.TestCaseName.Value).
				Str("Message", v.Message.Value).
				Str("File", v.Doc.URL.Value).
				Msg("Test failed")
		}
	*/
	return nil

}
func (a actionRunTest) runXCodebuildTest(ctx context.Context, path string, dest string) (string, error) {
	log.Info().
		Str("Destination", dest).
		Str("Output file", path).
		Msg("Running tests on destination")

	cmd := a.exec.CommandContext(ctx,
		xcode.XCodeBuild,
		a.xcode.GetArg(),
		a.xcode.GetProjectPath(),
		xcode.ActionClean,
		xcode.ActionTest,
		xcode.FlagScheme, "Swiftstraints iOS",
		xcode.FlagDestination, fmt.Sprintf("id=%s", dest),
		xcode.FlagResultBundlePath, path,
		"-showBuildTimingSummary",
		"CODE_SIGNING_ALLOWED=NO")

	pout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	perr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	go func() {
		f := output.NewFormatter(output.SimpleReporter{})
		f.Parse(pout)
		f.Parse(perr)
	}()

	//
	if err = cmd.Start(); err != nil {
		return "", err
	}

	if err = cmd.Wait(); err != nil {
		return "", err
	}

	return "", nil
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
func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
