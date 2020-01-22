package xcode

import (
	"fmt"
	"testing"

	"dothething/internal/utiltest"

	"github.com/stretchr/testify/assert"
)

var mockExec *utiltest.MockExec
var mockFileService *utiltest.MockFileService
var service XCodeListService

const XCODES = `/Applications/Xcode.app
/Applications/Xcode 10.3.1.app
/Invalid/path`

func setup() {
	mockExec = new(utiltest.MockExec)
	mockFileService = new(utiltest.MockFileService)
	service = XCodeListService{exec: mockExec, file: mockFileService}
}

func xcodePlist(version string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>CFBundleVersion</key>
			<string>%v</string>
			<key>CFBundleShortVersionString</key>
			<string>%v-snapshot</string>
		</dict>
	</plist>`, version, version)
}

func TestOpenFileContent(t *testing.T) {
	setup()

	// Expectation
	mockExec.
		On("Exec", MDFIND, []string{XCODE_BUNDLE_IDENTIFIER}).
		Return("Hello world", nil)

	wb, _ := service.spotlightSearch()
	assert.NotNil(t, wb)

	mockExec.AssertExpectations(t)
}

func TestResolveXcode(t *testing.T) {
	install := Install{
		Path:          "/Applications/Xcode.app",
		BundleVersion: "1.2.3",
		Version:       "1.2.3-snapshot",
	}

	t.Run("Should be able to resolve the install", func(t *testing.T) {
		setup()

		mockFileService.
			On("OpenAndReadFileContent", fmt.Sprintf("%v%v", install.Path, PLIST)).
			Return(xcodePlist("1.2.3"), nil)

		xc, err := service.resolveXcode(fmt.Sprintf(install.Path))

		assert.Nil(t, err)
		assert.EqualValues(t, xc, &install)
		mockFileService.AssertExpectations(t)
	})

	t.Run("Should fail to resolve invalid path", func(t *testing.T) {
		setup()

		mockFileService.
			On("OpenAndReadFileContent", fmt.Sprintf("%v%v", install.Path, PLIST)).
			Return(nil, fmt.Errorf("Error sample"))

		xc, err := service.resolveXcode(fmt.Sprintf(install.Path))

		assert.NotNil(t, err)
		assert.Nil(t, xc)
	})

	t.Run("Decoded error should be raise", func(t *testing.T) {
		setup()

		mockFileService.
			On("OpenAndReadFileContent", fmt.Sprintf("%v%v", install.Path, PLIST)).
			Return("invalid", nil)

		xc, err := service.resolveXcode(fmt.Sprintf(install.Path))
		assert.Nil(t, xc)
		assert.EqualError(t, err, "plist: type mismatch: tried to decode plist type `string' into value of type `xcode.infoPlist'")
	})
}

func TestSpotLightFailure(t *testing.T) {
	setup()

	mockExec.
		On("Exec", MDFIND, []string{XCODE_BUNDLE_IDENTIFIER}).
		Return(nil, fmt.Errorf("Error"))

	res, err := service.List()

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestList(t *testing.T) {
	setup()

	mockExec.
		On("Exec", MDFIND, []string{XCODE_BUNDLE_IDENTIFIER}).
		Return(XCODES, nil)

	mockFileService.On("IsDir", "/Applications/Xcode.app").Return(true, nil)
	mockFileService.On("IsDir", "/Applications/Xcode 10.3.1.app").Return(true, nil)
	mockFileService.On("IsDir", "/Invalid/path").Return(false, nil)

	mockFileService.
		On("OpenAndReadFileContent", "/Applications/Xcode.app"+PLIST).
		Return(xcodePlist("1.2.3"), nil)

	mockFileService.
		On("OpenAndReadFileContent", "/invalid/path"+PLIST).
		Return(xcodePlist("1.2.3"), nil)

	mockFileService.
		On("OpenAndReadFileContent", "/Applications/Xcode 10.3.1.app"+PLIST).
		Return(xcodePlist("10.3.1"), nil)

	res, err := service.List()
	assert.NoError(t, err)

	assert.EqualValues(t, res, []*Install{
		&Install{"/Applications/Xcode.app", "1.2.3", "1.2.3-snapshot"},
		&Install{"/Applications/Xcode 10.3.1.app", "10.3.1", "10.3.1-snapshot"},
	})
}
