package xcode

import (
	"context"
	"dothething/internal/utiltest"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type XCodeBuildServiceTestSuite struct {
	suite.Suite
	exec    *utiltest.MockExec
	subject XCodeBuildService
}

func (s *XCodeBuildServiceTestSuite) BeforeTest(suiteName string, testName string) {
	utiltest.SetupMockExec()
	s.exec = utiltest.Exec
	s.subject = NewService(s.exec, "/project/path/to/project.xcodeproj")
}

func (s *XCodeBuildServiceTestSuite) AfterTest(suiteName string, testName string) {
	utiltest.TearDownMockExec()
	s.exec = nil
}

func (s *XCodeBuildServiceTestSuite) TestNewService() {
	cases := []struct {
		path string
		flag string
	}{
		{path: "/path/to/project.xcodeproj", flag: FlagProject},
		{path: "/path/to/project.xcworkspace", flag: FlagWorkspace},
	}

	for _, tc := range cases {
		// when:
		res := NewService(s.exec, tc.path)

		// then:
		b, ok := res.(xcodeBuildService)
		s.Assert().True(ok, "Should be true")

		// and:
		s.Assert().EqualValues(b.arg, tc.flag, "Xcodebuild should be invoked with the right flag")
	}
}

func (s *XCodeBuildServiceTestSuite) TestRunShouldHandleErrors() {
	// setup:
	s.exec.On("ContextExec", mock.Anything,
		XCodeBuild,
		[]string{FlagShowDestinations, flagList}).
		Return(nil, errors.New("Fake error"))

	// when:
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	res, err := s.subject.Run(ctx, FlagShowDestinations, flagList)

	// then:
	s.Assert().Empty(res, "No result should be returned")
	s.Assert().EqualError(err, "Fake error", "The same error should be return")
}

func (s *XCodeBuildServiceTestSuite) TestRunShouldHandleTimeout() {
	// setup:
	s.exec.On("ContextExec", mock.Anything,
		XCodeBuild,
		[]string{FlagShowDestinations, flagList}).
		WaitUntil(time.After(time.Second*11)).
		Return(nil, errors.New("Fake error"))

	// when:
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	res, err := s.subject.Run(ctx, FlagShowDestinations, flagList)

	// then:
	s.Assert().Empty(res, "No result should be returned")
	s.Assert().EqualError(err, "context deadline exceeded", "The same error should be return")
}

func TestXCodeBuildServiceTestSuite(t *testing.T) {
	suite.Run(t, new(XCodeBuildServiceTestSuite))
}
