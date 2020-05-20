package commercio

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMnemonic(t *testing.T) {
	tests := []struct {
		name    string
		ep      func(b []byte) (n int, err error)
		wantErr bool
	}{
		{
			"read ok",
			nil,
			false,
		},
		{
			"entropy read fails, return from entropyProvider call",
			func(b []byte) (n int, err error) {
				return 0, errors.New("error!")
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ep != nil {
				entropyProvider = tt.ep
				defer func() {
					entropyProvider = rand.Read
				}()
			}

			res, err := NewMnemonic()

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, res)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, res)
		})
	}
}

func TestNewRSAKeypair(t *testing.T) {
	tests := []struct {
		name      string
		rp        func(random io.Reader, bits int) (*rsa.PrivateKey, error)
		wantErr   bool
		wantPanic bool
	}{
		{
			"all ok",
			nil,
			false,
			false,
		},
		{
			"rsa gen func return error",
			func(random io.Reader, bits int) (*rsa.PrivateKey, error) {
				return nil, errors.New("error!")
			},
			true,
			false,
		},
		{
			"rsa gen func return no error but no bytes have been generated",
			func(random io.Reader, bits int) (*rsa.PrivateKey, error) {
				return &rsa.PrivateKey{}, nil
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.rp != nil {
				rsaGenKeyProvider = tt.rp
				defer func() {
					rsaGenKeyProvider = rsa.GenerateKey
				}()
			}

			var pk string
			var pubk string
			var err error

			if tt.wantPanic {
				require.Panics(t, func() {
					pk, pubk, err = NewRSAKeypair()
				})
				return
			} else {
				pk, pubk, err = NewRSAKeypair()
			}

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, pk)
				require.Empty(t, pubk)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, pk)
			require.NotEmpty(t, pubk)
		})
	}
}
