package keeper

import (
	"github.com/pokt-network/posmint/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestKeeper_GetNodeFromPublicKey(t *testing.T) {
	ctx, vals, _, _, keeper := createTestInput(t, false)
	node, found := keeper.GetNodeFromPublicKey(ctx, crypto.PublicKey(vals[0].ConsPubKey.(ed25519.PubKeyEd25519)).String())
	assert.True(t, found)
	assert.Equal(t, vals[0], node)
}
