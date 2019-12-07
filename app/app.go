package app

import (
	"encoding/json"
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesKeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	bam "github.com/pokt-network/posmint/baseapp"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const appName = "pocket-core"

var (
	// NewBasicManager is in charge of setting up basic module elemnets
	ModuleBasics = module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	)
	// module account permissions
	moduleAccountPermissions = map[string][]string{
		auth.FeeCollectorName: nil,
	}
)

// NewPocketCoreApp is a constructor function for pocketCoreApp
func NewPocketCoreApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *pocketCoreApp {
	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()
	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetAppVersion("0.0.1")
	keys := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.TStoreKey, supply.StoreKey, params.StoreKey, pocketTypes.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(nodesTypes.TStoreKey, appsTypes.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &pocketCoreApp{
		BaseApp: bApp,
		cdc:     cdc,
		keys:    keys,
		tkeys:   tkeys,
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)
	// Set specific supspaces
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSupspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	nodesSubspace := app.paramsKeeper.Subspace(nodesTypes.DefaultParamspace)
	appsSubspace := app.paramsKeeper.Subspace(appsTypes.DefaultParamspace)
	pocketSubspace := app.paramsKeeper.Subspace(pocketTypes.DefaultParamspace)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		authSubspace,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		bankSupspace,
		bank.DefaultCodespace,
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		moduleAccountPermissions,
	)

	// The nodesKeeper keeper
	n := nodesKeeper.NewKeeper(
		app.cdc,
		keys[nodesTypes.StoreKey],
		app.bankKeeper,
		app.supplyKeeper,
		nodesSubspace,
		nodesTypes.DefaultCodespace,
	)

	// The apps keeper
	a := appsKeeper.NewKeeper(
		app.cdc,
		keys[appsTypes.StoreKey],
		app.bankKeeper,
		n,
		app.supplyKeeper,
		appsSubspace,
		appsTypes.DefaultCodespace,
	)

	keybase := GetKeybase()
	tendermintNode := GetTendermintNode()
	hostedBlockchains := GetHostedChains()
	passphrase := GetCoinbasePassphrase()

	// The NameserviceKeeper is the Keeper from the module for this tutorial
	// It handles interactions with the namestore
	app.pocketKeeper = pocketKeeper.NewPocketCoreKeeper(
		keys[pocketTypes.StoreKey],
		app.cdc,
		n,
		a,
		keybase,
		hostedBlockchains,
		pocketSubspace,
		passphrase,
	)

	app.mm = module.NewManager(
		auth.NewAppModule(app.accountKeeper, tendermintNode, keybase),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper, tendermintNode, keybase),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper, tendermintNode, keybase),
		nodes.NewAppModule(app.nodesKeeper, app.accountKeeper, app.supplyKeeper, tendermintNode, keybase),
		apps.NewAppModule(app.appsKeeper, app.supplyKeeper, app.nodesKeeper, tendermintNode, keybase),
		pocket.NewAppModule(app.pocketKeeper, app.nodesKeeper, app.appsKeeper),
	)

	app.mm.SetOrderBeginBlockers(appsTypes.ModuleName, nodesTypes.ModuleName, pocketTypes.ModuleName)
	app.mm.SetOrderEndBlockers(appsTypes.ModuleName, nodesTypes.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		nodesTypes.ModuleName,
		appsTypes.ModuleName,
		pocketTypes.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		supply.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
