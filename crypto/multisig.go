package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
)

var _ PublicKeyMultiSig = PublicKeyMultiSignature{}
var _ crypto.PubKey = PublicKeyMultiSignature{}
var _ PublicKey = PublicKeyMultiSignature{}

type PublicKeyMultiSignature struct {
	PublicKeys []PublicKey `json:"keys"`
}

func (pms PublicKeyMultiSignature) NewMultiKey(keys ...PublicKey) (PublicKeyMultiSig, error) {
	if keys == nil || len(keys) < 2 {
		return nil, errors.New("must have at least two public keys")
	}
	pms.PublicKeys = keys
	return pms, nil
}

func (pms PublicKeyMultiSignature) VerifyBytes(msg []byte, multiSignature []byte) bool {
	var multiSig MultiSig
	err := cdc.UnmarshalBinaryBare(multiSignature, &multiSig)
	if err != nil {
		return false
	}
	numOfSigs := multiSig.NumOfSigs()
	// ensure number of signatures == number of public keys
	if numOfSigs != len(pms.PublicKeys) {
		return false
	}
	for i := 0; i < numOfSigs; i++ {
		signature, found := multiSig.GetSignatureByIndex(i)
		if !found {
			return false
		}
		if !pms.PublicKeys[i].VerifyBytes(msg, signature) {
			return false
		}
	}
	return true
}

func (pms PublicKeyMultiSignature) Address() crypto.Address {
	return crypto.AddressHash(pms.Bytes())
}

func (pms PublicKeyMultiSignature) String() string {
	return hex.EncodeToString(pms.Bytes())
}

func (pms PublicKeyMultiSignature) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(pms)
}

func (pms PublicKeyMultiSignature) Keys() []PublicKey {
	return pms.PublicKeys
}

func (pms PublicKeyMultiSignature) Equals(other crypto.PubKey) bool {
	otherKey, sameType := other.(PublicKeyMultiSignature)
	if !sameType {
		return false
	}
	if len(pms.PublicKeys) != len(otherKey.PublicKeys) {
		return false
	}
	for i := 0; i < len(pms.PublicKeys); i++ {
		if !pms.PublicKeys[i].Equals(otherKey.PublicKeys[i]) {
			return false
		}
	}
	return true
}

func (pms PublicKeyMultiSignature) NewPublicKey(res []byte) (PublicKey, error) {
	err := cdc.UnmarshalBinaryBare(res, &pms)
	return pms, err
}

func (pms PublicKeyMultiSignature) PubKey() crypto.PubKey {
	return nil
}

func (pms PublicKeyMultiSignature) RawBytes() []byte {
	return pms.Bytes()
}

func (pms PublicKeyMultiSignature) RawString() string {
	return pms.String()
}

func (pms PublicKeyMultiSignature) PubKeyToPublicKey(crypto.PubKey) PublicKey {
	return nil
}

func (pms PublicKeyMultiSignature) Size() int {
	if len(pms.PublicKeys) != 0 {
		return pms.PublicKeys[0].Size()
	}
	return 0
}

type MultiSignature struct {
	Sigs [][]byte `json:"signatures"`
}

var _ MultiSig = MultiSignature{}

func (ms MultiSignature) AddSignature(sig []byte, key PublicKey, keys []PublicKey) (MultiSig, error) {
	index := getIndex(key, keys)
	if index == -1 {
		return nil, fmt.Errorf("provided key %s doesn't exist in the list of public keys", key.RawString())
	}
	return ms.AddSignatureByIndex(sig, index), nil
}

func (ms MultiSignature) AddSignatureByIndex(sig []byte, index int) MultiSig {
	// Signature already exists, just replace the value there
	sigsLen := len(ms.Sigs)
	if sigsLen-1 >= index {
		ms.Sigs[index] = sig
		return ms
	}
	// else add it to the list at the specific index
	for i := sigsLen; i < index-1; i++ {
		ms.Sigs = append(ms.Sigs, []byte{0})
	}
	ms.Sigs = append(ms.Sigs, sig)
	return ms
}

func (ms MultiSignature) NewMultiSignature() MultiSig {
	return MultiSignature{Sigs: make([][]byte, 0, 2)}
}

func (ms MultiSignature) Marshal() []byte {
	return cdc.MustMarshalBinaryBare(ms)
}

func (ms MultiSignature) Unmarshal(sig []byte) MultiSig {
	cdc.MustUnmarshalBinaryBare(sig, &ms)
	return ms
}

func (ms MultiSignature) String() string {
	return hex.EncodeToString(ms.Marshal())
}

func (ms MultiSignature) NumOfSigs() int {
	return len(ms.Sigs)
}

func (ms MultiSignature) Signatures() [][]byte {
	return ms.Sigs
}

func (ms MultiSignature) GetSignatureByIndex(i int) (sig []byte, found bool) {
	if len(ms.Sigs) < i {
		return nil, false
	}
	sig = ms.Sigs[i]
	if sig == nil {
		return sig, false
	}
	return sig, true
}

func (ms MultiSignature) GetSignatureByKey(pubKey PublicKey, keys []PublicKey) (sig []byte, found bool) {
	i := getIndex(pubKey, keys)
	return ms.GetSignatureByIndex(i)
}

func getIndex(pk PublicKey, keys []PublicKey) int {
	for i := 0; i < len(keys); i++ {
		if pk.Equals(keys[i]) {
			return i
		}
	}
	return -1
}
