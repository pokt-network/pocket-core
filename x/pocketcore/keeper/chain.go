package keeper

import (
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// "GenerateChain" - Generates the external network identifier for the non native chain with a few descriptive parameters:
// "ticker"  -> short marketplace name for the blockchain: ETH BTC etc.
// "netid"   -> the network identifier of the specific chain
// "version" -> the version of the client
// "client"  -> the name of the client `geth vs parity`
// "inter"   -> the mode to interface with the client (websockets/http)
// NOTE: this is just a guide to creating a non native chain, due to the nature of Pocket Network,
// a social consensus on network identifiers will be the source of truth.
func GenerateChain(ticker, netid, version, client, inter string) (string, sdk.Error) {
	return pc.NonNativeChain{
		Ticker:  ticker,
		Netid:   netid,
		Version: version,
		Client:  client,
		Inter:   inter,
	}.HashString()
}

// "GetHostedBlockchains" returns the non native chains hosted locally on this node
func (k Keeper) GetHostedBlockchains() *pc.HostedBlockchains {
	return k.hostedBlockchains
}
