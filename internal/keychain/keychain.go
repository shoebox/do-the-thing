package keychain

import (
	"bufio"
	"bytes"
	"context"
	"dothething/internal/api"
	"errors"
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
	FlagNonExtractable     = "-x"
)

var searchListRegexp = regexp.MustCompile(`^(?:\t|(?:\s)+)?"(.*)"`)

type keychain struct {
	filePath string
	*api.API
}

func NewKeyChain(api *api.API) (api.KeyChain, error) {
	return keychain{API: api, filePath: api.PathService.KeyChain()}, nil
}

func (k keychain) GetPath() string {
	return k.filePath
}

// Create will create a new temporary keychhain and add it to the search list
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
	_, err := k.API.Exec.CommandContext(ctx,
		SecurityUtil, ActionDeleteKeychain,
		k.filePath).Output()

	if err != nil {
		err = KeyChainError{msg: deleteError}
	}

	return err
}

// ImportCertificate Import one item into a keychain
func (k keychain) ImportCertificate(ctx context.Context, filePath, password string) error {
	log.Info().
		Str("FilePath", filePath).
		Msg("Importing Certificate")
	_, err := k.API.Exec.CommandContext(ctx,
		SecurityUtil,
		ActionImport,
		filePath,
		FlagKeychain, k.filePath, // Specify keychain into which item(s) will be imported.
		FlagPassphase, password, // Specify the unwrapping passphrase immediately.
		FlagAppPath, "/usr/bin/codesign", // Specify an application which may access the imported key;
		FlagNonExtractable).
		Output()

	if err != nil {
		err = CertificateImportError(err)
	}

	return err
}

// createKeychain Create keychain with provided password
func (k keychain) createKeychain(ctx context.Context, password string) error {
	if len(password) == 0 {
		return KeyChainError{msg: createError, err: errors.New("Keychain password should not be empty")}
	}

	_, err := k.API.Exec.CommandContext(ctx,
		SecurityUtil,
		ActionCreateKeychain,
		FlagPassword, password, // Use password as the password for the keychains being created.
		k.filePath).Output()

	if err != nil {
		err = KeyChainError{msg: createError, err: err}
	}

	return err
}

// configureKeychain : Set settings for keychain, or the default keychain if none is specified
func (k keychain) configureKeychain(ctx context.Context) error {
	// Omitting the timeout argument (-t) specified no-timeout
	_, err := k.API.Exec.
		CommandContext(ctx, SecurityUtil, ActionSettings, k.filePath).
		Output()

	if err != nil {
		err = KeyChainError{msg: configureError, err: err}
	}

	return err
}

// setPartitionList :  Sets the "partition list" for a key. The "partition list" is an extra
// parameter in the ACL which limits access to the key based on an application's code signature.
func (k keychain) setPartitionList(ctx context.Context, password string, description string) error {
	b, err := k.API.Exec.CommandContext(ctx,
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

	if err != nil {
		err = KeyChainError{msg: partitionError, err: err}
	}

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
	var res []string
	// Display the the keychain search list without any specified domain
	// TODO: Maybe necessary to define the domain later?
	b, err := k.API.Exec.
		CommandContext(ctx, SecurityUtil, ActionListKeyChains).
		Output()

	if err != nil {
		return res, err
	}

	return parseSearchList(b), nil
}

func parseSearchList(data []byte) []string {
	var res []string

	// Parse each line and try to parse the keychain p ath
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		txt := scanner.Text()
		l := isSearchListEntry(txt)
		if l != "" {
			res = append(res, l)
		}
	}

	return res
}

func isSearchListEntry(txt string) string {
	if searchListRegexp.MatchString(txt) {
		sm := searchListRegexp.FindStringSubmatch(txt)
		return sm[1]
	}

	return ""
}

// setSearchList will set the keychain search list
func (k keychain) setSearchList(ctx context.Context, list []string) error {
	// appending argumets
	args := append([]string{ActionListKeyChains, "-s"}, list...)

	//
	_, err := k.API.Exec.
		CommandContext(ctx, SecurityUtil, args...).
		Output()

	// TODO: Handle return

	return err
}
