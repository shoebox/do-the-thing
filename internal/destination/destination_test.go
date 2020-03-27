package destination

import (
	"dothething/internal/utiltest"
	"dothething/internal/xcode"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var projectSevice XCodeDestinationService
var execMock *utiltest.MockExec

var fakePath = "/path/to/project.xcodeproj"
var xMock *xcode.MockXCodeBuildService

func setupServiceTest() {
	xMock = new(xcode.MockXCodeBuildService)
	projectSevice = XCodeDestinationService{xcode: xMock}
}

func TestListDestinations(t *testing.T) {

	t.Run("Handle error", func(t *testing.T) {
		// setup:
		setupServiceTest()
		xMock.On("ShowDestinations", "toto").
			Return("", errors.New("test error"))

		// when:
		list, err := projectSevice.List("toto")
		assert.Nil(t, list)
		assert.EqualError(t, err, "Command execution failed")

		xMock.AssertExpectations(t)
	})

	t.Run("Should parse result", func(t *testing.T) {
		// setup:
		setupServiceTest()
		xMock.On("ShowDestinations", "toto").
			Return(`
		Available destinations for the "test" scheme:
			{ platform:iOS Simulator, id:20ADB361-71A8-4D73-8F5D-38241750CBF5, OS:13.3, name:iPad }

		Ineligible destinations for the "test" scheme:
			{ platform:iOS Simulator, id:D2D6C8CE-04B6-44E5-933F-63C29A5952C2, OS:13.3, name:iPad Air (3rd generation) }
			{ platform:iOS Simulator, id:dvtdevice-DVTiOSDeviceSimulatorPlaceholder-iphonesimulator:placeholder, name:Generic iOS Simulator Device }
				`, nil)

		// when:
		list, err := projectSevice.List("toto")

		// then:
		assert.NotEmpty(t, list)
		assert.EqualValues(t, list[0], Destination{
			Name:     "iPad",
			OS:       "13.3",
			Id:       "20ADB361-71A8-4D73-8F5D-38241750CBF5",
			Platform: "iOS Simulator",
		})
		assert.Nil(t, err)

		// and:
		xMock.AssertExpectations(t)
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
