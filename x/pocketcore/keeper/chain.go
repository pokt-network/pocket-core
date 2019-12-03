package keeper

import (
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) GenerateChain(ticker, netid, version, client, inter string) (string, sdk.Error) {
	return pc.KeyForChain(ticker, netid, version, client, inter)
}
