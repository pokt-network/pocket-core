package blockchain

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// BlockKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type accountKeeper struct {
	accKeeper auth.AccountKeeper
	CoinKeeper bank.Keeper
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the accountKeeper
func NewAccountKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) accountKeeper {
	return accountKeeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (ak accountKeeper) SendCoins(ctx sdk.Context, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, amount sdk.Coins) sdk.Error {
	return ak.CoinKeeper.SendCoins(ctx, fromAddress, toAddress, amount)
}

func (ak accountKeeper) Balance(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return ak.CoinKeeper.GetCoins(ctx, address)
}

//// todo may need to rope this into distribution for economic inflation control and monitoring
// todo need to look into module accounts
//func (ak accountKeeper) Mint(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) (newBalance sdk.Coins, err sdk.Error) {
//	//supply.Keeper().MintCoins()
//	//return
//}

func (ak accountKeeper) HasEnoughCoins(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) bool {
	return ak.CoinKeeper.HasCoins(ctx, address, amount)
}
