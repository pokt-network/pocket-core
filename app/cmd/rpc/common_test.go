package rpc

import (
	"context"
	"encoding/json"
	"github.com/pokt-network/pocket-core/app"
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
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	tmCfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/rpc/client"
	cTypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"testing"
	"time"
)

func NewInMemoryTendermintNode(t *testing.T, genesisState []byte) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	app.MakeCodec() // needed for queries and tx
	// create the in memory tendermint node and keybase
	tendermintNode, keybase = inMemTendermintNode(genesisState)
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
		inMemKB = nil
		return
	}
	return
}

func TestNewInMemory(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	defer cleanup()
}

var (
	memoryModAccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		nodesTypes.StakedPoolName: {supply.Burner, supply.Staking, supply.Minter},
		appsTypes.StakedPoolName:  {supply.Burner, supply.Staking, supply.Minter},
		nodesTypes.DAOPoolName:    {supply.Burner, supply.Staking, supply.Minter},
		nodesTypes.ModuleName:     {supply.Burner, supply.Staking, supply.Minter},
		appsTypes.ModuleName:      nil,
	}
	genState cfg.GenesisState
	memCDC   *codec.Codec
	inMemKB  keys.Keybase
	memCLI   client.Client
)

const (
	dummyChainsHash = "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"
	dummyChainsURL  = "https://foo.bar:8080"
	dummyServiceURL = "0.0.0.0:8081"
	defaultTMURI    = "tcp://localhost:26657"
)

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
		app.accountKeeper,
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
	app.pocketKeeper.Keybase = getInMemoryKeybase()
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
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper))
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

func newMemoryPCBaseApp(logger log.Logger, db dbm.DB, options ...func(*bam.BaseApp)) *memoryPCApp {
	bApp := bam.NewBaseApp("pocket-test", logger, db, auth.DefaultTxDecoder(memCodec()), options...)
	bApp.SetAppVersion("0.0.0")
	k := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.StoreKey, supply.StoreKey, params.StoreKey, pocketTypes.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(nodesTypes.TStoreKey, appsTypes.TStoreKey, pocketTypes.TStoreKey, params.TStoreKey)
	// Create the application
	return &memoryPCApp{
		BaseApp: bApp,
		cdc:     memCodec(),
		keys:    k,
		tkeys:   tkeys,
	}
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

func getInMemoryKeybase() keys.Keybase {
	if inMemKB == nil {
		inMemKB = keys.NewInMemory()
		_, err := inMemKB.Create("test")
		if err != nil {
			panic(err)
		}
		_, err = inMemKB.GetCoinbase()
		if err != nil {
			panic(err)
		}
	}
	return inMemKB
}

func inMemTendermintNode(genesisState []byte) (*node.Node, keys.Keybase) {
	kb := getInMemoryKeybase()
	cb, err := kb.GetCoinbase()
	if err != nil {
		panic(err)
	}
	pk, err := kb.ExportPrivateKeyObject(cb.GetAddress(), "test")
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
			AppState:   genesisState,
		}, nil
	}
	loggerFile, _ := os.Open(os.DevNull)
	c := config{
		TmConfig: getTestConfig(),
		Logger:   log.NewTMLogger(log.NewSyncWriter(loggerFile)),
	}
	db := dbm.NewMemDB()
	nodeKey := p2p.NodeKey{PrivKey: pk}
	privVal := cfg.GenFilePV(c.TmConfig.PrivValidatorKey, c.TmConfig.PrivValidatorState)
	privVal.Key.PrivKey = pk
	privVal.Key.PubKey = pk.PubKey()
	privVal.Key.Address = pk.PubKey().Address()
	creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *memoryPCApp {
		return newMemPCApp(logger, db, bam.SetPruning(store.PruneNothing))
	}
	//upgradePrivVal(c.TmConfig)
	dbProvider := func(*node.DBContext) (dbm.DB, error) {
		return db, nil
	}
	baseapp := creator(c.Logger, db, io.Writer(nil))
	tmNode, err := node.NewNode(
		c.TmConfig,
		privVal,
		&nodeKey,
		proxy.NewLocalClientCreator(baseapp),
		genDocProvider,
		dbProvider,
		node.DefaultMetricsProvider(c.TmConfig.Instrumentation),
		c.Logger.With("module", "node"),
	)
	if err != nil {
		panic(err)
	}
	baseapp.SetTendermintNode(tmNode)
	return tmNode, kb
}

func memCodec() *codec.Codec {
	if memCDC == nil {
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
	return memCDC
}

func getInMemoryTMClient() client.Client {
	if memCLI == nil || !memCLI.IsRunning() {
		memCLI = client.NewHTTP(defaultTMURI, "/websocket")
	}
	return memCLI
}

func subscribeTo(t *testing.T, eventType string) (cli client.Client, stopClient func(), eventChan <-chan cTypes.ResultEvent) {
	ctx, cancel := getBackgroundContext()
	cli = getInMemoryTMClient()
	if !cli.IsRunning() {
		_ = cli.Start()
	}
	stopClient = func() {
		err := cli.UnsubscribeAll(ctx, "helpers")
		if err != nil {
			t.Fatal(err)
		}
		err = cli.Stop()
		if err != nil {
			t.Fatal(err)
		}
		memCLI = nil
		cancel()
	}
	eventChan, err := cli.Subscribe(ctx, "helpers", types.QueryForEvent(eventType).String())
	if err != nil {
		panic(err)
	}
	return
}

func getBackgroundContext() (context.Context, func()) {
	return context.WithCancel(context.Background())
}

func getInMemHostedChains() pocketTypes.HostedBlockchains {
	return pocketTypes.HostedBlockchains{
		M: map[string]pocketTypes.HostedBlockchain{dummyChainsHash: {Hash: dummyChainsHash, URL: dummyChainsURL}},
	}
}

func getTestConfig() (tmConfg *tmCfg.Config) {
	tmConfg = tmCfg.TestConfig()
	tmConfg.RPC.ListenAddress = defaultTMURI
	return
}

func getUnstakedAccount(kb keys.Keybase) *keys.KeyPair {
	cb, err := kb.GetCoinbase()
	if err != nil {
		panic(err)
	}
	kps, err := kb.List()
	if err != nil {
		panic(err)
	}
	if len(kps) > 2 {
		panic("get unstaked account only works with the default 2 keypairs")
	}
	for _, kp := range kps {
		if kp.PublicKey != cb.PublicKey {
			return &kp
		}
	}
	return nil
}

func oneValTwoNodeGenesisState() []byte {
	kb := getInMemoryKeybase()
	kp1, err := kb.GetCoinbase()
	if err != nil {
		panic(err)
	}
	kp2, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	pubKey := kp1.PublicKey
	pubKey2 := kp2.PublicKey
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
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(1000000000000000)})
	res := memCodec().MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res
	// set coinbase as account holding coins
	rawAccounts := defaultGenesis[auth.ModuleName]
	var authGenState auth.GenesisState
	memCodec().MustUnmarshalJSON(rawAccounts, &authGenState)
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address: sdk.Address(pubKey.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000000))),
		PubKey:  pubKey,
	})
	// add second account
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address: sdk.Address(pubKey2.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000000))),
		PubKey:  pubKey,
	})
	res2 := memCodec().MustMarshalJSON(authGenState)
	defaultGenesis[auth.ModuleName] = res2
	// set default chain for module
	rawPocket := defaultGenesis[pocketTypes.ModuleName]
	var pocketGenesisState pocketTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawPocket, &pocketGenesisState)
	pocketGenesisState.Params.SupportedBlockchains = []string{dummyChainsHash}
	res3 := memCodec().MustMarshalJSON(pocketGenesisState)
	defaultGenesis[pocketTypes.ModuleName] = res3
	genState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
}

func fiveValidatorsOneAppGenesis() (genBz []byte, validators nodesTypes.Validators, app appsTypes.Application) {
	kb := getInMemoryKeybase()
	// create keypairs
	kp1, err := kb.GetCoinbase()
	if err != nil {
		panic(err)
	}
	kp2, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	// get public keys
	pubKey := kp1.PublicKey
	pubKey2 := crypto.GenerateEd25519PrivKey().PublicKey()
	pubKey3 := crypto.GenerateEd25519PrivKey().PublicKey()
	pubKey4 := crypto.GenerateEd25519PrivKey().PublicKey()
	pubKey5 := crypto.GenerateEd25519PrivKey().PublicKey()
	defaultGenesis := module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).DefaultGenesis()
	// setup validators
	rawPOS := defaultGenesis[nodesTypes.ModuleName]
	var posGenesisState nodesTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawPOS, &posGenesisState)
	// validator 1
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey.Address()),
			PublicKey:    pubKey,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(1000000000000000000)})
	// validator 2
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey2.Address()),
			PublicKey:    pubKey2,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 3
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey3.Address()),
			PublicKey:    pubKey3,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 4
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey4.Address()),
			PublicKey:    pubKey4,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 5
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey5.Address()),
			PublicKey:    pubKey5,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// marshal into json
	res := memCodec().MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res
	// setup applications
	rawApps := defaultGenesis[appsTypes.ModuleName]
	var appsGenesisState appsTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawApps, &appsGenesisState)
	// app 1
	appsGenesisState.Applications = append(appsGenesisState.Applications, appsTypes.Application{
		Address:                 kp2.GetAddress(),
		PublicKey:               kp2.PublicKey,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{dummyChainsHash},
		StakedTokens:            sdk.NewInt(10000000),
		MaxRelays:               sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	})
	res2 := memCodec().MustMarshalJSON(appsGenesisState)
	defaultGenesis[appsTypes.ModuleName] = res2
	// accounts
	rawAccounts := defaultGenesis[auth.ModuleName]
	var authGenState auth.GenesisState
	memCodec().MustUnmarshalJSON(rawAccounts, &authGenState)
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address: sdk.Address(pubKey.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000000))),
		PubKey:  pubKey,
	})
	res = memCodec().MustMarshalJSON(authGenState)
	defaultGenesis[auth.ModuleName] = res
	// setup supported blockchains
	rawPocket := defaultGenesis[pocketTypes.ModuleName]
	var pocketGenesisState pocketTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawPocket, &pocketGenesisState)
	pocketGenesisState.Params.SupportedBlockchains = []string{dummyChainsHash}
	res3 := memCodec().MustMarshalJSON(pocketGenesisState)
	defaultGenesis[pocketTypes.ModuleName] = res3
	genState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j, posGenesisState.Validators, appsGenesisState.Applications[0]
}

type config struct {
	TmConfig    *tmCfg.Config
	Logger      log.Logger
	TraceWriter string
}
