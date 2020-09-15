package keeper

import (
	"fmt"
	"os"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis - Init store state from genesis data
func (k Keeper) InitGenesis(ctx sdk.Ctx, data types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)
	// validate acl
	if err := k.GetACL(ctx).Validate(k.GetAllParamNames(ctx)); err != nil {
		k.Logger(ctx).Error(err.Error())
		os.Exit(1)
	}
	dao := k.GetDAOAccount(ctx)
	if dao == nil {
		k.Logger(ctx).Error(fmt.Errorf("%s module account has not been set", types.DAOAccountName).Error())
		os.Exit(1)
	}
	err := k.AuthKeeper.MintCoins(ctx, types.DAOAccountName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, data.DAOTokens)))
	if err != nil {
		k.Logger(ctx).Error(fmt.Errorf("unable to set dao tokens: %s", err.Error()).Error())
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (k Keeper) ExportGenesis(ctx sdk.Ctx) types.GenesisState {
	return types.NewGenesisState(k.GetParams(ctx), k.GetDAOTokens(ctx))
}
