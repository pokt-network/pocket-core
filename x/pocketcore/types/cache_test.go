package types

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"reflect"
	"testing"
)

func InitCacheTest() {
	InitCache("data", "data", db.MemDBBackend, db.MemDBBackend, 100, 100)
}

func TestAllEvidence_AddGetEvidence(t *testing.T) {
	InitCacheTest()
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum := hex.EncodeToString([]byte{0001})
	header := SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	proof := RelayProof{
		Entropy:            0,
		RequestHash:        header.HashString(), // fake
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
	SetProof(header, RelayEvidence, proof)
	assert.True(t, reflect.DeepEqual(GetProof(header, RelayEvidence, 0), proof))
}

func TestAllEvidence_DeleteEvidence(t *testing.T) {
	InitCacheTest()
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum := hex.EncodeToString([]byte{0001})
	header := SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	proof := RelayProof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		RequestHash:        header.HashString(), // fake
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	SetProof(header, RelayEvidence, proof)
	assert.True(t, reflect.DeepEqual(GetProof(header, RelayEvidence, 0), proof))
	GetProof(header, RelayEvidence, 0)
	DeleteEvidence(header, RelayEvidence)
	assert.Empty(t, GetProof(header, RelayEvidence, 0))
}

func TestAllEvidence_GetTotalProofs(t *testing.T) {
	InitCacheTest()
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum := hex.EncodeToString([]byte{0001})
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
	proof := RelayProof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		RequestHash:        header.HashString(), // fake
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	proof2 := RelayProof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		RequestHash:        header.HashString(), // fake
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	SetProof(header, RelayEvidence, proof)
	SetProof(header, RelayEvidence, proof2)
	SetProof(header2, RelayEvidence, proof2) // different header so shouldn't be counted
	assert.Equal(t, GetTotalProofs(header, RelayEvidence), int64(2))
}

func TestSetGetSession(t *testing.T) {
	InitCacheTest()
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	session2 := NewTestSession(t, hex.EncodeToString(Hash([]byte("bar"))))
	SetSession(session)
	s, found := GetSession(session.SessionHeader)
	assert.True(t, found)
	assert.Equal(t, s, session)
	s, found = GetSession(session2.SessionHeader)
	assert.False(t, found)
	SetSession(session2)
	s, found = GetSession(session2.SessionHeader)
	assert.True(t, found)
	assert.Equal(t, s, session2)
}

func TestIteratorValue(t *testing.T) {
	ClearSessionCache()
	InitCacheTest()
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	SetSession(session)
	sessIt := SessionIterator()
	assert.Equal(t, session, sessIt.Value())
}

func TestDeleteSession(t *testing.T) {
	InitCacheTest()
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	SetSession(session)
	DeleteSession(session.SessionHeader)
	_, found := GetSession(session.SessionHeader)
	assert.False(t, found)
}

func TestClearCache(t *testing.T) {
	InitCacheTest()
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	SetSession(session)
	ClearSessionCache()
	iter := SessionIterator()
	var count = 0
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		count++
	}
	assert.Zero(t, count)
}

func NewTestSession(t *testing.T, chain string) Session {
	appPubKey := getRandomPubKey()
	var vals []exported.ValidatorI
	for i := 0; i < 5; i++ {
		nodePubKey := getRandomPubKey()
		vals = append(vals, types.Validator{
			Address:      sdk.Address(nodePubKey.Address()),
			PublicKey:    nodePubKey,
			Status:       2,
			Chains:       []string{chain},
			StakedTokens: sdk.ZeroInt(),
		})
	}
	return Session{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey.RawString(),
			Chain:              chain,
			SessionBlockHeight: 1,
		},
		SessionKey:   appPubKey.RawBytes(), // fake
		SessionNodes: vals,
	}
}
