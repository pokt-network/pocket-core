package pocketcore

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetAllProofs(ctx, data.Proofs)
	keeper.SetAllUnverifiedProofs(ctx, data.Claims)
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Params: k.GetParams(ctx),
		Proofs: k.GetAllProofs(ctx),
		Claims: k.GetAllUnverifiedProofs(ctx),
	}
}
