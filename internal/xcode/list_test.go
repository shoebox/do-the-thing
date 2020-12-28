package xcode

import (
	"bytes"
	"fmt"
	"testing"

	"dothething/internal/api"
	"dothething/internal/util"

	"github.com/stretchr/testify/suite"
)

var install = api.Install{
	Path:          "/Applications/Xcode.app",
	BundleVersion: "1.2.3",
	Version:       "1.2.3-snapshot",
}

type XCodeListTestSuite struct {
	suite.Suite
	/*
		ctx             context.Context
		cmd             *utiltest.MockExecutorCmd
		mockExec        *utiltest.MockExecutor
		mockFileService *utiltest.MockFileService
		service         listService
	*/
	*listService
}

func (s *XCodeListTestSuite) SetupTest() {
	s.listService = new(listService)
	/*
		utiltest.SetupMockExec()
		s.cmd = new(utiltest.MockExecutorCmd)
		s.ctx = context.Background()
		s.mockExec = new(utiltest.MockExecutor)
		s.mockFileService = new(utiltest.MockFileService)
		s.service = listService{exec: s.mockExec, file: s.mockFileService}
	*/
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

func (s *XCodeListTestSuite) TestResolveInstall() {
	cases := []struct {
		i *api.Install
		p string
		e error
		c string
	}{
		{
			i: &api.Install{
				Path:          "/path/to/project.xcodeproj",
				BundleVersion: "10.0.1",
				Version:       "10.0.1-snapshot",
			},
			p: "/path/to/project.xcodeproj",
			c: xcodeContentPListFile("10.0.1"),
		},
		{
			i: nil,
			c: "toto",
			e: util.DecodingError{},
		},
	}

	for _, c := range cases {
		s.Run(c.c, func() {
			// when:
			i, e := s.listService.resolveInstall(c.p, bytes.NewReader([]byte(c.c)))

			// then:
			s.Assert().EqualValues(c.i, i)

			//
			s.Assert().Equal(c.e, e)
		})
	}
}
