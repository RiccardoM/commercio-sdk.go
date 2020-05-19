package commercio

import "errors"

var (
	// ErrNewSDK represents some kind of error that happened during the initialization phase of the SDK.
	ErrNewSDK = errors.New("could not initialize SDK")

	// ErrInvalidMessage represents some kind of error that happened during the adaptation of messages to the
	// format used for broadcasting.
	ErrInvalidMessage = errors.New("invalid transaction")

	// ErrInvalidPublicKey is used when a Bech32-encoded public key is invalid or malformed.
	ErrInvalidPublicKey = errors.New("invalid public key")

	// ErrInvalidAddress is used when a Bech32-encoded address is invalid or malformed.
	ErrInvalidAddress = errors.New("invalid address key")

	// ErrInvalidSignatureKey represents an error returned when the RSA signature key is invalid.
	ErrInvalidSignatureKey = errors.New("invalid signature key")

	// ErrInvalidTumblerKey represents an error returned when the RSA tumbler key is invalid.
	ErrInvalidTumblerKey = errors.New("invalid signature key")

	// ErrInvalidVerificationKey represents an error returned when the RSA verification key is invalid.
	ErrInvalidVerificationKey = errors.New("invalid verification key")

	// ErrProofCreation represents an error returned when the DidDocument proof creation fails.
	ErrProofCreation = errors.New("cannot create proof")

	// ErrInvalidAmount represents an error returned when the provided uint64 amount is invalid.
	ErrInvalidAmount = errors.New("invalid coin amount")

	// ErrInvalidPowerupParams represents an error returned when the provided power up params are invalid.
	ErrInvalidPowerupParams = errors.New("invalid powerup params")

	// ErrNotEnoughEntropy represents an error returned when there's not enough entropy to generate random data.
	ErrNotEnoughEntropy = errors.New("not enough entropy")

	// ErrEncryptionFailure represents an error returned when some error happens during the encryption process.
	ErrEncryptionFailure = errors.New("encryption failure")
)
