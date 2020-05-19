package commercio

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/cosmos/go-bip39"
)

// NewMnemonic return a cryptographically-secure wallet mnemonic.
func NewMnemonic() (string, error) {
	e := func(w error, ext error) (string, error) {
		return "", fmt.Errorf("%w, %s", w, ext.Error())
	}
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return e(ErrNotEnoughEntropy, err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return e(ErrNotEnoughEntropy, err)
	}

	return mnemonic, nil
}

// NewRSAKeypair returns a new PEM-encoded RSA 2048-bit keypair.
func NewRSAKeypair() (string, string, error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	pubk := &pk.PublicKey

	mpk, err := x509.MarshalPKCS8PrivateKey(pk)
	if err != nil {
		return "", "", err
	}

	mpubk, err := x509.MarshalPKIXPublicKey(pubk)
	if err != nil {
		return "", "", err
	}

	pubkStr := bytes.Buffer{}
	pkStr := bytes.Buffer{}

	pubkBlock := &pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   mpubk,
	}

	pkBlock := &pem.Block{
		Type:    "PRIVATE KEY",
		Headers: nil,
		Bytes:   mpk,
	}

	if err := pem.Encode(&pubkStr, pubkBlock); err != nil {
		return "", "", err
	}

	if err := pem.Encode(&pkStr, pkBlock); err != nil {
		return "", "", err
	}

	return pkStr.String(), pubkStr.String(), nil
}
