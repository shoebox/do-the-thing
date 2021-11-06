package path

import (
	"dothething/internal/api"
	"testing"

	"github.com/stretchr/testify/suite"
)

type pathServiceSuite struct {
	suite.Suite
	API     *api.API
	cfg     *api.Config
	subject pathService
}

func TestKeychainSuite(t *testing.T) {
	suite.Run(t, new(pathServiceSuite))
}

func (s *pathServiceSuite) BeforeTest(suiteName, testName string) {
	s.API = &api.API{
		Config: &api.Config{
			Path:          "/path/to/toto.xcodeproj",
			Scheme:        "schemeName",
			Configuration: "configName",
			Target:        "targetName",
		},
	}
	s.subject = pathService{s.API}
}

func (s *pathServiceSuite) AfterTest(suiteName, testName string) {
}

func (s *pathServiceSuite) TestBuildFolder() {
	// when:
	p := s.subject.buildFolder()

	// then:
	s.Assert().Equal("/path/to/Build", p)
}

func (s *pathServiceSuite) TestArchive() {
	// when:
	p := s.subject.Archive()

	// then:
	s.Assert().Equal("/path/to/Build/targetName-schemeName-configName.xcarchive", p)
}

func (s *pathServiceSuite) TestKeyChain() {
	// when:
	p := s.subject.KeyChain()

	// then:
	s.Assert().Equal("/path/to/Build/do-the-thing.keychain", p)
}

func (s *pathServiceSuite) TestExportPlist() {
	// when:
	p := s.subject.ExportPList()

	// then:
	s.Assert().Equal("/path/to/Build/targetName-schemeName-configName-export.plist", p)
}

func (s *pathServiceSuite) TestObjRoot() {
	// when:
	p := s.subject.ObjRoot()

	// then:
	s.Assert().Equal("OBJROOT=/path/to/Build/obj", p)
}

func (s *pathServiceSuite) TestSymRoot() {
	// when:
	p := s.subject.SymRoot()
	// then:
	s.Assert().Equal("SYMROOT=/path/to/Build/sym", p)
}

func (s *pathServiceSuite) TestXCResult() {
	// when:
	p := s.subject.XCResult()

	// then:
	s.Assert().Equal("/path/to/Build/targetName-schemeName-configName.xcresult", p)
}

func (s *pathServiceSuite) TestXCodeProject() {
	cases := []struct {
		Path     string
		Expected string
	}{
		{
			Path:     "/path/to/project.xcodeproj",
			Expected: "/path/to/project.xcodeproj",
		},
		{
			Path:     "/path/to/project.xcworkspace",
			Expected: "/path/to/project.xcodeproj",
		},
	}

	for _, c := range cases {
		s.API.Config.Path = c.Path

		// when:
		res := s.subject.XCodeProject()

		// then:
		s.Assert().EqualValues(c.Expected, res)
	}
}
