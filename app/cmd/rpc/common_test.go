package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	fp "path/filepath"
	"testing"
	"time"

	"github.com/tendermint/tendermint/privval"

	types2 "github.com/pokt-network/pocket-core/codec/types"

	"github.com/tendermint/tendermint/rpc/client/http"

	"github.com/pokt-network/pocket-core/app"
	bam "github.com/pokt-network/pocket-core/baseapp"
	"github.com/pokt-network/pocket-core/codec"
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
	cTypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var FS = string(fp.Separator)

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
	// chains := &pocketTypes.HostedBlockchains{M: make(map[string]pocketTypes.HostedBlockchain)}
	// chains.M[dummyChainsHash] = pocketTypes.HostedBlockchain{ID: dummyChainsHash, URL: dummyChainsURL }
	// init cache in memory
	pocketTypes.InitConfig(&pocketTypes.HostedBlockchains{
		M: make(map[string]pocketTypes.HostedBlockchain),
	}, tendermintNode.Logger, sdk.DefaultTestingPocketConfig())
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
		pocketTypes.ClearEvidence()
		pocketTypes.ClearSessionCache()
		inMemKB = nil
		//err = os.RemoveAll(tendermintNode.Config().DBPath)
		if err != nil {
			panic(err)
		}
		err = os.RemoveAll("data")
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
	return
}

func TestNewInMemory(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	defer cleanup()
}

var (
	memCDC  *codec.Codec
	inMemKB keys.Keybase
	memCLI  client.Client
)

const (
	dummyChainsHash = "0001"
	dummyChainsURL  = "http:127.0.0.1:8081"
	dummyServiceURL = "https://foo.bar:8081"
	defaultTMURI    = "tcp://localhost:26657"
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
	privVal := privval.GenFilePV(c.TmConfig.PrivValidatorKey, c.TmConfig.PrivValidatorState)
	privVal.Key.PrivKey = pk
	privVal.Key.PubKey = pk.PubKey()
	privVal.Key.Address = pk.PubKey().Address()
	pocketTypes.InitPVKeyFile(privVal.Key)

	creator := func(logger log.Logger, db dbm.DB, _ io.Writer) *app.PocketCoreApp {
		m := map[string]pocketTypes.HostedBlockchain{sdk.PlaceholderHash: {
			ID:  sdk.PlaceholderHash,
			URL: dummyChainsURL,
		}}
		p := app.NewPocketCoreApp(app.GenState, getInMemoryKeybase(), getInMemoryTMClient(), &pocketTypes.HostedBlockchains{M: m}, logger, db, bam.SetPruning(store.PruneNothing))
		return p
	}
	//upgradePrivVal(c.TmConfig)
	dbProvider := func(*node.DBContext) (dbm.DB, error) {
		return db, nil
	}
	txDB := dbm.NewMemDB()
	baseapp := creator(c.Logger, db, io.Writer(nil))
	tmNode, err := node.NewNode(baseapp,
		c.TmConfig,
		0,
		privVal,
		&nodeKey,
		proxy.NewLocalClientCreator(baseapp),
		sdk.NewTransactionIndexer(txDB),
		genDocProvider,
		dbProvider,
		node.DefaultMetricsProvider(c.TmConfig.Instrumentation),
		c.Logger.With("module", "node"),
	)
	if err != nil {
		panic(err)
	}
	baseapp.SetTxIndexer(tmNode.TxIndexer())
	baseapp.SetBlockstore(tmNode.BlockStore())
	baseapp.SetEvidencePool(tmNode.EvidencePool())
	baseapp.SetTendermintNode(tmNode)
	app.PCA = baseapp
	return tmNode, kb
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
		memCLI, _ = http.New(defaultTMURI, "/websocket")
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

func getTestConfig() (tmConfg *tmCfg.Config) {
	tmConfg = tmCfg.TestConfig()
	tmConfg.RPC.ListenAddress = defaultTMURI
	tmConfg.Consensus.CreateEmptyBlocks = true // Set this to false to only produce blocks when there are txs or when the AppHash changes
	tmConfg.Consensus.SkipTimeoutCommit = false
	tmConfg.Consensus.CreateEmptyBlocksInterval = time.Duration(50) * time.Millisecond
	tmConfg.Consensus.TimeoutCommit = time.Duration(50) * time.Millisecond
	tmConfg.TxIndex.Indexer = "kv"
	tmConfg.TxIndex.IndexKeys = "tx.hash,tx.height,message.sender"
	return
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
		gov.AppModuleBasic{},
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
	// set default governance in genesis
	var govGenesisState govTypes.GenesisState
	rawGov := defaultGenesis[govTypes.ModuleName]
	memCodec().MustUnmarshalJSON(rawGov, &govGenesisState)
	mACL := createTestACL(kp1)
	govGenesisState.Params.ACL = mACL
	govGenesisState.Params.DAOOwner = kp1.GetAddress()
	govGenesisState.Params.Upgrade = govTypes.NewUpgrade(10000, "2.0.0")
	res4 := memCodec().MustMarshalJSON(govGenesisState)
	defaultGenesis[govTypes.ModuleName] = res4
	pocketGenesisState.Params.SupportedBlockchains = []string{dummyChainsHash}
	// end genesis setup
	app.GenState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
}

var testACL govTypes.ACL

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
		acl.SetOwner("auth/FeeMultipliers", kp.GetAddress())
		acl.SetOwner("application/ApplicationStakeMinimum", kp.GetAddress())
		acl.SetOwner("pocketcore/ClaimExpiration", kp.GetAddress())
		acl.SetOwner("pocketcore/SessionNodeCount", kp.GetAddress())
		acl.SetOwner("pocketcore/MinimumNumberOfProofs", kp.GetAddress())
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
		acl.SetOwner("application/MaximumChains", kp.GetAddress())
		acl.SetOwner("pos/MaximumChains", kp.GetAddress())
		acl.SetOwner("pos/MaxJailedBlocks", kp.GetAddress())
		testACL = acl
	}
	return testACL
}

func fiveValidatorsOneAppGenesis() (genBz []byte, keys []crypto.PrivateKey, validators nodesTypes.Validators, application appsTypes.Application) {
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
	// application 1
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
	app.GenState = defaultGenesis
	j, _ := memCodec().MarshalJSONIndent(defaultGenesis, "", "    ")
	return j, kys, posGenesisState.Validators, appsGenesisState.Applications[0]
}

type config struct {
	TmConfig    *tmCfg.Config
	Logger      log.Logger
	TraceWriter string
}

func generateChainsJson(configFilePath string, chains []pocketTypes.HostedBlockchain) *pocketTypes.HostedBlockchains {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// ensure directory path made
		err = os.MkdirAll(configFilePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	chainsPath := configFilePath + FS + sdk.DefaultChainsName
	var jsonFile *os.File
	// if does not exist create one
	jsonFile, err := os.OpenFile(chainsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	// generate hosted chains from user input
	// create dummy input for the file
	res, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		panic(err)
	}
	// write to the file
	_, err = jsonFile.Write(res)
	if err != nil {
		panic(err)
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		panic(err)
	}
	m := make(map[string]pocketTypes.HostedBlockchain)
	for _, chain := range chains {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			panic(errors.New(fmt.Sprintf("invalid ID: %s in network identifier in %s file", chain.ID, app.GlobalConfig.PocketConfig.ChainsName)))
		}
		m[chain.ID] = chain
	}
	// return the map
	return &pocketTypes.HostedBlockchains{M: m}
}
