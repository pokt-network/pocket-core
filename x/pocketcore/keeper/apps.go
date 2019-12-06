package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// get all the apps from the world state
func (k Keeper) GetAllApps(ctx sdk.Context) []exported.ApplicationI {
	return k.appKeeper.GetAllApplications(ctx)
}

// get an app from the world state
func (k Keeper) GetApp(ctx sdk.Context, address sdk.ValAddress) (node exported.ApplicationI, found bool) {
	return k.appKeeper.GetApplication(ctx, address)
}

// get an app from a public key string
func (k Keeper) GetAppFromPublicKey(ctx sdk.Context, pubKey string) (app exported.ApplicationI, found bool) {
	pk, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return nil, false
	}
	return k.GetApp(ctx, pk.Address())
}

// get the chains of an application
func (k Keeper) GetAppChains(ctx sdk.Context, address sdk.ValAddress) (chains map[string]struct{}, found bool) {
	node, found := k.GetApp(ctx, address)
	if !found {
		return nil, false
	}
	return node.GetChains(), true
}

// see if the app has staked for a specific chain
func (k Keeper) AppChainsContains(ctx sdk.Context, address sdk.ValAddress, chain string) (contains bool, found bool) {
	chains, found := k.GetAppChains(ctx, address)
	if !found {
		return false, false
	}
	_, contains = chains[chain]
	return contains, true
}
