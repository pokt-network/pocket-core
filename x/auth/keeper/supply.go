package keeper

import (
	"fmt"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/exported"
	"github.com/pokt-network/pocket-core/x/auth/types"
)

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Ctx) (supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	b, _ := store.Get(types.SupplyKeyPrefix)
	if b == nil {
		ctx.Logger().Error(fmt.Sprintf("stored supply should not have been nil, at height: %d", ctx.BlockHeight()))
		return
	}
	supply, err := k.DecodeSupply(b)
	if err != nil {
		ctx.Logger().Error(fmt.Sprint(err.Error()))
		return
	}
	return supply
}

// SetSupply sets the Supply to store
func (k Keeper) SetSupply(ctx sdk.Ctx, supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.EncodeSupply(supply)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return
	}
	err = store.Set(types.SupplyKeyPrefix, bz)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return
	}
}

func (k Keeper) EncodeSupply(supply exported.SupplyI) ([]byte, error) {
	s, ok := supply.(types.Supply)
	if !ok {
		return nil, fmt.Errorf("%s", "supplyI must be of type Supply")
	}
	bz, err := k.Cdc.MarshalBinaryLengthPrefixed(&s)
	return bz, err
}

func (k Keeper) DecodeSupply(bz []byte) (exported.SupplyI, error) {
	var supply types.Supply
	err := k.Cdc.LegacyUnmarshalBinaryLengthPrefixed(bz, &supply)
	return supply, err
}
