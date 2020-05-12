package commercio

import "errors"

var (
	// ErrNewSDK represents some kind of error that happened during the initialization phase of the SDK.
	ErrNewSDK = errors.New("could not initialize SDK")

	// ErrInvalidMessage represents some kind of error that happened during the adaptation of messages to the
	// format used for broadcasting.
	ErrInvalidMessage = errors.New("invalid transaction")
)
