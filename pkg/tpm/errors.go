// Package tpm provides error definitions for tpm operations.
package tpm

import "errors"

var (
	// ErrTPMNotFound is returned when no TPM device is found.
	ErrTPMNotFound = errors.New("no TPM device found")
	// ErrPCRNotFound is returned when a requested PCR is not found in the TPM.
	ErrPCRNotFound = errors.New("requested PCR not found in TPM")
	// ErrTPMNotOpened is returned when attempting to use a TPM device that has not been opened.
	ErrTPMNotOpened = errors.New("TPM device not opened")
	// ErrInvalidNonce is returned when the provided nonce is invalid.
	ErrInvalidNonce = errors.New("invalid nonce")
	// ErrInvalidPCRs is returned when the provided PCRs are invalid.
	ErrInvalidPCRs = errors.New("invalid PCRs provided")
	// ErrNoSignatureCached is returned when no signature is cached and Quote() has not been called.
	ErrNoSignatureCached = errors.New("no signature cached; call Quote() first")
	// ErrNoPublicKeyCached is returned when no public key is cached and Quote() has not been called.
	ErrNoPublicKeyCached = errors.New("no public key cached; call Quote() first")
)
