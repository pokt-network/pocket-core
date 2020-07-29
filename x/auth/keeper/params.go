package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/types"
)

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Ctx, params types.Params) {
	k.subspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Ctx) (params types.Params) {
	k.subspace.GetParamSet(ctx, &params)
	return
}
