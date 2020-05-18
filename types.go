package commercio

import (
	"errors"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// messageEnclosure encloses a Cosmos message into its REST-accepted enclosure.
type messageEnclosure struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Address returns str as a Cosmos-compatible address, given str as a bech32-encoded string.
func Address(str string) (types.AccAddress, error) {
	return types.AccAddressFromBech32(str)
}

// Amount returns a Cosmos-compatible Commercio.network amount, expressed in ucommercio.
func Amount(amount uint64) (types.Coins, error) {
	if amount == 0 { // an uint64 can at most be zero!
		return nil, errors.New("amount cannot be zero")
	}

	c, err := types.ParseCoins(fmt.Sprintf("%ducommercio", amount))
	if err != nil {
		// we panic here because since we hardcode the "ucommercio" amount, something should go *really* wrong
		// for ParseCoins to return error, hence we must stop execution.
		panic(fmt.Errorf("could not convert well-known field to coins, %w", err))
	}

	return c, nil
}

//
// Commercio.network type exports
//

// MsgSend represents a message which sends some coins from an address to another
type MsgSend bank.MsgSend

// MsgShareDoc represents a message which shares a document from an account to another
type MsgShareDoc docs.MsgShareDocument
