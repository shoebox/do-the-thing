package path

import (
	"dothething/internal/api"
	"testing"

	"github.com/stretchr/testify/suite"
)

type pathServiceSuite struct {
	suite.Suite
	API     *api.APIMock
	cfg     *api.Config
	subject pathService
}

func TestKeychainSuite(t *testing.T) {
	suite.Run(t, new(pathServiceSuite))
}

func (s *pathServiceSuite) BeforeTest(suiteName, testName string) {
	s.API = new(api.APIMock)
	s.cfg = new(api.Config)
	s.cfg.Path = "/path/to/toto.xcodeproj"
	s.cfg.Scheme = "schemeName"
	s.cfg.Configuration = "configName"
	s.cfg.Target = "targetName"
	s.API.On("Config").Return(s.cfg)
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
