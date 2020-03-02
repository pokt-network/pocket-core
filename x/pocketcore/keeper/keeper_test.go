package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_GetHostedBlockchains(t *testing.T) {
	ethereum, err := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	bitcoin, err := types.NonNativeChain{
		Ticker:  "btc",
		Netid:   "1",
		Version: "0.19.0.1",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	eth := types.HostedBlockchain{
		Hash: ethereum,
		URL:  "https://www.google.com",
	}
	btc := types.HostedBlockchain{
		Hash: bitcoin,
		URL:  "https://www.google.com",
	}
	_, _, _, _, keeper, _ := createTestInput(t, false)
	hb := keeper.GetHostedBlockchains()
	assert.NotNil(t, hb)
	assert.True(t, hb.ContainsFromString(eth.Hash))
	assert.False(t, hb.ContainsFromString(btc.Hash))
}
