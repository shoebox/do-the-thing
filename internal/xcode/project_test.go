package xcode

import (
	"dothething/internal/utiltest"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var projectSevice *XCodeProjectService
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

func setupServiceTest() {
	execMock = new(utiltest.MockExec)
	projectSevice = NewProjectService(execMock)
}

func TestProjectResolution(t *testing.T) {
	t.Run("Should handle possible errors while invoking xcodebuild", func(t *testing.T) {
		setupServiceTest()

		execMock.
			On("Exec", xCodeBuild, []string{flagList, flagJSON}).
			Return(nil, errors.New("Error calling xcode"))

		pj, err := projectSevice.Parse()
		assert.Nil(t, pj)
		assert.EqualError(t, err, "Failed to call xcode API (Error : Error calling xcode)")

	})

	t.Run("Should properly call xcodedbuild and parse the result", func(t *testing.T) {
		setupServiceTest()

		execMock.
			On("Exec", xCodeBuild, []string{flagList, flagJSON}).
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

	t.Run("Should handling unmarshalling errors", func(t *testing.T) {
		setupServiceTest()

		execMock.
			On("Exec", xCodeBuild, []string{flagList, flagJSON}).
			Return("invalid json", nil)

		pj, err := projectSevice.Parse()
		assert.EqualError(t, err, ErrInvalidConfig.Error())
		assert.Nil(t, pj)
	})
}
