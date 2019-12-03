package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

func (k Keeper)GenerateChain(ticker, netid, version, client, inter string) (string, sdk.Error){
	return pc.KeyForChain(ticker, netid, version, client, inter)
}
