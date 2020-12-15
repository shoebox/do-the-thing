package keychain

import (
	"fmt"
)

const (
	configureError = "Failed to configure keychain"
	createError    = "Failed to create keychain"
	deleteError    = "Failed to delete keychain"
	fileError      = "Keychain file error"
	importError    = "Failed to import certificate int keychain"
	partitionError = "Failed to set partition list"
)

type KeyChainError struct {
	msg string
	err error
}

func (k KeyChainError) Error() string {
	if k.err != nil {
		return fmt.Sprintf("KeyChain error: %v (%v)", k.msg, k.err)
	}

	return k.msg
}

func CertificateImportError(err error) error {
	return KeyChainError{msg: importError, err: err}
}
