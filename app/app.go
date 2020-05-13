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
	cfg "github.com/pokt-network/posmint/config"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/gov"
	govKeeper "github.com/pokt-network/posmint/x/gov/keeper"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	dbm "github.com/tendermint/tm-db"
)

const (
	AppVersion = "RC-0.3.0"
)

// NewPocketCoreApp is a constructor function for PocketCoreApp
func NewPocketCoreApp(genState cfg.GenesisState, keybase keys.Keybase, tmClient client.Client, hostedChains *pocketTypes.HostedBlockchains, logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *PocketCoreApp {
	app := newPocketBaseApp(logger, db, baseAppOptions...)
	// setup subspaces
	authSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	nodesSubspace := sdk.NewSubspace(nodesTypes.DefaultParamspace)
	appsSubspace := sdk.NewSubspace(appsTypes.DefaultParamspace)
	pocketSubspace := sdk.NewSubspace(pocketTypes.DefaultParamspace)
	// The AuthKeeper handles address -> account lookups
	app.accountKeeper = auth.NewKeeper(
		app.cdc,
		app.keys[auth.StoreKey],
		authSubspace,
		moduleAccountPermissions,
	)
	// The nodesKeeper keeper handles pocket core nodes
	app.nodesKeeper = nodesKeeper.NewKeeper(
		app.cdc,
		app.keys[nodesTypes.StoreKey],
		app.accountKeeper,
		nodesSubspace,
		nodesTypes.DefaultCodespace,
	)
	// The apps keeper handles pocket core applications
	app.appsKeeper = appsKeeper.NewKeeper(
		app.cdc,
		app.keys[appsTypes.StoreKey],
		app.nodesKeeper,
		app.accountKeeper,
		appsSubspace,
		appsTypes.DefaultCodespace,
	)
	// The main pocket core
	app.pocketKeeper = pocketKeeper.NewKeeper(
		app.keys[pocketTypes.StoreKey],
		app.cdc,
		app.nodesKeeper,
		app.appsKeeper,
		hostedChains,
		pocketSubspace,
	)
	// The governance keeper
	app.govKeeper = govKeeper.NewKeeper(
		app.cdc,
		app.keys[pocketTypes.StoreKey],
		app.tkeys[pocketTypes.StoreKey],
		govTypes.DefaultCodespace,
		app.accountKeeper,
		authSubspace, nodesSubspace, appsSubspace, pocketSubspace,
	)
	// add the keybase to the pocket core keeper
	app.pocketKeeper.Keybase = keybase
	app.pocketKeeper.TmNode = tmClient
	// setup module manager
	app.mm = module.NewManager(
		auth.NewAppModule(app.accountKeeper),
		nodes.NewAppModule(app.nodesKeeper),
		apps.NewAppModule(app.appsKeeper),
		pocket.NewAppModule(app.pocketKeeper),
		gov.NewAppModule(app.govKeeper),
	)
	// setup the order of begin and end blockers
	app.mm.SetOrderBeginBlockers(nodesTypes.ModuleName, appsTypes.ModuleName, pocketTypes.ModuleName)
	app.mm.SetOrderEndBlockers(nodesTypes.ModuleName, appsTypes.ModuleName)
	// setup the order of Genesis
	app.mm.SetOrderInitGenesis(
		auth.ModuleName,
		nodesTypes.ModuleName,
		appsTypes.ModuleName,
		pocketTypes.ModuleName,
		gov.ModuleName,
	)
	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
	// The initChainer handles translating the genesis.json file into initial state for the network
	if genState == nil {
		app.SetInitChainer(app.InitChainer)
	} else {
		app.SetInitChainer(app.InitChainerWithGenesis)
	}
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper))
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	app.SetAppVersion(AppVersion)
	// load the latest persistent version of the store
	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}
