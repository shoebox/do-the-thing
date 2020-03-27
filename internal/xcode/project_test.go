package xcode

import (
	"dothething/internal/utiltest"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validResponse string = `{
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

const fakePath = "/path/to/project.xcodeproj"

var execMock *utiltest.MockExec
var projectService XCodeProjectService

func TestMain(m *testing.M) {
	execMock = new(utiltest.MockExec)
	projectService = XCodeProjectService{xcodeService: NewService(execMock, fakePath)}
}

func TestCases(t *testing.T) {
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
		execMock.
			On("ContextExec", XCodeBuild, []string{flagList, flagJSON, FlagProject, fakePath}).
			Return(tc.json, tc.execErr)

		_, err := projectService.Parse()

		assert.EqualValues(t, tc.expectedError, err)
		execMock.AssertExpectations(t)
	}
}

func TestProjectResolution(t *testing.T) {

	t.Run("Should properly call xcodedbuild and parse the result", func(t *testing.T) {
		execMock.
			On("ContextExec", XCodeBuild, []string{flagList, flagJSON, FlagProject, fakePath}).
			Return(validResponse, nil)

		pj, err := projectService.Parse()
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
