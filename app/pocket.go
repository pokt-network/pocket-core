package app

import (
	"encoding/json"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	nodesKeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	bam "github.com/pokt-network/posmint/baseapp"
	"github.com/pokt-network/posmint/codec"
	cfg "github.com/pokt-network/posmint/config"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
)

type pocketCoreApp struct {
	*bam.BaseApp
	cdc *codec.Codec
	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// Keepers for each module
	accountKeeper auth.AccountKeeper
	appsKeeper    appsKeeper.Keeper
	bankKeeper    bank.Keeper
	supplyKeeper  supply.Keeper
	nodesKeeper   nodesKeeper.Keeper
	paramsKeeper  params.Keeper
	pocketKeeper  pocketKeeper.Keeper

	// Module Manager
	mm *module.Manager
}

func (app *pocketCoreApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	genesisState := cfg.GenesisStateFromFile(app.cdc, GenesisFilepath)
	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *pocketCoreApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}
func (app *pocketCoreApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
func (app *pocketCoreApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *pocketCoreApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *pocketCoreApp) ExportAppState(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, err error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, err
	}
	return appState, nil
}
