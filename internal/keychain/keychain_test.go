package keychain

import (
	"dothething/internal/utiltest"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockExec *utiltest.MockExec
var subject KeyChainHandler
var tmpFile *os.File

func TestKeychainMain(t *testing.T) {
	mockExec = new(utiltest.MockExec)

	var err error
	tmpFile, err = ioutil.TempFile("", "test.*.keychain")
	assert.NoError(t, err)
	assert.NotNil(t, tmpFile)

	subject = KeyChainHandler{exec: mockExec, filePath: tmpFile.Name()}
	t.Cleanup(func() {
		fmt.Println("clean")
	})
}

func TestCreate(t *testing.T) {
	// setup:
	t.Run("Should handle error", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)

		mockExec.
			On("Exec", SecurityUtil, mock.Anything).
			Return("", errors.New("mock error"))

		// when:
		err := subject.Create("p4ssword")

		// then:
		assert.EqualError(t, err, "mock error")
	})
}

func TestImportCertificate(t *testing.T) {
	t.Run("Should invoke exec API", func(t *testing.T) {
		// setup:
		passwd := "p4ssword"
		filepath := "toto"
		TestKeychainMain(t)

		// Should create the keychain
		mockExec.
			On("Exec", SecurityUtil, []string{ActionImport,
				filepath,
				"-k", tmpFile.Name(),
				"-P", passwd,
				"-T", "/usr/bin/codesign"}).
			Return("", nil)

		// when:
		err := subject.ImportCertificate(filepath, passwd, "")

		// then:
		assert.NoError(t, err)

		// and:
		mockExec.AssertExpectations(t)
	})
}

func TestCreateKeyChain(t *testing.T) {

	t.Run("Password should not be empty", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)

		// when:
		err := subject.createKeychain("")

		// then:
		assert.EqualError(t, err, "Keychain password should not be empty")

		// and:
		mockExec.AssertExpectations(t)
	})

	t.Run("Should create keychain", func(t *testing.T) {
		// setup:
		passwd := "p4ssword"
		TestKeychainMain(t)

		// Should create the keychain
		mockExec.
			On("Exec", SecurityUtil, []string{ActionCreateKeychain,
				"-p", passwd,
				subject.filePath}).
			Return("", nil)

		// when:
		err := subject.createKeychain(passwd)

		// then:
		assert.NoError(t, err)

		// and:
		mockExec.AssertExpectations(t)
	})

	t.Run("Should handle error", func(t *testing.T) {
		// setup:
		passwd := "p4ssword"
		TestKeychainMain(t)

		mockExec.
			On("Exec", SecurityUtil, []string{ActionCreateKeychain,
				"-p", passwd,
				subject.filePath}).
			Return("", errors.New("mock error"))

		// when:
		err := subject.createKeychain(passwd)

		// then:
		assert.EqualError(t, err, "mock error")
	})
}

func TestKeyChainConfiguration(t *testing.T) {
	// setup:
	TestKeychainMain(t)

	// Should try to configure the keychain
	mockExec.
		On("Exec", SecurityUtil, []string{ActionSettings, tmpFile.Name()}).
		Return("", nil)

	// when:
	err := subject.configureKeychain()

	// then:
	assert.NoError(t, err)

	// and:
	mockExec.AssertExpectations(t)
}

func TestSetPartitionList(t *testing.T) {
	// setup:
	password := "p4ssword"
	identity := "identity"

	t.Run("Should handle errors", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil,
				[]string{ActionSetPartitionList,
					"-S", "apple:,apple-tool:,codesign:", // Partition ID
					"-s", // Match signing keys
					"-k", password,
					"-D", identity, // Match description string
					"-t", "private",
					tmpFile.Name()}).
			Return(nil, errors.New("error text"))

		// when:
		err := subject.setPartitionList(password, identity)

		// then:
		assert.EqualError(t, err, "error text")

		// and:
		mockExec.AssertExpectations(t)

	})

	t.Run("Should handle success", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil,
				[]string{ActionSetPartitionList,
					"-S", "apple:,apple-tool:,codesign:", // Partition ID
					"-s", // Match signing keys
					"-k", password,
					"-D", identity, // Match description string
					"-t", "private",
					tmpFile.Name()}).
			Return("", nil)

		// when:
		err := subject.setPartitionList(password, identity)

		// then:
		assert.NoError(t, err)

		// and:
		mockExec.AssertExpectations(t)
	})
}

func TestKeyChainDelete(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil, []string{ActionDeleteKeychain, tmpFile.Name()}).
			Return("", nil)

		// when:
		err := subject.Delete()

		// then:
		assert.NoError(t, err)

		// and:
		mockExec.AssertExpectations(t)
	})

	t.Run("Error case", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil, []string{ActionDeleteKeychain, tmpFile.Name()}).
			Return("", errors.New("Error text"))

		// when:
		err := subject.Delete()

		// then:
		assert.EqualError(t, err, "Error text")

		// and:
		mockExec.AssertExpectations(t)
	})
}

func TestSetSearchList(t *testing.T) {
	k1 := "/test/test1.keychain"
	k2 := "/test/test2/test3.keychain"

	t.Run("Should successfully set list", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil, []string{ActionListKeyChains, "-s", k1, k2}).
			Return("", nil)

		// when:
		err := subject.setSearchList([]string{k1, k2})

		// then:
		assert.NoError(t, err)

		// and:
		mockExec.AssertExpectations(t)
	})

	t.Run("Should handle errors", func(t *testing.T) {
		// setup
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil, []string{ActionListKeyChains, "-s", k1, k2}).
			Return("", errors.New("security: unknown command"))

		// when:
		err := subject.setSearchList([]string{k1, k2})

		// then:
		assert.EqualError(t, err, "security: unknown command")
	})
}

func TestGetSearchList(t *testing.T) {
	t.Run("Should return valid data", func(t *testing.T) {
		// setup:
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil, []string{ActionListKeyChains}).
			Return(`	
	"/Users/fake/Library/Keychains/login.keychain-db"
	"/Users/fake/toto"`, nil)

		// when:
		data, err := subject.getSearchList()

		// then:
		assert.NoError(t, err)
		assert.ObjectsAreEqual(data, []string{
			"/Users/fake/Library/Keychains/login.keychain-db",
			"/Users/fake/toto",
		})
	})

	t.Run("Should handle errors", func(t *testing.T) {
		// setup
		TestKeychainMain(t)
		mockExec.
			On("Exec", SecurityUtil, []string{ActionListKeyChains}).
			Return(nil, errors.New("security: unknown command"))

		// when:
		data, err := subject.getSearchList()

		// then:
		assert.EqualError(t, err, "security: unknown command")
		assert.Nil(t, data)
	})
}

/*
func TestImportCertificate(t *testing.T) {
	t.Run("Certificate password should be valid", func(t *testing.T) {
		// setup:
		data, err := os.Open("../../assets/Certificate.p12")
		assert.NoError(t, err)
		assert.NotNil(t, data)

		// when:
		err = subject.ImportCertificate(data, "hello:")

		// then:
		assert.EqualError(t, err, "Invalid")
	})

	t.Run("Importation API need to be invoked", func(t *testing.T) {
	})

}
*/
