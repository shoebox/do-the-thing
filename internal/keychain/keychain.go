package keychain

import (
	"bufio"
	"bytes"
	"dothething/internal/util"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	SecurityUtil = "security"

	ActionCreateKeychain   = "create-keychain"
	ActionDeleteKeychain   = "delete-keychain"
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
	file, err := ioutil.TempFile("", "do-the-thing.*.keychain")
	if err != nil {
		return nil, err
	}

	return file, err

	// k.exec.Exec()
}

func (k KeyChainHandler) ImportCertificate(file io.Reader, password string) error {
	return nil
}

func (k KeyChainHandler) createKeychain(password string) error {
	return nil
}

func (k KeyChainHandler) configureKeychain() error {
	// set-keychain-settings
	// set-key-partition-list
	return nil
}

func (k KeyChainHandler) deleteKeyChain() error {
	return nil
}

func (k KeyChainHandler) addKeyChainToSearchList() error {
	return nil
}

func (k KeyChainHandler) removeKeyChainToSearchList() error {
	return nil
}

func (k KeyChainHandler) getSearchList() ([]string, error) {
	b, err := k.exec.Exec(nil, SecurityUtil, ActionListKeyChains)
	if err != nil {
		return nil, err
	}

	res := []string{}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, `"`)
		line = strings.TrimSuffix(line, `"`)

		res = append(res, line)
	}

	return res, nil
}

func (k KeyChainHandler) setSearchList([]*os.File) error {
	return nil
}

/*
# create new empty keychain
(x)security create-keychain -p "${ADMIN_PASSWORD}" "${tempKeychain}"

# add keychain to user's keychain search list so they can access it
 security list-keychains -d user -s "${tempKeychain}" $(security list-keychains -d user | tr -d '"')

# removing relock timeout on keychain
(x) security set-keychain-settings "${tempKeychain}"

# import the certs
(x) security import foo.p12 -k "${tempKeychain}" -P "${CERT_PASSWORD}" -T "/usr/bin/codesign"

# tell os it's ok to access this identity from command line with tools shipped by apple (suppress codesign modal UI)
(x) security set-key-partition-list -S apple-tool:,apple: -s -k "$ADMIN_PASSWORD" -D "${identity}" -t private ${tempKeychain}

# set default keychain to temp keychain
security default-keychain -d user -s ${tempKeychain}

# unlock keychain
security unlock-keychain -p ${ADMIN_PASSWORD} ${tempKeychain}

# prove we added the code signing identity to the temp keychain
security find-identity -v -p codesigning

# do some codesign stuff

# clean up temp keychain we created
(x) security delete-keychain ${tempKeychain}
*/
