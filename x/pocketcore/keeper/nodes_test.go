package keeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_GetNodeFromPublicKey(t *testing.T) {
	ctx, vals, _, _, keeper := createTestInput(t, false)
	node, found := keeper.GetNodeFromPublicKey(ctx, vals[0].PublicKey.RawString())
	assert.True(t, found)
	assert.Equal(t, vals[0], node)
}
