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
	"github.com/pokt-network/posmint/crypto"
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
	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"testing"
	"time"
)

func NewInMemoryTendermintNode(t *testing.T) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	// make the codec
	makeMemCodec()
	// create the in memory tendermint node and keybase
	tendermintNode, keybase = inMemTendermintNode()
	// test assertions
	if tendermintNode == nil {
		panic("tendermintNode should not be nil")
	}
	if keybase == nil {
		panic("should not be nil")
	}
	assert.NotNil(t, tendermintNode)
	assert.NotNil(t, keybase)
	// start the in memory node
	err := tendermintNode.Start()
	// assert that it is not nil
	assert.Nil(t, err)
	// provide cleanup function
	cleanup = func() {
		err = tendermintNode.Stop()
		if err != nil {
			panic(err)
		}
		err = os.RemoveAll(tendermintNode.Config().DBPath)
		if err != nil {
			panic(err)
		}
	}
	return
}

func TestNewInMemory(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	time.Sleep(2 * time.Second)
}

type memoryPCApp struct {
	*bam.BaseApp
	cdc           *codec.Codec
	keys          map[string]*sdk.KVStoreKey
	tkeys         map[string]*sdk.TransientStoreKey
	accountKeeper auth.AccountKeeper
	appsKeeper    appsKeeper.Keeper
	bankKeeper    bank.Keeper
	supplyKeeper  supply.Keeper
	nodesKeeper   nodesKeeper.Keeper
	paramsKeeper  params.Keeper
	pocketKeeper  pocketKeeper.Keeper
	mm            *module.Manager
}

func newMemoryPCBaseApp(logger log.Logger, db dbm.DB, options ...func(*bam.BaseApp)) *memoryPCApp {
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(memCDC), options...)
	bApp.SetAppVersion(appVersion)
	k := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.StoreKey, supply.StoreKey, params.StoreKey, pocketTypes.StoreKey)
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
func newMemPCApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *memoryPCApp {
	app := newMemoryPCBaseApp(logger, db, baseAppOptions...)
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keys[params.StoreKey], app.tkeys[params.TStoreKey], params.DefaultCodespace)
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keys[auth.StoreKey],
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		app.ModuleAccountAddrs(),
	)
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		app.keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		memoryModAccPerms,
	)
	app.nodesKeeper = nodesKeeper.NewKeeper(
		app.cdc,
		app.keys[nodesTypes.StoreKey],
		app.bankKeeper,
		app.supplyKeeper,
		app.paramsKeeper.Subspace(nodesTypes.DefaultParamspace),
		nodesTypes.DefaultCodespace,
	)
	app.appsKeeper = appsKeeper.NewKeeper(
		app.cdc,
		app.keys[appsTypes.StoreKey],
		app.bankKeeper,
		app.nodesKeeper,
		app.supplyKeeper,
		app.paramsKeeper.Subspace(appsTypes.DefaultParamspace),
		appsTypes.DefaultCodespace,
	)
	app.pocketKeeper = pocketKeeper.NewPocketCoreKeeper(
		app.keys[pocketTypes.StoreKey],
		app.cdc,
		app.nodesKeeper,
		app.appsKeeper,
		getInMemHostedChains(),
		app.paramsKeeper.Subspace(pocketTypes.DefaultParamspace),
		"test",
	)
	app.pocketKeeper.Keybase = inMemKeybase
	app.pocketKeeper.TmNode = getInMemoryTMClient()
	app.mm = module.NewManager(
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		nodes.NewAppModule(app.nodesKeeper, app.accountKeeper, app.supplyKeeper),
		apps.NewAppModule(app.appsKeeper, app.supplyKeeper, app.nodesKeeper),
		pocket.NewAppModule(app.pocketKeeper, app.nodesKeeper, app.appsKeeper),
	)
	app.mm.SetOrderBeginBlockers(nodesTypes.ModuleName, appsTypes.ModuleName, pocketTypes.ModuleName)
	app.mm.SetOrderEndBlockers(nodesTypes.ModuleName, appsTypes.ModuleName)
	app.mm.SetOrderInitGenesis(
		auth.ModuleName,
		bank.ModuleName,
		nodesTypes.ModuleName,
		appsTypes.ModuleName,
		pocketTypes.ModuleName,
		supply.ModuleName,
	)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

func (app *memoryPCApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	return app.mm.InitGenesis(ctx, genState)
}

func (app *memoryPCApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *memoryPCApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *memoryPCApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *memoryPCApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range memoryModAccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

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
	memoryModAccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		nodesTypes.StakedPoolName: {supply.Burner, supply.Staking},
		appsTypes.StakedPoolName:  {supply.Burner, supply.Staking},
		nodesTypes.DAOPoolName:    {supply.Burner, supply.Staking},
		nodesTypes.ModuleName:     nil,
		appsTypes.ModuleName:      nil,
	}
	genState     cfg.GenesisState
	inMemKeybase keys.Keybase
	memCDC       *codec.Codec
)

func inMemTendermintNode() (*node.Node, keys.Keybase) {
	pk := ed25519.GenPrivKey()
	kb := keys.NewInMemory()
	inMemKeybase = kb
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
			AppState:   memGenesisState(kp.PublicKey),
		}, nil
	}
	err = kb.SetCoinbase(kp.GetAddress())
	if err != nil {
		panic(err)
	}
	loggerFile, err := os.Open(os.DevNull)
	c := config{
		TmConfig: getTestConfig(),
		Logger:   log.NewTMLogger(log.NewSyncWriter(loggerFile)),
	}
	db := dbm.NewMemDB()
	traceWriter, err := openTraceWriter(c.TraceWriter)
	if err != nil {
		panic(err)
	}
	nodeKey := p2p.NodeKey{PrivKey: pk}
	privVal := cfg.GenFilePV(c.TmConfig.PrivValidatorKey, c.TmConfig.PrivValidatorState)
	privVal.Key.PrivKey = pk
	privVal.Key.PubKey = pk.PubKey()
	privVal.Key.Address = pk.PubKey().Address()
	creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *memoryPCApp {
		return newMemPCApp(logger, db)
	}
	upgradePrivVal(c.TmConfig)
	dbProvider := func(*node.DBContext) (dbm.DB, error) {
		return db, nil
	}
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

func memGenesisState(pubKey crypto.PublicKey) []byte {
	defaultGenesis := module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).DefaultGenesis()
	// set coinbase as a validator
	rawPOS := defaultGenesis[nodesTypes.ModuleName]
	var posGenesisState nodesTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawPOS, &posGenesisState)
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey.Address()),
			PublicKey:    pubKey,
			Status:       sdk.Bonded,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	res := memCodec().MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res
	genState = defaultGenesis
	// set coinbase as account holding coins
	rawAccounts := defaultGenesis[auth.ModuleName]
	var authGenState auth.GenesisState
	memCodec().MustUnmarshalJSON(rawAccounts, &authGenState)
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address:       sdk.Address(pubKey.Address()),
		Coins:         sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000))),
		PubKey:        pubKey,
		AccountNumber: 0,
		Sequence:      0,
	})
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
}

func memCodec() *codec.Codec {
	if memCDC == nil {
		makeMemCodec()
	}
	return memCDC
}

func makeMemCodec() {
	memCDC = codec.New()
	module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).RegisterCodec(memCDC)
	sdk.RegisterCodec(memCDC)
	codec.RegisterCrypto(memCDC)
}

func getInMemoryTMClient() client.Client {
	return client.NewHTTP(defaultTMURI, "/websocket")
}

func getInMemHostedChains() pocketTypes.HostedBlockchains {
	return pocketTypes.HostedBlockchains{
		M: map[string]pocketTypes.HostedBlockchain{dummyChainsHash: {Hash: dummyChainsHash, URL: dummyChainsURL}},
	}
}

func getTestConfig() ( tmConfg *tmCfg.Config){
	tmConfg = tmCfg.TestConfig()
	return
}
