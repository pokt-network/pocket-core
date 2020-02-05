package app

import (
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
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// NewPocketCoreApp is a constructor function for pocketCoreApp
func NewPocketCoreApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *pocketCoreApp {
	app := newPocketBaseApp(logger, db, baseAppOptions...)
	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keys[params.StoreKey], app.tkeys[params.TStoreKey], params.DefaultCodespace)
	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keys[auth.StoreKey],
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		app.ModuleAccountAddrs(),
	)
	// The SupplyKeeper collects transaction fees
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		app.keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		moduleAccountPermissions,
	)
	// The nodesKeeper keeper handles pocket core nodes
	app.nodesKeeper = nodesKeeper.NewKeeper(
		app.cdc,
		app.keys[nodesTypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		app.supplyKeeper,
		app.paramsKeeper.Subspace(nodesTypes.DefaultParamspace),
		nodesTypes.DefaultCodespace,
	)
	// The apps keeper handles pocket core applications
	app.appsKeeper = appsKeeper.NewKeeper(
		app.cdc,
		app.keys[appsTypes.StoreKey],
		app.bankKeeper,
		app.nodesKeeper,
		app.supplyKeeper,
		app.paramsKeeper.Subspace(appsTypes.DefaultParamspace),
		appsTypes.DefaultCodespace,
	)
	// The main pocket core
	app.pocketKeeper = pocketKeeper.NewPocketCoreKeeper(
		app.keys[pocketTypes.StoreKey],
		app.cdc,
		app.nodesKeeper,
		app.appsKeeper,
		getHostedChains(),
		app.paramsKeeper.Subspace(pocketTypes.DefaultParamspace),
		getCoinbasePassphrase(),
	)
	// add the keybase to the pocket core keeper
	app.pocketKeeper.Keybase = MustGetKeybase()
	app.pocketKeeper.TmNode = getTMClient()
	// setup module manager
	app.mm = module.NewManager(
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		nodes.NewAppModule(app.nodesKeeper, app.accountKeeper, app.supplyKeeper),
		apps.NewAppModule(app.appsKeeper, app.supplyKeeper, app.nodesKeeper),
		pocket.NewAppModule(app.pocketKeeper, app.nodesKeeper, app.appsKeeper),
	)
	// setup the order of begin and end blockers
	app.mm.SetOrderBeginBlockers(nodesTypes.ModuleName, appsTypes.ModuleName, pocketTypes.ModuleName)
	app.mm.SetOrderEndBlockers(nodesTypes.ModuleName, appsTypes.ModuleName)
	// setup the order of Genesis
	app.mm.SetOrderInitGenesis(
		auth.ModuleName,
		bank.ModuleName,
		nodesTypes.ModuleName,
		appsTypes.ModuleName,
		pocketTypes.ModuleName,
		supply.ModuleName,
	)
	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper))
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	// load the latest persistent version of the store
	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}
