package keychain

import (
	"context"
	"errors"
	"fmt"
)

// Create will create a new temporary keychhain and add it to the search list
func (k keychain) Create(ctx context.Context, password string) error {
	if err := k.createKeychain(ctx, password); err != nil {
		return fmt.Errorf("failed to create keychain (Error: %v", err)
	}

	if err := k.configureKeychain(ctx); err != nil {
		return fmt.Errorf("failed to configure keychain (Error: %v", err)
	}

	err := k.addKeyChainToSearchList(ctx)
	if err != nil {
		return fmt.Errorf("failed to add the keychain to the search list (Error: %v", err)
	}

	return nil
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
		k.API.PathService.KeyChain()).Output()

	if err != nil {
		err = KeyChainError{msg: createError, err: err}
	}

	return err
}

// configureKeychain : Set settings for keychain, or the default keychain if none is specified
func (k keychain) configureKeychain(ctx context.Context) error {
	// Omitting the timeout argument (-t) specified no-timeout
	_, err := k.API.Exec.
		CommandContext(ctx, SecurityUtil, ActionSettings, k.API.PathService.KeyChain()).
		Output()

	if err != nil {
		err = KeyChainError{msg: configureError, err: err}
	}

	return err
}
