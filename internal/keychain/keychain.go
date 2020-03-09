package keychain

import (
	"bufio"
	"bytes"
	"dothething/internal/util"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	SecurityUtil = "security"

	ActionCreateKeychain   = "create-keychain"
	ActionDeleteKeychain   = "delete-keychain"
	ActionSettings         = "set-keychain-settings"
	ActionSetPartitionList = "set-key-partition-list"
	ActionListKeyChains    = "list-keychains"
)

type KeyChain interface {
	Create(password string) (*os.File, error)
}

type KeyChainHandler struct {
	exec         util.Exec
	keychainFile *os.File
}

func NewKeyChain(exec util.Exec) (KeyChain, error) {
	tmpFile, err := ioutil.TempFile("", "test.*.keychain")
	if err != nil {
		return nil, err
	}
	return KeyChainHandler{exec: exec, keychainFile: tmpFile}, nil
}

func (k KeyChainHandler) Create(password string) (*os.File, error) {
	return nil, nil

	// k.exec.Exec()
}

func (k KeyChainHandler) ImportCertificate(file io.Reader, password string) error {
	return nil
}

// createKeychain Create keychain with provided password
func (k KeyChainHandler) createKeychain(password string) error {
	if len(password) == 0 {
		return errors.New("Keychain password should not be empty")
	}

	res, err := k.exec.Exec(nil, SecurityUtil, ActionCreateKeychain,
		"-p", password, // Use password as the password for the keychains being created.
		k.keychainFile.Name())

	fmt.Println(res, err)

	if err != nil {
		return err
	}

	return nil
}

// configureKeychain : Set settings for keychain, or the default keychain if none is specified
func (k KeyChainHandler) configureKeychain() error {
	// Omitting the timeout argument (-t) specified no-timeout
	if _, err := k.exec.Exec(nil, SecurityUtil, ActionSettings, k.keychainFile.Name()); err != nil {
		return err
	}

	return nil
}

// setPartitionList :  Sets the "partition list" for a key. The "partition list" is an extra
// parameter in the ACL which limits access to the key based on an application's code signature.
func (k KeyChainHandler) setPartitionList(password string, description string) error {
	if _, err := k.exec.Exec(nil, SecurityUtil,
		ActionSetPartitionList,
		"-S", "apple:,apple-tool:,codesign:", // Partition ID
		"-s",           // Match keys that can sign
		"-k", password, // Password for keychain
		"-D", description, // Match description string
		"-t", "private", // We are looking for a private key
		k.keychainFile.Name()); err != nil {
		return err
	}

	return nil
}

// deleteKeyChain: Delete keychains and remove them from the search list
func (k KeyChainHandler) deleteKeyChain() error {
	if _, err := k.exec.Exec(nil,
		SecurityUtil, ActionDeleteKeychain,
		k.keychainFile.Name()); err != nil {
		return err
	}
	return nil
}

func (k KeyChainHandler) addKeyChainToSearchList() error {
	list, err := k.getSearchList()
	if err != nil {
		return err
	}

	return k.setSearchList(append(list, k.keychainFile.Name()))
}

func (k KeyChainHandler) removeKeyChainToSearchList() error {
	return nil
}

func (k KeyChainHandler) getSearchList() ([]string, error) {
	// Display the the keychain search list without any specified domain
	// TODO: Maybe necessary to define the domain later?
	b, err := k.exec.Exec(nil, SecurityUtil, ActionListKeyChains)
	if err != nil {
		return nil, err
	}

	res := []string{}

	// Parse each line and try to parse the keychain path
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, `"`)
		line = strings.TrimSuffix(line, `"`)

		// TODO: Should we add validation of the keychain path ??

		res = append(res, line)
	}

	return res, nil
}

func (k KeyChainHandler) setSearchList(list []string) error {
	args := []string{ActionListKeyChains, "-s"}
	args = append(args, list...)
	if _, err := k.exec.Exec(nil, SecurityUtil, args...); err != nil {
		return err
	}

	return nil
}
