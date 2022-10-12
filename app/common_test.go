package app

import (
	"context"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/tendermint/tendermint/privval"

	bam "github.com/pokt-network/pocket-core/baseapp"
	"github.com/pokt-network/pocket-core/codec"
	types2 "github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/store"
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
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	tmCfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/client/local"
	cTypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	dummyChainsHash = "0001"
)

var (
	tendermintTimeoutCommit = time.Duration(500) * time.Millisecond
)

type upgrades struct {
	codecUpgrade  codecUpgrade
	eight0Upgrade upgrade
}
type upgrade struct {
	height int64
}
type codecUpgrade struct {
	upgradeMod bool
	height     int64
	//after8     bool
}

// NewInMemoryTendermintNodeAmino will create a TM node with only one validator. LeanPocket is disabled.
func NewInMemoryTendermintNodeAmino(t *testing.T, genesisState []byte) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	return NewInMemoryTendermintNodeAminoWithValidators(t, genesisState, nil)
}

// NewInMemoryTendermintNodeAminoWithValidators will create a TM node with 'n' "validators".
// If "validators" is nil, LeanPOKT is disabled
func NewInMemoryTendermintNodeAminoWithValidators(t *testing.T, genesisState []byte, validators []crypto.PrivateKey) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	// create the in memory tendermint node and keybase
	tendermintNode, keybase = inMemTendermintNodeWithValidators(genesisState, validators)
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
	defaultConfig := sdk.DefaultTestingPocketConfig()
	if validators != nil {
		defaultConfig.PocketConfig.LeanPocket = true
	}

	pocketTypes.InitConfig(&pocketTypes.HostedBlockchains{
		M: make(map[string]pocketTypes.HostedBlockchain),
	}, tendermintNode.Logger, defaultConfig)
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
		pocketTypes.CleanPocketNodes()
		pocketTypes.StopEvidenceWorker()
		PCA = nil
		inMemKB = nil
		err := inMemDB.Close()
		if err != nil {
			panic(err)
		}
		cdc = nil
		memCDC = nil
		inMemDB = nil
		sdk.GlobalCtxCache = nil
		err = os.RemoveAll("data")
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
		codec.TestMode = 0
	}
	return
}

// NewInMemoryTendermintNodeProto will create a TM node with only one validator. LeanPocket is disabled.
func NewInMemoryTendermintNodeProto(t *testing.T, genesisState []byte) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	return NewInMemoryTendermintNodeProtoWithValidators(t, genesisState, nil)
}

// NewInMemoryTendermintNodeWithValidators will create a TM node with 'n' "validators".
// If "validators" is nil, this creates a pre-leanpokt TM node, else it will enable lean pocket
func NewInMemoryTendermintNodeProtoWithValidators(t *testing.T, genesisState []byte, validators []crypto.PrivateKey) (tendermintNode *node.Node, keybase keys.Keybase, cleanup func()) {
	// create the in memory tendermint node and keybase
	tendermintNode, keybase = inMemTendermintNodeWithValidators(genesisState, validators)
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
	defaultConfig := sdk.DefaultTestingPocketConfig()
	if validators != nil {
		defaultConfig.PocketConfig.LeanPocket = true
	}
	pocketTypes.InitConfig(&pocketTypes.HostedBlockchains{
		M: make(map[string]pocketTypes.HostedBlockchain),
	}, tendermintNode.Logger, defaultConfig)
	// start the in memory node
	err := tendermintNode.Start()
	if err != nil {
		panic(err)
	}
	// assert that it is not nil
	assert.Nil(t, err)
	// provide cleanup function
	cleanup = func() {
		codec.TestMode = 0

		err = tendermintNode.Stop()
		if err != nil {
			panic(err)
		}

		pocketTypes.CleanPocketNodes()
		pocketTypes.StopEvidenceWorker()

		PCA = nil
		inMemKB = nil
		err := inMemDB.Close()
		if err != nil {
			panic(err)
		}
		cdc = nil
		memCDC = nil
		inMemDB = nil
		sdk.GlobalCtxCache = nil
		err = os.RemoveAll("data")
		if err != nil {
			panic(err)
		}
		time.Sleep(3 * time.Second)
	}
	return
}

func TestNewInMemoryAminoWithValidators(t *testing.T) {
	gbz, validators, _, _ := generateGenesis(5, 5, 10)
	_, _, cleanup := NewInMemoryTendermintNodeAminoWithValidators(t, gbz, validators)
	defer cleanup()
}

func TestNewInMemoryProtoWithValidators(t *testing.T) {
	gbz, validators, _, _ := generateGenesis(5, 5, 10)
	_, _, cleanup := NewInMemoryTendermintNodeProtoWithValidators(t, gbz, validators)
	defer cleanup()
}

func TestNewInMemoryAmino(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNodeAmino(t, oneAppTwoNodeGenesis())
	defer cleanup()
}
func TestNewInMemoryProto(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNodeProto(t, oneAppTwoNodeGenesis())
	defer cleanup()
}

var (
	memCDC  *codec.Codec
	inMemKB keys.Keybase
	memCLI  client.Client
	inMemDB dbm.DB
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

func getInMemoryDB() dbm.DB {
	if inMemDB == nil {
		inMemDB = dbm.NewMemDB()
	}
	return inMemDB
}

// inMemTendermintNodeWithValidators will create a TM node with 'n' "validators".
// If "validators" is nil, LeanPokt is disabled and uses in memory CB as the sole validator for consensus
func inMemTendermintNodeWithValidators(genesisState []byte, validatorsPk []crypto.PrivateKey) (*node.Node, keys.Keybase) {
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
					MaxBytes:   pocketTypes.DefaultBlockByteSize,
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
	db := getInMemoryDB()
	traceWriter, err := openTraceWriter(c.TraceWriter)
	if err != nil {
		panic(err)
	}
	nodeKey := p2p.NodeKey{PrivKey: pk}
	var privVal *privval.FilePVLean
	if validatorsPk == nil {
		// only set cb as validator
		privVal = privval.GenFilePVLean(c.TmConfig.PrivValidatorKey, c.TmConfig.PrivValidatorState)
		privVal.Keys[0].PrivKey = pk
		privVal.Keys[0].PubKey = pk.PubKey()
		privVal.Keys[0].Address = pk.PubKey().Address()
		pocketTypes.CleanPocketNodes()
		pocketTypes.AddPocketNodeByFilePVKey(privVal.Keys[0], c.Logger)
	} else {
		// (LeanPOKT) Set multiple nodes as validators
		pocketTypes.CleanPocketNodes()
		// generating a stub of n validators
		privVal = privval.GenFilePVsLean(c.TmConfig.PrivValidatorKey, c.TmConfig.PrivValidatorState, uint(len(validatorsPk)))
		// replace the stub validators with the correct validators
		for i, pk := range validatorsPk {
			privVal.Keys[i].PrivKey = pk.PrivKey()
			privVal.Keys[i].PubKey = pk.PubKey()
			privVal.Keys[i].Address = pk.PubKey().Address()
			pocketTypes.AddPocketNode(pk, c.Logger)
		}
	}

	dbProvider := func(*node.DBContext) (dbm.DB, error) {
		return db, nil
	}
	app := GetApp(c.Logger, db, traceWriter)
	txDB := dbm.NewMemDB()
	tmNode, err := node.NewNode(app.BaseApp,
		c.TmConfig,
		0,
		privVal,
		&nodeKey,
		proxy.NewLocalClientCreator(app),
		sdk.NewTransactionIndexer(txDB),
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

func GetApp(logger log.Logger, db dbm.DB, traceWriter io.Writer) *PocketCoreApp {
	creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *PocketCoreApp {
		m := map[string]pocketTypes.HostedBlockchain{"0001": {
			ID:  sdk.PlaceholderHash,
			URL: sdk.PlaceholderURL,
		}}
		p := NewPocketCoreApp(GenState, getInMemoryKeybase(), getInMemoryTMClient(), &pocketTypes.HostedBlockchains{M: m, L: sync.RWMutex{}}, logger, db, false, 5000000, bam.SetPruning(store.PruneNothing))
		return p
	}
	return creator(logger, db, traceWriter)
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
	memCDC.SetUpgradeOverride(upgrade)
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
	eventChan, err := cli.Subscribe(ctx, "helpers", types.QueryForEvent(eventType).String(), 5)
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
	newTMConfig.Consensus.CreateEmptyBlocksInterval = tendermintTimeoutCommit
	newTMConfig.Consensus.TimeoutCommit = tendermintTimeoutCommit
	newTMConfig.P2P.MaxNumInboundPeers = 4
	newTMConfig.P2P.MaxNumOutboundPeers = 4
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

func oneAppTwoNodeGenesis() []byte {
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

func twoValTwoNodeGenesisState8() (genbz []byte, vals []nodesTypes.Validator) {
	kb := getInMemoryKeybase()
	kp1, err := kb.GetCoinbase()
	if err != nil {
		panic(err)
	}
	kp2, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	kp3, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	kp4, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	pubKey := kp1.PublicKey
	pubKey2 := kp2.PublicKey
	pubKey3 := kp3.PublicKey
	pubkey4 := kp4.PublicKey
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
		nodesTypes.Validator{
			Address:                 sdk.Address(pubKey.Address()),
			PublicKey:               pubKey,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{dummyChainsHash},
			ServiceURL:              sdk.PlaceholderServiceURL,
			StakedTokens:            sdk.NewInt(1000000000000000),
			UnstakingCompletionTime: time.Time{},
			OutputAddress:           kp3.GetAddress(),
		})
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{
			Address:                 sdk.Address(pubKey2.Address()),
			PublicKey:               pubKey2,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{dummyChainsHash},
			ServiceURL:              sdk.PlaceholderServiceURL,
			StakedTokens:            sdk.NewInt(1000000000),
			UnstakingCompletionTime: time.Time{},
			OutputAddress:           kp4.GetAddress(),
		})
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
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address: sdk.Address(pubKey3.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000000))),
		PubKey:  pubKey3,
	})
	// add second account
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address: sdk.Address(pubkey4.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000000))),
		PubKey:  pubkey4,
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
	return j, posGenesisState.Validators
}

func twoValTwoNodeGenesisState() (genbz []byte, vals []nodesTypes.Validator) {
	kb := getInMemoryKeybase()
	kp1, err := kb.GetCoinbase()
	if err != nil {
		panic(err)
	}
	kp2, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	kp3, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	kp4, err := kb.Create("test")
	if err != nil {
		panic(err)
	}
	pubKey := kp1.PublicKey
	pubKey2 := kp2.PublicKey
	pubKey3 := kp3.PublicKey
	pubkey4 := kp4.PublicKey
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
		nodesTypes.Validator{
			Address:                 sdk.Address(pubKey.Address()),
			PublicKey:               pubKey,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{dummyChainsHash},
			ServiceURL:              sdk.PlaceholderServiceURL,
			StakedTokens:            sdk.NewInt(1000000000000000),
			UnstakingCompletionTime: time.Time{},
			OutputAddress:           nil,
		})
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{
			Address:                 sdk.Address(pubKey2.Address()),
			PublicKey:               pubKey2,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{dummyChainsHash},
			ServiceURL:              sdk.PlaceholderServiceURL,
			StakedTokens:            sdk.NewInt(1000000000),
			UnstakingCompletionTime: time.Time{},
			OutputAddress:           nil,
		})
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
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address: sdk.Address(pubKey3.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000000))),
		PubKey:  pubKey3,
	})
	// add second account
	authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
		Address: sdk.Address(pubkey4.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000000))),
		PubKey:  pubkey4,
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
	return j, posGenesisState.Validators
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

// generateGenesis generates a genesis state of n validators, n servicers, and n apps
func generateGenesis(validators int, servicers int, appss int) ([]byte, []crypto.PrivateKey, []crypto.PrivateKey, []crypto.PrivateKey) {

	kb := getInMemoryKeybase()
	kp1, err := kb.GetCoinbase()

	if err != nil {
		panic(err)
	}

	validatorPks := []crypto.PrivateKey{}
	servicerPks := []crypto.PrivateKey{}
	appPks := []crypto.PrivateKey{}

	encryptPassPhrase := "test"
	for i := 0; i < validators; i++ {
		keyPair, err := kb.Create(encryptPassPhrase)
		if err != nil {
			panic(err)
		}
		pk, err := kb.ExportPrivateKeyObject(keyPair.GetAddress(), encryptPassPhrase)
		if err != nil {
			panic(err)
		}
		validatorPks = append(validatorPks, pk)
	}

	for i := 0; i < servicers; i++ {
		keyPair, err := kb.Create(encryptPassPhrase)
		if err != nil {
			panic(err)
		}
		pk, err := kb.ExportPrivateKeyObject(keyPair.GetAddress(), encryptPassPhrase)
		if err != nil {
			panic(err)
		}
		servicerPks = append(servicerPks, pk)
	}

	for i := 0; i < appss; i++ {
		keyPair, err := kb.Create(encryptPassPhrase)
		if err != nil {
			panic(err)
		}
		pk, err := kb.ExportPrivateKeyObject(keyPair.GetAddress(), encryptPassPhrase)
		if err != nil {
			panic(err)
		}
		appPks = append(appPks, pk)
	}

	defaultGenesis := module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		gov.AppModuleBasic{},
		nodes.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).DefaultGenesis()

	rawPOS := defaultGenesis[nodesTypes.ModuleName]
	var posGenesisState nodesTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawPOS, &posGenesisState)

	rawAccounts := defaultGenesis[auth.ModuleName]
	var authGenState auth.GenesisState
	memCodec().MustUnmarshalJSON(rawAccounts, &authGenState)

	MinStake := int64(10000000000)
	ValidatorStake := MinStake + 1000000

	posGenesisState.Params.StakeMinimum = MinStake
	posGenesisState.Params.MaxValidators = int64(validators)
	// validators kp
	for _, v := range validatorPks {
		posGenesisState.Validators = append(posGenesisState.Validators,
			nodesTypes.Validator{Address: sdk.Address(v.PublicKey().Address()),
				PublicKey:    v.PublicKey(),
				Status:       sdk.Staked,
				Chains:       []string{dummyChainsHash},
				ServiceURL:   sdk.PlaceholderServiceURL,
				StakedTokens: sdk.NewInt(ValidatorStake)})

		authGenState.Accounts = append(authGenState.Accounts, &auth.BaseAccount{
			Address: sdk.Address(v.PublicKey().Address()),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(ValidatorStake))),
			PubKey:  v.PublicKey(),
		})
	}

	for _, v := range servicerPks {
		posGenesisState.Validators = append(posGenesisState.Validators,
			nodesTypes.Validator{Address: sdk.Address(v.PublicKey().Address()),
				PublicKey:    v.PublicKey(),
				Status:       sdk.Staked,
				Chains:       []string{dummyChainsHash},
				ServiceURL:   sdk.PlaceholderServiceURL,
				StakedTokens: sdk.NewInt(MinStake)})
	}

	// validators+servicers
	res := memCodec().MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res

	// pokt holders
	res3 := memCodec().MustMarshalJSON(authGenState)
	defaultGenesis[auth.ModuleName] = res3

	// setup application
	rawApps := defaultGenesis[appsTypes.ModuleName]
	var appsGenesisState appsTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawApps, &appsGenesisState)
	for _, pk := range appPks {
		appsGenesisState.Applications = append(appsGenesisState.Applications, appsTypes.Application{
			Address:                 sdk.GetAddress(pk.PublicKey()),
			PublicKey:               pk.PublicKey(),
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{dummyChainsHash},
			StakedTokens:            sdk.NewInt(10000000),
			MaxRelays:               sdk.NewInt(100000),
			UnstakingCompletionTime: time.Time{},
		})
	}
	res2 := memCodec().MustMarshalJSON(appsGenesisState)
	defaultGenesis[appsTypes.ModuleName] = res2

	// set default chain for module
	rawPocket := defaultGenesis[pocketTypes.ModuleName]
	var pocketGenesisState pocketTypes.GenesisState
	memCodec().MustUnmarshalJSON(rawPocket, &pocketGenesisState)
	pocketGenesisState.Params.SessionNodeCount = int64(validators + servicers)
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
	return j, validatorPks, servicerPks, appPks
}

//
//func TestGatewayChecker(t *testing.T) {
//	startheight := 14681
//	iterations := 30
//	blocks := 96 // ~ 24 hours
//	oldTotalSupply := 0
//	// Code below
//	type Supply struct {
//		Total string `json:"total"`
//	}
//	type Result struct {
//		Inflation    int `json:"inflation"`
//		Day          int `json:"days_ago"`
//		Height       int `json:"height"`
//		DeviationPer int `json:"dev_perc"`
//	}
//	var results []Result
//	var supply Supply
//	var sum int
//	for i := 0; i <= iterations; i++ {
//		jsonStr := `{"height":` + strconv.Itoa(startheight) + `}`
//		req, _ := http2.NewRequest("POST", "http://localhost:8081/v1/query/supply", bytes.NewBuffer([]byte(jsonStr)))
//		client := http2.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			panic(err)
//		}
//		bd, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			panic(err)
//		}
//		err = json.Unmarshal(bd, &supply)
//		if err != nil {
//			panic(err)
//		}
//		total, err := strconv.Atoi(supply.Total)
//		if err != nil {
//			panic(err)
//		}
//		if oldTotalSupply != 0 {
//			results = append(results, Result{
//				Inflation: oldTotalSupply - total,
//				Day:       i,
//				Height:    startheight,
//			})
//			sum += oldTotalSupply - total
//		}
//		oldTotalSupply = total
//		startheight = startheight - blocks
//	}
//	avg := sum / iterations
//	for _, result := range results {
//		result.DeviationPer = (100 * (result.Inflation - avg)) / avg
//		bz, err := json.MarshalIndent(result, "", "  ")
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(string(bz))
//	}
//}
