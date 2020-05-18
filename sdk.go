package commercio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/commercionetwork/commercionetwork/app"
	"github.com/commercionetwork/sacco.go"
)

const (
	// default hrp
	hrp = "did:com:"

	// default derivation path
	derivationPath = sacco.CosmosDerivationPath
)

var (
	// DefaultConfig represents the default configuration for a commercio.network SDK instance.
	DefaultConfig = SDKConfig{
		DerivationPath: derivationPath,
		hrp:            hrp,
		LCDEndpoint:    "http://localhost:1317",
		Mode:           TxModeSync,
	}
)

// TxMode represents the mode used by the SDK to broadcast the transaction.
// TxMode can be either:
//  - `sync`: the LCD will do basic validity checks on the messages, will not wait for the message to be included in a block; it'll always return no error.
// 	- `async`: like `sync`, but no checks are performed.
//  - `block`: like `sync`, but it will wait for the message to be included in a block; it could return error.
type TxMode string

// make TxMode behave like a sacco TxMode
func (tm TxMode) asSaccoMode() sacco.TxMode { return sacco.TxMode(tm) }

const (
	// TxModeSync represents the `sync` transaction mode.
	TxModeSync = "sync"

	// TxModeAsync represents the `async` transaction mode.
	TxModeAsync = "async"

	// TxModeBlock represents the `block` transaction mode.
	TxModeBlock = "block"
)

// SDKConfig allows callers to customize default behaviors the commercio.network SDK assumes.
type SDKConfig struct {
	// DerivationPath represents the derivation path used while performing crypto-related operations.
	DerivationPath string

	// hrp represents the human-readable part, placed before the address encoded in Bech32.
	hrp string

	// LCDEndpoint is the commercio.network REST LCD server endpoint, where transaction will be broadcasted.
	LCDEndpoint string

	// Mode is the TxMode to be used while performing transaction-related operations.
	Mode TxMode
}

// validate checks that each and every field of sc are complying with the specification (no empty fields).
func (sc SDKConfig) validate() error {
	if sc.DerivationPath == "" {
		return errors.New("missing derivation path")
	}

	if sc.LCDEndpoint == "" {
		return errors.New("missing LCD endpoint")
	}

	_, err := url.Parse(sc.LCDEndpoint)
	if err != nil || !(strings.HasPrefix(sc.LCDEndpoint, "http://") || strings.HasPrefix(sc.LCDEndpoint, "https://")) {
		return errors.New("malformed LCD endpoint")
	}

	if sc.Mode != TxModeSync && sc.Mode != TxModeAsync && sc.Mode != TxModeBlock {
		return errors.New("invalid transaction mode")
	}

	return nil
}

// SDK represents the entrypoint for the commercio.network SDK.
type SDK struct {
	wallet      *sacco.Wallet
	config      SDKConfig
	typeMapping typeMapping
}

// NewSDK returns a new instance of SDK initialized by given mnemonic and config.
func NewSDK(mnemonic string, config SDKConfig) (*SDK, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("%w, %s", ErrNewSDK, err.Error())
	}

	config.hrp = hrp

	w, err := sacco.FromMnemonic(config.hrp, mnemonic, config.DerivationPath)
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrNewSDK, err.Error())
	}

	codec := app.MakeCodec()

	return &SDK{wallet: w, config: config, typeMapping: generateTypeMappings(codec)}, nil
}

// SendTransaction sends all the messages contained in rawMsgs through the pre-defined LCD, then returns the transaction
// hash.
func (sdk *SDK) SendTransaction(rawMsgs ...interface{}) (string, error) {
	txp, err := sdk.genTx(rawMsgs...)
	if err != nil {
		return "", err
	}

	return sdk.wallet.SignAndBroadcast(txp, sdk.config.LCDEndpoint, sdk.config.Mode.asSaccoMode())
}

func (sdk *SDK) genTx(rawMsgs ...interface{}) (sacco.TransactionPayload, error) {
	if len(rawMsgs) == 0 {
		return sacco.TransactionPayload{}, errors.New("no message provided")
	}

	msgs := make([]json.RawMessage, len(rawMsgs))

	for i := 0; i < len(rawMsgs); i++ {
		var err error

		enclosure := messageEnclosure{
			Type:  sdk.typeMapping.cosmosType(rawMsgs[i]),
			Value: rawMsgs[i],
		}

		msgs[i], err = json.Marshal(enclosure)
		if err != nil {
			return sacco.TransactionPayload{}, fmt.Errorf("%w, message #%d: %s", ErrInvalidMessage, i, err.Error())
		}
	}

	fee := sacco.Fee{
		Amount: make([]sacco.Coin, len(msgs)),
		Gas:    "200000",
	}

	for i := range msgs {
		fee.Amount[i] = sacco.Coin{
			Denom:  "ucommercio",
			Amount: "10000",
		}
	}

	return sacco.TransactionPayload{
		Message: msgs,
		Fee:     fee,
	}, nil
}
