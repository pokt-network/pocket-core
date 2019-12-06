package keeper

import (
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// generate the network identifier for the non native chain
func (k Keeper) GenerateChain(ticker, netid, version, client, inter string) (string, sdk.Error) {
	return pc.NonNativeChain{
		Ticker:  ticker,
		Netid:   netid,
		Version: version,
		Client:  client,
		Inter:   inter,
	}.HashString()
}
