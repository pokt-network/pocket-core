package types

import (
	"encoding/hex"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"reflect"
	"testing"
)

func InitCacheTest() {
	logger := log.NewNopLogger()
	testingConfig := sdk.DefaultTestingPocketConfig()
	testingConfig.PocketConfig.LeanPocket = true
	InitConfig(&HostedBlockchains{
		M: make(map[string]HostedBlockchain),
	}, logger, testingConfig)
	// init cache in memory

	// init needed maps for cache
	servicerPk := GetRandomPrivateKey()
	InitNodeWithCacheLean(servicerPk)

}

func TestMain(m *testing.M) {
	InitCacheTest()
	m.Run()
	err := os.RemoveAll("data")
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}

func TestIsUniqueProof(t *testing.T) {
	h := SessionHeader{
		ApplicationPubKey:  "0",
		Chain:              "0001",
		SessionBlockHeight: 0,
	}
	e, _ := GetEvidence(h, RelayEvidence, sdk.NewInt(100000))
	p := RelayProof{
		Entropy: 1,
	}
	p1 := RelayProof{
		Entropy: 2,
	}
	assert.True(t, IsUniqueProof(p, e), "p is unique")
	e.AddProof(p)
	SetEvidence(e)
	e, err := GetEvidence(h, RelayEvidence, sdk.ZeroInt())
	assert.Nil(t, err)
	assert.False(t, IsUniqueProof(p, e), "p is no longer unique")
	assert.True(t, IsUniqueProof(p1, e), "p is unique")
}

func TestAllEvidence_AddGetEvidence(t *testing.T) {
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
	SetProof(header, RelayEvidence, proof, sdk.NewInt(100000))
	assert.True(t, reflect.DeepEqual(GetProof(header, RelayEvidence, 0), proof))
}


func TestAllEvidence_DeleteEvidence(t *testing.T) {
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
	SetProof(header, RelayEvidence, proof, sdk.NewInt(100000))
	assert.True(t, reflect.DeepEqual(GetProof(header, RelayEvidence, 0), proof))
	GetProof(header, RelayEvidence, 0)
	_ = DeleteEvidence(header, RelayEvidence)
	assert.Empty(t, GetProof(header, RelayEvidence, 0))
}

func TestAllEvidence_DeleteEvidenceLean(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := GetPrivateKeyLean().PublicKey()
	address := sdk.GetAddress(servicerPubKey)
	servicerPubKeyString := servicerPubKey.RawString()
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
		ServicerPubKey:     servicerPubKeyString,
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
	SetProofLean(header, RelayEvidence, proof, sdk.NewInt(100000), &address)
	assert.True(t, reflect.DeepEqual(GetProofLean(header, RelayEvidence, 0, &address), proof))
	GetProofLean(header, RelayEvidence, 0, &address)
	_ = DeleteEvidenceLean(header, RelayEvidence, &address)
	assert.Empty(t, GetProofLean(header, RelayEvidence, 0, &address))
}

func TestAllEvidence_GetTotalProofs(t *testing.T) {
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
	SetProof(header, RelayEvidence, proof, sdk.NewInt(100000))
	SetProof(header, RelayEvidence, proof2, sdk.NewInt(100000))
	SetProof(header2, RelayEvidence, proof2, sdk.NewInt(100000)) // different header so shouldn't be counted
	_, totalRelays := GetTotalProofs(header, RelayEvidence, sdk.NewInt(100000))
	assert.Equal(t, totalRelays, int64(2))
}

func TestAllEvidence_GetTotalProofsLean(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := GetPrivateKeyLean().PublicKey()
	address := sdk.GetAddress(servicerPubKey)
	servicerPubKeyString := servicerPubKey.RawString()
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
		ServicerPubKey:     servicerPubKeyString,
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
		ServicerPubKey:     servicerPubKeyString,
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
	SetProofLean(header, RelayEvidence, proof, sdk.NewInt(100000), &address)
	SetProofLean(header, RelayEvidence, proof2, sdk.NewInt(100000), &address)
	SetProofLean(header2, RelayEvidence, proof2, sdk.NewInt(100000), &address) // different header so shouldn't be counted
	_, totalRelays := GetTotalProofsLean(header, RelayEvidence, sdk.NewInt(100000), &address)
	assert.Equal(t, totalRelays, int64(2))
}

func TestSetGetSession(t *testing.T) {
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	session2 := NewTestSession(t, hex.EncodeToString(Hash([]byte("bar"))))
	SetSession(session)
	s, found := GetSession(session.SessionHeader)
	assert.True(t, found)
	assert.Equal(t, s, session)
	_, found = GetSession(session2.SessionHeader)
	assert.False(t, found)
	SetSession(session2)
	s, found = GetSession(session2.SessionHeader)
	assert.True(t, found)
	assert.Equal(t, s, session2)
}

func TestSetGetSessionLean(t *testing.T) {
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	session2 := NewTestSession(t, hex.EncodeToString(Hash([]byte("bar"))))

	randomAddr := sdk.GetAddress(GetPrivateKeyLean().PublicKey())
	SetSessionLean(session, &randomAddr)

	s, found := GetSessionLean(session.SessionHeader, &randomAddr)
	assert.True(t, found)
	assert.Equal(t, s, session)
	_, found = GetSessionLean(session2.SessionHeader, &randomAddr)
	assert.False(t, found)
	SetSession(session2)
	s, found = GetSessionLean(session2.SessionHeader, &randomAddr)
	assert.True(t, found)
	assert.Equal(t, s, session2)
}

func TestDeleteSession(t *testing.T) {
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	SetSession(session)
	DeleteSession(session.SessionHeader)
	_, found := GetSession(session.SessionHeader)
	assert.False(t, found)
}

func TestDeleteSessionLean(t *testing.T) {
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	randomAddr := sdk.GetAddress(GetPrivateKeyLean().PublicKey())
	SetSessionLean(session, &randomAddr)
	DeleteSessionLean(session.SessionHeader, &randomAddr)
	_, found := GetSessionLean(session.SessionHeader, &randomAddr)
	assert.False(t, found)
}

func TestClearCache(t *testing.T) {
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

func TestClearCacheLean(t *testing.T) {
	session := NewTestSession(t, hex.EncodeToString(Hash([]byte("foo"))))
	randomAddr := sdk.GetAddress(GetPrivateKeyLean().PublicKey())
	SetSessionLean(session, &randomAddr)
	ClearSessionCacheLean(&randomAddr)
	iter := SessionIteratorLean(&randomAddr)
	var count = 0
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		count++
	}
	assert.Zero(t, count)
}

func NewTestSession(t *testing.T, chain string) Session {
	appPubKey := getRandomPubKey()
	var vals []sdk.Address
	for i := 0; i < 5; i++ {
		nodePubKey := getRandomPubKey()
		vals = append(vals, sdk.Address(nodePubKey.Address()))
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
