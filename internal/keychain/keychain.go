package keychain

import (
	"bufio"
	"bytes"
	"context"
	"dothething/internal/util"
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog/log"
)

const (
	SecurityUtil = "security"

	ActionCreateKeychain   = "create-keychain"
	ActionDeleteKeychain   = "delete-keychain"
	ActionImport           = "import"
	ActionSettings         = "set-keychain-settings"
	ActionSetPartitionList = "set-key-partition-list"
	ActionListKeyChains    = "list-keychains"
	FlagAppPath            = "-T"
	FlagKeychain           = "-k"
	FlagPartitionList      = "-S"
	FlagPassword           = "-p"
	FlagPassphase          = "-P"
)

type KeyChain interface {
	Create(ctx context.Context, password string) error
	Delete(ctx context.Context) error
	ImportCertificate(ctx context.Context, filePathh string, password string, identity string) error
	GetPath() string
}

type keychain struct {
	executor util.Executor
	filePath string
}

func NewKeyChain(executor util.Executor) (KeyChain, error) {
	tmpDir, err := ioutil.TempDir("", "do-the-thing-*")
	if err != nil {
		return nil, err
	}

	return keychain{
		executor: executor,
		filePath: filepath.Join(tmpDir, "do-the-thing.keychain"),
	}, nil
}

func (k keychain) GetPath() string {
	return k.filePath
}

// Create will create a new temporary keychhain and add it to the
// search list
func (k keychain) Create(ctx context.Context, password string) error {
	if err := k.createKeychain(ctx, password); err != nil {
		return err
	}

	if err := k.configureKeychain(ctx); err != nil {
		return err
	}

	return k.addKeyChainToSearchList(ctx)
}

// Delete will delete the keychain and remove them from the search list
func (k keychain) Delete(ctx context.Context) error {
	_, err := k.executor.CommandContext(ctx,
		SecurityUtil, ActionDeleteKeychain,
		k.filePath).Output()

	return err
}

// ImportCertificate Import one item into a keychain
func (k keychain) ImportCertificate(ctx context.Context, filePath, password, identity string) error {
	_, err := k.executor.CommandContext(ctx,
		SecurityUtil,
		ActionImport,
		filePath,
		FlagKeychain, k.filePath, // Specify keychain into which item(s) will be imported.
		FlagPassphase, password, // Specify the unwrapping passphrase immediately.
		FlagAppPath, "/usr/bin/codesign"). // Specify an application which may access the imported key;
		Output()

	return err
}

// createKeychain Create keychain with provided password
func (k keychain) createKeychain(ctx context.Context, password string) error {
	if len(password) == 0 {
		return errors.New("Keychain password should not be empty")
	}

	_, err := k.executor.CommandContext(ctx,
		SecurityUtil,
		ActionCreateKeychain,
		FlagPassword, password, // Use password as the password for the keychains being created.
		k.filePath).Output()

	return err
}

// configureKeychain : Set settings for keychain, or the default keychain if none is specified
func (k keychain) configureKeychain(ctx context.Context) error {
	// Omitting the timeout argument (-t) specified no-timeout
	_, err := k.executor.CommandContext(ctx, SecurityUtil, ActionSettings, k.filePath).Output()

	return err
}

// setPartitionList :  Sets the "partition list" for a key. The "partition list" is an extra
// parameter in the ACL which limits access to the key based on an application's code signature.
func (k keychain) setPartitionList(ctx context.Context, password string, description string) error {
	b, err := k.executor.CommandContext(ctx,
		SecurityUtil,
		ActionSetPartitionList,
		FlagPartitionList, "apple:,apple-tool:,codesign:", // Partition ID
		"-s",           // Match keys that can sign
		"-k", password, // Password for keychain
		"-D", description, // Match description string
		"-t", "private", // We are looking for a private key
		k.filePath).
		Output()

	log.Debug().Msg(string(b))

	return err
}

func (k keychain) addKeyChainToSearchList(ctx context.Context) error {
	list, err := k.getSearchList(ctx)
	if err != nil {
		return err
	}

	return k.setSearchList(ctx, append(list, k.filePath))
}

func (k keychain) getSearchList(ctx context.Context) ([]string, error) {
	// Display the the keychain search list without any specified domain
	// TODO: Maybe necessary to define the domain later?
	b, err := k.executor.
		CommandContext(ctx, SecurityUtil, ActionListKeyChains).
		Output()
	if err != nil {
		return nil, err
	}

	return parseSearchList(b), nil
}

func parseSearchList(data []byte) []string {
	var res []string
	r := regexp.MustCompile(`^(?:\t|(?:\s)+)?"(.*)"`)

	// Parse each line and try to parse the keychain path
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		txt := scanner.Text()
		if r.MatchString(txt) {
			sm := r.FindStringSubmatch(txt)
			res = append(res, sm[1])
		}
	}

	return res
}

func (k keychain) setSearchList(ctx context.Context, list []string) error {
	args := append([]string{ActionListKeyChains, "-s"}, list...)
	_, err := k.executor.
		CommandContext(ctx, SecurityUtil, args...).
		Output()
	return err
}
