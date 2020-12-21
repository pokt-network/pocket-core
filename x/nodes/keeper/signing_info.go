package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// GetValidatorSigningInfo - Retrieve signing information for the validator by address
func (k Keeper) GetValidatorSigningInfo(ctx sdk.Ctx, addr sdk.Address) (info types.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := store.Get(types.KeyForValidatorSigningInfo(addr))
	if bz == nil {
		found = false
		return
	}
	_ = k.Cdc.UnmarshalBinaryLengthPrefixed(bz, &info)
	found = true
	return
}

// SetValidatorSigningInfo - Store signing information for the validator by address
func (k Keeper) SetValidatorSigningInfo(ctx sdk.Ctx, addr sdk.Address, info types.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := k.Cdc.MarshalBinaryLengthPrefixed(&info)
	_ = store.Set(types.KeyForValidatorSigningInfo(addr), bz)
}

func (k Keeper) DeleteValidatorSigningInfo(ctx sdk.Ctx, addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForValidatorSigningInfo(addr))
}

func (k Keeper) ResetValidatorSigningInfo(ctx sdk.Ctx, addr sdk.Address) {
	signInfo, found := k.GetValidatorSigningInfo(ctx, addr)
	if !found {
		ctx.Logger().Error(fmt.Sprintf("error in ResetValidatorSigningInfo: signing info not found, at height: %d", ctx.BlockHeight()))
		signInfo = types.ValidatorSigningInfo{
			Address:     addr,
			StartHeight: ctx.BlockHeight(),
		}
	}
	signInfo.ResetSigningInfo()
	k.SetValidatorSigningInfo(ctx, addr, signInfo)
	k.clearValidatorMissed(ctx, addr)
}

// IterateAndExecuteOverValSigningInfo - Goes over signing info validators and executes handler
func (k Keeper) IterateAndExecuteOverValSigningInfo(ctx sdk.Ctx, handler func(addr sdk.Address, info types.ValidatorSigningInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter, _ := sdk.KVStorePrefixIterator(store, types.ValidatorSigningInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		address, err := types.GetValidatorSigningInfoAddress(iter.Key())
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("unable to execute over validator %s error: %v, at height: %d", iter.Key(), err, ctx.BlockHeight()).Error())
		}
		var info types.ValidatorSigningInfo
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &info)
		if handler(address, info) {
			break
		}
	}
}

// valMissedAt - Check if validator is missed
func (k Keeper) valMissedAt(ctx sdk.Ctx, addr sdk.Address, index int64) (missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := store.Get(types.GetValMissedBlockKey(addr, index))
	if bz == nil { // lazy: treat empty key as not missed
		missed = false
		return
	}
	b := sdk.Bool(missed)
	_ = k.Cdc.UnmarshalBinaryLengthPrefixed(bz, &b)
	return bool(b)
}

// SetValidatorMissedAt - Store missed validator
func (k Keeper) SetValidatorMissedAt(ctx sdk.Ctx, addr sdk.Address, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	b := sdk.Bool(missed)
	bz, _ := k.Cdc.MarshalBinaryLengthPrefixed(&b)
	_ = store.Set(types.GetValMissedBlockKey(addr, index), bz)
}

// clearValidatorMissed - Remove all missed validators from store
func (k Keeper) clearValidatorMissed(ctx sdk.Ctx, addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	iter, _ := sdk.KVStorePrefixIterator(store, types.GetValMissedBlockPrefixKey(addr))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		_ = store.Delete(iter.Key())
	}
}

// IterateAndExecuteOverMissedArray - Stored by *validator* address (not operator address)
func (k Keeper) IterateAndExecuteOverMissedArray(ctx sdk.Ctx,
	address sdk.Address, handler func(index int64, missed bool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	index := int64(0)
	// Array may be sparse
	for ; index < k.SignedBlocksWindow(ctx); index++ {
		var missed bool
		bz, _ := store.Get(types.GetValMissedBlockKey(address, index))
		if bz == nil {
			continue
		}
		b := sdk.Bool(missed)
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(bz, &b)
		if handler(index, missed) {
			break
		}
	}
}
