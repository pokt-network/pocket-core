package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/willf/bloom"
)

func TestMultiAppend(t *testing.T) {
	h1 := merkleHash([]byte("a"))
	h2 := merkleHash([]byte("b"))
	sum := uint64ToBytes(0, 1000000)
	m := runtime.MemStats{}
	m2 := runtime.MemStats{}
	//var start time.Time
	runtime.ReadMemStats(&m)
	//start = time.Now()
	x := merkleHash(append(append(h1, h2...), sum...))
	//fmt.Println(time.Since(start))
	runtime.ReadMemStats(&m2)
	//fmt.Println(m2.Alloc - m.Alloc)
	m = runtime.MemStats{}
	m2 = runtime.MemStats{}
	runtime.ReadMemStats(&m)
	//start = time.Now()
	dest := make([]byte, MerkleHashLength*2+16)
	y := merkleHash(MultiAppend(dest, h1, h2, sum))
	//fmt.Println(time.Since(start))
	runtime.ReadMemStats(&m2)
	//fmt.Println(m2.Alloc - m.Alloc)
	assert.Equal(t, x, y)
}

func TestEvidence_GenerateMerkleRoot(t *testing.T) {
	ClearEvidence()
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
		EvidenceType: RelayEvidence,
		NumOfProofs:  5,
		Proofs: []Proof{
			&RelayProof{
				Entropy:            3238283,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            34939492,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            12383,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            96384,
				RequestHash:        validAAT.HashString(), // fake
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
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
	root := i.GenerateMerkleRoot(0)
	assert.NotNil(t, root.Hash)
	assert.NotEmpty(t, root.Hash)
	assert.Nil(t, HashVerification(hex.EncodeToString(root.Hash)))
	assert.True(t, root.isValidRange())
	assert.Zero(t, root.Range.Lower)
	assert.NotZero(t, root.Range.Upper)

	iter := EvidenceIterator()
	// Make sure its stored in order!
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		e := iter.Value()
		assert.Equal(t, i, e)
		newRoot := e.GenerateMerkleRoot(0)
		assert.Equal(t, root, newRoot)
	}
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
		EvidenceType: RelayEvidence,
	}
	index := 4
	proof, leaf := i.GenerateMerkleProof(0, index)
	assert.Len(t, proof.HashRanges, 3)
	assert.Contains(t, i.Proofs, leaf)
	assert.Equal(t, proof.Target.Hash, merkleHash(leaf.Bytes()))
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
			&RelayProof{
				Entropy:            83,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            3492332332249492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            121212123232323383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            23121223232396384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            963223233238481322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
		},
		EvidenceType: RelayEvidence,
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
			&RelayProof{
				Entropy:            82398289423,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            34932332249492,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            1212121232383,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            23192932384,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            2993223481322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            993223423981322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            90333981322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
				Entropy:            2398123322,
				SessionBlockHeight: 1,
				ServicerPubKey:     nodePubKey.RawString(),
				RequestHash:        validAAT.HashString(), // fake
				Blockchain:         ethereum,
				Token:              validAAT,
				Signature:          "",
			},
			&RelayProof{
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
	root := i.GenerateMerkleRoot(0)
	proofs, leaf := i.GenerateMerkleProof(0, index)
	// validate level count on claim by total relays
	res, _ := proofs.Validate(0, root, leaf, len(proofs.HashRanges))
	assert.True(t, res)
	index2 := 0
	root2 := i2.GenerateMerkleRoot(0)
	proofs2, leaf2 := i2.GenerateMerkleProof(0, index2)
	res, _ = proofs2.Validate(0, root2, leaf2, len(proofs2.HashRanges))
	assert.True(t, res)
	// wrong root
	res, _ = proofs.Validate(0, root2, leaf, len(proofs.HashRanges))
	assert.False(t, res)
	// wrong leaf provided
	res, _ = proofs.Validate(0, root, leaf2, len(proofs.HashRanges))
	assert.False(t, res)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func Test_sortAndStructure(t *testing.T) {
	type args struct {
		hr []HashRange
		p  []Proof
	}
	lol := make([]Proof, 0)
	rand.Seed(time.Now().UnixNano())
	sum := 0
	for i := 1; i < 100000; i++ {
		sum += i
		lol = append(lol, RelayProof{
			RequestHash:        RandStringBytes(9),
			Entropy:            rand.Int63n(1000000000000),
			SessionBlockHeight: 1,
			ServicerPubKey:     RandStringBytes(32),
			Blockchain:         "0001",
			Token:              AAT{},
			Signature:          RandStringBytes(64),
		})
	}
	// get the # of proofs
	numberOfProofs := len(lol)
	// initialize the hashRange
	hashRanges := make([]HashRange, numberOfProofs)
	tests := []struct {
		name string
		args args
	}{
		{"sortAndStructure Consistency Test", args{
			hr: hashRanges,
			p:  lol,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := true
			for i := 0; i < 1; i++ {
				gotSortedHR, gotProof := sortAndStructure(tt.args.p)
				gotSortedHR2, gotProof2 := sortAndStructure(tt.args.p)
				assert.Equal(t, len(gotSortedHR), len(gotSortedHR2))
				assert.Equal(t, cap(gotSortedHR), cap(gotSortedHR2))
				if !reflect.DeepEqual(gotSortedHR, gotSortedHR2) {
					fmt.Println("HashRanges Not Equal")
					assert.Equal(t, gotSortedHR, gotSortedHR2)
					jgotSortedHR, _ := json.Marshal(gotSortedHR)
					jgotSortedHR2, _ := json.Marshal(gotSortedHR2)
					fmt.Println(string(jgotSortedHR))
					fmt.Println(string(jgotSortedHR2))
					t.FailNow()
				}
				if !reflect.DeepEqual(gotProof, gotProof2) {
					t.FailNow()
				}
				gotSortedHR3, gotProof3 := structureProofs(gotProof2)
				if !reflect.DeepEqual(gotSortedHR3, gotSortedHR2) {
					fmt.Println("HashRanges Not Equal")
					assert.Equal(t, gotSortedHR3, gotSortedHR2)
					jgotSortedHR3, _ := json.Marshal(gotSortedHR3)
					jgotSortedHR2, _ := json.Marshal(gotSortedHR2)
					fmt.Println(string(jgotSortedHR3))
					fmt.Println(string(jgotSortedHR2))
					t.FailNow()
				}
				if !reflect.DeepEqual(gotProof3, gotProof2) {
					t.FailNow()
				}
			}
			assert.True(t, result)
		})
	}
}

type benchmarkArgs struct {
	hr []HashRange
	p  []Proof
}

func Benchmark_sortAndStructure(b *testing.B) {
	lol := make([]Proof, 0)
	rand.Seed(time.Now().UnixNano())
	sum := 0
	for i := 1; i < 1000000; i++ {
		sum += i
		lol = append(lol, RelayProof{
			RequestHash:        RandStringBytes(9),
			Entropy:            rand.Int63n(1000000000000),
			SessionBlockHeight: 1,
			ServicerPubKey:     RandStringBytes(32),
			Blockchain:         "0001",
			Token:              AAT{},
			Signature:          RandStringBytes(64),
		})
	}
	// get the # of proofs
	numberOfProofs := len(lol)
	// initialize the hashRange
	hashRanges := make([]HashRange, numberOfProofs)
	tests := []struct {
		name string
		args benchmarkArgs
		f    func(proofs []Proof) ([]HashRange, []Proof)
	}{
		{
			name: "custom_qsort_A",
			args: benchmarkArgs{
				hr: hashRanges,
				p:  lol,
			},
			f: sortAndStructure,
		},
	}
	b.StopTimer()
	b.ResetTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				tt.f(tt.args.p)
				b.StopTimer()
			}
		})
	}
}
