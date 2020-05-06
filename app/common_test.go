package app

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

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
	authKeeper "github.com/pokt-network/posmint/x/auth/keeper"
	"github.com/pokt-network/posmint/x/gov"
	govKeeper "github.com/pokt-network/posmint/x/gov/keeper"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
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
)

const (
	dummyChainsHash = "00"
)

func NewInMemoryTendermintNode(t *testing.T, genesisState []byte) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
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
	// init cache in memory
	pocketTypes.InitCache("data", "data", dbm.MemDBBackend, dbm.MemDBBackend, 100, 100, "pocket_evidence", "session")
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
		err = os.RemoveAll("data")
		if err != nil {
			panic(err)
		}
		pocketTypes.ClearEvidence()
		pocketTypes.ClearSessionCache()
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
		nodesTypes.StakedPoolName: {auth.Burner, auth.Minter, auth.Staking},
		appsTypes.StakedPoolName:  {auth.Burner, auth.Minter, auth.Staking},
		nodesTypes.ModuleName:     {auth.Burner, auth.Minter, auth.Staking},
		govTypes.DAOAccountName:   {auth.Burner, auth.Staking, auth.Minter},
		appsTypes.ModuleName:      nil,
	}
	genState cfg.GenesisState
	memCDC   *codec.Codec
	inMemKB  keys.Keybase
	memCLI   client.Client
)

type memoryPCApp struct {
	*bam.BaseApp
	cdc           *codec.Codec
	keys          map[string]*sdk.KVStoreKey
	tkeys         map[string]*sdk.TransientStoreKey
	accountKeeper authKeeper.Keeper
	appsKeeper    appsKeeper.Keeper
	nodesKeeper   nodesKeeper.Keeper
	govKeeper     govKeeper.Keeper
	pocketKeeper  pocketKeeper.Keeper
	mm            *module.Manager
}

// NewPocketCoreApp is a constructor function for pocketCoreApp
func newMemPCApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *memoryPCApp {
	app := newMemoryPCBaseApp(logger, db, baseAppOptions...)
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
		getInMemHostedChains(),
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
	app.pocketKeeper.Keybase = getInMemoryKeybase()
	app.pocketKeeper.TmNode = getInMemoryTMClient()
	app.mm = module.NewManager(
		auth.NewAppModule(app.accountKeeper),
		nodes.NewAppModule(app.nodesKeeper),
		apps.NewAppModule(app.appsKeeper),
		pocket.NewAppModule(app.pocketKeeper),
		gov.NewAppModule(app.govKeeper),
	)
	app.mm.SetOrderBeginBlockers(nodesTypes.ModuleName, appsTypes.ModuleName, pocketTypes.ModuleName)
	app.mm.SetOrderEndBlockers(nodesTypes.ModuleName, appsTypes.ModuleName)
	app.mm.SetOrderInitGenesis(
		auth.ModuleName,
		nodesTypes.ModuleName,
		appsTypes.ModuleName,
		pocketTypes.ModuleName,
		gov.ModuleName,
	)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
	app.SetInitChainer(app.InitChainer)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper))
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
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(memCodec()), options...)
	bApp.SetAppVersion(AppVersion)
	k := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, nodesTypes.StoreKey, appsTypes.StoreKey, gov.StoreKey, pocketTypes.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(nodesTypes.TStoreKey, appsTypes.TStoreKey, pocketTypes.TStoreKey, gov.TStoreKey)
	// Create the application
	return &memoryPCApp{
		BaseApp: bApp,
		cdc:     memCodec(),
		keys:    k,
		tkeys:   tkeys,
	}
}

func (app *memoryPCApp) InitChainer(ctx sdk.Ctx, req abci.RequestInitChain) abci.ResponseInitChain {
	return app.mm.InitGenesis(ctx, genState)
}

func (app *memoryPCApp) BeginBlocker(ctx sdk.Ctx, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *memoryPCApp) EndBlocker(ctx sdk.Ctx, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *memoryPCApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *memoryPCApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range memoryModAccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *memoryPCApp) ExportAppState(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, err error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()}).WithAppVersion("0.0.0")
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

func getRandomPrivateKey() crypto.Ed25519PrivateKey {
	return crypto.Ed25519PrivateKey{}.GenPrivateKey().(crypto.Ed25519PrivateKey)
}

func inMemTendermintNode(genesisState []byte) (*node.Node, keys.Keybase) {
	// TODO add a second validator
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
		Logger:   log.NewTMLogger(loggerFile),
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
	pocketTypes.InitPVKeyFile(privVal.Key)

	creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *memoryPCApp {
		return newMemPCApp(logger, db, bam.SetPruning(store.PruneNothing))
	}
	upgradePrivVal(c.TmConfig)
	dbProvider := func(*node.DBContext) (dbm.DB, error) {
		return db, nil
	}
	app := creator(c.Logger, db, traceWriter)
	tmNode, err := node.NewNode(
		c.TmConfig,
		privVal,
		&nodeKey,
		proxy.NewLocalClientCreator(app),
		genDocProvider,
		dbProvider,
		node.DefaultMetricsProvider(c.TmConfig.Instrumentation),
		c.Logger.With("module", "node"),
	)
	if err != nil {
		panic(err)
	}
	testPCA = app
	app.SetTendermintNode(tmNode)
	return tmNode, kb
}

var testPCA *memoryPCApp

func memCodec() *codec.Codec {

	if memCDC == nil {
		memCDC = codec.New()
		module.NewBasicManager(
			apps.AppModuleBasic{},
			auth.AppModuleBasic{},
			gov.AppModuleBasic{},
			nodes.AppModuleBasic{},
			pocket.AppModuleBasic{},
		).RegisterCodec(memCDC)
		sdk.RegisterCodec(memCDC)
		codec.RegisterCrypto(memCDC)
	}
	return memCDC
}

func getInMemoryTMClient() client.Client {
	if memCLI == nil || !memCLI.IsRunning() {
		memCLI = client.NewHTTP(tmCfg.TestConfig().RPC.ListenAddress, "/websocket")
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

func getInMemHostedChains() *pocketTypes.HostedBlockchains {
	return &pocketTypes.HostedBlockchains{
		M: map[string]pocketTypes.HostedBlockchain{dummyChainsHash: {ID: dummyChainsHash, URL: PlaceholderURL}},
	}
}

func getTestConfig() (newTMConfig *tmCfg.Config) {
	newTMConfig = tmCfg.DefaultConfig()
	// setup tendermint node config
	newTMConfig.SetRoot("data")
	newTMConfig.FastSyncMode = false
	newTMConfig.NodeKey = "data" + FS + DefaultNKName
	newTMConfig.PrivValidatorKey = "data" + FS + DefaultPVKName
	newTMConfig.PrivValidatorState = "data" + FS + DefaultPVSName
	newTMConfig.RPC.ListenAddress = DefaultListenAddr + "36657"
	newTMConfig.P2P.ListenAddress = DefaultListenAddr + "36656" // Node listen address. (0.0.0.0:0 means any interface, any port)
	newTMConfig.Consensus = tmCfg.TestConsensusConfig()
	newTMConfig.Consensus.CreateEmptyBlocks = true // Set this to false to only produce blocks when there are txs or when the AppHash changes
	newTMConfig.Consensus.SkipTimeoutCommit = false
	newTMConfig.Consensus.CreateEmptyBlocksInterval = time.Duration(10) * time.Millisecond
	newTMConfig.Consensus.TimeoutCommit = time.Duration(10) * time.Millisecond
	newTMConfig.P2P.MaxNumInboundPeers = 40
	newTMConfig.P2P.MaxNumOutboundPeers = 10
	pocketTypes.InitClientBlockAllowance(10000)
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
		gov.AppModuleBasic{},
		nodes.AppModuleBasic{},
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
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(1000000000000000)})
	res := memCodec().MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res

	// setup application
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
	res3 := memCodec().MustMarshalJSON(authGenState)
	defaultGenesis[auth.ModuleName] = res3
	// set default chain for module
	rawPocket := defaultGenesis[pocketTypes.ModuleName]
	var pocketGenesisState pocketTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawPocket, &pocketGenesisState)
	pocketGenesisState.Params.SupportedBlockchains = []string{"00"}
	res4 := memCodec().MustMarshalJSON(pocketGenesisState)
	defaultGenesis[pocketTypes.ModuleName] = res4
	// set default governance in genesis
	var govGenesisState govTypes.GenesisState
	rawGov := defaultGenesis[govTypes.ModuleName]
	memCodec().MustUnmarshalJSON(rawGov, &govGenesisState)
	nMACL := createTestACL(kp1)
	govGenesisState.Params.Upgrade = govTypes.NewUpgrade(10000, "2.0.0")
	govGenesisState.Params.ACL = nMACL
	govGenesisState.Params.DAOOwner = kp1.GetAddress()
	govGenesisState.DAOTokens = sdk.NewInt(1000)
	res5 := memCodec().MustMarshalJSON(govGenesisState)
	defaultGenesis[govTypes.ModuleName] = res5
	// end genesis setup
	genState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
}

var testACL govTypes.ACL

func resetTestACL() {
	testACL = nil
}

func createTestACL(kp keys.KeyPair) govTypes.ACL {
	if testACL == nil {
		acl := govTypes.ACL{}
		acl = make([]govTypes.ACLPair, 0)
		acl.SetOwner("auth/MaxMemoCharacters", kp.GetAddress())
		acl.SetOwner("auth/TxSigLimit", kp.GetAddress())
		acl.SetOwner("gov/daoOwner", kp.GetAddress())
		acl.SetOwner("gov/acl", kp.GetAddress())
		acl.SetOwner("pos/StakeDenom", kp.GetAddress())
		acl.SetOwner("pocketcore/SupportedBlockchains", kp.GetAddress())
		acl.SetOwner("pos/DowntimeJailDuration", kp.GetAddress())
		acl.SetOwner("pos/SlashFractionDoubleSign", kp.GetAddress())
		acl.SetOwner("pos/SlashFractionDowntime", kp.GetAddress())
		acl.SetOwner("application/ApplicationStakeMinimum", kp.GetAddress())
		acl.SetOwner("pocketcore/ClaimExpiration", kp.GetAddress())
		acl.SetOwner("pocketcore/SessionNodeCount", kp.GetAddress())
		acl.SetOwner("pocketcore/ReplayAttackBurnMultiplier", kp.GetAddress())
		acl.SetOwner("pos/MaxValidators", kp.GetAddress())
		acl.SetOwner("pos/ProposerPercentage", kp.GetAddress())
		acl.SetOwner("application/StabilityAdjustment", kp.GetAddress())
		acl.SetOwner("application/AppUnstakingTime", kp.GetAddress())
		acl.SetOwner("application/ParticipationRateOn", kp.GetAddress())
		acl.SetOwner("pos/MaxEvidenceAge", kp.GetAddress())
		acl.SetOwner("pos/MinSignedPerWindow", kp.GetAddress())
		acl.SetOwner("pos/StakeMinimum", kp.GetAddress())
		acl.SetOwner("pos/UnstakingTime", kp.GetAddress())
		acl.SetOwner("pos/RelaysToTokensMultiplier", kp.GetAddress())
		acl.SetOwner("application/BaseRelaysPerPOKT", kp.GetAddress())
		acl.SetOwner("pocketcore/ClaimSubmissionWindow", kp.GetAddress())
		acl.SetOwner("pos/DAOAllocation", kp.GetAddress())
		acl.SetOwner("pos/SignedBlocksWindow", kp.GetAddress())
		acl.SetOwner("pos/BlocksPerSession", kp.GetAddress())
		acl.SetOwner("application/MaxApplications", kp.GetAddress())
		acl.SetOwner("gov/daoOwner", kp.GetAddress())
		acl.SetOwner("gov/upgrade", kp.GetAddress())
		testACL = acl
	}
	return testACL
}

func twoValTwoNodeGenesisState() []byte {
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
		gov.AppModuleBasic{},
		nodes.AppModuleBasic{},
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
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(1000000000000000)})
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey2.Address()),
			PublicKey:    pubKey2,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(1000000000)})
	posGenesisState.Params.UnstakingTime = time.Nanosecond
	posGenesisState.Params.SessionBlockFrequency = 5
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
	// set default governance in genesis
	var govGenesisState govTypes.GenesisState
	rawGov := defaultGenesis[govTypes.ModuleName]
	memCodec().MustUnmarshalJSON(rawGov, &govGenesisState)
	nMACL := createTestACL(kp1)
	govGenesisState.Params.Upgrade = govTypes.NewUpgrade(10000, "2.0.0")
	govGenesisState.Params.ACL = nMACL
	govGenesisState.Params.DAOOwner = kp1.GetAddress()
	govGenesisState.DAOTokens = sdk.NewInt(1000)
	res4 := memCodec().MustMarshalJSON(govGenesisState)
	defaultGenesis[govTypes.ModuleName] = res4
	// end genesis setup
	genState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
}

func fiveValidatorsOneAppGenesis() (genBz []byte, keys []crypto.PrivateKey, validators nodesTypes.Validators, app appsTypes.Application) {
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
	pk1, err := kb.ExportPrivateKeyObject(kp1.GetAddress(), "test")
	if err != nil {
		panic(err)
	}
	pk2, err := kb.ExportPrivateKeyObject(kp2.GetAddress(), "test")
	var kys []crypto.PrivateKey
	kys = append(kys, pk1, pk2, crypto.GenerateEd25519PrivKey(), crypto.GenerateEd25519PrivKey(), crypto.GenerateEd25519PrivKey())
	// get public kys
	pubKey := kp1.PublicKey
	pubKey2 := kp2.PublicKey
	pubKey3 := kys[2].PublicKey()
	pubKey4 := kys[3].PublicKey()
	pubKey5 := kys[4].PublicKey()
	defaultGenesis := module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		gov.AppModuleBasic{},
		nodes.AppModuleBasic{},
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
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(1000000000000000000)})
	// validator 2
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey2.Address()),
			PublicKey:    pubKey2,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 3
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey3.Address()),
			PublicKey:    pubKey3,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 4
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey4.Address()),
			PublicKey:    pubKey4,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 5
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey5.Address()),
			PublicKey:    pubKey5,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   PlaceholderServiceURL,
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
	pocketGenesisState.Params.ClaimSubmissionWindow = 10
	res3 := memCodec().MustMarshalJSON(pocketGenesisState)
	defaultGenesis[pocketTypes.ModuleName] = res3
	// set default governance in genesis
	var govGenesisState govTypes.GenesisState
	rawGov := defaultGenesis[govTypes.ModuleName]
	memCodec().MustUnmarshalJSON(rawGov, &govGenesisState)
	nMACL := createTestACL(kp1)
	govGenesisState.Params.Upgrade = govTypes.NewUpgrade(10000, "2.0.0")
	govGenesisState.Params.ACL = nMACL
	govGenesisState.Params.DAOOwner = kp1.GetAddress()
	govGenesisState.DAOTokens = sdk.NewInt(1000)
	res4 := memCodec().MustMarshalJSON(govGenesisState)
	defaultGenesis[govTypes.ModuleName] = res4
	// end genesis setup
	genState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j, kys, posGenesisState.Validators, appsGenesisState.Applications[0]
}
