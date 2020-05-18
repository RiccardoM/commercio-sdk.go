package commercio

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_setCosmosConfig(t *testing.T) {
	// setCosmosConfig() gets called when the package initializes, no need to call it already,
	// just check that the required values are there
	config := types.GetConfig()
	require.Equal(t, "did:com:", config.GetBech32AccountAddrPrefix())
	require.Equal(t, "did:com:pub", config.GetBech32AccountPubPrefix())
	require.Equal(t, "did:com:valcons", config.GetBech32ConsensusAddrPrefix())
	require.Equal(t, "did:com:valconspub", config.GetBech32ConsensusPubPrefix())
	require.Equal(t, "did:com:valoper", config.GetBech32ValidatorAddrPrefix())
	require.Equal(t, "did:com:valoperpub", config.GetBech32ValidatorPubPrefix())
}
