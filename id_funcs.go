package commercio

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"crypto"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/hdkeychain"
	id "github.com/commercionetwork/commercionetwork/x/id/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fastjson"
)

// BuildDidDocument creates a DidDocument for the account associated to sdk, given its publick key, signature and
// verification RSA public keys.
func (sdk *SDK) BuildDidDocument(pubKeyString string, signatureKey, verificationKey io.Reader) (DidDocument, error) {
	e := func(w error, ext error) (DidDocument, error) {
		return DidDocument{}, fmt.Errorf("%w, %s", w, ext.Error())
	}

	uAddr := sdk.wallet.Address

	_, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, pubKeyString)
	if err != nil {
		return e(ErrInvalidPublicKey, err)
	}

	wacc, err := types.AccAddressFromBech32(uAddr)
	if err != nil {
		return e(ErrInvalidAddress, err)
	}

	sk, _, err := readKey(signatureKey, typePublicKey)
	if err != nil {
		return e(ErrInvalidSignatureKey, err)
	}

	vk, _, err := readKey(verificationKey, typePublicKey)
	if err != nil {
		return e(ErrInvalidVerificationKey, err)
	}

	didDocument := DidDocument{
		Context: "https://www.w3.org/ns/did/v1",
		ID:      wacc,
		PubKeys: id.PubKeys{
			id.PubKey{
				ID:           uAddr + "#keys-1",
				Type:         "RsaVerificationKey2018",
				Controller:   wacc,
				PublicKeyPem: vk,
			},
			id.PubKey{
				ID:           uAddr + "#keys-2",
				Type:         "RsaSignatureKey2018",
				Controller:   wacc,
				PublicKeyPem: sk,
			},
		},
	}

	oProof := Proof{
		Type:               "EcdsaSecp256k1VerificationKey2019",
		Created:            time.Now(),
		ProofPurpose:       "authentication",
		Controller:         uAddr,
		VerificationMethod: pubKeyString,
	}

	data, err := json.Marshal(didDocument)
	if err != nil {
		return e(ErrProofCreation, err)
	}

	wex, err := sdk.wallet.ExportWithPrivateKey()
	if err != nil {
		return e(ErrProofCreation, err)
	}

	kc, err := hdkeychain.NewKeyFromString(fastjson.GetString([]byte(wex), "private_key"))
	if err != nil {
		return e(ErrProofCreation, err)
	}

	ec, err := kc.ECPrivKey()
	if err != nil {
		return e(ErrProofCreation, err)
	}

	sum := sha256.Sum256(data[:])
	signature, err := ec.Sign(sum[:])
	if err != nil {
		return e(ErrProofCreation, err)
	}

	oProof.SignatureValue = base64.StdEncoding.EncodeToString(serializeSig(signature))

	ddProof := id.Proof(oProof)
	didDocument.Proof = &ddProof

	return didDocument, nil
}

// PowerUpParams are parameters used by BuildPowerupRequests during its lifecycle.
type PowerUpParams struct {
	// PubKey is the bech32-encoded public key of the account which sends the Power-up request.
	PubKey string

	// TumblerKey is the tumbler RSA PKIX public key.
	TumblerKey io.Reader

	// SignatureKey is the user RSA PKCS8 private key, used to build the message proof.
	//
	// It must be the private part of the keypair created and associated with the account which sends
	// the transaction.
	SignatureKey io.Reader

	// Amount represents the amount the user signing the transaction wants to send to its Pairwise address.
	Amount uint64

	// PairwiseAddress is the address to which the user signing the transaction wants to send Amount tokens.
	PairwiseAddress types.AccAddress
}

// validate checks that p is valid and can be used.
func (p PowerUpParams) validate() error {
	if p.PubKey == "" {
		return errors.New("missing pubkey")
	}

	if p.TumblerKey == nil {
		return errors.New("tumbler key reader must not be nil")
	}

	if p.SignatureKey == nil {
		return errors.New("verification key reader must not be nil")
	}

	if p.Amount == 0 {
		return errors.New("amount cannot be zero")
	}

	if p.PairwiseAddress == nil {
		return errors.New("pairwise address must not be nil")
	}

	return nil
}

// requestPowerupProof is the proof payload that will be used to calculate the Powerup Request proof.
type requestPowerupProof struct {
	SenderDid   types.AccAddress `json:"sender_did"`
	PairwiseDid types.AccAddress `json:"pairwise_did"`
	Timestamp   int64            `json:"timestamp"`
	Signature   string           `json:"signature"`
}

// BuildPowerupRequest creates a MsgRequestDidPowerUp based on params.
func (sdk *SDK) BuildPowerupRequest(params PowerUpParams) (MsgRequestDidPowerUp, error) {
	e := func(w error, ext error) (MsgRequestDidPowerUp, error) {
		return MsgRequestDidPowerUp{}, fmt.Errorf("%w, %s", w, ext.Error())
	}

	if err := params.validate(); err != nil {
		return e(ErrInvalidPowerupParams, err)
	}

	uAddr := sdk.wallet.Address

	_, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, params.PubKey)
	if err != nil {
		return e(ErrInvalidPublicKey, err)
	}

	wacc, err := types.AccAddressFromBech32(uAddr)
	if err != nil {
		return e(ErrInvalidAddress, err)
	}

	coinsAmount, err := Amount(params.Amount)
	if err != nil {
		return e(ErrInvalidAmount, err)
	}

	request := MsgRequestDidPowerUp{
		Claimant: wacc,
		Amount:   coinsAmount,
		ID:       uuid.NewV4().String(),
	}

	proof := requestPowerupProof{
		SenderDid:   wacc,
		PairwiseDid: params.PairwiseAddress,
		Timestamp:   time.Now().Unix(),
	}

	sigPayload := proof.SenderDid.String() + proof.PairwiseDid.String() + strconv.FormatInt(proof.Timestamp, 10)
	payloadHash := sha256.Sum256([]byte(sigPayload))

	_, rawKey, err := readKey(params.SignatureKey, typePrivateKey)
	if err != nil {
		return e(ErrInvalidSignatureKey, err)
	}

	privKey, ok := rawKey.(*rsa.PrivateKey)
	if !ok {
		// something's seriously wrong here, since the readKey function should already return stuff that *can* be
		// casted to *rsa.PrivateKey, panic!
		panic(fmt.Errorf("readKey parsed the private SignatureKey, but somehow it's not a *rsa.PrivateKey"))
	}

	// sign the proof hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, payloadHash[:])
	if err != nil {
		log.Fatal(err)
	}

	// encode in base64
	proof.Signature = base64.StdEncoding.EncodeToString(signature)

	// encode to json
	cdc := codec.New()
	proofJSON, _ := cdc.MarshalJSON(proof)

	/*
		proof now contains the blob we will encrypt with the tumbler public key
	*/
	_, tumblerRawKey, err := readKey(params.TumblerKey, typePublicKey)
	if err != nil {
		return e(ErrInvalidTumblerKey, err)
	}

	tumblerKey, ok := tumblerRawKey.(*rsa.PublicKey)
	if !ok {
		// something's seriously wrong here, since the readKey function should already return stuff that *can* be
		// casted to *rsa.PublicKey, panic!
		panic(fmt.Errorf("readKey parsed the public tumbler key, but somehow it's not a *rsa.PublicKey"))
	}

	// AES-GCM encryption of proofJSON
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return e(ErrNotEnoughEntropy, fmt.Errorf("cannot generate AES nonce, %w", err))
	}

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return e(ErrNotEnoughEntropy, fmt.Errorf("cannot generate AES encryption key, %w", err))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return e(ErrEncryptionFailure, err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return e(ErrEncryptionFailure, err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, proofJSON, nil)

	finalc := bytes.Buffer{}
	finalc.Write(nonce)
	finalc.Write(ciphertext)

	// convert it in base64
	epb64 := base64.StdEncoding.EncodeToString(finalc.Bytes())

	request.Proof = epb64

	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, tumblerKey, key)
	if err != nil {
		return e(ErrEncryptionFailure, err)
	}

	// convert it in base64
	keyb64 := base64.StdEncoding.EncodeToString(encryptedKey)

	request.ProofKey = keyb64

	return request, nil
}

// keyType is used to define whether readKey should validate a private or public RSA key.
type keyType int

const (
	typePrivateKey keyType = iota
	typePublicKey
)

// readKey reads an RSA PKIX public key or a RSA PKCS8 private key depending on kt, validates it and returns its
// PEM representation along with the associated interface type.
// Will return error if the key is invalid.
func readKey(r io.Reader, kt keyType) (string, interface{}, error) {
	key, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil, err
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return "", nil, errors.New("no valid PEM data found")
	}

	var pk interface{}

	switch kt {
	case typePublicKey:
		pk, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return "", nil, fmt.Errorf("invalid public key: %w", err)
		}
	case typePrivateKey:
		pk, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", nil, fmt.Errorf("invalid private key key: %w", err)
		}
	}

	return string(key), pk, nil

}

// serializeSig serializes a btcec.Signature in the (R || S) format.
func serializeSig(sig *btcec.Signature) []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}
