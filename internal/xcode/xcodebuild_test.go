package xcode

import (
	"dothething/internal/api"
	"dothething/internal/utiltest"
	"testing"

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
		// Exec:   x.mockExecutor,
	}
	x.API.Config.Path = pathWorkspace
	x.subject = xcodeBuildService{x.API}
}

/*
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
*/
