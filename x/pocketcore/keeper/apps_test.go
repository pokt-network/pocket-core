package keeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetApp(t *testing.T) {
	ctx, _, apps, _, keeper, _ := createTestInput(t, false)
	a, found := keeper.GetApp(ctx, apps[0].Address)
	assert.True(t, found)
	assert.Equal(t, a, apps[0])
	randomAddr := getRandomValidatorAddress()
	a, found = keeper.GetApp(ctx, randomAddr)
	assert.False(t, found)
}

func TestGetAppFromPublicKey(t *testing.T) {
	ctx, _, apps, _, keeper, _ := createTestInput(t, false)
	pk := apps[0].PublicKey.RawString()
	a, found := keeper.GetAppFromPublicKey(ctx, pk)
	assert.True(t, found)
	assert.Equal(t, a, apps[0])
	randomPubKey := getRandomPubKey().String()
	a, found = keeper.GetAppFromPublicKey(ctx, randomPubKey)
	assert.False(t, found)
}
