package app

import (
	"context"
	types2 "github.com/pokt-network/pocket-core/codec/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/client/local"
	tmStore "github.com/tendermint/tendermint/store"
	"io"
	"os"
	"testing"
	"time"

	bam "github.com/pokt-network/pocket-core/baseapp"
	"github.com/pokt-network/pocket-core/codec"
	cfg "github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/store"
	storeTypes "github.com/pokt-network/pocket-core/store/types"

	// sdk "github.com/pokt-network/pocket-core/types"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/gov"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmCfg "github.com/tendermint/tendermint/config"
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
	dummyChainsHash = "0001"
)

func BeforeEach(t *testing.T) {
	pocketTypes.ClearSessionCache()
	pocketTypes.ClearEvidence()
	sdk.GlobalCtxCache.Purge()
}

type upgrades struct {
	codecUpgrade
}
type codecUpgrade struct {
	upgradeMod bool
	height     int64
}

func NewInMemoryTendermintNodeAmino(t *testing.T, genesisState []byte) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	// create the in memory tendermint node and keybase
	tendermintNode, keybase = inMemTendermintNode(genesisState, false)
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
	pocketTypes.InitConfig(&pocketTypes.HostedBlockchains{
		M: make(map[string]pocketTypes.HostedBlockchain),
	}, tendermintNode.Logger, sdk.DefaultTestingPocketConfig())
	// start the in memory node
	err := tendermintNode.Start()
	if err != nil {
		panic(err)
	}
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
		PCA = nil
		inMemKB = nil
	}
	return
}
func NewInMemoryTendermintNodeProto(t *testing.T, genesisState []byte) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	// create the in memory tendermint node and keybase
	tendermintNode, keybase = inMemTendermintNode(genesisState, true)
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
	pocketTypes.InitConfig(&pocketTypes.HostedBlockchains{
		M: make(map[string]pocketTypes.HostedBlockchain),
	}, tendermintNode.Logger, sdk.DefaultTestingPocketConfig())
	// start the in memory node
	err := tendermintNode.Start()
	if err != nil {
		panic(err)
	}
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
		PCA = nil
		inMemKB = nil
	}
	return
}

func TestNewInMemoryAmino(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNodeAmino(t, oneValTwoNodeGenesisState())
	defer cleanup()
}
func TestNewInMemoryProto(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNodeProto(t, oneValTwoNodeGenesisState())
	defer cleanup()
}

var (
	memCDC  *codec.Codec
	inMemKB keys.Keybase
	memCLI  client.Client
)

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

func inMemTendermintNode(genesisState []byte, protoCodec bool) (*node.Node, keys.Keybase) {
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

	//creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *PocketCoreApp {
	//	m := map[string]pocketTypes.HostedBlockchain{"0001": {
	//		ID:  sdk.PlaceholderHash,
	//		URL: sdk.PlaceholderURL,
	//	}}
	//	p := NewPocketCoreApp(GenState, getInMemoryKeybase(), getInMemoryTMClient(), &pocketTypes.HostedBlockchains{M: m}, logger, db, bam.SetPruning(store.PruneNothing))
	//	return p
	//}
	dbProvider := func(*node.DBContext) (dbm.DB, error) {
		return db, nil
	}
	app := creator(c.Logger, db, traceWriter)
	tmNode, err := node.NewNode(app,
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
	PCA = app
	app.SetTxIndexer(tmNode.TxIndexer())
	app.SetBlockstore(tmNode.BlockStore())
	app.SetEvidencePool(tmNode.EvidencePool())
	app.pocketKeeper.TmNode = local.New(tmNode)
	app.SetTendermintNode(tmNode)
	return tmNode, kb
}

func GetApp(protoCdc bool, logger log.Logger, db dbm.DB, traceWriter io.Writer) *PocketCoreApp {
	creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *PocketCoreApp {
		m := map[string]pocketTypes.HostedBlockchain{"0001": {
			ID:  sdk.PlaceholderHash,
			URL: sdk.PlaceholderURL,
		}}
		p := NewPocketCoreApp(GenState, getInMemoryKeybase(), getInMemoryTMClient(), &pocketTypes.HostedBlockchains{M: m}, logger, db, bam.SetPruning(store.PruneNothing))
		return p
	}
	app := creator(logger, db, traceWriter)
	if protoCdc {
		ctx := new(Ctx)
		ctx.On("IsAfterUpgradeHeight").Return(true)
		app.appsKeeper.UpgradeCodec(ctx)
		app.accountKeeper.UpgradeCodec(ctx)
		app.govKeeper.UpgradeCodec(ctx)
		app.pocketKeeper.UpgradeCodec(ctx)
		app.nodesKeeper.UpgradeCodec(ctx)
	}
	return app

}

func memCodec() *codec.Codec {
	if memCDC == nil {
		memCDC = codec.NewCodec(types2.NewInterfaceRegistry())
		module.NewBasicManager(
			apps.AppModuleBasic{},
			auth.AppModuleBasic{},
			gov.AppModuleBasic{},
			nodes.AppModuleBasic{},
			pocket.AppModuleBasic{},
		).RegisterCodec(memCDC)
		sdk.RegisterCodec(memCDC)
		crypto.RegisterAmino(memCDC.AminoCodec().Amino)
	}
	return memCDC
}

func memCodecMod(upgrade bool) *codec.Codec {
	if memCDC == nil {
		memCDC = codec.NewCodec(types2.NewInterfaceRegistry())
		module.NewBasicManager(
			apps.AppModuleBasic{},
			auth.AppModuleBasic{},
			gov.AppModuleBasic{},
			nodes.AppModuleBasic{},
			pocket.AppModuleBasic{},
		).RegisterCodec(memCDC)
		sdk.RegisterCodec(memCDC)
		crypto.RegisterAmino(memCDC.AminoCodec().Amino)
	}
	memCDC.SetAfterUpgradeMod(upgrade)
	return memCDC
}

func getInMemoryTMClient() client.Client {
	if memCLI == nil || !memCLI.IsRunning() {
		memCLI, _ = http.New(tmCfg.TestConfig().RPC.ListenAddress, "/websocket")
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

func getTestConfig() (newTMConfig *tmCfg.Config) {
	newTMConfig = tmCfg.DefaultConfig()
	// setup tendermint node config
	newTMConfig.SetRoot("data")
	newTMConfig.FastSyncMode = false
	newTMConfig.NodeKey = "data" + FS + sdk.DefaultNKName
	newTMConfig.PrivValidatorKey = "data" + FS + sdk.DefaultPVKName
	newTMConfig.PrivValidatorState = "data" + FS + sdk.DefaultPVSName
	newTMConfig.RPC.ListenAddress = sdk.DefaultListenAddr + "36657"
	newTMConfig.P2P.ListenAddress = sdk.DefaultListenAddr + "36656" // Node listen address. (0.0.0.0:0 means any interface, any port)
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
			ServiceURL:   sdk.PlaceholderServiceURL,
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
	pocketGenesisState.Params.SupportedBlockchains = []string{"0001"}
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
	GenState = defaultGenesis
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
		acl.SetOwner("application/ApplicationStakeMinimum", kp.GetAddress())
		acl.SetOwner("application/AppUnstakingTime", kp.GetAddress())
		acl.SetOwner("application/BaseRelaysPerPOKT", kp.GetAddress())
		acl.SetOwner("application/MaxApplications", kp.GetAddress())
		acl.SetOwner("application/MaximumChains", kp.GetAddress())
		acl.SetOwner("application/ParticipationRateOn", kp.GetAddress())
		acl.SetOwner("application/StabilityAdjustment", kp.GetAddress())
		acl.SetOwner("auth/MaxMemoCharacters", kp.GetAddress())
		acl.SetOwner("auth/TxSigLimit", kp.GetAddress())
		acl.SetOwner("auth/FeeMultipliers", kp.GetAddress())
		acl.SetOwner("gov/acl", kp.GetAddress())
		acl.SetOwner("gov/daoOwner", kp.GetAddress())
		acl.SetOwner("gov/upgrade", kp.GetAddress())
		acl.SetOwner("pocketcore/ClaimExpiration", kp.GetAddress())
		acl.SetOwner("pocketcore/ClaimSubmissionWindow", kp.GetAddress())
		acl.SetOwner("pocketcore/MinimumNumberOfProofs", kp.GetAddress())
		acl.SetOwner("pocketcore/ReplayAttackBurnMultiplier", kp.GetAddress())
		acl.SetOwner("pocketcore/SessionNodeCount", kp.GetAddress())
		acl.SetOwner("pocketcore/SupportedBlockchains", kp.GetAddress())
		acl.SetOwner("pos/BlocksPerSession", kp.GetAddress())
		acl.SetOwner("pos/DAOAllocation", kp.GetAddress())
		acl.SetOwner("pos/DowntimeJailDuration", kp.GetAddress())
		acl.SetOwner("pos/MaxEvidenceAge", kp.GetAddress())
		acl.SetOwner("pos/MaximumChains", kp.GetAddress())
		acl.SetOwner("pos/MaxJailedBlocks", kp.GetAddress())
		acl.SetOwner("pos/MaxValidators", kp.GetAddress())
		acl.SetOwner("pos/MinSignedPerWindow", kp.GetAddress())
		acl.SetOwner("pos/ProposerPercentage", kp.GetAddress())
		acl.SetOwner("pos/RelaysToTokensMultiplier", kp.GetAddress())
		acl.SetOwner("pos/SignedBlocksWindow", kp.GetAddress())
		acl.SetOwner("pos/SlashFractionDoubleSign", kp.GetAddress())
		acl.SetOwner("pos/SlashFractionDowntime", kp.GetAddress())
		acl.SetOwner("pos/StakeDenom", kp.GetAddress())
		acl.SetOwner("pos/StakeMinimum", kp.GetAddress())
		acl.SetOwner("pos/UnstakingTime", kp.GetAddress())
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
			ServiceURL:   sdk.PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(1000000000000000)})
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey2.Address()),
			PublicKey:    pubKey2,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   sdk.PlaceholderServiceURL,
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
	GenState = defaultGenesis
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
	if err != nil {
		panic(err)
	}
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
			ServiceURL:   sdk.PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(1000000000000000000)})
	// validator 2
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey2.Address()),
			PublicKey:    pubKey2,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   sdk.PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 3
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey3.Address()),
			PublicKey:    pubKey3,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   sdk.PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 4
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey4.Address()),
			PublicKey:    pubKey4,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   sdk.PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	// validator 5
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey5.Address()),
			PublicKey:    pubKey5,
			Status:       sdk.Staked,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   sdk.PlaceholderServiceURL,
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
	GenState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j, kys, posGenesisState.Validators, appsGenesisState.Applications[0]
}

type Ctx struct {
	mock.Mock
}

// BlockGasMeter provides a mock function with given fields:
func (_m *Ctx) BlockGasMeter() storeTypes.GasMeter {
	ret := _m.Called()

	var r0 storeTypes.GasMeter
	if rf, ok := ret.Get(0).(func() storeTypes.GasMeter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.GasMeter)
		}
	}

	return r0
}

// BlockHeader provides a mock function with given fields:
func (_m *Ctx) BlockHeader() abcitypes.Header {
	ret := _m.Called()

	var r0 abcitypes.Header
	if rf, ok := ret.Get(0).(func() abcitypes.Header); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(abcitypes.Header)
	}

	return r0
}

// BlockHeight provides a mock function with given fields:
func (_m *Ctx) BlockHeight() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// BlockStore provides a mock function with given fields:
func (_m *Ctx) BlockStore() *tmStore.BlockStore {
	ret := _m.Called()

	var r0 *tmStore.BlockStore
	if rf, ok := ret.Get(0).(func() *tmStore.BlockStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tmStore.BlockStore)
		}
	}

	return r0
}

// BlockTime provides a mock function with given fields:
func (_m *Ctx) BlockTime() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// CacheContext provides a mock function with given fields:
func (_m *Ctx) CacheContext() (sdk.Context, func()) {
	ret := _m.Called()

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func() sdk.Context); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	var r1 func()
	if rf, ok := ret.Get(1).(func() func()); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	return r0, r1
}

// ChainID provides a mock function with given fields:
func (_m *Ctx) ChainID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ConsensusParams provides a mock function with given fields:
func (_m *Ctx) ConsensusParams() *abcitypes.ConsensusParams {
	ret := _m.Called()

	var r0 *abcitypes.ConsensusParams
	if rf, ok := ret.Get(0).(func() *abcitypes.ConsensusParams); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcitypes.ConsensusParams)
		}
	}

	return r0
}

// Context provides a mock function with given fields:
func (_m *Ctx) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// EventManager provides a mock function with given fields:
func (_m *Ctx) EventManager() *sdk.EventManager {
	ret := _m.Called()

	var r0 *sdk.EventManager
	if rf, ok := ret.Get(0).(func() *sdk.EventManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sdk.EventManager)
		}
	}

	return r0
}

// GasMeter provides a mock function with given fields:
func (_m *Ctx) GasMeter() storeTypes.GasMeter {
	ret := _m.Called()

	var r0 storeTypes.GasMeter
	if rf, ok := ret.Get(0).(func() storeTypes.GasMeter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.GasMeter)
		}
	}

	return r0
}

// IsCheckTx provides a mock function with given fields:
func (_m *Ctx) IsCheckTx() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsZero provides a mock function with given fields:
func (_m *Ctx) IsZero() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsZero provides a mock function with given fields:
func (_m *Ctx) IsAfterUpgradeHeight() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// KVStore provides a mock function with given fields: key
func (_m *Ctx) KVStore(key storeTypes.StoreKey) storeTypes.KVStore {
	ret := _m.Called(key)

	var r0 storeTypes.KVStore
	if rf, ok := ret.Get(0).(func(storeTypes.StoreKey) storeTypes.KVStore); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.KVStore)
		}
	}

	return r0
}

// Logger provides a mock function with given fields:
func (_m *Ctx) Logger() log.Logger {
	ret := _m.Called()

	var r0 log.Logger
	if rf, ok := ret.Get(0).(func() log.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.Logger)
		}
	}

	return r0
}

// MinGasPrices provides a mock function with given fields:
func (_m *Ctx) MinGasPrices() sdk.DecCoins {
	ret := _m.Called()

	var r0 sdk.DecCoins
	if rf, ok := ret.Get(0).(func() sdk.DecCoins); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sdk.DecCoins)
		}
	}

	return r0
}

// MultiStore provides a mock function with given fields:
func (_m *Ctx) MultiStore() storeTypes.MultiStore {
	ret := _m.Called()

	var r0 storeTypes.MultiStore
	if rf, ok := ret.Get(0).(func() storeTypes.MultiStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.MultiStore)
		}
	}

	return r0
}

// MustGetPrevCtx provides a mock function with given fields: height
func (_m *Ctx) MustGetPrevCtx(height int64) sdk.Context {
	ret := _m.Called(height)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(int64) sdk.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// PrevCtx provides a mock function with given fields: height
func (_m *Ctx) PrevCtx(height int64) (sdk.Context, error) {
	ret := _m.Called(height)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(int64) sdk.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransientStore provides a mock function with given fields: key
func (_m *Ctx) TransientStore(key storeTypes.StoreKey) storeTypes.KVStore {
	ret := _m.Called(key)

	var r0 storeTypes.KVStore
	if rf, ok := ret.Get(0).(func(storeTypes.StoreKey) storeTypes.KVStore); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.KVStore)
		}
	}

	return r0
}

// TxBytes provides a mock function with given fields:
func (_m *Ctx) TxBytes() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// Value provides a mock function with given fields: key
func (_m *Ctx) Value(key interface{}) interface{} {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// VoteInfos provides a mock function with given fields:
func (_m *Ctx) VoteInfos() []abcitypes.VoteInfo {
	ret := _m.Called()

	var r0 []abcitypes.VoteInfo
	if rf, ok := ret.Get(0).(func() []abcitypes.VoteInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]abcitypes.VoteInfo)
		}
	}

	return r0
}

// WithBlockGasMeter provides a mock function with given fields: meter
func (_m *Ctx) WithBlockGasMeter(meter storeTypes.GasMeter) sdk.Context {
	ret := _m.Called(meter)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(storeTypes.GasMeter) sdk.Context); ok {
		r0 = rf(meter)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithBlockHeader provides a mock function with given fields: header
func (_m *Ctx) WithBlockHeader(header abcitypes.Header) sdk.Context {
	ret := _m.Called(header)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(abcitypes.Header) sdk.Context); ok {
		r0 = rf(header)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithBlockHeight provides a mock function with given fields: height
func (_m *Ctx) WithBlockHeight(height int64) sdk.Context {
	ret := _m.Called(height)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(int64) sdk.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithBlockStore provides a mock function with given fields: bs
func (_m *Ctx) WithBlockStore(bs *tmStore.BlockStore) sdk.Context {
	ret := _m.Called(bs)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(*tmStore.BlockStore) sdk.Context); ok {
		r0 = rf(bs)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithBlockTime provides a mock function with given fields: newTime
func (_m *Ctx) WithBlockTime(newTime time.Time) sdk.Context {
	ret := _m.Called(newTime)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(time.Time) sdk.Context); ok {
		r0 = rf(newTime)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithChainID provides a mock function with given fields: chainID
func (_m *Ctx) WithChainID(chainID string) sdk.Context {
	ret := _m.Called(chainID)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(string) sdk.Context); ok {
		r0 = rf(chainID)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithConsensusParams provides a mock function with given fields: params
func (_m *Ctx) WithConsensusParams(params *abcitypes.ConsensusParams) sdk.Context {
	ret := _m.Called(params)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(*abcitypes.ConsensusParams) sdk.Context); ok {
		r0 = rf(params)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithContext provides a mock function with given fields: ctx
func (_m *Ctx) WithContext(ctx context.Context) sdk.Context {
	ret := _m.Called(ctx)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(context.Context) sdk.Context); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithEventManager provides a mock function with given fields: em
func (_m *Ctx) WithEventManager(em *sdk.EventManager) sdk.Context {
	ret := _m.Called(em)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(*sdk.EventManager) sdk.Context); ok {
		r0 = rf(em)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithGasMeter provides a mock function with given fields: meter
func (_m *Ctx) WithGasMeter(meter storeTypes.GasMeter) sdk.Context {
	ret := _m.Called(meter)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(storeTypes.GasMeter) sdk.Context); ok {
		r0 = rf(meter)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithIsCheckTx provides a mock function with given fields: isCheckTx
func (_m *Ctx) WithIsCheckTx(isCheckTx bool) sdk.Context {
	ret := _m.Called(isCheckTx)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(bool) sdk.Context); ok {
		r0 = rf(isCheckTx)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithLogger provides a mock function with given fields: logger
func (_m *Ctx) WithLogger(logger log.Logger) sdk.Context {
	ret := _m.Called(logger)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(log.Logger) sdk.Context); ok {
		r0 = rf(logger)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithMinGasPrices provides a mock function with given fields: gasPrices
func (_m *Ctx) WithMinGasPrices(gasPrices sdk.DecCoins) sdk.Context {
	ret := _m.Called(gasPrices)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(sdk.DecCoins) sdk.Context); ok {
		r0 = rf(gasPrices)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithMultiStore provides a mock function with given fields: ms
func (_m *Ctx) WithMultiStore(ms storeTypes.MultiStore) sdk.Context {
	ret := _m.Called(ms)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(storeTypes.MultiStore) sdk.Context); ok {
		r0 = rf(ms)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithProposer provides a mock function with given fields: addr
func (_m *Ctx) WithProposer(addr sdk.Address) sdk.Context {
	ret := _m.Called(addr)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(sdk.Address) sdk.Context); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithTxBytes provides a mock function with given fields: txBytes
func (_m *Ctx) WithTxBytes(txBytes []byte) sdk.Context {
	ret := _m.Called(txBytes)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func([]byte) sdk.Context); ok {
		r0 = rf(txBytes)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithValue provides a mock function with given fields: key, value
func (_m *Ctx) WithValue(key interface{}, value interface{}) sdk.Context {
	ret := _m.Called(key, value)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) sdk.Context); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

// WithValue provides a mock function with given fields: key, value
func (_m *Ctx) AppVersion() string {
	return ""
}

// WithVoteInfos provides a mock function with given fields: voteInfo
func (_m *Ctx) WithVoteInfos(voteInfo []abcitypes.VoteInfo) sdk.Context {
	ret := _m.Called(voteInfo)

	var r0 sdk.Context
	if rf, ok := ret.Get(0).(func([]abcitypes.VoteInfo) sdk.Context); ok {
		r0 = rf(voteInfo)
	} else {
		r0 = ret.Get(0).(sdk.Context)
	}

	return r0
}

func (_m *Ctx) BlockHash(cdc *codec.Codec) ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).([]byte)
	}

	var r1 error
	if rf1, ok := ret.Get(1).(func() error); ok {
		r1 = rf1()
	} else {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
