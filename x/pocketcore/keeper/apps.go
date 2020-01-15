package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// get all the apps from the world state
func (k Keeper) GetAllApps(ctx sdk.Context) []exported.ApplicationI {
	return k.appKeeper.AllApplications(ctx)
}

// get an app from the world state
func (k Keeper) GetApp(ctx sdk.Context, address sdk.Address) (a exported.ApplicationI, found bool) {
	a = k.appKeeper.Application(ctx, address)
	if a == nil {
		return a, false
	}
	return a, true
}

// get an app from a public key string
func (k Keeper) GetAppFromPublicKey(ctx sdk.Context, pubKey string) (app exported.ApplicationI, found bool) {
	pk, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return nil, false
	}
	return k.GetApp(ctx, pk.Address())
}
