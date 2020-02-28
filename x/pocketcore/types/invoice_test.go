package types

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetAllInvoices(t *testing.T) {
	assert.NotNil(t, GetAllInvoices().M)
}

func TestAllInvoices_AddGetInvoice(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	header := SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	proof := Proof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	err = GetAllInvoices().AddToInvoice(header, proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.True(t, reflect.DeepEqual(GetAllInvoices().GetProof(header, 0), proof))
}

func TestAllInvoices_DeleteInvoice(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	header := SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	proof := Proof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	err = GetAllInvoices().AddToInvoice(header, proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.True(t, reflect.DeepEqual(GetAllInvoices().GetProof(header, 0), proof))
	GetAllInvoices().GetProof(header, 0)
	GetAllInvoices().DeleteInvoice(header)
	assert.Empty(t, GetAllInvoices().GetProof(header, 0))
}

func TestAllInvoices_GetProofs(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	header := SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	proof := Proof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	proof2 := Proof{
		Entropy:            1, // just for testing equality
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	err = GetAllInvoices().AddToInvoice(header, proof)
	err = GetAllInvoices().AddToInvoice(header, proof2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	proofs := GetAllInvoices().GetProofs(header)
	assert.NotNil(t, proofs)
	assert.Len(t, proofs, 2)
	assert.Equal(t, proofs[0], proof)
	assert.Equal(t, proofs[1], proof2)
}

func TestAllInvoices_GetTotalRelays(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	header := SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	header2 := SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 101,
	}
	proof := Proof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	proof2 := Proof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	err = GetAllInvoices().AddToInvoice(header, proof)
	err = GetAllInvoices().AddToInvoice(header, proof2)
	err = GetAllInvoices().AddToInvoice(header2, proof2) // different header so shouldn't be counted
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, GetAllInvoices().GetTotalRelays(header), int64(2))
}
