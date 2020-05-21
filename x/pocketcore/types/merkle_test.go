package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/willf/bloom"
	"math"
	"reflect"
	"testing"
)

func TestEvidence_GenerateMerkleRoot(t *testing.T) {
	appPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	clientPrivateKey := GetRandomPrivateKey()
	clientPublicKey := clientPrivateKey.PublicKey().RawString()
	nodePubKey := getRandomPubKey()
	ethereum := hex.EncodeToString([]byte{01})
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
	i := Evidence{
		Bloom: *bloom.New(10000, 4),
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		NumOfProofs: 5,
		Proofs: []Proof{
			RelayProof{
				Entropy:            3238283,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            34939492,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            12383,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            96384,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            96384812,
				RequestHash:        validAAT.HashString(), // fake
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

func TestEvidence_GenerateMerkleProof(t *testing.T) {
	appPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	clientPrivateKey := GetRandomPrivateKey()
	clientPublicKey := clientPrivateKey.PublicKey().RawString()
	nodePubKey := getRandomPubKey()
	ethereum := hex.EncodeToString([]byte{01})
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
	i := Evidence{
		Bloom: *bloom.New(10000, 4),
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		NumOfProofs: 5,
		Proofs: []Proof{
			RelayProof{
				Entropy:            3238283,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            34939492,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            12383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            96384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            96384812,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
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

func TestEvidence_VerifyMerkleProof(t *testing.T) {
	appPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	clientPrivateKey := GetRandomPrivateKey()
	clientPublicKey := clientPrivateKey.PublicKey().RawString()
	nodePubKey := getRandomPubKey()
	ethereum := hex.EncodeToString([]byte{01})
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
	i := Evidence{
		Bloom: *bloom.New(10000, 4),
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		NumOfProofs: 5,
		Proofs: []Proof{
			RelayProof{
				Entropy:            83,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            3492332332249492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            121212123232323383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            23121223232396384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            963223233238481322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
		},
	}
	i2 := Evidence{
		Bloom: *bloom.New(10000, 4),
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		NumOfProofs: 9,
		Proofs: []Proof{
			RelayProof{
				Entropy:            82398289423,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            34932332249492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            1212121232383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            23192932384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            2993223481322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            993223423981322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            90333981322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            2398123322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			RelayProof{
				Entropy:            99322342381322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
		},
	}
	index := 4
	root := i.GenerateMerkleRoot()
	proofs, cousinIndex := i.GenerateMerkleProof(index)
	res, _ := proofs.Validate(root, i.Proofs[index], i.Proofs[cousinIndex], int64(len(i.Proofs)))
	assert.True(t, res)
	index2 := 0
	root2 := i2.GenerateMerkleRoot()
	proofs2, cousinIndex2 := i2.GenerateMerkleProof(index2)
	res, _ = proofs2.Validate(root2, i2.Proofs[index2], i2.Proofs[cousinIndex2], int64(len(i2.Proofs)))
	assert.True(t, res)
	// wrong root
	res, _ = proofs.Validate(root2, i.Proofs[index], i.Proofs[cousinIndex], int64(len(i.Proofs)))
	assert.False(t, res)
	// wrong cousin provided
	res, _ = proofs.Validate(root, i.Proofs[index], i.Proofs[cousinIndex2], int64(len(i.Proofs)))
	assert.False(t, res)
	// wrong leaf provided
	res, _ = proofs.Validate(root, i.Proofs[index2], i.Proofs[cousinIndex], int64(len(i.Proofs)))
	assert.False(t, res)
	// wrong tree size
	res, _ = proofs.Validate(root, i.Proofs[index], i.Proofs[cousinIndex], int64(len(i2.Proofs)))
	assert.False(t, res)
}
