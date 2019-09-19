package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// wrapper around go-secp256k1 public key
type PublicKey struct {
	ecdsaPubKey ecdsa.PublicKey
}

// wrapper around go-secp256k1 private key
type PrivateKey struct {
	ecdsaPrivKey *ecdsa.PrivateKey
}

// Private Key functions
// TODO: Create a function to import private key from []byte or hex dump.

// Generates a new private key
func NewPrivateKey() (PrivateKey, error) {
	ecdsaPrivKey, err := ecdsa.GenerateKey(S256(), rand.Reader)
	if err != nil {
		return PrivateKey{}, err
	}
	return PrivateKey{ecdsaPrivKey}, nil
}

// Returns the private key public key
func (privKey *PrivateKey) GetPublicKey() PublicKey {
	return PublicKey{privKey.ecdsaPrivKey.PublicKey}
}

// Dumps the PrivateKey on a []bytes
func (privKey *PrivateKey) Bytes() []byte {
	if privKey == nil {
		return nil
	}
	return paddedBigBytes(privKey.ecdsaPrivKey.D, privKey.ecdsaPrivKey.Params().BitSize/8)
}

// Signs a message, the message has to be [32]byte
func (privKey *PrivateKey) Sign(message []byte) ([]byte, error) {
	if privKey == nil {
		return nil, InvalidPrivateKeyError
	}
	return secp256k1.Sign(message, privKey.Bytes())
}

// Public Key functions

// Dumps the PublicKey on a []bytes
func (pubKey *PublicKey) Bytes() []byte {
	if pubKey == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pubKey.ecdsaPubKey.X, pubKey.ecdsaPubKey.Y)
}

// Creates a PublicKey instance from bytes
func publicKeyFromBytes(bytes []byte) PublicKey {
	s256curve := S256()
	x, y := s256curve.ScalarBaseMult(bytes)
	//s256curve.ScalarMult
	ecdsaPublicKey := ecdsa.PublicKey{Curve: s256curve, X: x, Y: y}
	return PublicKey{ecdsaPublicKey}
}

// KeyPair functions

// Creates a new key pair
func NewKeypair() (PrivateKey, PublicKey, error) {
	var privKey PrivateKey
	var pubKey PublicKey
	privKey, err := NewPrivateKey()
	if err != nil {
		return privKey, pubKey, err
	}
	pubKey = privKey.GetPublicKey()

	return privKey, pubKey, nil
}

// Signature Functions

// "GetPublicKeyBytesFromSignature" confirms the public key from the signature.
func GetPublicKeyBytesFromSignature(messageHash, signature []byte) ([]byte, error) {
	// RecoverPubkey returns the public key of the signer.
	// msg must be the 32-byte Hash of the message to be signed.
	// sig must be a 65-byte compact ECDSA signature containing the
	// recovery id as the last element.
	return secp256k1.RecoverPubkey(messageHash, signature)
}

// "VerifySignature" verifies the signature with the public key.
func VerifySignature(pubKey PublicKey, messageHash []byte, signature []byte) bool {
	// VerifySignature checks that the given pubkey created signature over message.
	// The signature should be in [R || S] format. !!! Remove V (last byte)
	return secp256k1.VerifySignature(pubKey.Bytes(), messageHash, signature[:len(signature)-1])
}

// "VerifySignature" verifies the signature with the public key bytes.
func VerifySignatureWithPubKeyBytes(pubKeyBytes []byte, messageHash []byte, signature []byte) bool {
	// VerifySignatureWithPubKeyBytes checks that the given pubkey bytes created signature over message.
	// The signature should be in [R || S] format. !!! Remove V (last byte)
	return secp256k1.VerifySignature(pubKeyBytes, messageHash, signature[:len(signature)-1])
}
