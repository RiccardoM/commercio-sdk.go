package commercio

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type failingReader struct{}

func (f failingReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("error!")
}

func Test_readKey(t *testing.T) {
	tests := []struct {
		name      string
		keyReader io.Reader
		keyType   keyType
		want      string
		ifaceNil  bool
		wantErr   bool
	}{
		{
			"cannot read from reader",
			failingReader{},
			typePublicKey,
			"",
			true,
			true,
		},
		{
			"reader returns empty string",
			strings.NewReader(""),
			typePublicKey,
			"",
			true,
			true,
		},
		{
			"reader returns bogus data",
			strings.NewReader("aaaaa"),
			typePublicKey,
			"",
			true,
			true,
		},
		{
			"data is not typePublicKey",
			strings.NewReader("-----BEGIN CERTIFICATE REQUEST-----\nMIICXTCCAUUCAQAwGDEWMBQGA1UEAwwNaGkgd2lraXBlZGlhITCCASIwDQYJKoZI\nhvcNAQEBBQADggEPADCCAQoCggEBAMTwzCYD+iLlDwTu5Y43aQH9q1LF3kgot8I4\n9ZgbFhDmCE4YlLhZKO4hieK6z8z+IfZjfapn01rzuzvTHESj5bSSU6AcEsKSOgTQ\nuB+KKn4mgngyBrJwWjr4IZ9XkGsCLAP2/wkyJC2ire6FuTSQ00YGhKf1B3WbIBbn\n5i1rvZXnYxlheWlNSmxx54q4gTwcd/V4nS4BThYA/ypATjHS/gfQ650cOQzRK/Jh\nWfAbfnETYUpD6MCgZAIbaBuYvYpQEGqQ4niTvtSd07RHKnewcPFqJhMV86qN4HQY\n4ZBNzQcF/2aCGHYyRniKznSDNijT2kaAz/L7ORqh+90qH/BLnKsCAwEAAaAAMA0G\nCSqGSIb3DQEBCwUAA4IBAQAqV5g9AZGXEbM97ouTGDJqFNP2QjO9ZK9J3BOUTrFO\ntMUrVWj+ixhC6vXD3o5uVL/fg6OlmK+13gsBpzg2mq72TBrZsNOK4+O0XvltIvSx\n0H5tf1NYwuHxFgHDqgs/fQBOKFTadebJZHbPBtMrqlnenKYJiVb5YSWBZ7JKRCK7\nVSgwNxxAMnSCNI0xF3EjZ1bjQkM8xGhnwe+n/RAd5Q2pMLIrquMoGMTUYLOq1xSB\nsGTp8iLWbbWPl6gC1hcSMpFsbdyjMCWs+a2R2F8QnahrRfvpgFEndvzA2EvqHIoR\nBHE1ChD7l691PxZP1eKA1I4AzZno5sb6SWyd8+pqY0oG\n-----END CERTIFICATE REQUEST-----"),
			typePublicKey,
			"",
			true,
			true,
		},
		{
			"data is not typePrivateKey",
			strings.NewReader("-----BEGIN CERTIFICATE REQUEST-----\nMIICXTCCAUUCAQAwGDEWMBQGA1UEAwwNaGkgd2lraXBlZGlhITCCASIwDQYJKoZI\nhvcNAQEBBQADggEPADCCAQoCggEBAMTwzCYD+iLlDwTu5Y43aQH9q1LF3kgot8I4\n9ZgbFhDmCE4YlLhZKO4hieK6z8z+IfZjfapn01rzuzvTHESj5bSSU6AcEsKSOgTQ\nuB+KKn4mgngyBrJwWjr4IZ9XkGsCLAP2/wkyJC2ire6FuTSQ00YGhKf1B3WbIBbn\n5i1rvZXnYxlheWlNSmxx54q4gTwcd/V4nS4BThYA/ypATjHS/gfQ650cOQzRK/Jh\nWfAbfnETYUpD6MCgZAIbaBuYvYpQEGqQ4niTvtSd07RHKnewcPFqJhMV86qN4HQY\n4ZBNzQcF/2aCGHYyRniKznSDNijT2kaAz/L7ORqh+90qH/BLnKsCAwEAAaAAMA0G\nCSqGSIb3DQEBCwUAA4IBAQAqV5g9AZGXEbM97ouTGDJqFNP2QjO9ZK9J3BOUTrFO\ntMUrVWj+ixhC6vXD3o5uVL/fg6OlmK+13gsBpzg2mq72TBrZsNOK4+O0XvltIvSx\n0H5tf1NYwuHxFgHDqgs/fQBOKFTadebJZHbPBtMrqlnenKYJiVb5YSWBZ7JKRCK7\nVSgwNxxAMnSCNI0xF3EjZ1bjQkM8xGhnwe+n/RAd5Q2pMLIrquMoGMTUYLOq1xSB\nsGTp8iLWbbWPl6gC1hcSMpFsbdyjMCWs+a2R2F8QnahrRfvpgFEndvzA2EvqHIoR\nBHE1ChD7l691PxZP1eKA1I4AzZno5sb6SWyd8+pqY0oG\n-----END CERTIFICATE REQUEST-----"),
			typePrivateKey,
			"",
			true,
			true,
		},
		{
			"data is typePrivateKey",
			strings.NewReader("-----BEGIN PRIVATE KEY-----\nMIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEA/cBbcJSfO4QFfXKO\nsCGlC06y2oWw2zxa61BmaTsNA7TIDlQD6cjzZ2PRyRRisw5iuiBdz04pA8CF+FZu\nQwbHJwIDAQABAkAPNXBFlyLUFl2d3zfeJqYVv2nI3ypyeXOZlwAMXpWxGw41IMnm\n5LzziTRdbVBDf+wdFfLDe+70bsWGCG4awF3BAiEA/xvGFG0toMTfEeRDtQNJL5dn\nEcoj3lq8IxD2bc+YdhsCIQD+o16NfLIDXdDg1Ir2oKNVvfZPDh0ZGtPFIedKzTVz\n5QIgZ63W+/g/QgahDjlyFwAF33St6/n2R+kiazH6pThoox8CIQCCAn2DNdhZuaut\nLzeoRko+u9enc2hN6hGXxACog2+4NQIgcSK9PYyuRqNS8IfcHWbnMUkvuxsQrjte\np+OJ7MHMC4E=\n-----END PRIVATE KEY-----"),
			typePrivateKey,
			"-----BEGIN PRIVATE KEY-----\nMIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEA/cBbcJSfO4QFfXKO\nsCGlC06y2oWw2zxa61BmaTsNA7TIDlQD6cjzZ2PRyRRisw5iuiBdz04pA8CF+FZu\nQwbHJwIDAQABAkAPNXBFlyLUFl2d3zfeJqYVv2nI3ypyeXOZlwAMXpWxGw41IMnm\n5LzziTRdbVBDf+wdFfLDe+70bsWGCG4awF3BAiEA/xvGFG0toMTfEeRDtQNJL5dn\nEcoj3lq8IxD2bc+YdhsCIQD+o16NfLIDXdDg1Ir2oKNVvfZPDh0ZGtPFIedKzTVz\n5QIgZ63W+/g/QgahDjlyFwAF33St6/n2R+kiazH6pThoox8CIQCCAn2DNdhZuaut\nLzeoRko+u9enc2hN6hGXxACog2+4NQIgcSK9PYyuRqNS8IfcHWbnMUkvuxsQrjte\np+OJ7MHMC4E=\n-----END PRIVATE KEY-----",
			false,
			false,
		},
		{
			"data is typePublicKey",
			strings.NewReader("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyz9KTGQpQIoACbF2LKy4\n3QhENi/7/vtQl/iANUhxHo315b+0+u+imrhqG0iADQsMN/x7QequCLmFWXtowT35\nhU/YE5jnf3jO0AkH9/nWk5OBLTiF7Hr0R6zLRw6V/MKgaq8HvG54hJqqDE/BKI9c\nNxc6ucS+8xQSq3PZm+KRjZqOs8nVyG6lW+P59OMBY6nFDyLD4Ym7CkZ2uPSyDphz\n35fQeSTLSzjD0iebix4WV7w7zV1UbDDy0qKL64hKx0gTa7F1Kf4L3WUfCnVQppcq\n7w9WxaM22bxPin+XYE2B3j64QX1JC7ybI9O1DvpM9CSVXhpiIbJ5m8BsPoM4GAdi\nHwIDAQAB\n-----END PUBLIC KEY-----"),
			typePublicKey,
			"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyz9KTGQpQIoACbF2LKy4\n3QhENi/7/vtQl/iANUhxHo315b+0+u+imrhqG0iADQsMN/x7QequCLmFWXtowT35\nhU/YE5jnf3jO0AkH9/nWk5OBLTiF7Hr0R6zLRw6V/MKgaq8HvG54hJqqDE/BKI9c\nNxc6ucS+8xQSq3PZm+KRjZqOs8nVyG6lW+P59OMBY6nFDyLD4Ym7CkZ2uPSyDphz\n35fQeSTLSzjD0iebix4WV7w7zV1UbDDy0qKL64hKx0gTa7F1Kf4L3WUfCnVQppcq\n7w9WxaM22bxPin+XYE2B3j64QX1JC7ybI9O1DvpM9CSVXhpiIbJ5m8BsPoM4GAdi\nHwIDAQAB\n-----END PUBLIC KEY-----",
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, i, e := readKey(tt.keyReader, tt.keyType)

			require.Equal(t, tt.want, s)
			require.Equal(t, tt.ifaceNil, i == nil)

			if tt.wantErr {
				require.Error(t, e)
				return
			}

			require.NoError(t, e)
		})
	}
}

func TestPowerUpParams_validate(t *testing.T) {
	addr, err := Address("did:com:1rv8jkqulyf5j55pcjte7v8fg6h0gxcerw8a042")
	require.NoError(t, err)

	tests := []struct {
		name    string
		params  PowerUpParams
		wantErr bool
	}{
		{
			"missing pubkey",
			PowerUpParams{},
			true,
		},
		{
			"missing tumbler key",
			PowerUpParams{
				PubKey: "pk",
			},
			true,
		},
		{
			"missing signature key",
			PowerUpParams{
				PubKey:     "pk",
				TumblerKey: &bytes.Buffer{},
			},
			true,
		},
		{
			"missing amount",
			PowerUpParams{
				PubKey:       "pk",
				TumblerKey:   &bytes.Buffer{},
				SignatureKey: &bytes.Buffer{},
			},
			true,
		},
		{
			"missing pairwise address",
			PowerUpParams{
				PubKey:       "pk",
				TumblerKey:   &bytes.Buffer{},
				SignatureKey: &bytes.Buffer{},
				Amount:       42,
			},
			true,
		},
		{
			"all ok",
			PowerUpParams{
				PubKey:          "pk",
				TumblerKey:      &bytes.Buffer{},
				SignatureKey:    &bytes.Buffer{},
				Amount:          42,
				PairwiseAddress: addr,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				require.Error(t, tt.params.validate())
			} else {
				require.NoError(t, tt.params.validate())
			}
		})
	}
}

func TestSDK_BuildPowerupRequest(t *testing.T) {
	sdk, err := NewSDK("first purse atom language viable marble switch industry pill prevent drive develop prison art hard useless search shoulder promote rapid split wrestle balcony focus", DefaultSDKConfig)
	require.NoError(t, err)

	addr, err := Address("did:com:1rv8jkqulyf5j55pcjte7v8fg6h0gxcerw8a042")
	require.NoError(t, err)

	tests := []struct {
		name    string
		params  PowerUpParams
		wantErr bool
	}{
		{
			"invalid params",
			PowerUpParams{},
			true,
		},
		{
			"invalid public key",
			PowerUpParams{
				PubKey:          "pk",
				TumblerKey:      &bytes.Buffer{},
				SignatureKey:    &bytes.Buffer{},
				Amount:          42,
				PairwiseAddress: addr,
			},
			true,
		},
		{
			"invalid signature key reader",
			PowerUpParams{
				PubKey:          "did:com:pub1addwnpepqdr89xxl6pwpj87tzsycmlr035tcmpvcc7xadz5vr2nq9nmcu5hp7xytrlz",
				TumblerKey:      &bytes.Buffer{},
				SignatureKey:    &bytes.Buffer{},
				Amount:          42,
				PairwiseAddress: addr,
			},
			true,
		},
		{
			"invalid tumbler key reader",
			PowerUpParams{
				PubKey:          "did:com:pub1addwnpepqdr89xxl6pwpj87tzsycmlr035tcmpvcc7xadz5vr2nq9nmcu5hp7xytrlz",
				TumblerKey:      &bytes.Buffer{},
				SignatureKey:    strings.NewReader("-----BEGIN PRIVATE KEY-----\nMIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEA/cBbcJSfO4QFfXKO\nsCGlC06y2oWw2zxa61BmaTsNA7TIDlQD6cjzZ2PRyRRisw5iuiBdz04pA8CF+FZu\nQwbHJwIDAQABAkAPNXBFlyLUFl2d3zfeJqYVv2nI3ypyeXOZlwAMXpWxGw41IMnm\n5LzziTRdbVBDf+wdFfLDe+70bsWGCG4awF3BAiEA/xvGFG0toMTfEeRDtQNJL5dn\nEcoj3lq8IxD2bc+YdhsCIQD+o16NfLIDXdDg1Ir2oKNVvfZPDh0ZGtPFIedKzTVz\n5QIgZ63W+/g/QgahDjlyFwAF33St6/n2R+kiazH6pThoox8CIQCCAn2DNdhZuaut\nLzeoRko+u9enc2hN6hGXxACog2+4NQIgcSK9PYyuRqNS8IfcHWbnMUkvuxsQrjte\np+OJ7MHMC4E=\n-----END PRIVATE KEY-----"),
				Amount:          42,
				PairwiseAddress: addr,
			},
			true,
		},
		{
			"all ok",
			PowerUpParams{
				PubKey:          "did:com:pub1addwnpepqdr89xxl6pwpj87tzsycmlr035tcmpvcc7xadz5vr2nq9nmcu5hp7xytrlz",
				TumblerKey:      strings.NewReader("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyz9KTGQpQIoACbF2LKy4\n3QhENi/7/vtQl/iANUhxHo315b+0+u+imrhqG0iADQsMN/x7QequCLmFWXtowT35\nhU/YE5jnf3jO0AkH9/nWk5OBLTiF7Hr0R6zLRw6V/MKgaq8HvG54hJqqDE/BKI9c\nNxc6ucS+8xQSq3PZm+KRjZqOs8nVyG6lW+P59OMBY6nFDyLD4Ym7CkZ2uPSyDphz\n35fQeSTLSzjD0iebix4WV7w7zV1UbDDy0qKL64hKx0gTa7F1Kf4L3WUfCnVQppcq\n7w9WxaM22bxPin+XYE2B3j64QX1JC7ybI9O1DvpM9CSVXhpiIbJ5m8BsPoM4GAdi\nHwIDAQAB\n-----END PUBLIC KEY-----"),
				SignatureKey:    strings.NewReader("-----BEGIN PRIVATE KEY-----\nMIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEA/cBbcJSfO4QFfXKO\nsCGlC06y2oWw2zxa61BmaTsNA7TIDlQD6cjzZ2PRyRRisw5iuiBdz04pA8CF+FZu\nQwbHJwIDAQABAkAPNXBFlyLUFl2d3zfeJqYVv2nI3ypyeXOZlwAMXpWxGw41IMnm\n5LzziTRdbVBDf+wdFfLDe+70bsWGCG4awF3BAiEA/xvGFG0toMTfEeRDtQNJL5dn\nEcoj3lq8IxD2bc+YdhsCIQD+o16NfLIDXdDg1Ir2oKNVvfZPDh0ZGtPFIedKzTVz\n5QIgZ63W+/g/QgahDjlyFwAF33St6/n2R+kiazH6pThoox8CIQCCAn2DNdhZuaut\nLzeoRko+u9enc2hN6hGXxACog2+4NQIgcSK9PYyuRqNS8IfcHWbnMUkvuxsQrjte\np+OJ7MHMC4E=\n-----END PRIVATE KEY-----"),
				Amount:          42,
				PairwiseAddress: addr,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, e := sdk.BuildPowerupRequest(tt.params)

			if tt.wantErr {
				require.Error(t, e)
				require.Equal(t, MsgRequestDidPowerUp{}, p)
				return
			}

			require.NoError(t, e)
			require.NotEqual(t, MsgRequestDidPowerUp{}, p)
		})
	}
}

func TestSDK_BuildDidDocument(t *testing.T) {
	sdk, err := NewSDK("first purse atom language viable marble switch industry pill prevent drive develop prison art hard useless search shoulder promote rapid split wrestle balcony focus", DefaultSDKConfig)
	require.NoError(t, err)

	tests := []struct {
		name         string
		pubkey       string
		sigKeyReader io.Reader
		verKeyReader io.Reader
		wantErr      bool
	}{
		{
			"missing pubkey",
			"",
			&bytes.Buffer{},
			&bytes.Buffer{},
			true,
		},
		{
			"missing signature key",
			"did:com:pub1addwnpepqdr89xxl6pwpj87tzsycmlr035tcmpvcc7xadz5vr2nq9nmcu5hp7xytrlz",
			&bytes.Buffer{},
			&bytes.Buffer{},
			true,
		},
		{
			"missing verification key",
			"did:com:pub1addwnpepqdr89xxl6pwpj87tzsycmlr035tcmpvcc7xadz5vr2nq9nmcu5hp7xytrlz",
			strings.NewReader("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyz9KTGQpQIoACbF2LKy4\n3QhENi/7/vtQl/iANUhxHo315b+0+u+imrhqG0iADQsMN/x7QequCLmFWXtowT35\nhU/YE5jnf3jO0AkH9/nWk5OBLTiF7Hr0R6zLRw6V/MKgaq8HvG54hJqqDE/BKI9c\nNxc6ucS+8xQSq3PZm+KRjZqOs8nVyG6lW+P59OMBY6nFDyLD4Ym7CkZ2uPSyDphz\n35fQeSTLSzjD0iebix4WV7w7zV1UbDDy0qKL64hKx0gTa7F1Kf4L3WUfCnVQppcq\n7w9WxaM22bxPin+XYE2B3j64QX1JC7ybI9O1DvpM9CSVXhpiIbJ5m8BsPoM4GAdi\nHwIDAQAB\n-----END PUBLIC KEY-----"),
			&bytes.Buffer{},
			true,
		},
		{
			"all ok",
			"did:com:pub1addwnpepqdr89xxl6pwpj87tzsycmlr035tcmpvcc7xadz5vr2nq9nmcu5hp7xytrlz",
			strings.NewReader("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyz9KTGQpQIoACbF2LKy4\n3QhENi/7/vtQl/iANUhxHo315b+0+u+imrhqG0iADQsMN/x7QequCLmFWXtowT35\nhU/YE5jnf3jO0AkH9/nWk5OBLTiF7Hr0R6zLRw6V/MKgaq8HvG54hJqqDE/BKI9c\nNxc6ucS+8xQSq3PZm+KRjZqOs8nVyG6lW+P59OMBY6nFDyLD4Ym7CkZ2uPSyDphz\n35fQeSTLSzjD0iebix4WV7w7zV1UbDDy0qKL64hKx0gTa7F1Kf4L3WUfCnVQppcq\n7w9WxaM22bxPin+XYE2B3j64QX1JC7ybI9O1DvpM9CSVXhpiIbJ5m8BsPoM4GAdi\nHwIDAQAB\n-----END PUBLIC KEY-----"),
			strings.NewReader("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyz9KTGQpQIoACbF2LKy4\n3QhENi/7/vtQl/iANUhxHo315b+0+u+imrhqG0iADQsMN/x7QequCLmFWXtowT35\nhU/YE5jnf3jO0AkH9/nWk5OBLTiF7Hr0R6zLRw6V/MKgaq8HvG54hJqqDE/BKI9c\nNxc6ucS+8xQSq3PZm+KRjZqOs8nVyG6lW+P59OMBY6nFDyLD4Ym7CkZ2uPSyDphz\n35fQeSTLSzjD0iebix4WV7w7zV1UbDDy0qKL64hKx0gTa7F1Kf4L3WUfCnVQppcq\n7w9WxaM22bxPin+XYE2B3j64QX1JC7ybI9O1DvpM9CSVXhpiIbJ5m8BsPoM4GAdi\nHwIDAQAB\n-----END PUBLIC KEY-----"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, e := sdk.BuildDidDocument(tt.pubkey, tt.sigKeyReader, tt.verKeyReader)

			if tt.wantErr {
				require.Error(t, e)
				require.Equal(t, DidDocument{}, d)
				return
			}

			require.NotEqual(t, DidDocument{}, d)
			require.NoError(t, e)
		})
	}
}
