package keychain

import (
	"context"
	"dothething/internal/api"
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
	*api.API
}

func NewKeyChain(api *api.API) (api.KeyChain, error) {
	return keychain{API: api}, nil
}

// Delete will delete the keychain and remove them from the search list
func (k keychain) Delete(ctx context.Context) error {
	b, err := k.securityCmd(
		ctx,
		ActionDeleteKeychain,
		[]string{k.API.PathService.KeyChain()},
	).Output()

	if err != nil {
		err = KeyChainError{msg: deleteError}
	}

	log.Info().
		Bytes("Result", b).
		Msg("Deletion result")

	return err
}

// ImportCertificate Import one item into a keychain
func (k keychain) ImportCertificate(ctx context.Context, filePath, password, commonName string) error {
	log.Info().
		Str("FilePath", filePath).
		Msg("Importing Certificate")
	b, err := k.securityCmd(
		ctx,
		ActionImport,
		[]string{
			filePath,
			FlagKeychain, k.API.PathService.KeyChain(), // Specify keychain into which item(s) will be imported.
			FlagPassphase, password, // Specify the unwrapping passphrase immediately.
			FlagAppPath, "/usr/bin/codesign", // Specify an application which may access the imported key;
			FlagNonExtractable,
		},
	).Output()

	log.Info().Bytes("Result", b).Msg("Importation result")

	if err != nil {
		return CertificateImportError(err)
	}

	return k.setPartitionList(ctx, "dothething")
}

// setPartitionList :  Sets the "partition list" for a key. The "partition list" is an extra
// parameter in the ACL which limits access to the key based on an application's code signature.
func (k keychain) setPartitionList(ctx context.Context, password string) error {
	log.Debug().Msg("Set partition list")
	b, err := k.securityCmd(
		ctx,
		ActionSetPartitionList,
		[]string{
			FlagPartitionList, "apple:,apple-tool:,codesign:", // Partition ID
			"-s",           // Match keys that can sign
			"-k", password, // Password for keychain
			"-t", "private", // We are looking for a private key
			k.API.PathService.KeyChain(),
		},
	).Output()

	log.Info().Msg(string(b))

	if err != nil {
		err = KeyChainError{msg: partitionError, err: err}
	}

	return err
}
