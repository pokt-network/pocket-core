package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
	"math/big"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

func SHA3FromBytes(b []byte) []byte {
	hasher := sha3.New256()
	hasher.Write(b)
	return hasher.Sum(nil)
}

func SHA3FromString(s string) []byte {
	hasher := sha3.New256()
	hasher.Write([]byte(s))
	return hasher.Sum(nil)
}

// A wrapper for elliptic curve functions

// "Sign" is cyrptographic DSA function that signs a [32]byte.
func Sign(messageHash, secretKey []byte) ([]byte, error) {
	// Sign creates a recoverable ECDSA signature.
	// The produced signature is in the 65-byte [R || S || V] format where V is 0 or 1.
	//
	// The caller is responsible for ensuring that msg cannot be chosen
	// directly by an attacker. It is usually preferable to use a cryptographic
	// Hash function on any input before handing it to this function.
	return secp256k1.Sign(messageHash, secretKey)
}

// FromECDSA exports a private key into a binary dump.
func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

// "FromECDSAPub' exports a public key into a binary dump.
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

// "GetPublicKeyFromSignature" confirms the public key from the signature.
func GetPublicKeyFromSignature(messageHash, signature []byte) ([]byte, error) {
	// RecoverPubkey returns the public key of the signer.
	// msg must be the 32-byte Hash of the message to be signed.
	// sig must be a 65-byte compact ECDSA signature containing the
	// recovery id as the last element.
	return secp256k1.RecoverPubkey(messageHash, signature)
}

// "VerifySignature" verifies the signature with the public key.
func VerifySignature(publicKey, messageHash, signature []byte) bool {
	// VerifySignature checks that the given pubkey created signature over message.
	// The signature should be in [R || S] format. !!! Remove V (last byte)
	return secp256k1.VerifySignature(publicKey, messageHash, signature)
}

// "UncompressedPublicKey" uncompresses the public key from 33 bytes.
func UncompressPublicKey(publicKey []byte) (x, y *big.Int) {
	// DecompressPubkey parses a public key in the 33-byte compressed format.
	// It returns non-nil coordinates if the public key is valid.
	return secp256k1.DecompressPubkey(publicKey)
}

// "CompressPublicKey" converts a public key to 33-byte format.
func CompressPublicKey(x, y *big.Int) []byte {
	// CompressPubkey encodes a public key to 33-byte compressed format.
	return secp256k1.CompressPubkey(x, y)
}

// "NewPrivateKey" creates a private key.
func NewPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
}

// "NewPublicKey" creates an EC public key.
func NewPublicKey(key *ecdsa.PrivateKey) ecdsa.PublicKey {
	return key.PublicKey
}

// UnmarshalPubkey converts bytes to an EC public key.
func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(secp256k1.S256(), pub)
	if x == nil {
		return nil, InvalidPublicKeyError
	}
	return &ecdsa.PublicKey{Curve: secp256k1.S256(), X: x, Y: y}, nil
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}

// UnmarshalPrivateKey creates a private key with the given D value.
func UnmarshalPrivateKey(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, true)
}

// ToECDSAUnsafe blindly converts a binary blob to a private key. It should almost
// never be used unless you are sure the input is valid and want to avoid hitting
// errors due to bad origin encoding (0 prefixes cut off).
func ToECDSAUnsafe(d []byte) *ecdsa.PrivateKey {
	priv, _ := toECDSA(d, false)
	return priv
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, InvalidHexStringError
	}
	return UnmarshalPrivateKey(b)
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, InvalidPrivateKeyError
	}
	return priv, nil
}
