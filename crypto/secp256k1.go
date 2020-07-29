package crypto

import (
	"encoding/hex"
	"encoding/json"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"strings"
)

type (
	Secp256k1PublicKey  secp256k1.PubKeySecp256k1
	Secp256k1PrivateKey secp256k1.PrivKeySecp256k1
)

var (
	_ PublicKey  = Secp256k1PublicKey{}
	_ PrivateKey = Secp256k1PrivateKey{}

	_ crypto.PubKey  = Secp256k1PublicKey{}
	_ crypto.PrivKey = Secp256k1PrivateKey{}
)

const (
	Secp256k1PrivateKeySize = 32
	Secp256k1PublicKeySize  = secp256k1.PubKeySecp256k1Size
)

func (Secp256k1PublicKey) NewPublicKey(b []byte) (PublicKey, error) {
	var bz [Secp256k1PublicKeySize]byte
	copy(bz[:], b)
	pubkey := secp256k1.PubKeySecp256k1(bz)
	pk := Secp256k1PublicKey(pubkey)
	return pk, nil
}

func (Secp256k1PublicKey) PubKeyToPublicKey(key crypto.PubKey) PublicKey {
	return Secp256k1PublicKey(key.(secp256k1.PubKeySecp256k1))
}

func (pub Secp256k1PublicKey) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(pub)
	if err != nil {
		panic(err)
	}
	return bz
}

func (pub Secp256k1PublicKey) PubKey() crypto.PubKey {
	return secp256k1.PubKeySecp256k1(pub)
}

func (pub Secp256k1PublicKey) RawBytes() []byte {
	pkBytes := [Secp256k1PublicKeySize]byte(pub)
	return pkBytes[:]
}

func (pub Secp256k1PublicKey) String() string {
	return hex.EncodeToString(pub.Bytes())
}

func (pub Secp256k1PublicKey) RawString() string {
	return hex.EncodeToString(pub.RawBytes())
}

func (pub Secp256k1PublicKey) Address() crypto.Address {
	return secp256k1.PubKeySecp256k1(pub).Address()
}

func (pub Secp256k1PublicKey) VerifyBytes(msg []byte, sig []byte) bool {
	return secp256k1.PubKeySecp256k1(pub).VerifyBytes(msg, sig)
}

func (pub Secp256k1PublicKey) Equals(other crypto.PubKey) bool {
	return secp256k1.PubKeySecp256k1(pub).Equals(secp256k1.PubKeySecp256k1(other.(Secp256k1PublicKey)))
}

func (pub Secp256k1PublicKey) Size() int {
	return Secp256k1PublicKeySize
}

func (pub Secp256k1PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pub.RawString())
}

func (pub *Secp256k1PublicKey) UnmarshalJSON(data []byte) error {
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
func (Secp256k1PrivateKey) PrivateKeyFromBytes(b []byte) (PrivateKey, error) {
	var bz [Secp256k1PrivateKeySize]byte
	copy(bz[:], b)
	pri := secp256k1.PrivKeySecp256k1(bz)
	pk := Secp256k1PrivateKey(pri)

	return pk, nil
}

func (priv Secp256k1PrivateKey) RawBytes() []byte {
	pkBytes := [Secp256k1PrivateKeySize]byte(priv)
	return pkBytes[:]
}

func (priv Secp256k1PrivateKey) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(priv)
	if err != nil {
		panic(err)
	}
	return bz
}

func (priv Secp256k1PrivateKey) String() string {
	return hex.EncodeToString(priv.Bytes())
}

func (priv Secp256k1PrivateKey) Equals(other crypto.PrivKey) bool {
	return secp256k1.PrivKeySecp256k1(priv).Equals(secp256k1.PrivKeySecp256k1(other.(Secp256k1PrivateKey)))
}

func (priv Secp256k1PrivateKey) RawString() string {
	return hex.EncodeToString(priv.RawBytes())
}

func (priv Secp256k1PrivateKey) PublicKey() PublicKey {
	return Secp256k1PublicKey(secp256k1.PrivKeySecp256k1(priv).PubKey().(secp256k1.PubKeySecp256k1))
}

func (priv Secp256k1PrivateKey) PubKey() crypto.PubKey {
	return secp256k1.PubKeySecp256k1(secp256k1.PrivKeySecp256k1(priv).PubKey().(secp256k1.PubKeySecp256k1))
}

func (priv Secp256k1PrivateKey) Sign(msg []byte) ([]byte, error) {
	return secp256k1.PrivKeySecp256k1(priv).Sign(msg)
}

func (priv Secp256k1PrivateKey) Size() int {
	return Secp256k1PrivateKeySize
}

func (priv Secp256k1PrivateKey) PrivKey() crypto.PrivKey {
	return secp256k1.PrivKeySecp256k1(priv)
}

func (Secp256k1PrivateKey) PrivKeyToPrivateKey(key crypto.PrivKey) PrivateKey {
	return Secp256k1PrivateKey(key.(secp256k1.PrivKeySecp256k1))
}

func (Secp256k1PrivateKey) GenPrivateKey() PrivateKey {
	return Secp256k1PrivateKey(secp256k1.GenPrivKey())
}
