package keeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetApp(t *testing.T) {
	ctx, _, apps, _, keeper, _, _ := createTestInput(t, false)
	a, found := keeper.GetApp(ctx, apps[0].Address)
	assert.True(t, found)
	assert.Equal(t, a, apps[0])
	randomAddr := getRandomValidatorAddress()
	_, found = keeper.GetApp(ctx, randomAddr)
	assert.False(t, found)
}

func TestGetAppFromPublicKey(t *testing.T) {
	ctx, _, apps, _, keeper, _, _ := createTestInput(t, false)
	pk := apps[0].PublicKey.RawString()
	a, found := keeper.GetAppFromPublicKey(ctx, pk)
	assert.True(t, found)
	assert.Equal(t, a, apps[0])
	randomPubKey := getRandomPubKey().String()
	_, found = keeper.GetAppFromPublicKey(ctx, randomPubKey)
	assert.False(t, found)
}
