package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"math"
	"reflect"
	"testing"
)

func TestInvoice_GenerateMerkleRoot(t *testing.T) {
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	clientPrivateKey := getRandomPrivateKey()
	clientPublicKey := clientPrivateKey.PublicKey().RawString()
	nodePubKey := getRandomPubKey()
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
	validAAT := AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPubKey,
		ClientPublicKey:      clientPublicKey,
		ApplicationSignature: "",
	}
	appSig, er := appPrivateKey.Sign(validAAT.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validAAT.ApplicationSignature = hex.EncodeToString(appSig)
	i := Invoice{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		TotalRelays: 5,
		Proofs: []RelayProof{
			{
				Entropy:            3238283,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            34939492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            12383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            96384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            96384812,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
		},
	}
	root := i.GenerateMerkleRoot()
	assert.NotNil(t, root.Hash)
	assert.NotEmpty(t, root.Hash)
	assert.Nil(t, HashVerification(hex.EncodeToString(root.Hash)))
	assert.NotNil(t, root.Sum)
	assert.NotZero(t, root.Sum)
}

func TestInvoice_GenerateMerkleProof(t *testing.T) {
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	clientPrivateKey := getRandomPrivateKey()
	clientPublicKey := clientPrivateKey.PublicKey().RawString()
	nodePubKey := getRandomPubKey()
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
	validAAT := AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPubKey,
		ClientPublicKey:      clientPublicKey,
		ApplicationSignature: "",
	}
	appSig, er := appPrivateKey.Sign(validAAT.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validAAT.ApplicationSignature = hex.EncodeToString(appSig)
	i := Invoice{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		TotalRelays: 5,
		Proofs: []RelayProof{
			{
				Entropy:            3238283,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            34939492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            12383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            96384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            96384812,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
		},
	}
	index := 4
	proofs, cousinIndex := i.GenerateMerkleProof(index)
	assert.NotPanics(t, func() {
		if reflect.DeepEqual(proofs[0], proofs[1]) {
			t.Fatalf("Equal MerkleProofs")
		}
	})
	assert.Len(t, proofs[0].HashSums, 3)
	assert.Len(t, proofs[1].HashSums, 3)
	assert.True(t, math.Abs(float64(cousinIndex-index)) < 3)
}

func TestInvoice_VerifyMerkleProof(t *testing.T) {
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	clientPrivateKey := getRandomPrivateKey()
	clientPublicKey := clientPrivateKey.PublicKey().RawString()
	nodePubKey := getRandomPubKey()
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
	validAAT := AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPubKey,
		ClientPublicKey:      clientPublicKey,
		ApplicationSignature: "",
	}
	appSig, er := appPrivateKey.Sign(validAAT.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validAAT.ApplicationSignature = hex.EncodeToString(appSig)
	i := Invoice{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		TotalRelays: 5,
		Proofs: []RelayProof{
			{
				Entropy:            83,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            3492332332249492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            121212123232323383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            23121223232396384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            963223233238481322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
		},
	}
	i2 := Invoice{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		TotalRelays: 9,
		Proofs: []RelayProof{
			{
				Entropy:            82398289423,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            34932332249492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            1212121232383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            23192932384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            2993223481322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            993223423981322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            90333981322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            2398123322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			{
				Entropy:            99322342381322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
		},
	}
	index := 4
	root := i.GenerateMerkleRoot()
	proofs, cousinIndex := i.GenerateMerkleProof(index)
	assert.True(t, proofs.Validate(root, i.Proofs[index], i.Proofs[cousinIndex], int64(len(i.Proofs))))
	index2 := 0
	root2 := i2.GenerateMerkleRoot()
	proofs2, cousinIndex2 := i2.GenerateMerkleProof(index2)
	assert.True(t, proofs2.Validate(root2, i2.Proofs[index2], i2.Proofs[cousinIndex2], int64(len(i2.Proofs))))
	// wrong root
	assert.False(t, proofs.Validate(root2, i.Proofs[index], i.Proofs[cousinIndex], int64(len(i.Proofs))))
	// wrong cousin provided
	assert.False(t, proofs.Validate(root, i.Proofs[index], i.Proofs[cousinIndex2], int64(len(i.Proofs))))
	// wrong leaf provided
	assert.False(t, proofs.Validate(root, i.Proofs[index2], i.Proofs[cousinIndex], int64(len(i.Proofs))))
	// wrong tree size
	assert.False(t, proofs.Validate(root, i.Proofs[index], i.Proofs[cousinIndex], int64(len(i2.Proofs))))
}
