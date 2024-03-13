package types

import (
	sha "crypto"
	"encoding/hex"
	"math/big"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/privval"
	_ "golang.org/x/crypto/sha3"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
)

var (
	Hasher                  = sha.SHA3_256
	HashLength              = sha.SHA3_256.Size()
	NetworkIdentifierLength = 4
	AddrLength              = tmhash.TruncatedSize
	globalPVKeyFile         = privval.FilePVKey{}
)

// "NetworkIdentifierVerification"- Verify the netID format (hex string)
func NetworkIdentifierVerification(hash string) sdk.Error {
	// decode string into bz
	h, err := hex.DecodeString(hash)
	if err != nil {
		return NewHexDecodeError(ModuleName, err)
	}
	hashLen := len(h)
	// ensure Length isn't 0
	if hashLen == 0 {
		return NewEmptyHashError(ModuleName)
	}
	// ensure Length
	if hashLen > NetworkIdentifierLength {
		return NewInvalidHashLengthError(ModuleName)
	}
	return nil
}

// SignatureVerification verifies that:
// - The privKey associated with publicKey generated sigHex
// - sigHex signed msgHex
func SignatureVerification(publicKeyHex, msgHex, sigHex string) sdk.Error {
	// Decode the signature from hex
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return NewSigDecodeError(ModuleName)
	}
	// Ensure Length is valid
	if len(sigBytes) != crypto.Ed25519SignatureSize {
		return NewInvalidSignatureSizeError(ModuleName)
	}
	// Decode public key from hex
	publicKey, err := crypto.NewPublicKey(publicKeyHex)
	if err != nil {
		return NewPubKeyDecodeError(ModuleName)
	}
	// Decode message from hex
	msgBytes, err := hex.DecodeString(msgHex)
	if err != nil {
		return NewMsgDecodeError(ModuleName)
	}
	// verify that publicKey signed msgBytes and generated sigBytes
	if ok := publicKey.VerifyBytes(msgBytes, sigBytes); !ok {
		return NewInvalidSignatureError(ModuleName)
	}
	return nil
}

// "PubKeyVerification" - Verifies the public key format (hex string)
func PubKeyVerification(pk string) sdk.Error {
	// decode the bz
	pkBz, err := hex.DecodeString(pk)
	if err != nil {
		return NewPubKeyDecodeError(ModuleName)
	}
	// ensure Length
	if len(pkBz) != crypto.Ed25519PubKeySize {
		return NewPubKeySizeError(ModuleName)
	}
	return nil
}

// "HashVerification" - Verifies the merkleHash format (hex string)
func HashVerification(hash string) sdk.Error {
	// decode the merkleHash
	h, err := hex.DecodeString(hash)
	if err != nil {
		return NewHexDecodeError(ModuleName, err)
	}
	hLen := len(h)
	// ensure Length isn't 0
	if hLen == 0 {
		return NewEmptyHashError(ModuleName)
	}
	// ensure Length
	if hLen != HashLength {
		return NewInvalidHashLengthError(ModuleName)
	}
	return nil
}

// "AddressVerification" - Verifies the address format (hex strign)
func AddressVerification(addr string) sdk.Error {
	// decode the address
	address, err := hex.DecodeString(addr)
	if err != nil {
		return NewHexDecodeError(ModuleName, err)
	}
	addrLen := len(address)
	// ensure Length isn't 0
	if addrLen == 0 {
		return NewEmptyAddressError(ModuleName)
	}
	// ensure Length
	if addrLen != AddrLength {
		return NewAddressInvalidLengthError(ModuleName)
	}
	return nil
}

// "ID"- Converts []byte to hashed []byte
func Hash(b []byte) []byte {
	hasher := Hasher.New()
	hasher.Write(b) //nolint:golint,errcheck
	return hasher.Sum(nil)
}

func PseudorandomSelection(max sdk.BigInt, hash []byte) (index sdk.BigInt) {
	// merkleHash for show and convert back to decimal
	intHash := sdk.NewIntFromBigInt(new(big.Int).SetBytes(hash[:8]))
	// mod the selection
	return intHash.Mod(max)
}
