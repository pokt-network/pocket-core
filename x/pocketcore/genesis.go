package pocketcore

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/keeper"
	"github.com/pokt-network/posmint/x/pos/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis sets up the module based on the genesis state
// First TM block is at height 1, so state updates applied from
// genesis.json are in block 0.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, supplyKeeper types.SupplyKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	return types.GenesisState{}
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data types.GenesisState) error {
	return nil
}
