package unittest

import (
	"context"
	"crypto/rand"
	"dothething/internal/util"
	"dothething/internal/xcode"
	"dothething/internal/xcode/output"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

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
	exec  util.Exec
	xcode xcode.XCodeBuildService
}

type failedTest struct {
	name   string
	msg    string
	target string
}

func NewActionRun(service xcode.XCodeBuildService, exec util.Exec) ActionRunTest {
	return actionRunTest{xcode: service, exec: exec}
}

func (a actionRunTest) Run(ctx context.Context, dest string) error {
	// Creating a temp folder to contains the test results
	path := tempFileName("dothething", ".xcresult")

	// Run test via xcodebuild
	res, err := a.runXCodebuildTest(ctx, path, dest)
	if err != nil {
		return err
	}
	// fmt.Println(res)

	output.Parse(strings.NewReader(res))
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

		/*
			a.xcode.RunAction(ctx,
				"clean",
				xcode.ActionTest,
				xcode.FlagScheme, "test",
				xcode.FlagDestination, fmt.Sprintf("id=%s", dest),
				xcode.FlagResultBundlePath, path,
				"-parallel-testing-enabled",
				fmt.Sprintf("-maximum-parallel-testing-worker=%v", runtime.NumCPU()),
				"-showBuildTimingSummary")
		*/

	return "", nil
}

func (a actionRunTest) decodeXCResultFile(ctx context.Context, path string) ([]byte, error) {
	return a.exec.ContextExec(ctx,
		xCRun, resultTool,
		actionGet,
		flagFormat, formatJson,
		flagPath, path,
	)
}

// TempFileName generates a temporary filename for use in testing or whatever
func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
