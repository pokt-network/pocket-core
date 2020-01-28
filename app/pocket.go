package app

import (
	"encoding/json"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodesKeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
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
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

// pocket core is an extension of baseapp
type pocketCoreApp struct {
	// extends baseapp
	*bam.BaseApp
	// the codec (uses amino)
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

// new pocket core base
func newPocketBaseApp(logger log.Logger, db db.DB, options ...func(*bam.BaseApp)) *pocketCoreApp {
	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), options...)
	// set version of the baseapp
	bApp.SetAppVersion(appVersion)
	// setup the key value store keys
	k := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.StoreKey, supply.StoreKey, params.StoreKey, pocketTypes.StoreKey)
	// setup the transient store keys
	tkeys := sdk.NewTransientStoreKeys(nodesTypes.TStoreKey, appsTypes.TStoreKey, pocketTypes.TStoreKey, params.TStoreKey)
	// Create the application
	return &pocketCoreApp{
		BaseApp: bApp,
		cdc:     cdc,
		keys:    k,
		tkeys:   tkeys,
	}
}

// inits from genesis
func (app *pocketCoreApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	genesisState := cfg.GenesisStateFromFile(app.cdc, genesisPath())
	return app.mm.InitGenesis(ctx, genesisState)
}

// setups all of the begin blockers for each module
func (app *pocketCoreApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// setups all of the end blockers for each module
func (app *pocketCoreApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// loads the hight from the store
func (app *pocketCoreApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the pcInstance's module account addresses.
func (app *pocketCoreApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// exports the app state to json
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

var (
	// module account permissions
	moduleAccountPermissions = map[string][]string{
		auth.FeeCollectorName:     nil,
		nodesTypes.StakedPoolName: {supply.Burner, supply.Minter, supply.Staking},
		appsTypes.StakedPoolName:  {supply.Burner, supply.Minter, supply.Staking},
		nodesTypes.DAOPoolName:    {supply.Burner, supply.Minter, supply.Staking},
		nodesTypes.ModuleName:     nil,
		appsTypes.ModuleName:      nil,
	}
)

const (
	appName    = "pocket-core"
	appVersion = "0.0.1"
)
