package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/types"
	"reflect"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type PublicKey interface {
	PubKey() crypto.PubKey
	Bytes() []byte
	RawBytes() []byte
	String() string
	RawString() string
	Address() crypto.Address
	Equals(other crypto.PubKey) bool
	VerifyBytes(msg []byte, sig []byte) bool
	PubKeyToPublicKey(crypto.PubKey) PublicKey
	Size() int
}

type PrivateKey interface {
	Bytes() []byte
	RawBytes() []byte
	String() string
	RawString() string
	PrivKey() crypto.PrivKey
	PubKey() crypto.PubKey
	Equals(other crypto.PrivKey) bool
	PublicKey() PublicKey
	Sign(msg []byte) ([]byte, error)
	PrivKeyToPrivateKey(crypto.PrivKey) PrivateKey
	GenPrivateKey() PrivateKey
	Size() int
}

type PublicKeyMultiSig interface {
	Address() crypto.Address
	String() string
	Bytes() []byte
	Equals(other crypto.PubKey) bool
	VerifyBytes(msg []byte, multiSignature []byte) bool
	PubKey() crypto.PubKey
	RawBytes() []byte
	RawString() string
	PubKeyToPublicKey(crypto.PubKey) PublicKey
	Size() int
	// new methods
	NewMultiKey(keys ...PublicKey) (PublicKeyMultiSig, error)
	Keys() []PublicKey
}

type MultiSig interface {
	AddSignature(sig []byte, key PublicKey, keys []PublicKey) (MultiSig, error)
	AddSignatureByIndex(sig []byte, index int) MultiSig
	Marshal() []byte
	Unmarshal([]byte) MultiSig
	NewMultiSignature() MultiSig
	String() string
	NumOfSigs() int
	Signatures() [][]byte
	GetSignatureByIndex(i int) (sig []byte, found bool)
}

func NewPublicKey(hexString string) (PublicKey, error) {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return NewPublicKeyBz(b)
}

func NewPublicKeyBz(b []byte) (PublicKey, error) {
	x := len(b)
	if x == Ed25519PubKeySize {
		return Ed25519PublicKey{}.NewPublicKey(b)
	} else if x == Secp256k1PublicKeySize {
		return Secp256k1PublicKey{}.NewPublicKey(b)
	} else if pk, err := PublicKeyMultiSignature.NewPublicKey(PublicKeyMultiSignature{}, b); err == nil {
		return pk, err
	} else {
		return nil, fmt.Errorf("unsupported public key type, length of: %d", x)
	}
}

func PubKeyToPublicKey(key crypto.PubKey) (PublicKey, error) {
	k := key
	switch k.(type) {
	case secp256k1.PubKeySecp256k1:
		return Secp256k1PublicKey{}.PubKeyToPublicKey(key), nil
	case ed25519.PubKeyEd25519:
		return Ed25519PublicKey(key.(ed25519.PubKeyEd25519)), nil
	case Ed25519PublicKey:
		return key.(Ed25519PublicKey), nil
	case Secp256k1PublicKey:
		return key.(Secp256k1PublicKey), nil
	default:
		return nil, errors.New("error converting pubkey to public key -> unsupported public key type")
	}
}

func NewPrivateKeyBz(b []byte) (PrivateKey, error) {
	switch len(b) {
	case Ed25519PrivKeySize:
		return Ed25519PrivateKey{}.PrivateKeyFromBytes(b)
	case Secp256k1PrivateKeySize:
		return Secp256k1PrivateKey{}.PrivateKeyFromBytes(b)
	default:
		return nil, errors.New("unsupported private key type")
	}
}

func PrivKeyToPrivateKey(key crypto.PrivKey) (PrivateKey, error) {
	k := key
	switch k.(type) {
	case secp256k1.PrivKeySecp256k1:
		return Secp256k1PrivateKey{}.PrivKeyToPrivateKey(key), nil
	case ed25519.PrivKeyEd25519:
		return Ed25519PrivateKey{}.PrivKeyToPrivateKey(key), nil
	case Secp256k1PrivateKey:
		return key.(Secp256k1PrivateKey), nil
	case Ed25519PrivateKey:
		return key.(Ed25519PrivateKey), nil
	default:
		return nil, errors.New("error converting privkey to private key -> unsupported private key type")
	}
}

func PrivKeyFromBytes(privKeyBytes []byte) (privKey PrivateKey, err error) {
	err = cdc.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

func PubKeyFromBytes(pubKeyBytes []byte) (pubKey PublicKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}

func GenerateSecp256k1PrivKey() PrivateKey {
	return Secp256k1PrivateKey{}.GenPrivateKey()
}

func GenerateEd25519PrivKey() PrivateKey {
	return Ed25519PrivateKey{}.GenPrivateKey()
}

// should prevent unknown keys from being in consensus
func CheckConsensusPubKey(pubKey crypto.PubKey) (abci.PubKey, error) {
	switch pk := pubKey.(type) {
	case ed25519.PubKeyEd25519:
		return abci.PubKey{
			Type: types.ABCIPubKeyTypeEd25519,
			Data: pk[:],
		}, nil
	case secp256k1.PubKeySecp256k1:
		return abci.PubKey{
			Type: types.ABCIPubKeyTypeSecp256k1,
			Data: pk[:],
		}, nil
	default:
		return abci.PubKey{}, fmt.Errorf("unknown pubkey type: %v %v", pubKey, reflect.TypeOf(pubKey))
	}
}
