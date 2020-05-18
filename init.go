package commercio

import "github.com/cosmos/cosmos-sdk/types"

// We use init() to configure commercionetwork's Cosmos settings, otherwise address/codec won't work
func init() {
	setCosmosConfig()
}

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
