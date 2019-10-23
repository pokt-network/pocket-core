package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// accountKeeper handles all account based state access/modifier
type poktKeeper struct {
	supplyKeeper supply.Keeper
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
}

// newAccountKeeper creates new instances of the accountKeeper
func newPoktKeeper(sk supply.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) poktKeeper {
	return poktKeeper{
		supplyKeeper: sk,
		storeKey:     storeKey,
		cdc:          cdc,
	}
}

// Mints sdk.Coins from pocket core module
func (pk poktKeeper) Mint(ctx sdk.Context, amount sdk.Coins, address sdk.AccAddress) sdk.Result {
	mintErr := pk.supplyKeeper.MintCoins(ctx, ModuleName, amount)
	if mintErr != nil {
		return mintErr.Result()
	}
	sendErr := pk.supplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, address, amount)
	if sendErr != nil {
		return sendErr.Result()
	}
	logString := amount.String() + " was sucessfully minted to " + address.String()
	return sdk.Result{
		Log: logString,
	}
}
