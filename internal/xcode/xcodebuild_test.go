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

type XCodebuildSuite struct {
	suite.Suite
	cmd      *utiltest.MockExecutorCmd
	ctx      context.Context
	cancel   context.CancelFunc
	subject  BuildService
	executor *utiltest.MockExecutor
	path     string
}

func TestXCodebuildSuite(t *testing.T) {
	suite.Run(t, new(XCodebuildSuite))
}

func (s *XCodebuildSuite) BeforeTest(suiteName, testName string) {
	s.executor = new(utiltest.MockExecutor)
	s.path = "/path/to/project.xcworkspace"
	s.subject = NewService(s.executor, s.path)
	s.cmd = new(utiltest.MockExecutorCmd)
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 60*time.Second)
}

func (s *XCodebuildSuite) AfterTest(suiteName, testName string) {
	s.cancel()
}

func (s *XCodebuildSuite) TestNewServiceFlag() {
	cases := []struct {
		path string
		flag string
	}{
		{path: "/path/to/project.xcodeproj", flag: FlagProject},
		{path: "/path/to/project.xcworkspace", flag: FlagWorkspace},
	}

	for _, tc := range cases {
		// when:
		res := NewService(s.executor, tc.path)

		// then:
		b, ok := res.(xcodeBuildService)
		s.Assert().True(ok, "Should be true")

		// and:
		s.Assert().EqualValues(b.arg, tc.flag, "Xcodebuild should be invoked with the right flag")
	}
}

func (s *XCodebuildSuite) TestListShouldHandleError() {
	// setup:
	s.cmd.On("Output").Return("", errors.New("Fake error"))

	s.executor.On("CommandContext",
		mock.Anything,
		Build,
		[]string{FlagList, FlagJSON, "-workspace", s.path}).
		Return(s.cmd)

	// whhen
	res, err := s.subject.List(s.ctx)

	// then:
	s.Assert().Empty(res, "No result should be returned")
	s.Assert().EqualError(err, "Error -1 - Unknown error", "The same error should be return")
}

func (s *XCodebuildSuite) TestListShouldReturnResult() {
	// setup:
	txt := "Destination list text"
	s.cmd.On("Output").Return(txt, nil)

	s.executor.On("CommandContext",
		mock.Anything,
		Build,
		[]string{FlagList, FlagJSON, "-workspace", s.path}).
		Return(s.cmd)

	// whhen
	res, err := s.subject.List(s.ctx)

	// then:
	s.Assert().EqualValues(res, txt)
	s.Assert().NoError(err)
}

func (s *XCodebuildSuite) TestShowDestinations() {
	// setup:
	txt := "Destination list text"
	s.cmd.On("Output").Return(txt, nil)

	s.executor.On("CommandContext",
		s.ctx,
		Build,
		[]string{
			FlagShowDestinations,
			FlagWorkspace,
			s.path,
			FlagScheme,
			"test",
		}).
		Return(s.cmd)

	// whhen
	res, err := s.subject.ShowDestinations(s.ctx, "test")

	// then:
	s.Assert().EqualValues(res, txt)
	s.Assert().NoError(err)
}
