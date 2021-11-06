package destination

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/utiltest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var dest = api.Destination{ID: "123-456-789"}

type DestinationTestSuite struct {
	API          *api.API
	cancel       func()
	ctx          context.Context
	mockCmd      *utiltest.MockExecutorCmd
	mockExecutor *utiltest.MockExecutor
	// mockXcodeBuild *xcode.XCodeBuildMock
	subject destinationService
	suite.Suite
}

// suite testing entrypoint
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DestinationTestSuite))
}

func (s *DestinationTestSuite) BeforeTest(suiteName, testName string) {
	s.mockCmd = new(utiltest.MockExecutorCmd)
	s.mockExecutor = new(utiltest.MockExecutor)
	// s.mockXcodeBuild = new(xcode.XCodeBuildMock)

	s.API = new(api.API)
	// s.API.Exec = s.mockExecutor
	//s.API.BuildService = s.mockXcodeBuild

	s.ctx, s.cancel = context.WithTimeout(context.Background(), 1*time.Second)
	s.subject = destinationService{s.API}
}

func (s *DestinationTestSuite) AfterTest(suiteName, testName string) {
	s.cancel()
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
	s.Assert().EqualValues(api.Destination{
		Name:     "iPad",
		OS:       "13.3",
		ID:       "20ADB361-71A8-4D73-8F5D-38241750CBF5",
		Platform: "iOS Simulator",
	}, res[0])
}

func (s *DestinationTestSuite) TestDestinationSectionDetection() {
	cases := []struct {
		line     string
		expected bool
		start    bool
		res      []api.Destination
	}{
		{line: `Available destinations for the "test" scheme:`, expected: true},
		{line: `Toto`, expected: false},
		{line: `Available destinations for the "hello world" scheme:`, expected: true},
	}

	for _, c := range cases {
		// when:
		s.subject.parseLine(c.line, &c.start, &c.res)

		// then:
		s.Assert().Equal(c.expected, c.start)
	}
}

func TestFillStruct(t *testing.T) {
	var fakeID string = "fake-id"
	var fakePlatform string = "fake-platform"
	var fakeOS string = "fake-os"

	cases := []struct {
		name string
		d    api.Destination
		m    map[string]string
	}{
		{
			name: "Should populate multiple values",
			m:    map[string]string{"Platform": fakePlatform, "ID": fakeID},
			d:    api.Destination{ID: fakeID, Platform: fakePlatform},
		},
		{
			name: "Should handle invalid fields",
			m:    map[string]string{"toto": "osef", "OS": fakeOS},
			d:    api.Destination{OS: fakeOS},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// setup:
			dest := api.Destination{}

			// when:
			fillStruct(c.m, &dest)

			// then:
			assert.EqualValues(t, c.d, dest)
		})
	}
}
