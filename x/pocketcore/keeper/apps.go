package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) GetAllApps(ctx sdk.Context) []exported.ApplicationI {
	return k.appKeeper.GetAllApplications(ctx)
}

func (k Keeper) GetApp(ctx sdk.Context, address sdk.ValAddress) (node exported.ApplicationI, found bool) {
	return k.appKeeper.GetApplication(ctx, address)
}

func (k Keeper) GetAppFromPublicKey(ctx sdk.Context, pubKey string) (app exported.ApplicationI, found bool) {
	appAddr, err := k.AddressFromPubKeyString(pubKey)
	if err != nil {
		return nil, false
	}
	return k.GetApp(ctx, appAddr)
}

func (k Keeper) GetAppChains(ctx sdk.Context, address sdk.ValAddress) (chains map[string]struct{}, found bool) {
	node, found := k.GetApp(ctx, address)
	if !found {
		return nil, false
	}
	return node.GetChains(), true
}

func (k Keeper) AppChainsContains(ctx sdk.Context, address sdk.ValAddress, chain string) (contains bool, found bool) {
	chains, found := k.GetAppChains(ctx, address)
	if !found {
		return false, false
	}
	_, contains = chains[chain]
	return contains, true
}
