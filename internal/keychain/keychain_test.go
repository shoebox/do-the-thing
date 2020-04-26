package keychain

import (
	"context"
	"dothething/internal/utiltest"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

var errText = "mock error"

type keychainTestSuite struct {
	ctx    context.Context
	cancel context.CancelFunc
	cmd    *utiltest.MockExecutorCmd
	suite.Suite
	subject  keychain
	executor *utiltest.MockExecutor2
}

func TestKeychainSuite(t *testing.T) {
	suite.Run(t, new(keychainTestSuite))
}

func (s *keychainTestSuite) BeforeTest(suiteName, testName string) {
	s.executor = new(utiltest.MockExecutor2)
	s.cmd = new(utiltest.MockExecutorCmd)
	s.subject = keychain{s.executor, "/path/to/file.keychain"}
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 60*time.Second)
}

func (s *keychainTestSuite) AfterTest(suiteName, testName string) {
	s.cancel()
}

func (s *keychainTestSuite) TestCreateShouldHandleErrors() {
	errText := "mock error"

	// setup:
	s.configureCommandOutput(
		SecurityUtil,
		[]string{ActionCreateKeychain, FlagPassword, "p4ssword", s.subject.GetPath()},
		"",
		&errText)

	// when:
	err := s.subject.Create(s.ctx, "p4ssword")

	// then:
	s.Assert().EqualError(err, errText)
}

func (s *keychainTestSuite) TestImportCertificateShouldHandleErrors() {
	// setup:
	passwd := "p4ssword"

	s.configureCommandOutput(
		SecurityUtil,
		[]string{
			ActionImport,
			"toto/file.p12",
			"-k", s.subject.GetPath(),
			"-P", passwd,
			"-T", "/usr/bin/codesign",
		},
		"",
		&errText)

	// when:
	err := s.subject.ImportCertificate(s.ctx, "toto/file.p12", passwd, "")

	// then:
	s.Assert().EqualError(err, errText)
}

func (s *keychainTestSuite) TestKeyChainConfiguration() {
	// setup
	s.configureCommandOutput(
		SecurityUtil,
		[]string{ActionSettings, s.subject.GetPath()},
		"",
		&errText)

	// when:
	err := s.subject.configureKeychain(s.ctx)

	// then:
	s.Assert().EqualError(err, errText)
}

func (s *keychainTestSuite) TestKeyChainCreatePasswordShouldNotBeEmpty() {
	// when:
	err := s.subject.createKeychain(s.ctx, "")

	// then:
	s.Assert().EqualError(err, "Keychain password should not be empty")
}

func (s *keychainTestSuite) TestSetPartitionList() {
	// setup
	s.configureCommandOutput(
		SecurityUtil,
		[]string{
			ActionSetPartitionList,
			FlagPartitionList, "apple:,apple-tool:,codesign:", // Partition ID
			"-s",             // Match keys that can sign
			"-k", "p4ssword", // Password for keychain
			"-D", "description", // Match description string
			"-t", "private", // We are looking for a private key
			"/path/to/file.keychain",
		},
		"",
		&errText)

	// when:
	err := s.subject.setPartitionList(s.ctx, "p4ssword", "description")
	s.Assert().EqualError(err, errText)
}

func (s *keychainTestSuite) TestDelete() {
	// setup
	file := "/path/to/file.keychain"
	s.configureCommandOutput(
		SecurityUtil,
		[]string{ActionDeleteKeychain, file},
		"",
		&errText)

	// when:
	err := s.subject.Delete(s.ctx)

	// then:
	s.Assert().EqualError(err, errText)
}

func (s *keychainTestSuite) configureCommandOutput(name string,
	args []string,
	result string,
	errTxt *string) {

	var err error
	if errTxt != nil {
		err = errors.New(*errTxt)
	}

	s.cmd.
		On("Output").
		Return(result, err)

	s.executor.
		On("CommandContext", s.ctx, name, args).
		Return(s.cmd)
}

func (s *keychainTestSuite) TestGetSearchListShouldReturnValidData() {
	// setup:
	s.configureCommandOutput(SecurityUtil,
		[]string{ActionListKeyChains}, `
    "/Users/fake/Library/Keychains/login.keychain-db"
    "/Users/fake/toto"`,
		nil)

	// when:
	data, err := s.subject.getSearchList(s.ctx)

	// then:
	s.Assert().NoError(err)
	s.Assert().EqualValues([]string{
		"/Users/fake/Library/Keychains/login.keychain-db",
		"/Users/fake/toto",
	}, data)
}

func (s *keychainTestSuite) TestGetSearchListShouldHandleErrors() {
	// setup:
	s.configureCommandOutput(SecurityUtil,
		[]string{ActionListKeyChains}, "", &errText)

	// when:
	data, err := s.subject.getSearchList(s.ctx)

	// then:
	s.Assert().EqualError(err, errText)
	s.Assert().Empty(data)
}

func (s *keychainTestSuite) TestSetSearchList() {
	// setup:
	k1 := "/test/test1.keychain"
	k2 := "/test/test2/test3.keychain"
	list := []string{k1, k2}

	s.configureCommandOutput(SecurityUtil,
		[]string{ActionListKeyChains, "-s", k1, k2},
		"",
		&errText)

	// when:
	err := s.subject.setSearchList(s.ctx, list)

	// then:
	s.Assert().EqualError(err, errText)
}
