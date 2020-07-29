package crypto

import (
	ed255192 "crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strings"
)

type (
	Ed25519PublicKey  ed25519.PubKeyEd25519
	Ed25519PrivateKey ed25519.PrivKeyEd25519
)

var (
	_ PublicKey      = Ed25519PublicKey{}
	_ PrivateKey     = Ed25519PrivateKey{}
	_ crypto.PubKey  = Ed25519PublicKey{}
	_ crypto.PrivKey = Ed25519PrivateKey{}
)

const (
	Ed25519PrivKeySize   = ed255192.PrivateKeySize
	Ed25519PubKeySize    = ed25519.PubKeyEd25519Size
	Ed25519SignatureSize = ed25519.SignatureSize
)

func (Ed25519PublicKey) NewPublicKey(b []byte) (PublicKey, error) {
	var bz [Ed25519PubKeySize]byte
	copy(bz[:], b)
	pubkey := ed25519.PubKeyEd25519(bz)
	pk := Ed25519PublicKey(pubkey)
	return pk, nil
}

func (Ed25519PublicKey) PubKeyToPublicKey(key crypto.PubKey) PublicKey {
	return Ed25519PublicKey(key.(ed25519.PubKeyEd25519))
}

func (pub Ed25519PublicKey) PubKey() crypto.PubKey {
	return ed25519.PubKeyEd25519(pub)
}

func (pub Ed25519PublicKey) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(pub)
	if err != nil {
		panic(err)
	}
	return bz
}

func (pub Ed25519PublicKey) RawBytes() []byte {
	pkBytes := [Ed25519PubKeySize]byte(pub)
	return pkBytes[:]
}

func (pub Ed25519PublicKey) String() string {
	return hex.EncodeToString(pub.Bytes())
}

func (pub Ed25519PublicKey) RawString() string {
	return hex.EncodeToString(pub.RawBytes())
}

func (pub Ed25519PublicKey) Address() crypto.Address {
	return ed25519.PubKeyEd25519(pub).Address()
}

func (pub Ed25519PublicKey) VerifyBytes(msg []byte, sig []byte) bool {
	return ed25519.PubKeyEd25519(pub).VerifyBytes(msg, sig)
}

func (pub Ed25519PublicKey) Equals(other crypto.PubKey) bool {
	return ed25519.PubKeyEd25519(pub).Equals(ed25519.PubKeyEd25519(other.(Ed25519PublicKey)))
}

func (pub Ed25519PublicKey) Size() int {
	return Ed25519PubKeySize
}

func (pub Ed25519PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pub.RawString())
}

func (pub *Ed25519PublicKey) UnmarshalJSON(data []byte) error {
	hexstring := strings.Trim(string(data[:]), "\"")

	bytes, err := hex.DecodeString(hexstring)
	if err != nil {
		return err
	}
	pk, err := NewPublicKeyBz(bytes)
	if err != nil {
		return err
	}
	err = cdc.UnmarshalBinaryBare(pk.Bytes(), pub)

	if err != nil {
		return err
	}

	return nil
}

func (Ed25519PrivateKey) PrivateKeyFromBytes(b []byte) (PrivateKey, error) {
	var bz [Ed25519PrivKeySize]byte
	copy(bz[:], b)
	pri := ed25519.PrivKeyEd25519(bz)
	pk := Ed25519PrivateKey(pri)

	return pk, nil
}

func (priv Ed25519PrivateKey) RawBytes() []byte {
	pkBytes := [Ed25519PrivKeySize]byte(priv)
	return pkBytes[:]
}

func (priv Ed25519PrivateKey) RawString() string {
	return hex.EncodeToString(priv.RawBytes())
}

func (priv Ed25519PrivateKey) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(priv)
	if err != nil {
		panic(err)
	}
	return bz
}

func (priv Ed25519PrivateKey) String() string {
	return hex.EncodeToString(priv.Bytes())
}

func (priv Ed25519PrivateKey) PublicKey() PublicKey {
	return Ed25519PublicKey(ed25519.PrivKeyEd25519(priv).PubKey().(ed25519.PubKeyEd25519))
}

func (priv Ed25519PrivateKey) PubKey() crypto.PubKey {
	return ed25519.PrivKeyEd25519(priv).PubKey().(ed25519.PubKeyEd25519)
}

func (priv Ed25519PrivateKey) Equals(other crypto.PrivKey) bool {
	return ed25519.PrivKeyEd25519(priv).Equals(ed25519.PrivKeyEd25519(other.(Ed25519PrivateKey)))
}

func (priv Ed25519PrivateKey) Sign(msg []byte) ([]byte, error) {
	return ed25519.PrivKeyEd25519(priv).Sign(msg)
}

func (priv Ed25519PrivateKey) Size() int {
	return Ed25519PrivKeySize
}

func (priv Ed25519PrivateKey) PrivKey() crypto.PrivKey {
	return ed25519.PrivKeyEd25519(priv)
}

func (Ed25519PrivateKey) PrivKeyToPrivateKey(key crypto.PrivKey) PrivateKey {
	return Ed25519PrivateKey(key.(ed25519.PrivKeyEd25519))
}

func (Ed25519PrivateKey) GenPrivateKey() PrivateKey {
	return Ed25519PrivateKey(ed25519.GenPrivKey())
}
