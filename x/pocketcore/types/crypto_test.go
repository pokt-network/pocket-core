package types

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	ed255192 "golang.org/x/crypto/ed25519"
	"testing"
)

func TestSignatureVerification(t *testing.T) {
	testData := []byte("test")
	badVerifyData := hex.EncodeToString([]byte("bad"))
	privateKey := GetRandomPrivateKey()
	goodVerifyPubKey := privateKey.PublicKey().RawString()
	badVerifyPublicKey := getRandomPubKey().RawString()
	signature, err := privateKey.Sign(testData)
	if err != nil {
		t.Fatalf(err.Error())
	}
	tests := []struct {
		name         string
		signData     string
		verifyData   string
		verifyPubKey string
		signature    string
		hasError     bool
	}{
		{
			name:         "Bad verify data",
			signData:     hex.EncodeToString(testData),
			verifyData:   badVerifyData,
			verifyPubKey: goodVerifyPubKey,
			signature:    hex.EncodeToString(signature),
			hasError:     true,
		},
		{
			name:         "Bad verify publicKey",
			signData:     hex.EncodeToString(testData),
			verifyData:   hex.EncodeToString(testData),
			verifyPubKey: badVerifyPublicKey,
			signature:    hex.EncodeToString(signature),
			hasError:     true,
		},
		{
			name:         "Valid verify",
			signData:     hex.EncodeToString(testData),
			verifyData:   hex.EncodeToString(testData),
			verifyPubKey: goodVerifyPubKey,
			signature:    hex.EncodeToString(signature),
			hasError:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, SignatureVerification(tt.verifyPubKey, tt.verifyData, tt.signature) != nil, tt.hasError)
		})
	}
}

func TestPubKeyVerification(t *testing.T) {
	privateKeyBytes := [ed255192.PrivateKeySize]byte(GetRandomPrivateKey())
	privateKey := hex.EncodeToString(privateKeyBytes[:])
	pubKeyWrongECBytes := [secp256k1.PubKeySecp256k1Size]byte(secp256k1.GenPrivKey().PubKey().(secp256k1.PubKeySecp256k1))
	pkWrongEC := hex.EncodeToString(pubKeyWrongECBytes[:])
	pk := getRandomPubKey().RawString()
	tests := []struct {
		name     string
		key      string
		hasError bool
	}{
		{
			name:     "Empty",
			key:      "",
			hasError: true,
		},
		{
			name:     "Private key instead of public key",
			key:      privateKey,
			hasError: true,
		},
		{
			name:     "Wrong eliptic curve algorithm",
			key:      pkWrongEC,
			hasError: true,
		},
		{
			name:     "Correct Public Key",
			key:      pk,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, PubKeyVerification(tt.key) != nil, tt.hasError)
		})
	}
}

func TestAddressVerification(t *testing.T) {
	badAddress := hex.EncodeToString([]byte("asdflkjaflkjasdf"))
	addressBadEncoding := base64.StdEncoding.EncodeToString(getRandomValidatorAddress())
	address := getRandomValidatorAddress().String()
	tests := []struct {
		name     string
		addr     string
		hasError bool
	}{
		{
			name:     "Bad addr",
			addr:     badAddress,
			hasError: true,
		},
		{
			name:     "Bad addr encoding",
			addr:     addressBadEncoding,
			hasError: true,
		},
		{
			name:     "Good addr encoding",
			addr:     address,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, AddressVerification(tt.addr) != nil, tt.hasError)
		})
	}
}

func TestHash(t *testing.T) {
	testData := []byte("test")
	badHash := hex.EncodeToString([]byte("testasldfjalsdjfnvaklsdj"))
	hash := hex.EncodeToString(Hash(testData))
	badEncoding := base64.StdEncoding.EncodeToString(Hash(testData))
	tests := []struct {
		name     string
		hash     string
		hasError bool
	}{
		{
			name:     "Bad addr",
			hash:     badHash,
			hasError: true,
		},
		{
			name:     "Bad addr encoding",
			hash:     badEncoding,
			hasError: true,
		},
		{
			name:     "Good addr encoding",
			hash:     hash,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, HashVerification(tt.hash) != nil, tt.hasError)
		})
	}
}
