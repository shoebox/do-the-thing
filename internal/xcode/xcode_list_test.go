package xcode

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"dothething/internal/utiltest"

	"github.com/stretchr/testify/mock"
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
	ctx             context.Context
	cmd             *utiltest.MockExecutorCmd
	mockExec        *utiltest.MockExecutor2
	mockFileService *utiltest.MockFileService
	service         listService
}

func (s *XCodeListTestSuite) SetupTest() {
	utiltest.SetupMockExec()
	s.cmd = new(utiltest.MockExecutorCmd)
	s.ctx = context.Background()
	s.mockExec = new(utiltest.MockExecutor2)
	s.mockFileService = new(utiltest.MockFileService)
	s.service = listService{exec: s.mockExec, file: s.mockFileService}
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
	s.mockReturn(
		MdFind,
		[]string{XCodeBundleIdentifier},
		"Hello world",
		nil,
	)

	// when:
	wb, _ := s.service.spotlightSearch(s.ctx)

	// then:
	s.Assert().NotNil(wb)

	// and:
	s.mockExec.AssertExpectations(s.T())
}

func (s *XCodeListTestSuite) mockReturn(cmd string, arg []string, res string, err error) {
	s.cmd.
		On("Output").
		Return(res, err)

	s.mockExec.
		On("CommandContext",
			mock.Anything,
			MdFind, []string{XCodeBundleIdentifier}).
		Return(s.cmd)
}

func (s *XCodeListTestSuite) TestSpotLightFailure() {
	s.mockReturn(
		MdFind,
		[]string{XCodeBundleIdentifier},
		"",
		fmt.Errorf("Error"),
	)

	res, err := s.service.List(s.ctx)

	s.Assert().Error(err)
	s.Assert().Nil(res)
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

func (s *XCodeListTestSuite) TestParseSpotLightSearchResults() {
	// setup:
	s.cmd.
		On("Output").
		Return("result", errors.New("toto"))

	s.mockExec.On("CommandContext",
		mock.Anything,
		MdFind,
		[]string{XCodeBundleIdentifier}).Return(s.cmd)

	// when:
	cases := []struct {
		path    string
		valid   bool
		version string
	}{
		{path: "/Applications/Xcode.app", valid: true, version: "13.1.1"},
		{path: "/Applications/Xcode-10.3.1.app", valid: true, version: "10.3.1"},
		{path: "/invalid/path", valid: false},
	}

	// when:
	for _, c := range cases {
		b := s.mockFileService.
			On("OpenAndReadFileContent", fmt.Sprintf("%v%v", c.path, ContentPListFile))

		if c.valid {
			b.Return(xcodeContentPListFile(c.version), nil)
		} else {
			b.Return(nil, errors.New("error text"))
		}
	}

	// then:
	for _, c := range cases {
		i, err := s.service.resolveXcode(c.path)

		if c.valid {
			s.Assert().EqualValues(i, &Install{
				Path:          c.path,
				BundleVersion: c.version,
				Version:       fmt.Sprintf("%v-snapshot", c.version),
			})
		} else {
			s.Assert().EqualError(err, "error text")
		}
	}
}

func (s *XCodeListTestSuite) TestParseSpotLightEntryShouldHandleError() {
	// setup:
	s.mockFileService.On("IsDir", "/toto/tata.app").Return(nil, errors.New("toto"))

	// when:
	_, err := s.service.parseSpotlightEntry("/toto/tata.app")

	// then:
	s.Assert().EqualError(err, "toto")
}

func (s *XCodeListTestSuite) TestParseSpotLightEntryShouldParseResult() {
	version := "10.3.1"
	path := "/Applications/Xcode-10.3.1.app"

	// setup:
	s.mockFileService.On("IsDir", path).
		Return(true, nil)

	s.mockFileService.
		On("OpenAndReadFileContent", fmt.Sprintf("%v%v", path, ContentPListFile)).
		Return(xcodeContentPListFile("10.3.1"), nil)

	// when:
	i, err := s.service.parseSpotlightEntry(path)

	// then:
	s.Assert().EqualValues(&Install{
		Path:          path,
		BundleVersion: version,
		Version:       fmt.Sprintf("%v-snapshot", version),
	}, i)
	s.Assert().NoError(err)
}

func (s *XCodeListTestSuite) TestListShouldHandleError() {
	// setup:
	s.cmd.
		On("Output").
		Return("result", errors.New("toto"))

	s.mockExec.On("CommandContext",
		mock.Anything,
		MdFind,
		[]string{XCodeBundleIdentifier}).Return(s.cmd)

	// when:
	installs, err := s.service.List(s.ctx)

	// then:
	s.Assert().EqualError(err, "toto")
	s.Assert().Empty(installs)
}

func (s XCodeListTestSuite) TestListShouldParseResults() {
	// setup:
	path := "/Applications/Xcode-10.3.1.app"
	version := "10.3.1"

	// setup:
	s.mockFileService.On("IsDir", path).
		Return(true, nil)

	s.mockFileService.
		On("OpenAndReadFileContent", fmt.Sprintf("%v%v", path, ContentPListFile)).
		Return(xcodeContentPListFile("10.3.1"), nil)

	s.cmd.
		On("Output").
		Return(path, nil)

	s.mockExec.On("CommandContext",
		mock.Anything,
		MdFind,
		[]string{XCodeBundleIdentifier}).Return(s.cmd)

	// when:
	installs, err := s.service.List(s.ctx)

	// then:
	s.Assert().NoError(err)
	s.Assert().EqualValues([]*Install{&Install{
		Path:          path,
		BundleVersion: version,
		Version:       fmt.Sprintf("%v-snapshot", version),
	}}, installs)
	s.Assert().NoError(err)
}
