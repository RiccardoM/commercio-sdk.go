package commercio

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/commercionetwork/sacco.go"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestSDKConfig_validate(t *testing.T) {
	tests := []struct {
		name    string
		config  SDKConfig
		wantErr bool
	}{
		{
			"default config must not error",
			DefaultSDKConfig,
			false,
		},
		{
			"missing derivation path",
			SDKConfig{},
			true,
		},
		{
			"missing lcd endpoint",
			SDKConfig{
				DerivationPath: sacco.CosmosDerivationPath,
			},
			true,
		},
		{
			"lcd url doesn't begin with http:// or https://",
			SDKConfig{
				DerivationPath: sacco.CosmosDerivationPath,
				LCDEndpoint:    "aaa.com",
			},
			true,
		},
		{
			"missing tx mode",
			SDKConfig{
				DerivationPath: sacco.CosmosDerivationPath,
				LCDEndpoint:    "http://aaa.com",
			},
			true,
		},
		{
			"invalid tx mode",
			SDKConfig{
				DerivationPath: sacco.CosmosDerivationPath,
				LCDEndpoint:    "http://aaa.com",
				Mode:           TxMode("asd"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				require.Error(t, tt.config.validate())
				return
			}

			require.NoError(t, tt.config.validate())
		})
	}
}

func TestNewSDK(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		config   SDKConfig
		wantErr  bool
	}{
		{
			"missing mnemonic",
			"",
			DefaultSDKConfig,
			true,
		},
		{
			"invalid config",
			"mnemonic",
			SDKConfig{},
			true,
		},
		{
			"well-formed mnemonic and config",
			"first purse atom language viable marble switch industry pill prevent drive develop prison art hard useless search shoulder promote rapid split wrestle balcony focus",
			DefaultSDKConfig,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := NewSDK(tt.mnemonic, tt.config)

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

func TestSDK_genTx(t *testing.T) {
	sdk, err := NewSDK("first purse atom language viable marble switch industry pill prevent drive develop prison art hard useless search shoulder promote rapid split wrestle balcony focus", DefaultSDKConfig)
	require.NoError(t, err)

	tests := []struct {
		name    string
		rawMsgs []interface{}
		wantErr bool
	}{
		{
			"no raw messages",
			nil,
			true,
		},
		{
			"message marshaling error",
			[]interface{}{make(chan int)},
			true,
		},
		{
			"everything is fine",
			[]interface{}{"hello"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := sdk.genTx(tt.rawMsgs...)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, sacco.TransactionPayload{}, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.rawMsgs), len(res.Message))
		})
	}
}

func TestSDK_SendTransaction(t *testing.T) {
	sdk, err := NewSDK("first purse atom language viable marble switch industry pill prevent drive develop prison art hard useless search shoulder promote rapid split wrestle balcony focus", DefaultSDKConfig)
	require.NoError(t, err)

	okResponder := httpmock.NewJsonResponderOrPanic(http.StatusOK, sacco.TxResponse{
		Code:   0,
		TxHash: "ok!",
	})

	errResponder := httpmock.NewJsonResponderOrPanic(http.StatusForbidden, sacco.Error{
		Error: "error!",
	})

	tests := []struct {
		name      string
		msgs      []interface{}
		responder httpmock.Responder
		wantErr   bool
	}{
		{
			"trigger error in genTx",
			nil,
			errResponder,
			true,
		},
		{
			"error from the LCD endpoint",
			[]interface{}{MsgSend{}},
			errResponder,
			true,
		},
		{
			"no error from the LCD endpoint",
			[]interface{}{MsgSend{}},
			okResponder,
			false,
		},
	}
	for _, tt := range tests {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, "http://localhost:1317/node_info", httpmock.NewJsonResponderOrPanic(http.StatusOK, sacco.NodeInfo{}))
		httpmock.RegisterResponder(http.MethodPost, "http://localhost:1317/txs", tt.responder)
		httpmock.RegisterRegexpResponder(http.MethodGet, regexp.MustCompile("http://localhost:1317/auth/accounts/(.+)"), httpmock.NewJsonResponderOrPanic(http.StatusOK, sacco.AccountData{Result: sacco.AccountDataResult{Value: sacco.AccountDataValue{Address: "address"}}}))

		t.Run(tt.name, func(t *testing.T) {
			var res string
			var err error

			if tt.msgs == nil {
				res, err = sdk.SendTransaction()
			} else {
				res, err = sdk.SendTransaction(tt.msgs)
			}

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, "", res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "ok!", res) // we can assume that because we're forcing okResponder
		})
	}
}
