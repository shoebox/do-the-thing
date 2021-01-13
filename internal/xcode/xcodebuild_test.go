package xcode

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/utiltest"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const pathWorkspace = "/path/to/project.xcworkspace"
const pathProject = "/path/to/project.xcodeproj"

type XCodebuildSuite struct {
	suite.Suite
	*api.API
	subject      xcodeBuildService
	mockCmd      *utiltest.MockExecutorCmd
	mockExecutor *utiltest.MockExecutor
}

func TestXCodebuildSuite(t *testing.T) {
	suite.Run(t, new(XCodebuildSuite))
}

func (x *XCodebuildSuite) BeforeTest(suiteName, testName string) {
	x.mockCmd = new(utiltest.MockExecutorCmd)
	x.mockExecutor = new(utiltest.MockExecutor)

	x.API = &api.API{
		Config: new(api.Config),
		Exec:   x.mockExecutor,
	}
	x.API.Config.Path = pathWorkspace
	x.subject = xcodeBuildService{x.API}
}

func (x *XCodebuildSuite) TestXCodeBuildArgumentDependingOfProjectType() {
	cases := []struct {
		path     string
		arg      string
		expected string
	}{
		{path: "~/dev/test.xcodeproj", expected: FlagProject},
		{path: "~/dev/test.xcworkspace", expected: FlagWorkspace},
		{path: "~/dev/", expected: FlagProject},
		{path: pathProject, expected: FlagProject},
		{path: pathWorkspace, expected: FlagWorkspace},
	}

	for _, c := range cases {
		x.Run(c.path, func() {
			// setup:
			x.API.Config.Path = c.path

			// when:
			res := x.subject.GetArg()

			// then:
			x.Assert().EqualValues(c.expected, res)
		})
	}
}

func (x *XCodebuildSuite) TestList() {
	cases := []struct {
		name           string
		output         string
		err            error
		expectedResult string
		expectedError  error
		path           string
		flag           string
	}{
		{path: pathWorkspace, flag: FlagWorkspace, name: "error", output: "", err: errors.New("fake error"), expectedResult: "", expectedError: NewError(-1)},
		{path: pathProject, flag: FlagProject, name: "error", output: "", err: errors.New("fake error"), expectedResult: "", expectedError: NewError(-1)},
		{path: pathWorkspace, flag: FlagWorkspace, name: "basic", output: "hello-world", expectedResult: "hello-world", expectedError: nil},
	}

	for _, c := range cases {
		x.Run(c.name, func() {
			// setup:
			x.BeforeTest("XCodebuildSuite", c.name)
			x.API.Config.Path = c.path
			x.mockCmd.On("CombinedOutput").Return(c.output, c.err)
			x.mockExecutor.On("CommandContext",
				mock.Anything,
				Cmd,
				[]string{FlagList, FlagJSON, c.flag, c.path}).
				Return(x.mockCmd)

			// when:
			res, err := x.subject.List(context.Background())

			// then:
			x.Assert().EqualValues(c.expectedError, err)
			x.Assert().EqualValues(c.expectedResult, res)
		})
	}
}

func (x *XCodebuildSuite) TestShowDestinations() {
	txt := "Destination list text"

	// setup:
	x.mockExecutor.MockCommandContext(
		Cmd,
		[]string{
			FlagShowDestinations,
			FlagWorkspace,
			pathWorkspace,
			FlagScheme,
			"test",
		},
		txt,
		nil,
	)
	// whhen
	res, err := x.subject.ShowDestinations(context.Background(), "test")

	// then:
	x.Assert().EqualValues(res, txt)
	x.Assert().NoError(err)
}
