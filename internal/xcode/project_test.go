package xcode

import (
	"dothething/internal/utiltest"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var projectSevice XCodeProjectService
var execMock *utiltest.MockExec

var validResponse string = `{
				"project" : {
					"configurations" : [
						"Production Debug",
						"Production Release",
						"Staging Debug",
						"Staging Release",
						"Test Debug",
						"Test Release"
					],
					"name" : "Sample",
					"schemes" : [
						"ActionPopoverButton",
						"Sample",
						"SampleTests Prod",
						"SampleTests Stag",
						"SampleTests Test",
						"TopShelfExtension"
					],
					"targets" : [
						"ActionPopoverButton",
						"Sample",
						"SampleTests",
						"TopShelfExtension"
					]
				}
			}`

var fakePath = "/path/to/project.xcodeproj"

func setupServiceTest() {
	execMock = new(utiltest.MockExec)
	projectSevice = XCodeProjectService{exec: execMock,
		arg:  flagProject,
		path: fakePath}
}

func TestProjectServiceCreate(t *testing.T) {
	t.Run("For workspaces", func(t *testing.T) {
		// when:
		pj := NewProjectService(execMock, "/path/test/toto.xcworkspace")

		//
		pj2, ok := pj.(XCodeProjectService)
		assert.True(t, ok)

		// then:
		assert.EqualValues(t, "-workspace", pj2.arg)
	})

	t.Run("For project", func(t *testing.T) {
		// when:
		pj := NewProjectService(execMock, "/path/test/toto.xcodeproj")

		//
		pj2, ok := pj.(XCodeProjectService)
		assert.True(t, ok)

		// then:
		assert.EqualValues(t, "-project", pj2.arg)
	})
}

func TestProjectResolution(t *testing.T) {
	params := []struct {
		execErr       error
		expectedError error
		json          string
		path          string
	}{
		{
			execErr:       nil,
			json:          "invalid json",
			expectedError: ErrInvalidConfig,
		},
		{
			execErr:       errors.New("Error calling xcode"),
			json:          "invalid json",
			expectedError: errors.New("Failed to call xcode API (Error : Error calling xcode)"),
		},
		{
			execErr:       errors.New("Error calling xcode"),
			expectedError: errors.New("Failed to call xcode API (Error : Error calling xcode)"),
			json:          "invalid json",
			path:          "/r/t",
		},
	}

	for _, tc := range params {
		setupServiceTest()

		execMock.
			On("Exec", xCodeBuild, []string{flagList, flagJSON, flagProject, fakePath}).
			Return(tc.json, tc.execErr)

		_, err := projectSevice.Parse()

		assert.EqualValues(t, tc.expectedError, err)
		execMock.AssertExpectations(t)
	}

	t.Run("Should properly call xcodedbuild and parse the result", func(t *testing.T) {
		setupServiceTest()

		execMock.
			On("Exec", xCodeBuild, []string{flagList, flagJSON, flagProject, fakePath}).
			Return(validResponse, nil)

		pj, err := projectSevice.Parse()
		assert.NoError(t, err)
		assert.NotNil(t, pj)

		execMock.AssertExpectations(t)

		assert.EqualValues(t, []string{
			"ActionPopoverButton",
			"Sample",
			"SampleTests",
			"TopShelfExtension",
		}, pj.Targets)

		assert.EqualValues(t, "Sample", pj.Name)

		assert.EqualValues(t, []string{
			"Production Debug",
			"Production Release",
			"Staging Debug",
			"Staging Release",
			"Test Debug",
			"Test Release"}, pj.Configurations)

		assert.EqualValues(t, []string{
			"ActionPopoverButton",
			"Sample",
			"SampleTests Prod",
			"SampleTests Stag",
			"SampleTests Test",
			"TopShelfExtension"}, pj.Schemes)
	})
}

func TestListDestinations(t *testing.T) {
	t.Run("Handle error", func(t *testing.T) {
		// setup:
		setupServiceTest()
		execMock.
			On("ContextExec", mock.Anything,
				xCodeBuild,
				[]string{
					flagShowDestinations,
					flagProject, fakePath,
					flagScheme, "toto",
				}). // Build the scheme specified by schemename.
			Return(nil, errors.New("test error"))

		// when:
		list, err := projectSevice.ListDestinations("toto")
		assert.Nil(t, list)
		assert.EqualError(t, err, "test error")
	})

	t.Run("Should handle timeout", func(t *testing.T) {
		// setup:
		setupServiceTest()
		execMock.
			On("ContextExec", mock.Anything,
				xCodeBuild,
				[]string{
					flagShowDestinations,
					flagProject, fakePath,
					flagScheme, "toto",
				}). // Build the scheme specified by schemename.
			WaitUntil(time.After(20*time.Second)).
			Return(nil, errors.New("test error"))

		// when:
		list, err := projectSevice.ListDestinations("toto")
		assert.Nil(t, list)
		assert.EqualError(t, err, ErrDestinationResolutionFailed.Error())
	})

	t.Run("Should parse result", func(t *testing.T) {
		// setup:
		setupServiceTest()
		execMock.
			On("ContextExec", mock.Anything,
				xCodeBuild,
				[]string{
					flagShowDestinations,
					flagProject, fakePath,
					flagScheme, "toto",
				}). // Build the scheme specified by schemename.
			Return(` 
	Available destinations for the "test" scheme:
		{ platform:iOS Simulator, id:20ADB361-71A8-4D73-8F5D-38241750CBF5, OS:13.3, name:iPad }

	Ineligible destinations for the "test" scheme:
		{ platform:iOS Simulator, id:D2D6C8CE-04B6-44E5-933F-63C29A5952C2, OS:13.3, name:iPad Air (3rd generation) }
		{ platform:iOS Simulator, id:dvtdevice-DVTiOSDeviceSimulatorPlaceholder-iphonesimulator:placeholder, name:Generic iOS Simulator Device }
			`, nil)

		// when:
		list, err := projectSevice.ListDestinations("toto")

		// then:
		assert.NotEmpty(t, list)
		assert.EqualValues(t, list[0], Destination{
			Name:     "iPad",
			OS:       "13.3",
			ID:       "20ADB361-71A8-4D73-8F5D-38241750CBF5",
			Platform: "iOS Simulator",
		})
		assert.Nil(t, err)

		// and:
		execMock.AssertExpectations(t)
	})
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
		assert.EqualValues(t, dest, Destination{ID: "fake-id", Platform: "fake-platform"})
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
