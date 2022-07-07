package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func TestPocketNodeAdd(t *testing.T) {
	key := GetRandomPrivateKey()
	address := sdk.GetAddress(key.PublicKey())
	AddPocketNode(key, log.NewNopLogger())
	_, ok := GlobalPocketNodes[address.String()]
	assert.True(t, ok)
}

func TestPocketNodeGetByAddress(t *testing.T) {
	key := GetRandomPrivateKey()
	address := sdk.GetAddress(key.PublicKey())
	AddPocketNode(key, log.NewNopLogger())
	node, err := GetPocketNodeByAddress(&address)
	assert.Nil(t, err)
	assert.NotNil(t, node)
}

func TestPocketNodeGet(t *testing.T) {
	key := GetRandomPrivateKey()
	AddPocketNode(key, log.NewNopLogger())
	node := GetPocketNode()
	assert.NotNil(t, node)
}

func TestPocketNodeCleanCache(t *testing.T) {
	key := GetRandomPrivateKey()
	AddPocketNode(key, log.NewNopLogger())
	CleanPocketNodes()
	assert.Nil(t, GlobalSessionCache)
	assert.Nil(t, GlobalEvidenceCache)
	assert.EqualValues(t, 0, len(GlobalPocketNodes))
}

func TestPocketNodeInitCache(t *testing.T) {
	CleanPocketNodes()
	key := GetRandomPrivateKey()
	testingConfig := sdk.DefaultTestingPocketConfig()
	AddPocketNode(key, log.NewNopLogger())
	InitPocketNodeCaches(testingConfig, log.NewNopLogger())
	address := sdk.GetAddress(key.PublicKey())
	node, err := GetPocketNodeByAddress(&address)
	assert.NotNil(t, GlobalSessionCache)
	assert.NotNil(t, GlobalEvidenceCache)
	assert.EqualValues(t, 1, len(GlobalPocketNodes))
	assert.Nil(t, err)
	assert.NotNil(t, node.EvidenceStore)
	assert.NotNil(t, node.SessionStore)
}

func TestPocketNodeInitCaches(t *testing.T) {
	CleanPocketNodes()
	key := GetRandomPrivateKey()
	key2 := GetRandomPrivateKey()
	logger := log.NewNopLogger()
	testingConfig := sdk.DefaultTestingPocketConfig()
	testingConfig.PocketConfig.LeanPocket = true
	AddPocketNode(key, logger)
	AddPocketNode(key2, logger)
	InitPocketNodeCaches(testingConfig, logger)
	assert.NotNil(t, GlobalSessionCache)
	assert.NotNil(t, GlobalEvidenceCache)
	assert.EqualValues(t, 2, len(GlobalPocketNodes))
	addresses := []sdk.Address{sdk.GetAddress(key.PublicKey()), sdk.GetAddress(key2.PublicKey())}
	for _, address := range addresses {
		node, err := GetPocketNodeByAddress(&address)
		assert.Nil(t, err)
		assert.NotNil(t, node.EvidenceStore)
		assert.NotNil(t, node.SessionStore)
	}
}
