package commercio

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
)

func Test_generateTypeMappings(t *testing.T) {
	tests := []struct {
		name    string
		want    typeMapping
		cdcFunc func() *codec.Codec
	}{
		{
			"no structs registered",
			typeMapping{},
			func() *codec.Codec {
				return codec.New()
			},
		},
		{
			"a struct registered",
			typeMapping{
				"Str": "st/Str",
			},
			func() *codec.Codec {
				type Str struct{}
				cdc := codec.New()
				cdc.RegisterConcrete(Str{}, "st/Str", nil)
				return cdc
			},
		},
		{
			"an interface registered",
			typeMapping{},
			func() *codec.Codec {
				type Str interface {
					Dummy()
				}
				cdc := codec.New()
				cdc.RegisterInterface((*Str)(nil), nil)
				return cdc
			},
		},
		{
			"a struct and an interface registered",
			typeMapping{
				"Str": "st/Str",
			},
			func() *codec.Codec {
				type StrI interface {
					Dummy()
				}

				type Str struct{}

				cdc := codec.New()
				cdc.RegisterConcrete(Str{}, "st/Str", nil)
				cdc.RegisterInterface((*StrI)(nil), nil)
				return cdc
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc := tt.cdcFunc()

			require.Equal(t, tt.want, generateTypeMappings(cdc))
		})
	}
}

func Test_typeMapping_cosmosType(t *testing.T) {
	tests := []struct {
		name     string
		cdcFunc  func() *codec.Codec
		typeFunc func() interface{}
		want     string
	}{
		{
			"nil type",
			func() *codec.Codec {
				return codec.New()
			},
			func() interface{} {
				return nil
			},
			"",
		},
		{
			"registered type",
			func() *codec.Codec {
				type Str struct{}
				cdc := codec.New()
				cdc.RegisterConcrete(Str{}, "st/Str", nil)
				return cdc
			},
			func() interface{} {
				type Str struct{}
				return Str{}
			},
			"st/Str",
		},
		{
			"unregistered type",
			func() *codec.Codec {
				type Str struct{}
				cdc := codec.New()
				cdc.RegisterConcrete(Str{}, "st/Str", nil)
				return cdc
			},
			func() interface{} {
				type Strr struct{}
				return Strr{}
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc := tt.cdcFunc()
			m := generateTypeMappings(cdc)
			tp := tt.typeFunc()
			require.Equal(t, tt.want, m.cosmosType(tp))
		})
	}
}
