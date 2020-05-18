package commercio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	tests := []struct {
		name    string
		addrStr string
		wantErr bool
	}{
		{
			"not bech32 string",
			"aaa",
			true,
		},
		{
			"a commercio bech32 address",
			"did:com:1rv8jkqulyf5j55pcjte7v8fg6h0gxcerw8a042",
			false,
		},
		{
			"a valid bech32, but hrp is not did:com:",
			"cosmos1s5afhd6gxevu37mkqcvvsj8qeylhn0rz46zdlq",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Address(tt.addrStr)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
		})
	}
}

func TestAmount(t *testing.T) {
	tests := []struct {
		name       string
		coinAmount uint64
		wantErr    bool
	}{
		{
			"zero coins",
			0,
			true,
		},
		{
			"some coins",
			42,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Amount(tt.coinAmount)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
		})
	}
}
