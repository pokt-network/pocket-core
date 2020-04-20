package pocketcore

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// "InitGenesis" - Initializes the state with a genesis state object
func InitGenesis(ctx sdk.Ctx, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	// set the params in store
	keeper.SetParams(ctx, data.Params)
	// set the receipt objects in store
	keeper.SetReceipts(ctx, data.Receipts)
	// set the claim objects in store
	keeper.SetClaims(ctx, data.Claims)
	return []abci.ValidatorUpdate{}
}

// "ExportGenesis" - Exports the state in a genesis state object
func ExportGenesis(ctx sdk.Ctx, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Params:   k.GetParams(ctx),
		Receipts: k.GetAllReceipts(ctx),
		Claims:   k.GetAllClaims(ctx),
	}
}
