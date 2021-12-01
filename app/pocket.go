package app

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/libs/os"

	bam "github.com/pokt-network/pocket-core/baseapp"
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/gov"
	govKeeper "github.com/pokt-network/pocket-core/x/gov/keeper"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	nodesKeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

// pocket core is an extension of baseapp
type PocketCoreApp struct {
	// extends baseapp
	*bam.BaseApp
	// the codec (uses amino)
	cdc *codec.Codec
	// Keys to access the substores
	Keys  map[string]*sdk.KVStoreKey
	Tkeys map[string]*sdk.TransientStoreKey
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
func NewPocketBaseApp(logger log.Logger, db db.DB, cache bool, iavlCacheSize int64, options ...func(*bam.BaseApp)) *PocketCoreApp {
	cdc = Codec()
	bam.SetABCILogging(GlobalConfig.PocketConfig.ABCILogging)
	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, cache, iavlCacheSize, auth.DefaultTxDecoder(cdc), cdc, options...)
	// set version of the baseapp
	bApp.SetAppVersion(AppVersion)
	// setup the key value store Keys
	k := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.StoreKey, gov.StoreKey, pocketTypes.StoreKey)
	// setup the transient store Keys
	tkeys := sdk.NewTransientStoreKeys(nodesTypes.TStoreKey, appsTypes.TStoreKey, pocketTypes.TStoreKey, gov.TStoreKey)
	// add params Keys too
	// Create the application
	return &PocketCoreApp{
		BaseApp: bApp,
		cdc:     cdc,
		Keys:    k,
		Tkeys:   tkeys,
	}
}

// inits from genesis
func (app *PocketCoreApp) InitChainer(ctx sdk.Ctx, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	switch GlobalGenesisType {
	case MainnetGenesisType:
		genesisState = GenesisStateFromJson(mainnetGenesis)
	case TestnetGenesisType:
		genesisState = GenesisStateFromJson(testnetGenesis)
	default:
		genesisState = GenesisStateFromFile(cdc, GlobalConfig.PocketConfig.DataDir+FS+sdk.ConfigDirName+FS+GlobalConfig.PocketConfig.GenesisName)
	}
	return app.mm.InitGenesis(ctx, genesisState)
}

var GenState GenesisState

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

// ModuleAccountAddrs returns all the pcInstance's module account addresses.
func (app *PocketCoreApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

type GenesisState map[string]json.RawMessage

func GenesisStateFromFile(cdc *codec.Codec, genFile string) GenesisState {
	if !os.FileExists(genFile) {
		panic(fmt.Errorf("%s does not exist, run `init` first", genFile))
	}
	genDoc := GenesisFileToGenDoc(genFile)
	return GenesisStateFromGenDoc(cdc, *genDoc)
}

func GenesisFileToGenDoc(genFile string) *tmtypes.GenesisDoc {
	if !os.FileExists(genFile) {
		panic(fmt.Errorf("%s does not exist, run `init` first", genFile))
	}
	genDoc, err := tmtypes.GenesisDocFromFile(genFile)
	if err != nil {
		panic(err)
	}
	return genDoc
}

func GenesisStateFromGenDoc(cdc *codec.Codec, genDoc tmtypes.GenesisDoc) (genesisState map[string]json.RawMessage) {
	if err := cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		panic(err)
	}
	return genesisState
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

func (app *PocketCoreApp) ExportState(height int64, chainID string) (string, error) {
	j, err := app.ExportAppState(height, false, nil)
	if err != nil {
		return "", err
	}
	if chainID == "" {
		chainID = "<Input New ChainID>"
	}
	j, _ = Codec().MarshalJSONIndent(types.GenesisDoc{
		ChainID: chainID,
		ConsensusParams: &types.ConsensusParams{
			Block: types.BlockParams{
				MaxBytes:   4000000,
				MaxGas:     -1,
				TimeIotaMs: 1,
			},
			Evidence: types.EvidenceParams{
				MaxAge: 1000000,
			},
			Validator: types.ValidatorParams{
				PubKeyTypes: []string{"ed25519"},
			},
		},
		Validators: nil,
		AppHash:    nil,
		AppState:   j,
	}, "", "    ")
	return SortJSON(j), err
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
