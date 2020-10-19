package types

import (
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/exported"
)

// "GetAppFromPublicKey" - Retrieves an application from the app store, using the appKeeper (a link to the apps module)
// using a hex string public key
func GetAppFromPublicKey(ctx sdk.Ctx, appsKeeper AppsKeeper, pubKey string) (app exported.ApplicationI, found bool) {
	pk, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return nil, false
	}
	return GetApp(ctx, appsKeeper, pk.Address().Bytes())
}

// "GetApp" - Retrieves an application from the app store, using the appKeeper (a link to the apps module)
func GetApp(ctx sdk.Ctx, appsKeeper AppsKeeper, address sdk.Address) (a exported.ApplicationI, found bool) {
	a = appsKeeper.Application(ctx, address)
	if a == nil {
		return a, false
	}
	return a, true
}
