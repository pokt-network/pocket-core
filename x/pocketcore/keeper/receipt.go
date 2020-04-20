package keeper

import (
	"encoding/hex"
	"fmt"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// "SetReceipt" - Sets the receipt object for a certain address in the state storage
func (k Keeper) SetReceipt(ctx sdk.Ctx, address sdk.Address, p pc.Receipt) error {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// marshal the receipt object into amino bz
	bz := k.cdc.MustMarshalBinaryBare(p)
	// generate the key for the receipt
	key, err := pc.KeyForReceipt(ctx, address, p.SessionHeader, p.EvidenceType)
	if err != nil {
		return err
	}
	// set kv into store
	store.Set(key, bz)
	return nil
}

// "GetReceipt" - Retrieves the receipt object for a certain address in the state storage
func (k Keeper) GetReceipt(ctx sdk.Ctx, address sdk.Address, header pc.SessionHeader, evidenceType pc.EvidenceType) (receipt pc.Receipt, found bool) {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// generate the key for the receipt
	key, err := pc.KeyForReceipt(ctx, address, header, evidenceType)
	if err != nil {
		ctx.Logger().Error("There was a problem creating a key for the receipt:\n" + err.Error())
		return pc.Receipt{}, false
	}
	// get the bytes from the store
	res := store.Get(key)
	if res == nil {
		return pc.Receipt{}, false
	}
	// unmarshal bytes into amino-json
	k.cdc.MustUnmarshalBinaryBare(res, &receipt)
	return receipt, true
}

// "SetReceipts" - Sets many receipt objects in the store
func (k Keeper) SetReceipts(ctx sdk.Ctx, receipts []pc.Receipt) {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// loop through all of the receipts
	for _, receipt := range receipts {
		// get the address
		addr, err := hex.DecodeString(receipt.ServicerAddress)
		if err != nil {
			panic(fmt.Sprintf("an error occured setting the receipts:\n%v", err))
		}
		// marshal the receipt into json-amino
		bz := k.cdc.MustMarshalBinaryBare(receipt)
		// generate the key for the receipt
		key, err := pc.KeyForReceipt(ctx, addr, receipt.SessionHeader, receipt.EvidenceType)
		if err != nil {
			panic(fmt.Sprintf("an error occured setting the receipts:\n%v", err))
		}
		// set it in the store
		store.Set(key, bz)
	}
}

// "GetReceipts" - Retrieves all the receipt objects for a certain address
func (k Keeper) GetReceipts(ctx sdk.Ctx, address sdk.Address) (receipts []pc.Receipt, err error) {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// generate the key for the address
	key, err := pc.KeyForReceipts(address)
	if err != nil {
		return nil, err
	}
	// iterate through all of the receipts
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// unmarshal into a new receipt object
		var receipt pc.Receipt
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &receipt)
		// append the new receipt object to the list
		receipts = append(receipts, receipt)
	}
	return
}

// "GetAllReceipts" - Retrieves all the receipt objects in the storage
func (k Keeper) GetAllReceipts(ctx sdk.Ctx) (receipts []pc.Receipt) {
	// get the store
	store := ctx.KVStore(k.storeKey)
	// iterate through the  objects
	iterator := sdk.KVStorePrefixIterator(store, pc.ReceiptKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// unmarshal into a new receipt object
		var receipt pc.Receipt
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &receipt)
		// append the new receipt object to the list
		receipts = append(receipts, receipt)
	}
	return
}
