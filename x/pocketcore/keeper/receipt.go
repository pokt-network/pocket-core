package keeper

import (
	"encoding/hex"
	"fmt"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Receipts (stored proof of work completed)
// set the verified proof of work (receipt)
func (k Keeper) SetReceipt(ctx sdk.Ctx, address sdk.Address, p pc.Receipt) error {
	ctx.Logger().Info(fmt.Sprintf("GetReceipt(address= %v, header= %+v) \n", address.String(), p))
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(p)
	key, err := pc.KeyForReceipt(ctx, address, p.SessionHeader)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

// retrieve the verified proof of work (receipt)
func (k Keeper) GetReceipt(ctx sdk.Ctx, address sdk.Address, header pc.SessionHeader) (receipt pc.Receipt, found bool) {
	ctx.Logger().Info(fmt.Sprintf("GetReceipt(address= %v, header= %+v) \n", address.String(), header))
	store := ctx.KVStore(k.storeKey)
	key, err := pc.KeyForReceipt(ctx, address, header)
	if err != nil {
		ctx.Logger().Error("There was a problem creating a key for the receipt:\n" + err.Error())
		return pc.Receipt{}, false
	}
	res := store.Get(key)
	if res == nil {
		return pc.Receipt{}, false
	}
	k.cdc.MustUnmarshalBinaryBare(res, &receipt)
	return receipt, true
}

// set verified proof of work (receipts) in world state
func (k Keeper) SetReceipts(ctx sdk.Ctx, receipts []pc.Receipt) {
	ctx.Logger().Info(fmt.Sprintf("SetReceipts(receipts %v) \n", receipts))
	store := ctx.KVStore(k.storeKey)
	for _, receipt := range receipts {
		addrbz, err := hex.DecodeString(receipt.ServicerAddress)
		if err != nil {
			panic(fmt.Sprintf("an error occured setting the receipts:\n%v", err))
		}
		bz := k.cdc.MustMarshalBinaryBare(receipt)
		key, err := pc.KeyForReceipt(ctx, addrbz, receipt.SessionHeader)
		if err != nil {
			panic(fmt.Sprintf("an error occured setting the receipts:\n%v", err))
		}
		store.Set(key, bz)
	}
}

// get all verified proof of work (receipts) for this address
func (k Keeper) GetReceipts(ctx sdk.Ctx, address sdk.Address) (receipts []pc.Receipt, err error) {
	ctx.Logger().Info(fmt.Sprintf("GetReceipts(address %v) \n", address.String()))
	store := ctx.KVStore(k.storeKey)
	key, err := pc.KeyForReceipts(address)
	if err != nil {
		return nil, err
	}
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.Receipt
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		receipts = append(receipts, summary)
	}
	return
}

// get all receipts for this address
func (k Keeper) GetAllReceipts(ctx sdk.Ctx) (receipts []pc.Receipt) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.ReceiptKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.Receipt
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		receipts = append(receipts, summary)
	}
	return
}
