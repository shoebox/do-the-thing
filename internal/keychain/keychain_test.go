package keychain

const (
	k1 = "/test/test1.keychain"
	k2 = "/test/test2/test3.keychain"
)

var (
	errText      = "mock error"
	passwd       = "p4ssword"
	keychainPath = "/path/to/file.keychain"
)

/*
type keychainTestSuite struct {
	suite.Suite
	API     *api.API
	pm      *api.PathMock
	exec    *utiltest.MockExecutor
	subject keychain
}

func TestKeychainSuite(t *testing.T) {
	suite.Run(t, new(keychainTestSuite))
}

func (s *keychainTestSuite) BeforeTest(suiteName, testName string) {
	s.exec = new(utiltest.MockExecutor)
	s.pm = new(api.PathMock)

	s.API = &api.API{
		Exec:        &s.exec,
		PathService: s.pm,
	}
	s.subject = keychain{API: s.API}

	// setup:
	s.pm.On("KeyChain").Return(k1)

}

func (s *keychainTestSuite) AfterTest(suiteName, testName string) {
}

func (s *keychainTestSuite) TestNewKeychain() {
	// when:
	k, err := NewKeyChain(s.API)

	// then:
	s.Assert().NoError(err)

	// and:
	s.Assert().NotNil(k)
}

func (s *keychainTestSuite) TestCreateShouldHandleErrors() {
	errText := "mock error"

	// setup:
	s.exec.MockCommandContext(
		SecurityUtil,
		[]string{ActionCreateKeychain, FlagPassword, "p4ssword", s.pm.KeyChain()},
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
		k1,
	},
		"",
		errors.New(errText),
	)

	// when:
	s.pm.On("Keychain").Return(keychainPath)
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
		FlagKeychain, k1,
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
	s.exec.MockCommandContextError(SecurityUtil, []string{ActionDeleteKeychain, k1}, errors.New("any"))

	// when:
	err := s.subject.Delete(context.Background())

	// then:
	s.Assert().EqualError(err, deleteError)
}
*/
