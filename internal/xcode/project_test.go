package xcode

import (
	"dothething/internal/utiltest"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
