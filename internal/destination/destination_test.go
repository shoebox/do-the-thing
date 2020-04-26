package destination

import (
	"context"
	"dothething/internal/utiltest"
	"dothething/internal/xcode"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DestinationTestSuite struct {
	suite.Suite
	cancel   context.CancelFunc
	ctx      context.Context
	cmd      *utiltest.MockExecutorCmd
	dest     Destination
	subject  destinationService
	executor *utiltest.MockExecutor2
	xcode    xcode.XCodeBuildService
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DestinationTestSuite))
}

func (s *DestinationTestSuite) BeforeTest(suiteName, testName string) {
	s.cmd = new(utiltest.MockExecutorCmd)
	s.dest = Destination{Id: "123-456-789"}
	s.executor = new(utiltest.MockExecutor2)
	s.xcode = xcode.NewService(s.executor, "/path/to/project.pbxproj")
	s.subject = destinationService{s.xcode, s.executor}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	s.ctx = ctx
	s.cancel = cancel
}

func (s *DestinationTestSuite) AfterTest(suiteName, testName string) {
	s.cancel()
}

func (s *DestinationTestSuite) TestBootShouldHandleResults() {
	// setup:
	s.cmd.
		On("Output").
		Return("result", nil)

	s.executor.
		On("CommandContext",
			mock.Anything,
			xcRun, []string{simCtl, actionBootStatus, s.dest.Id, flagBoot}).
		Return(s.cmd)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// when:
	err := s.subject.Boot(ctx, s.dest)

	// then:
	s.Assert().NoError(err)
}

func (s *DestinationTestSuite) TestBootShouldHandleError() {
	// setup:
	s.cmd.
		On("Output").
		Return("result", errors.New("Error text"))

	s.executor.
		On("CommandContext",
			s.ctx,
			xcRun, []string{simCtl, actionBootStatus, s.dest.Id, flagBoot}).
		Return(s.cmd)

	// when:
	err := s.subject.Boot(s.ctx, s.dest)

	// then:
	s.Assert().EqualError(err, "Error text")
}

func (s *DestinationTestSuite) TestShutdown() {
	// setup:
	s.cmd.
		On("Output").
		Return("result", errors.New("Error text"))

	s.executor.
		On("CommandContext",
			s.ctx,
			xcRun, []string{simCtl, actionShutdown, s.dest.Id}).
		Return(s.cmd)

	// when:
	err := s.subject.ShutDown(s.ctx, s.dest)

	// then:
	s.Assert().EqualValues(errors.New("Error text"), err)
}

func (s *DestinationTestSuite) TestList() {
	// setup:
	s.cmd.
		On("Output").
		Return("result", errors.New("toto"))

	s.executor.
		On("CommandContext",
			s.ctx,
			xcode.XCodeBuild,
			[]string{
				xcode.FlagShowDestinations,
				xcode.FlagProject,
				"/path/to/project.pbxproj",
				xcode.FlagScheme,
				"test",
			}).
		Return(s.cmd)

	// when:
	dest, err := s.subject.List(s.ctx, "test")

	// then:
	s.Assert().EqualError(err, "Error -1 - Unknown error")
	s.Assert().Empty(dest)
}

func (s *DestinationTestSuite) TestDestinationsParsing() {
	data := `
			Available destinations for the "test" scheme:
				{ platform:iOS Simulator, id:20ADB361-71A8-4D73-8F5D-38241750CBF5, OS:13.3, name:iPad }

			Ineligible destinations for the "test" scheme:
				{ platform:iOS Simulator, id:D2D6C8CE-04B6-44E5-933F-63C29A5952C2, OS:13.3, name:iPad Air (3rd generation) }
				{ platform:iOS Simulator, id:dvtdevice-DVTiOSDeviceSimulatorPlaceholder-iphonesimulator:placeholder, name:Generic iOS Simulator Device }`

	// when:
	res := s.subject.parseDestinations(data)

	// then:
	s.Assert().EqualValues(Destination{
		Name:     "iPad",
		OS:       "13.3",
		Id:       "20ADB361-71A8-4D73-8F5D-38241750CBF5",
		Platform: "iOS Simulator",
	}, res[0])
}

func TestFillStruct(t *testing.T) {
	t.Run("Should populate multiple values", func(t *testing.T) {
		// setup:
		dest := Destination{}

		m := map[string]string{}
		m["platform"] = "fake-platform"
		m["id"] = "fake-id"

		// when:
		fillStruct(m, &dest)

		// then:
		assert.EqualValues(t, dest, Destination{Id: "fake-id", Platform: "fake-platform"})
	})

	t.Run("Should handle invalid fields", func(t *testing.T) {
		// setup:
		dest := Destination{}
		m := map[string]string{}
		m["toto"] = "tutu"
		m["OS"] = "fake-os"

		// when:
		fillStruct(m, &dest)

		// then:
		assert.EqualValues(t, dest, Destination{OS: "fake-os"})
	})
}
