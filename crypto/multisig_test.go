package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func MultiSigSetup(t *testing.T) (pubKey PublicKeyMultiSig, privateKeys []PrivateKey) {
	pubKey = PublicKeyMultiSignature{}
	privateKeys = append(privateKeys, getRandomPrivateKey(t), getRandomPrivateKey(t))
	pubKey, err := pubKey.NewMultiKey(privateKeys[0].PublicKey(), privateKeys[1].PublicKey())
	if err != nil {
		t.Fatal(err)
	}
	if pubKey == nil {
		t.Fatal("nil pk")
	}
	return
}

func TestMultiSignature_AddSignatureVerifyBytes(t *testing.T) {
	msg := []byte("foo")
	pubKey, privKeys := MultiSigSetup(t)
	ms := (&MultiSignature{}).NewMultiSignature()
	sig1, err := privKeys[0].Sign(msg)
	if err != nil {
		t.Fatal(err)
	}
	sig2, err := privKeys[1].Sign(msg)
	if err != nil {
		t.Fatal(err)
	}
	ms, err = ms.AddSignature(sig1, pubKey.Keys()[0], pubKey.Keys())
	assert.Nil(t, err)
	ms = ms.AddSignatureByIndex(sig1, 0)
	assert.Nil(t, err)
	ms, err = ms.AddSignature(sig2, pubKey.Keys()[1], pubKey.Keys())
	assert.Nil(t, err)
	assert.True(t, pubKey.VerifyBytes(msg, ms.Marshal()))
	// wrong signature
	msWS := ms
	msWS, err = msWS.AddSignature(sig2, pubKey.Keys()[0], pubKey.Keys())
	assert.Nil(t, err)
	assert.False(t, pubKey.VerifyBytes(msg, msWS.Marshal()))
	// wrong message
	assert.Nil(t, err)
	assert.False(t, pubKey.VerifyBytes([]byte("bar"), ms.Marshal()))
	// extra signature
	msES := ms
	msES = msES.AddSignatureByIndex(sig1, 3)
	assert.Nil(t, err)
	assert.False(t, pubKey.VerifyBytes(msg, msES.Marshal()))
	// empty signatures
	msNS := MultiSignature{}
	assert.Nil(t, err)
	assert.False(t, pubKey.VerifyBytes(msg, msNS.Marshal()))
}
