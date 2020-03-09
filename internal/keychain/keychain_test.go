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

func TestMain(t *testing.T) {
	mockExec = new(utiltest.MockExec)
	var err error
	tmpFile, err = ioutil.TempFile("", "test.*.keychain")
	assert.NoError(t, err)

	subject = KeyChainHandler{exec: mockExec, keychainFile: tmpFile}
	t.Cleanup(func() {
		fmt.Println("clean")
	})
}

func TestCreateKeyChain(t *testing.T) {
	t.Run("Password should not be empty", func(t *testing.T) {
		// when:
		err := subject.createKeychain("")

		// then:
		assert.EqualError(t, err, "Keychain password should not be empty")

		// and:
		mockExec.AssertExpectations(t)
	})

	t.Run("Should create keychain", func(t *testing.T) {
		// Should create the keychain
		mockExec.
			On(SecurityUtil, ActionCreateKeychain, "-p", "p4ssword", mock.Anything).
			Return()

		// when:
		err := subject.createKeychain("p4ssword")

		// then:
		assert.NoError(t, err)

		// and:
		mockExec.AssertExpectations(t)
	})
}

func TestKeyChainConfiguration(t *testing.T) {
	// Should try to configure the keychain
	mockExec.
		On("security", "set-keychain-settings", mock.Anything).
		Return()

	mockExec.
		On(SecurityUtil,
			ActionSetPartitionList,
			"-S",
			"apple-tool:,apple:",
			"-s",
			"-k",
			"p4ssword",
			mock.Anything)

	// when:
	err := subject.configureKeychain()

	// then:
	assert.NoError(t, err)

	// and:
	mockExec.AssertExpectations(t)
}

func TestKeyChainDelete(t *testing.T) {
	//
	mockExec.
		On(SecurityUtil, ActionDeleteKeychain, mock.Anything).
		Return()

	// when:
	err := subject.deleteKeyChain()

	// then:
	assert.NoError(t, err)

	// and:
	mockExec.AssertExpectations(t)
}

func TestAddKeyChainToSearchList(t *testing.T) {
}

func TestGetSearchList(t *testing.T) {
	t.Run("Should return valid data", func(t *testing.T) {
		// setup:
		TestMain(t)
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
		TestMain(t)
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

func TestImportCertificate(t *testing.T) {
	t.Run("Certificate password should be valid", func(t *testing.T) {
		// setup:
		data, err := os.Open("../../assets/Certificate.p12")
		assert.NoError(t, err)
		assert.NotNil(t, data)

		// when:
		err = subject.ImportCertificate(data, "hello:")

		// then:
		assert.EqualError(t, err, "Invalidd")
	})

	t.Run("Importation API need to be invoked", func(t *testing.T) {
	})

}
