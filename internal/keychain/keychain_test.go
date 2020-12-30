package keychain

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/utiltest"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	k1 = "/test/test1.keychain"
	k2 = "/test/test2/test3.keychain"
)

var (
	errText      = "mock error"
	passwd       = "p4ssword"
	keychainPath = "/path/to/file.keychain"
)

type keychainTestSuite struct {
	suite.Suite
	API     *api.APIMock
	pm      *api.PathMock
	exec    *utiltest.MockExecutor
	subject keychain
}

func TestKeychainSuite(t *testing.T) {
	suite.Run(t, new(keychainTestSuite))
}

func (s *keychainTestSuite) BeforeTest(suiteName, testName string) {
	s.exec = new(utiltest.MockExecutor)
	s.API = new(api.APIMock)
	s.pm = new(api.PathMock)
	s.API.On("Exec").Return(s.exec)
	s.API.On("PathService").Return(s.pm)
	s.subject = keychain{API: s.API}
}

func (s *keychainTestSuite) AfterTest(suiteName, testName string) {
}

func (s *keychainTestSuite) TestNewKeychain() {
	// setup:
	p := "/path/to/k.keychain"
	s.pm.On("KeyChain").Return(p)

	// when:
	k, err := NewKeyChain(s.API)

	// then:
	s.Assert().NoError(err)

	// and:
	s.Assert().Equal(k.GetPath(), p)
}

func (s *keychainTestSuite) TestCreateShouldHandleErrors() {
	errText := "mock error"

	// setup:
	s.exec.MockCommandContext(
		SecurityUtil,
		[]string{ActionCreateKeychain, FlagPassword, "p4ssword", s.subject.GetPath()},
		"",
		errors.New(errText))

	// when:
	err := s.subject.Create(context.Background(), passwd)

	// then:
	s.Assert().EqualValues(err,
		KeyChainError{
			msg: createError,
			err: errors.New(errText),
		})
}

func (s *keychainTestSuite) TestKeyChainCreatePasswordShouldNotBeEmpty() {
	// when:
	err := s.subject.createKeychain(context.Background(), "")

	// then:
	s.Assert().EqualValues(err,
		KeyChainError{
			msg: createError,
			err: errors.New("Keychain password should not be empty"),
		})
}

func (s *keychainTestSuite) TestKeyChainConfiguration() {
	// setup
	s.subject.filePath = k1
	s.exec.MockCommandContext(SecurityUtil,
		[]string{ActionSettings, k1},
		"",
		errors.New(errText))

	// when:
	err := s.subject.configureKeychain(context.Background())

	// then:
	s.Assert().EqualValues(err,
		KeyChainError{
			msg: configureError,
			err: errors.New(errText),
		})
}

func (s *keychainTestSuite) TestSetPartitionList() {
	// setup
	s.exec.MockCommandContext(SecurityUtil, []string{
		ActionSetPartitionList,
		FlagPartitionList, "apple:,apple-tool:,codesign:", // Partition ID
		"-s",         // Match keys that can sign
		"-k", passwd, // Password for keychain
		"-D", "description", // Match description string
		"-t", "private", // We are looking for a private key
		keychainPath,
	},
		"",
		errors.New(errText),
	)

	// when:
	s.subject.filePath = keychainPath
	err := s.subject.setPartitionList(context.Background(), passwd, "description")

	// then:
	s.Assert().EqualValues(err,
		KeyChainError{
			msg: partitionError,
			err: errors.New(errText),
		})
}

func (s *keychainTestSuite) TestImportCertificateShouldHandleErrors() {
	// setup:

	s.exec.MockCommandContext(SecurityUtil, []string{
		ActionImport,
		"toto/file.p12",
		FlagKeychain, s.subject.GetPath(),
		FlagPassphase, passwd,
		FlagAppPath, "/usr/bin/codesign",
		FlagNonExtractable,
	}, "", errors.New(errText))

	// when: when trying to import a certificate who should fail
	err := s.subject.ImportCertificate(context.Background(), "toto/file.p12", passwd)

	// then: error should be reported
	s.Assert().EqualValues(KeyChainError{msg: importError, err: errors.New(errText)}, err)
}

func (s *keychainTestSuite) TestDeleteKeychainShouldHandleError() {
	// setup
	file := "/path/to/file.keychain"
	s.subject.filePath = file
	s.exec.MockCommandContextError(SecurityUtil, []string{ActionDeleteKeychain, file}, errors.New("any"))

	// when:
	err := s.subject.Delete(context.Background())

	// then:
	s.Assert().EqualError(err, deleteError)
}

func (s *keychainTestSuite) TestGetSearchListInputOutputAndParsing() {
	ctx := context.Background()

	// setup:
	cases := []struct {
		expected []string
		output   string
		err      error
	}{
		{
			output:   `    "/Users/user.name/Library/Keychains/login.keychain-db"`,
			expected: []string{"/Users/user.name/Library/Keychains/login.keychain-db"},
			err:      nil,
		},
		{
			output: `    "/Users/user.name/Library/Keychains/login.keychain-db"
    "/Users/user.name/Library/Keychains/test2.keychain-db"`,
			expected: []string{
				"/Users/user.name/Library/Keychains/login.keychain-db",
				"/Users/user.name/Library/Keychains/test2.keychain-db",
			},
			err: nil,
		},
		{
			output: "",
			err:    errors.New("toto"),
		},
	}

	for _, c := range cases {
		// setup:
		s.BeforeTest("keychainTestSuite", "TestSetSearchListShouldHandleErrors")
		s.exec.MockCommandContext(SecurityUtil, []string{ActionListKeyChains}, c.output, c.err)

		// when:
		r, e := s.subject.getSearchList(ctx)

		// and:
		s.Assert().EqualValues(c.expected, r)
		s.Assert().EqualValues(c.err, e)
	}
}

func (s *keychainTestSuite) TestSetSearchListShouldHandleErrors() {
	// setup:
	errText := "mock error"

	s.exec.MockCommandContext(
		SecurityUtil,
		[]string{
			ActionListKeyChains,
			"-s",
			k1,
			k2,
		},
		"",
		errors.New(errText),
	)

	// when:
	err := s.subject.setSearchList(context.Background(), []string{k1, k2})

	// then:
	s.Assert().EqualError(err, errText)
}

func (s *keychainTestSuite) TestParseSearchListLine() {
	cases := []struct {
		line     string
		expected string
	}{
		{
			line:     `    "/Users/user.name/Library/Keychains/login.keychain-db"`,
			expected: `/Users/user.name/Library/Keychains/login.keychain-db`,
		},
		{
			line:     `    "/private/var/folders/m1/s05mrl8s4fbfw8zf4hw04tjwh740pl/T/do-the-thing-856917238/do-the-thing.keychain"`,
			expected: "/private/var/folders/m1/s05mrl8s4fbfw8zf4hw04tjwh740pl/T/do-the-thing-856917238/do-the-thing.keychain",
		},
		{
			line:     "/path/to",
			expected: "",
		},
		{
			line:     `    /path/to`,
			expected: "",
		},
	}

	for _, c := range cases {
		s.Run(c.line, func() {
			// when:
			res := isSearchListEntry(c.line)

			// then:
			s.Assert().EqualValues(res, c.expected)
		})
	}
}
