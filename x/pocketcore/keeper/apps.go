package keeper

import (
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/exported"
)

// "GetApp" - Retrieves an application from the app store, using the appKeeper (a link to the apps module)
func (k Keeper) GetApp(ctx sdk.Ctx, address sdk.Address) (a exported.ApplicationI, found bool) {
	a = k.appKeeper.Application(ctx, address)
	if a == nil {
		return a, false
	}
	return a, true
}

// "GetAppFromPublicKey" - Retrieves an application from the app store, using the appKeeper (a link to the apps module)
// using a hex string public key
func (k Keeper) GetAppFromPublicKey(ctx sdk.Ctx, pubKey string) (app exported.ApplicationI, found bool) {
	pk, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return nil, false
	}
	return k.GetApp(ctx, sdk.Address(pk.Address()))
}
