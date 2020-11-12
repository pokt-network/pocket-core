package auth

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/keeper"
	"github.com/pokt-network/pocket-core/x/auth/types"
	"log"
)

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Ctx, k keeper.Keeper) types.GenesisState {
	params := k.GetParams(ctx)
	accounts := k.GetAllAccountsExport(ctx)
	supply := k.GetSupply(ctx)
	return types.NewGenesisState(params, accounts, supply.GetTotal())
}

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Ctx, k keeper.Keeper, data types.GenesisState) {
	// check for duplicate keys in fee multi
	keys := make(map[string]struct{})
	for _, feeM := range data.Params.FeeMultiplier.FeeMultis {
		if _, found := keys[feeM.Key]; found {
			log.Fatal(fmt.Sprintf("cannot have duplicate message types in feeMultiplierParams: key=%s already found", feeM.Key))
		}
		keys[feeM.Key] = struct{}{}
	}
	k.SetParams(ctx, data.Params)
	for _, account := range data.Accounts {
		k.SetAccount(ctx, account)
	}
	// manually set the total supply based on accounts if not provided
	if data.Supply.Empty() {
		var totalSupply sdk.Coins
		k.IterateAccounts(ctx,
			func(acc Account) (stop bool) {
				totalSupply = totalSupply.Add(acc.GetCoins())
				return false
			},
		)

		data.Supply = totalSupply
	}
	k.SetSupply(ctx, types.NewSupply(data.Supply))
}
