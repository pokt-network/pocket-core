package types

import (
	sha "crypto"
	"encoding/hex"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"golang.org/x/crypto/sha3"
)

var (
	HashLength = sha.SHA3_256.Size()
)

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

func HashVerification(hash string) sdk.Error {
	if len(hash)==0 {
		return NewEmptyHashError(ModuleName)
	}
	if len(hash)!= HashLength {
		return NewInvalidHashLengthError(ModuleName)
	}
}

// Converts []byte to SHA3-256 hashed []byte
func SHA3FromBytes(b []byte) []byte {
	hasher := sha3.New256()
	hasher.Write(b)
	return hasher.Sum(nil)
}

// Converts string to SHA3-256 hashed []byte
func SHA3FromString(s string) []byte {
	hasher := sha3.New256()
	hasher.Write([]byte(s))
	return hasher.Sum(nil)
}
