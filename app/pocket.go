package app

import (
	"encoding/json"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodesKeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	bam "github.com/pokt-network/pocket-core/baseapp"
	"github.com/pokt-network/pocket-core/codec"
	cfg "github.com/pokt-network/pocket-core/config"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/gov"
	govKeeper "github.com/pokt-network/pocket-core/x/gov/keeper"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	db "github.com/tendermint/tm-db"
)

// pocket core is an extension of baseapp
type PocketCoreApp struct {
	// extends baseapp
	*bam.BaseApp
	// the codec (uses amino)
	cdc *codec.Codec
	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey
	// Keepers for each module
	accountKeeper auth.Keeper
	appsKeeper    appsKeeper.Keeper
	nodesKeeper   nodesKeeper.Keeper
	govKeeper     govKeeper.Keeper
	pocketKeeper  pocketKeeper.Keeper
	// Module Manager
	mm *module.Manager
}

// new pocket core base
func newPocketBaseApp(logger log.Logger, db db.DB, options ...func(*bam.BaseApp)) *PocketCoreApp {
	Codec()
	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), options...)
	// set version of the baseapp
	bApp.SetAppVersion(AppVersion)
	// setup the key value store keys
	k := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.StoreKey, gov.StoreKey, pocketTypes.StoreKey)
	// setup the transient store keys
	tkeys := sdk.NewTransientStoreKeys(nodesTypes.TStoreKey, appsTypes.TStoreKey, pocketTypes.TStoreKey, gov.TStoreKey)
	// add params keys too
	// Create the application
	return &PocketCoreApp{
		BaseApp: bApp,
		cdc:     cdc,
		keys:    k,
		tkeys:   tkeys,
	}
}

// inits from genesis
func (app *PocketCoreApp) InitChainer(ctx sdk.Ctx, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState cfg.GenesisState
	switch GlobalGenesisType {
	case MainnetGenesisType:
		genesisState = GenesisStateFromJson(mainnetGenesis)
	case TestnetGenesisType:
		genesisState = GenesisStateFromJson(testnetGenesis)
	default:
		genesisState = cfg.GenesisStateFromFile(cdc, GlobalConfig.PocketConfig.DataDir+FS+ConfigDirName+FS+GlobalConfig.PocketConfig.GenesisName)
	}
	return app.mm.InitGenesis(ctx, genesisState)
}

var GenState cfg.GenesisState

// inits from genesis
func (app *PocketCoreApp) InitChainerWithGenesis(ctx sdk.Ctx, req abci.RequestInitChain) abci.ResponseInitChain {
	return app.mm.InitGenesis(ctx, GenState)
}

// setups all of the begin blockers for each module
func (app *PocketCoreApp) BeginBlocker(ctx sdk.Ctx, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// setups all of the end blockers for each module
func (app *PocketCoreApp) EndBlocker(ctx sdk.Ctx, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// loads the hight from the store
func (app *PocketCoreApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the pcInstance's module account addresses.
func (app *PocketCoreApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// exports the app state to json
func (app *PocketCoreApp) ExportAppState(height int64, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, err error) {
	// as if they could withdraw from the start of the next block
	ctx, err := app.NewContext(height)
	if err != nil {
		return nil, err
	}
	genState := app.mm.ExportGenesis(ctx)
	appState, err = Codec().MarshalJSONIndent(genState, "", "    ")
	if err != nil {
		return nil, err
	}
	return appState, nil
}

func (app *PocketCoreApp) NewContext(height int64) (sdk.Ctx, error) {
	store := app.Store()
	blockStore := app.BlockStore()
	ctx := sdk.NewContext(store, abci.Header{}, false, app.Logger()).WithBlockStore(blockStore)
	return ctx.PrevCtx(height)
}

func (app *PocketCoreApp) GetClient() client.Client {
	return app.pocketKeeper.TmNode
}

var (
	// module account permissions
	moduleAccountPermissions = map[string][]string{
		auth.FeeCollectorName:     {auth.Burner, auth.Minter, auth.Staking},
		nodesTypes.StakedPoolName: {auth.Burner, auth.Minter, auth.Staking},
		appsTypes.StakedPoolName:  {auth.Burner, auth.Minter, auth.Staking},
		govTypes.DAOAccountName:   {auth.Burner, auth.Minter, auth.Staking},
		nodesTypes.ModuleName:     {auth.Burner, auth.Minter, auth.Staking},
		appsTypes.ModuleName:      nil,
	}
)

const (
	appName = "pocket-core"
)
