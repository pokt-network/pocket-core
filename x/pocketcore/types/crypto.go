package types

import (
	sha "crypto"
	"encoding/hex"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	_ "golang.org/x/crypto/sha3"
)

var (
	Hasher     = sha.SHA3_256
	HashLength = sha.SHA3_256.Size()
	AddrLength = tmhash.TruncatedSize
)

// verify the signature using strings
func SignatureVerification(publicKey, msgHex, sigHex string) sdk.Error {
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return NewSigDecodeError(ModuleName)
	}
	if len(sig) != crypto.SignatureSize {
		return NewInvalidSignatureSizeError(ModuleName)
	}
	pk, err := crypto.NewPublicKey(publicKey)
	if err != nil {
		return NewPubKeyDecodeError(ModuleName)
	}
	msg, err := hex.DecodeString(msgHex)
	if err != nil {
		return NewMsgDecodeError(ModuleName)
	}
	if ok := pk.VerifySignature(msg, sig); !ok {
		return NewInvalidSignatureError(ModuleName)
	}
	return nil
}

// verify the public key format
func PubKeyVerification(pk string) sdk.Error {
	pkBz, err := hex.DecodeString(pk)
	if err != nil {
		return NewPubKeyDecodeError(ModuleName)
	}
	if len(pkBz) != crypto.PubKeySize {
		return NewPubKeySizeError(ModuleName)
	}
	return nil
}

// verify the hash format
func HashVerification(hash string) sdk.Error {
	if len(hash) == 0 {
		return NewEmptyHashError(ModuleName)
	}
	if len(hash) != HashLength {
		return NewInvalidHashLengthError(ModuleName)
	}
	return nil
}

func AddressVerification(address string) sdk.Error {
	if len(address) == 0 {
		return NewEmptyAddressError(ModuleName)
	}
	if len(address) != AddrLength {
		return NewAddressInvalidLengthError(ModuleName)
	}
	return nil
}

// Converts []byte to SHA3-256 hashed []byte
func Hash(b []byte) []byte {
	hasher := Hasher.New()
	hasher.Write(b)
	return hasher.Sum(nil)
}
