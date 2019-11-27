package types

import (
	"encoding/hex"
	"github.com/pokt-network/posmint/crypto"
	"golang.org/x/crypto/sha3"
)

func SignatureVerification(publicKey, msgHex, sigHex string) error {
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return err
	}
	if len(sig) != crypto.SignatureSize {
		return InvalidSignatureSizeError
	}
	pk, err := crypto.NewPublicKey(publicKey)
	if err != nil {
		return err
	}
	msg, err := hex.DecodeString(msgHex)
	if err != nil {
		return err
	}
	if ok := pk.VerifySignature(msg, sig); !ok {
		return InvalidSignatureError
	}
	return nil
}

func PubKeyVerification(pk string) error {
	pkBz, err := hex.DecodeString(pk)
	if err != nil {
		return err
	}
	if len(pkBz) != crypto.PubKeySize {
		return PubKeySizeError
	}
	return nil
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
