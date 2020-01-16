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
	cfg "github.com/pokt-network/posmint/config"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	tmCfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"testing"
	"time"
)

var memCDC *codec.Codec

// pocket core is an extension of baseapp
type memoryPCApp struct {
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
func newMemoryPCBaseApp(logger log.Logger, db dbm.DB, options ...func(*bam.BaseApp)) *memoryPCApp {
	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(memCDC), options...)
	// set version of the baseapp
	bApp.SetAppVersion(appVersion)
	// setup the key value store keys
	k := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.StoreKey, supply.StoreKey, params.StoreKey, pocketTypes.StoreKey)
	// setup the transient store keys
	tkeys := sdk.NewTransientStoreKeys(nodesTypes.TStoreKey, appsTypes.TStoreKey, pocketTypes.TStoreKey, params.TStoreKey)
	// Create the application
	return &memoryPCApp{
		BaseApp: bApp,
		cdc:     memCDC,
		keys:    k,
		tkeys:   tkeys,
	}
}

// NewPocketCoreApp is a constructor function for pocketCoreApp
func NewMemoryPCApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *memoryPCApp {
	app := newMemoryPCBaseApp(logger, db, baseAppOptions...)
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
		memoryModAccPerms,
	)
	// The nodesKeeper keeper handles pocket core nodes
	app.nodesKeeper = nodesKeeper.NewKeeper(
		app.cdc,
		app.keys[nodesTypes.StoreKey],
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

// inits from genesis
func (app *memoryPCApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	return app.mm.InitGenesis(ctx, genState)
}

// setups all of the begin blockers for each module
func (app *memoryPCApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// setups all of the end blockers for each module
func (app *memoryPCApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// loads the hight from the store
func (app *memoryPCApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the pcInstance's module account addresses.
func (app *memoryPCApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range memoryModAccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// exports the app state to json
func (app *memoryPCApp) ExportAppState(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, err error) {
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
	memoryModAccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		nodesTypes.StakedPoolName: {supply.Burner, supply.Staking},
		appsTypes.StakedPoolName:  {supply.Burner, supply.Staking},
		nodesTypes.DAOPoolName:    {supply.Burner, supply.Staking},
		nodesTypes.ModuleName:     nil,
		appsTypes.ModuleName:      nil,
	}
)

var genState cfg.GenesisState

func InMemoryTendermintNode() (*node.Node, keys.Keybase) {
	pk := ed25519.GenPrivKey()
	kb := keys.NewInMemory()
	kp, err := kb.ImportPrivateKeyObject(pk, "test")
	if err != nil {
		panic(err)
	}
	genDocProvider := func() (*types.GenesisDoc, error) {
		return &types.GenesisDoc{
			GenesisTime: time.Time{},
			ChainID:     "pocket-test",
			ConsensusParams: &types.ConsensusParams{
				Block: types.BlockParams{
					MaxBytes:   15000,
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
			AppState:   newMemDefaultGenesisState(kp.PubKey),
		}, nil
	}
	err = kb.SetCoinbase(kp.GetAddress())
	if err != nil {
		panic(err)
	}
	c := config{
		TmConfig: tmCfg.DefaultConfig(),
		Logger:   log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}
	// setup the database
	db := dbm.NewMemDB()
	// open the tracewriter
	traceWriter, err := openTraceWriter(c.TraceWriter)
	if err != nil {
		panic(err)
	}
	// load the node key
	nodeKey := p2p.NodeKey{PrivKey: pk}
	privVal := cfg.GenFilePV(c.TmConfig.PrivValidatorKey, c.TmConfig.PrivValidatorState)
	privVal.Key.PrivKey = pk
	privVal.Key.PubKey = pk.PubKey()
	privVal.Key.Address = pk.PubKey().Address()
	// app creator function
	creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *memoryPCApp {
		return NewMemoryPCApp(logger, db)
	}
	// upgrade the privVal file
	upgradePrivVal(c.TmConfig)
	dbProvider := func(*node.DBContext) (dbm.DB, error) {
		return db, nil
	}
	// create & start tendermint node
	tmNode, err := node.NewNode(
		c.TmConfig,
		privVal,
		&nodeKey,
		proxy.NewLocalClientCreator(creator(c.Logger, db, traceWriter)),
		genDocProvider,
		dbProvider,
		node.DefaultMetricsProvider(c.TmConfig.Instrumentation),
		c.Logger.With("module", "node"),
	)
	if err != nil {
		panic(err)
	}
	return tmNode, kb
}

func newMemDefaultGenesisState(pubKey crypto.PubKey) []byte {
	defaultGenesis := module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).DefaultGenesis()
	rawPOS := defaultGenesis[nodesTypes.ModuleName]
	var posGenesisState nodesTypes.GenesisState
	MemCodec().MustUnmarshalJSON(rawPOS, &posGenesisState)
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey.Address()),
			ConsPubKey:   pubKey,
			Status:       sdk.Bonded,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	res := MemCodec().MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res
	genState = defaultGenesis
	j, _ := MemCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
}

func MemCodec() *codec.Codec {
	if memCDC == nil {
		MakeMemCodec()
	}
	return memCDC
}

func MakeMemCodec() {
	// create a new codec
	memCDC = codec.New()
	// register all of the app module types
	module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).RegisterCodec(memCDC)
	// register the sdk types
	sdk.RegisterCodec(memCDC)
	// register the crypto types
	codec.RegisterCrypto(memCDC)
}

func TestNewInMemoryTest(t *testing.T) {
	MakeMemCodec()
	tmNode, kb := InMemoryTendermintNode()
	assert.NotNil(t, tmNode)
	assert.NotNil(t, kb)
	err := tmNode.Start()
	assert.Nil(t, err)
	time.Sleep(2*time.Second)
	err = tmNode.Stop()
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(tmNode.Config().DBPath)
	if err != nil {
		panic(err)
	}
}
