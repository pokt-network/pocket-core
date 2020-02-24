package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// get an app from the world state
func (k Keeper) GetApp(ctx sdk.Context, address sdk.Address) (a exported.ApplicationI, found bool) {
	ctx.Logger().Info(fmt.Sprintf("GetApp(Address = %v)", address.String()))
	a = k.appKeeper.Application(ctx, address)
	if a == nil {
		return a, false
	}
	return a, true
}

// get an app from a public key string
func (k Keeper) GetAppFromPublicKey(ctx sdk.Context, pubKey string) (app exported.ApplicationI, found bool) {
	ctx.Logger().Info(fmt.Sprintf("GetApp(PubKey = %v) \n", pubKey))
	pk, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return nil, false
	}
	return k.GetApp(ctx, sdk.Address(pk.Address()))
}
