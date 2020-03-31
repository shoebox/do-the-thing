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
)

var projectSevice destinationService
var execMock *utiltest.MockExec
var xMock *xcode.MockXCodeBuildService
var dest Destination

func setupServiceTest() {
	dest = Destination{Id: "123-456-789"}
	xMock = new(xcode.MockXCodeBuildService)
	execMock = new(utiltest.MockExec)
	projectSevice = destinationService{xcode: xMock, exec: execMock}
}

func TestBoot(t *testing.T) {
	var flagtests = []struct {
		name         string
		resToReturn  []byte
		errToReturn  error
		waitTime     int
		stringResult string
		expectedErr  error
		timeout      int
	}{
		{name: "Should handle error",
			waitTime:    0,
			timeout:     10,
			errToReturn: errors.New("test error"),
			expectedErr: errors.New("test error"),
		},
		{name: "Should handle timeout",
			waitTime:     2,
			timeout:      1,
			stringResult: "toto",
			errToReturn:  errors.New("test error"),
			expectedErr:  context.DeadlineExceeded,
		},
		{name: "Should handle result",
			waitTime:     0,
			timeout:      1,
			stringResult: "toto",
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			// setup:
			setupServiceTest()
			execMock.
				On("ContextExec", mock.Anything, xcRun, []string{simCtl, actionBootStatus, dest.Id, flagBoot}).
				WaitUntil(time.After(time.Duration(tt.waitTime)*time.Second)).
				Return(tt.stringResult, tt.errToReturn)

			// when:
			ctx, cancel := context.WithTimeout(context.Background(),
				time.Duration(tt.timeout)*time.Second)
			defer cancel() // The cancel should be deferred so resources are cleaned up

			err := projectSevice.Boot(ctx, dest)
			assert.EqualValues(t, tt.expectedErr, err)
		})
	}
}

func TestShutdown(t *testing.T) {
	// setup:
	ctx := context.Background()
	execMock.
		On("ContextExec", ctx, xcRun, []string{simCtl, actionShutdown, dest.Id}).
		Return(nil, errors.New("test error"))

	// when:
	err := projectSevice.ShutDown(ctx, dest)

	// then:
	assert.EqualValues(t, errors.New("test error"), err)
}

func TestXcRun(t *testing.T) {
	t.Run("Should handle error", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		// setup:
		setupServiceTest()
		execMock.
			On("ContextExec", ctx, xcRun, []string{simCtl, actionBootStatus, dest.Id}).
			Return(nil, errors.New("test error"))

		// when:
		errc := make(chan error, 1)
		resc := make(chan string, 1)
		projectSevice.xcRun(ctx, resc, errc, actionBootStatus, dest.Id)

		// then:
		err := <-errc
		res := <-resc
		assert.Empty(t, res, "No result should be returned")
		assert.EqualError(t, err, "test error", "The mocked error should be returned")
	})

	t.Run("Should handle success", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		// setup:
		setupServiceTest()
		execMock.
			On("ContextExec", ctx, xcRun, []string{simCtl, actionBootStatus, dest.Id}).
			Return("Hello world", nil)

		// when:
		errc := make(chan error, 1)
		resc := make(chan string, 1)
		projectSevice.xcRun(ctx, resc, errc, actionBootStatus, dest.Id)

		// then:
		err := <-errc
		res := <-resc
		assert.EqualValues(t, res, "Hello world", "Valid result should be returned")
		assert.NoError(t, err, "No error should be returned")
	})
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
