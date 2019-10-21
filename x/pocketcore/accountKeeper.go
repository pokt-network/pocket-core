package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// accountKeeper handles all account based state access/modifier
type accountKeeper struct {
	accKeeper  auth.AccountKeeper
	CoinKeeper bank.Keeper
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
}

// newAccountKeeper creates new instances of the accountKeeper
func newAccountKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) accountKeeper {
	return accountKeeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// get the balance by address
func (ak accountKeeper) Balance(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return ak.CoinKeeper.GetCoins(ctx, address)
}

// check if an address has enough coins
func (ak accountKeeper) HasEnoughCoins(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) bool {
	return ak.CoinKeeper.HasCoins(ctx, address, amount)
}
