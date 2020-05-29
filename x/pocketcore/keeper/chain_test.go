package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_GetHostedBlockchains(t *testing.T) {
	ethereum := hex.EncodeToString([]byte{01})
	bitcoin := hex.EncodeToString([]byte{02})
	eth := types.HostedBlockchain{
		ID:  ethereum,
		URL: "https://www.google.com:443",
	}
	btc := types.HostedBlockchain{
		ID:  bitcoin,
		URL: "https://www.google.com:443",
	}
	_, _, _, _, keeper, _, _ := createTestInput(t, false)
	hb := keeper.GetHostedBlockchains()
	assert.NotNil(t, hb)
	assert.True(t, hb.Contains(eth.ID))
	assert.False(t, hb.Contains(btc.ID))
}
