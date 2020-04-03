package xcode

import (
	"fmt"
	"testing"

	"dothething/internal/utiltest"

	"github.com/stretchr/testify/suite"
)

const XCODES = `/Applications/Xcode.app
/Applications/Xcode 10.3.1.app
/Invalid/path`

var install = Install{
	Path:          "/Applications/Xcode.app",
	BundleVersion: "1.2.3",
	Version:       "1.2.3-snapshot",
}

type XCodeListTestSuite struct {
	suite.Suite
	mockExec        *utiltest.MockExec
	mockFileService *utiltest.MockFileService
	service         XCodeListService
}

func (s *XCodeListTestSuite) SetupTest() {
	utiltest.SetupMockExec()
	s.mockExec = utiltest.Exec
	s.mockFileService = new(utiltest.MockFileService)
	s.service = XCodeListService{exec: s.mockExec, file: s.mockFileService}
}

func xcodeContentPListFile(version string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD ContentPListFile 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>CFBundleVersion</key>
			<string>%v</string>
			<key>CFBundleShortVersionString</key>
			<string>%v-snapshot</string>
		</dict>
	</plist>`, version, version)
}

func TestXCodeListTestSuite(t *testing.T) {
	suite.Run(t, new(XCodeListTestSuite))
}

func (s *XCodeListTestSuite) TestOpenFileContent() {
	// setup:
	s.mockExec.
		On("Exec", MdFind, []string{XCodeBundleIdentifier}).
		Return("Hello world", nil)

	wb, _ := s.service.spotlightSearch()
	s.Assert().NotNil(wb)

	s.mockExec.AssertExpectations(s.T())
}

func (s *XCodeListTestSuite) TestShouldBeAbleToResolveTheInstall() {
	s.mockFileService.
		On("OpenAndReadFileContent", fmt.Sprintf("%v%v", install.Path, ContentPListFile)).
		Return(xcodeContentPListFile("1.2.3"), nil)

	xc, err := s.service.resolveXcode(fmt.Sprintf(install.Path))

	s.Assert().Nil(err)
	s.Assert().EqualValues(xc, &install)

	s.mockExec.AssertExpectations(s.T())
}

func (s *XCodeListTestSuite) TesTshouldFailToResolveInvalidPath() {
	s.mockFileService.
		On("OpenAndReadFileContent", fmt.Sprintf("%v%v", install.Path, ContentPListFile)).
		Return(nil, fmt.Errorf("Error sample"))

	xc, err := s.service.resolveXcode(fmt.Sprintf(install.Path))

	s.Assert().NotNil(err)
	s.Assert().Nil(xc)
}

func (s *XCodeListTestSuite) TestResolveXcodeDecodingErrorsHandling() {
	s.mockFileService.
		On("OpenAndReadFileContent", fmt.Sprintf("%v%v", install.Path, ContentPListFile)).
		Return("invalid", nil)

	xc, err := s.service.resolveXcode(fmt.Sprintf(install.Path))
	s.Assert().Nil(xc)
	s.Assert().EqualError(err, "plist: type mismatch: tried to decode plist type `string' into value of type `xcode.infoPlist'")
}

func (s *XCodeListTestSuite) TestSpotLightFailure() {
	s.mockExec.
		On("Exec", MdFind, []string{XCodeBundleIdentifier}).
		Return(nil, fmt.Errorf("Error"))

	res, err := s.service.List()

	s.Assert().Error(err)
	s.Assert().Nil(res)
}

func (s *XCodeListTestSuite) TestList() {
	s.mockExec.
		On("Exec", MdFind, []string{XCodeBundleIdentifier}).
		Return(XCODES, nil)

	s.mockFileService.On("IsDir", "/Applications/Xcode.app").Return(true, nil)
	s.mockFileService.On("IsDir", "/Applications/Xcode 10.3.1.app").Return(true, nil)
	s.mockFileService.On("IsDir", "/Invalid/path").Return(false, nil)

	s.mockFileService.
		On("OpenAndReadFileContent", "/Applications/Xcode.app"+ContentPListFile).
		Return(xcodeContentPListFile("1.2.3"), nil)

	s.mockFileService.
		On("OpenAndReadFileContent", "/invalid/path"+ContentPListFile).
		Return(xcodeContentPListFile("1.2.3"), nil)

	s.mockFileService.
		On("OpenAndReadFileContent", "/Applications/Xcode 10.3.1.app"+ContentPListFile).
		Return(xcodeContentPListFile("10.3.1"), nil)

	res, err := s.service.List()
	s.Assert().NoError(err)

	s.Assert().EqualValues(res, []*Install{
		&Install{"/Applications/Xcode.app", "1.2.3", "1.2.3-snapshot"},
		&Install{"/Applications/Xcode 10.3.1.app", "10.3.1", "10.3.1-snapshot"},
	})
}
