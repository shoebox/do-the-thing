package destination

import (
	"context"
	"dothething/internal/utiltest"
	"dothething/internal/xcode"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DestinationTestSuite struct {
	suite.Suite
	dest     Destination
	execMock *utiltest.MockExec
	subject  destinationService
	xMock    *xcode.MockXCodeBuildService
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DestinationTestSuite))
}

func (s *DestinationTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Println("SetupTest")
	s.dest = Destination{Id: "123-456-789"}
	s.xMock = new(xcode.MockXCodeBuildService)
	s.execMock = new(utiltest.MockExec)
	s.subject = destinationService{xcode: s.xMock, exec: s.execMock}
}

func (s *DestinationTestSuite) TestBootShouldHandleResults() {
	// setup:
	s.execMock.
		On("ContextExec",
			mock.Anything,
			xcRun, []string{simCtl, actionBootStatus, s.dest.Id, flagBoot}).
		Return("toto", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// when:
	err := s.subject.Boot(ctx, s.dest)

	// then:
	s.Assert().NoError(err)
}

func (s *DestinationTestSuite) TestBootShouldHandleError() {
	// setup:
	s.execMock.
		On("ContextExec",
			mock.Anything,
			xcRun, []string{simCtl, actionBootStatus, s.dest.Id, flagBoot}).
		Return("", errors.New("test error"))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// when:
	err := s.subject.Boot(ctx, s.dest)

	// then:
	s.Assert().EqualValues(errors.New("test error"), err)
}

func (s *DestinationTestSuite) TestBootShouldHandleTimeout() {
	// setup:
	s.execMock.
		On("ContextExec",
			mock.Anything,
			xcRun, []string{simCtl, actionBootStatus, s.dest.Id, flagBoot}).
		WaitUntil(time.After(time.Duration(10)*time.Second)).
		Return("toto", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// when:
	err := s.subject.Boot(ctx, s.dest)

	// then:
	s.Assert().EqualValues(context.DeadlineExceeded, err)
}

func (s *DestinationTestSuite) TestShutdown() {
	// setup:
	ctx := context.Background()
	s.execMock.
		On("ContextExec", ctx, xcRun, []string{simCtl, actionShutdown, s.dest.Id}).
		Return(nil, errors.New("test error"))

	// when:
	err := s.subject.ShutDown(ctx, s.dest)

	// then:
	s.Assert().EqualValues(errors.New("test error"), err)
}

func (s *DestinationTestSuite) TestXcRunErrorHandling() {

	ctx := context.Background()
	defer ctx.Done()

	// setup:
	s.execMock.
		On("ContextExec", ctx, xcRun, []string{simCtl, actionBootStatus, s.dest.Id}).
		Return(nil, errors.New("test error"))

	// when:
	errc := make(chan error, 1)
	resc := make(chan string, 1)
	s.subject.xcRun(ctx, resc, errc, actionBootStatus, s.dest.Id)

	// then:
	err := <-errc
	res := <-resc
	s.Assert().Empty(res, "No result should be returned")
	s.Assert().EqualError(err, "test error", "The mocked error should be returned")
}

func (s *DestinationTestSuite) TestXcRun() {
	ctx := context.Background()
	defer ctx.Done()

	// setup:
	s.execMock.
		On("ContextExec", ctx, xcRun, []string{simCtl, actionBootStatus, s.dest.Id}).
		Return("Hello world", nil)

	// when:
	errc := make(chan error, 1)
	resc := make(chan string, 1)
	s.subject.xcRun(ctx, resc, errc, actionBootStatus, s.dest.Id)

	// then:
	err := <-errc
	res := <-resc
	s.Assert().EqualValues("Hello world", res, "Valid result should be returned")
	s.Assert().NoError(err, "No error should be returned")
}

func (s *DestinationTestSuite) TestShowDestinationsError() {
	ctx := context.Background()

	s.xMock.On("ShowDestinations", "toto").
		Return("", errors.New("error test"))

	// when:
	list, err := s.subject.List(ctx, "toto")

	// then:
	s.Assert().Empty(list)
	s.Assert().EqualError(err, ErrDestinationResolutionFailed.Error())

}

func (s *DestinationTestSuite) TestShowDestinationsResult() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	s.xMock.On("ShowDestinations", "toto").
		Return(`
			Available destinations for the "test" scheme:
				{ platform:iOS Simulator, id:20ADB361-71A8-4D73-8F5D-38241750CBF5, OS:13.3, name:iPad }

			Ineligible destinations for the "test" scheme:
				{ platform:iOS Simulator, id:D2D6C8CE-04B6-44E5-933F-63C29A5952C2, OS:13.3, name:iPad Air (3rd generation) }
				{ platform:iOS Simulator, id:dvtdevice-DVTiOSDeviceSimulatorPlaceholder-iphonesimulator:placeholder, name:Generic iOS Simulator Device }
					`, nil)

	// when:
	list, err := s.subject.List(ctx, "toto")

	// then:
	s.Assert().NotEmpty(list)
	s.Assert().EqualValues(list[0], Destination{
		Name:     "iPad",
		OS:       "13.3",
		Id:       "20ADB361-71A8-4D73-8F5D-38241750CBF5",
		Platform: "iOS Simulator",
	})
	s.Assert().Nil(err)
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
