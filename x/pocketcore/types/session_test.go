package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSessionKey(t *testing.T) {
	appPubKey := getRandomPubKey()
	ctx := newContext(t, false).WithAppVersion("0.0.0")
	blockhash := hex.EncodeToString(ctx.BlockHeader().LastBlockId.Hash)
	ethereum := hex.EncodeToString([]byte{01})
	bitcoin := hex.EncodeToString([]byte{02})
	key1, err := NewSessionKey(appPubKey.RawString(), ethereum, blockhash)
	assert.Nil(t, err)
	assert.NotNil(t, key1)
	assert.NotEmpty(t, key1)
	assert.Nil(t, HashVerification(hex.EncodeToString(key1)))
	key2, err := NewSessionKey(appPubKey.RawString(), bitcoin, blockhash)
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
	realKey := SessionKey(merkleHash([]byte("validKey")))
	assert.NotNil(t, fakeKey1.Validate())
	assert.NotNil(t, fakeKey2.Validate())
	assert.Nil(t, realKey.Validate())
}

//func TestNewSessionNodes(t *testing.T) {
//	fakeSessionKey, err := hex.DecodeString("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey1, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab81")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey2, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab82")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey3, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab83")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey4, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab84")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey5, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab85")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey6, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab86")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey7, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab87")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey8, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab88")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey9, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab89")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey10, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab8A")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey11, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab8B")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	fakePubKey12, err := crypto.NewPublicKey("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab8C")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	ethereum := hex.EncodeToString([]byte{01})
//	var allNodes []exported.ValidatorI
//	node12 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey12.Address()),
//		PublicKey:               fakePubKey12,
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node1 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey1.Address()),
//		PublicKey:               (fakePubKey1),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node2 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey2.Address()),
//		PublicKey:               (fakePubKey2),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node3 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey3.Address()),
//		PublicKey:               (fakePubKey3),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node4 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey4.Address()),
//		PublicKey:               (fakePubKey4),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node5 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey5.Address()),
//		PublicKey:               (fakePubKey5),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node6 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey6.Address()),
//		PublicKey:               (fakePubKey6),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node7 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey7.Address()),
//		PublicKey:               (fakePubKey7),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node8 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey8.Address()),
//		PublicKey:               (fakePubKey8),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node9 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey9.Address()),
//		PublicKey:               (fakePubKey9),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node10 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey10.Address()),
//		PublicKey:               (fakePubKey10),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	node11 := nodesTypes.Validator{
//		Address:                 sdk.Address(fakePubKey11.Address()),
//		PublicKey:               (fakePubKey11),
//		Jailed:                  false,
//		Status:                  sdk.Staked,
//		Chains:                  []string{ethereum},
//		ServiceUrl:              "https://www.google.com:443",
//		StakedTokens:            sdk.NewInt(100000),
//		UnstakingCompletionTime: time.Time{},
//	}
//	allNodes = make([]exported.ValidatorI, 12)
//	allNodes[0] = node12
//	allNodes[1] = node1
//	allNodes[2] = node2
//	allNodes[3] = node3
//	allNodes[4] = node4
//	allNodes[5] = node5
//	allNodes[6] = node6
//	allNodes[7] = node7
//	allNodes[8] = node8
//	allNodes[9] = node9
//	allNodes[10] = node10
//	allNodes[11] = node11
//	k := MockPosKeeper{Validators: allNodes}
//	sessionNodes, err := NewSessionNodes(newContext(t, false).WithAppVersion("0.0.0"), newContext(t, false).WithAppVersion("0.0.0"), k, ethereum, fakeSessionKey, 5)
//	assert.Nil(t, err)
//	assert.Len(t, sessionNodes, 5)
//	assert.Contains(t, sessionNodes, allNodes[0].(nodesTypes.Validator))
//	assert.Contains(t, sessionNodes, allNodes[1].(nodesTypes.Validator))
//	assert.NotContains(t, sessionNodes, allNodes[2].(nodesTypes.Validator))
//	assert.NotContains(t, sessionNodes, allNodes[3].(nodesTypes.Validator))
//	assert.Contains(t, sessionNodes, allNodes[4].(nodesTypes.Validator))
//	assert.NotContains(t, sessionNodes, allNodes[5].(nodesTypes.Validator))
//	assert.NotContains(t, sessionNodes, allNodes[6].(nodesTypes.Validator))
//	assert.Contains(t, sessionNodes, allNodes[7].(nodesTypes.Validator))
//	assert.Contains(t, sessionNodes, allNodes[8].(nodesTypes.Validator))
//	assert.NotContains(t, sessionNodes, allNodes[9].(nodesTypes.Validator))
//	assert.NotContains(t, sessionNodes, allNodes[10].(nodesTypes.Validator))
//	assert.NotContains(t, sessionNodes, allNodes[11].(nodesTypes.Validator))
//	assert.True(t, sessionNodes.Contains(node12))
//	assert.True(t, sessionNodes.Contains(node8))
//	assert.True(t, sessionNodes.Contains(node7))
//	assert.True(t, sessionNodes.Contains(node4))
//	assert.True(t, sessionNodes.Contains(node1))
//	assert.False(t, sessionNodes.Contains(node2))
//	assert.Nil(t, sessionNodes.Validate(5))
//	assert.NotNil(t, SessionNodes(make([]exported.ValidatorI, 5)).Validate(5))
//}
