package types

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
	"time"
)

func TestNewSessionKey(t *testing.T) {
	appPubKey := crypto.PublicKey(getRandomPubKey())
	ctx := newContext(t, false)
	blockhash := hex.EncodeToString(ctx.BlockHeader().LastBlockId.Hash)
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	bitcoin, err := NonNativeChain{
		Ticker:  "btc",
		Netid:   "1",
		Version: "0.19.0.1",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	key1, err := NewSessionKey(appPubKey.String(), ethereum, blockhash)
	assert.Nil(t, err)
	assert.NotNil(t, key1)
	assert.NotEmpty(t, key1)
	assert.Nil(t, HashVerification(hex.EncodeToString(key1)))
	key2, err := NewSessionKey(appPubKey.String(), bitcoin, blockhash)
	assert.Nil(t, err)
	assert.NotNil(t, key2)
	assert.NotEmpty(t, key2)
	assert.Nil(t, HashVerification(hex.EncodeToString(key2)))
	assert.Equal(t, len(key1), len(key2))
	assert.NotEqual(t, key1, key2)
}

func TestSessionKey_Validate(t *testing.T) {
	fakeKey1 := SessionKey([]byte("fakekey"))
	fakeKey2 := SessionKey([]byte(""))
	realKey := SessionKey(hash([]byte("validKey")))
	assert.NotNil(t, fakeKey1.Validate())
	assert.NotNil(t, fakeKey2.Validate())
	assert.Nil(t, realKey.Validate())
}

func TestNewSessionNodes(t *testing.T) {
	fakeSessionKey, err := hex.DecodeString("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fakePubKey1, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab81")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey2, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab82")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey3, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab83")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey4, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab84")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey5, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab85")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey6, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab86")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey7, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab87")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey8, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab88")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey9, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab89")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey10, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab8A")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey11, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab8B")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fakePubKey12, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab8C")
	if err != nil {
		t.Fatalf(err.Error())
	}
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
	var allNodes []exported.ValidatorI
	node12 := nodesTypes.Validator{
		Address:                 fakePubKey12.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey12),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node1 := nodesTypes.Validator{
		Address:                 fakePubKey1.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey1),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node2 := nodesTypes.Validator{
		Address:                 fakePubKey2.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey2),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node3 := nodesTypes.Validator{
		Address:                 fakePubKey3.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey3),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node4 := nodesTypes.Validator{
		Address:                 fakePubKey4.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey4),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node5 := nodesTypes.Validator{
		Address:                 fakePubKey5.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey5),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node6 := nodesTypes.Validator{
		Address:                 fakePubKey6.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey6),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node7 := nodesTypes.Validator{
		Address:                 fakePubKey7.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey7),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node8 := nodesTypes.Validator{
		Address:                 fakePubKey8.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey8),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node9 := nodesTypes.Validator{
		Address:                 fakePubKey9.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey9),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node10 := nodesTypes.Validator{
		Address:                 fakePubKey10.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey10),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	node11 := nodesTypes.Validator{
		Address:                 fakePubKey11.Address(),
		ConsPubKey:              ed25519.PubKeyEd25519(fakePubKey11),
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  []string{ethereum},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	allNodes = make([]exported.ValidatorI, 12)
	allNodes[0] = node12
	allNodes[1] = node1
	allNodes[2] = node2
	allNodes[3] = node3
	allNodes[4] = node4
	allNodes[5] = node5
	allNodes[6] = node6
	allNodes[7] = node7
	allNodes[8] = node8
	allNodes[9] = node9
	allNodes[10] = node10
	allNodes[11] = node11
	sessionNodes, err := NewSessionNodes(ethereum, fakeSessionKey, allNodes, 5)
	assert.Nil(t, err)
	assert.Len(t, sessionNodes, 5)
	assert.NotContains(t, sessionNodes, allNodes[0].(nodesTypes.Validator))
	assert.Contains(t, sessionNodes, allNodes[1].(nodesTypes.Validator))
	assert.Contains(t, sessionNodes, allNodes[2].(nodesTypes.Validator))
	assert.Contains(t, sessionNodes, allNodes[3].(nodesTypes.Validator))
	assert.Contains(t, sessionNodes, allNodes[4].(nodesTypes.Validator))
	assert.Contains(t, sessionNodes, allNodes[5].(nodesTypes.Validator))
	assert.NotContains(t, sessionNodes, allNodes[6].(nodesTypes.Validator))
	assert.NotContains(t, sessionNodes, allNodes[7].(nodesTypes.Validator))
	assert.NotContains(t, sessionNodes, allNodes[8].(nodesTypes.Validator))
	assert.NotContains(t, sessionNodes, allNodes[9].(nodesTypes.Validator))
	assert.NotContains(t, sessionNodes, allNodes[10].(nodesTypes.Validator))
	assert.NotContains(t, sessionNodes, allNodes[11].(nodesTypes.Validator))
	assert.True(t, sessionNodes.Contains(node1))
	assert.True(t, sessionNodes.Contains(node2))
	assert.True(t, sessionNodes.Contains(node3))
	assert.True(t, sessionNodes.Contains(node4))
	assert.True(t, sessionNodes.Contains(node5))
	assert.False(t, sessionNodes.Contains(node6))
	assert.False(t, sessionNodes.Contains(node7))
	assert.False(t, sessionNodes.Contains(node8))
	assert.False(t, sessionNodes.Contains(node9))
	assert.False(t, sessionNodes.Contains(node10))
	assert.False(t, sessionNodes.Contains(node11))
	assert.False(t, sessionNodes.Contains(node12))
	assert.Nil(t, sessionNodes.Validate(5))
	assert.NotNil(t, SessionNodes(make([]exported.ValidatorI, 5)).Validate(5))
	nodesbz, err := sessionNodes.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, nodesbz)
	sn := SessionNodes{}
	sn, err = sn.UnmarshalJSON(nodesbz)
	assert.Nil(t, err)
	assert.NotNil(t, sn)
	assert.Equal(t, sessionNodes, sn)
}
