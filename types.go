package commercio

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// setCosmosConfig sets up the Commercio.network HRP for Cosmos-SDK function calling.
func setCosmosConfig() {
	// this bit comes directly from commercio.network app/app.go.
	bech32MainPrefix := "did:com:"

	prefixValidator := "val"
	prefixConsensus := "cons"
	prefixPublic := "pub"
	prefixOperator := "oper"

	bech32PrefixAccAddr := bech32MainPrefix
	bech32PrefixAccPub := bech32MainPrefix + prefixPublic
	bech32PrefixValAddr := bech32MainPrefix + prefixValidator + prefixOperator
	bech32PrefixValPub := bech32MainPrefix + prefixValidator + prefixOperator + prefixPublic
	bech32PrefixConsAddr := bech32MainPrefix + prefixValidator + prefixConsensus
	bech32PrefixConsPub := bech32MainPrefix + prefixValidator + prefixConsensus + prefixPublic

	config := types.GetConfig()
	config.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
	config.Seal()
}

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
func Amount(amount uint64) types.Coins {
	c, err := types.ParseCoins(fmt.Sprintf("%ducommercio", amount))
	if err != nil {
		panic(fmt.Errorf("could not convert well-known field to coins, %w", err))
	}

	return c
}

//
// Commercio.network type exports
//

// MsgSend represents a message which sends some coins from an address to another
type MsgSend bank.MsgSend

// MsgShareDoc represents a message which shares a document from an account to another
type MsgShareDoc docs.MsgShareDocument
