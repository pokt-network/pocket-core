package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	log2 "log"
	"os"
	fp "path/filepath"
	"strings"
	"syscall"
	"time"

	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/baseapp"
	"github.com/pokt-network/posmint/codec"
	cfg "github.com/pokt-network/posmint/config"
	"github.com/pokt-network/posmint/crypto"
	kb "github.com/pokt-network/posmint/crypto/keys"
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/gov"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/spf13/cobra"
	con "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli/flags"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/rpc/client"
	tmType "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	DefaultDDName                   = ".pocket"
	DefaultKeybaseName              = "pocket-keybase"
	DefaultPVKName                  = "priv_val_key.json"
	DefaultPVSName                  = "priv_val_state.json"
	DefaultNKName                   = "node_key.json"
	DefaultChainsName               = "chains.json"
	DefaultGenesisName              = "genesis.json"
	DefaultRPCPort                  = "8081"
	DefaultSessionDBType            = dbm.GoLevelDBBackend
	DefaultEvidenceDBType           = dbm.GoLevelDBBackend
	DefaultSessionDBName            = "session"
	DefaultEvidenceDBName           = "pocket_evidence"
	DefaultTMURI                    = "tcp://localhost:26657"
	DefaultMaxSessionCacheEntries   = 100
	DefaultMaxEvidenceCacheEntries  = 100
	DefaultListenAddr               = "tcp://0.0.0.0:"
	DefaultClientBlockSyncAllowance = 10
	DefaultJSONSortRelayResponses   = true
	DefaultDBBackend                = string(dbm.GoLevelDBBackend)
	DefaultTxIndexer                = "kv"
	DefaultTxIndexTags              = "tx.hash,tx.height,message.sender,transfer.recipient"
	ConfigDirName                   = "config"
	ConfigFileName                  = "config.json"
	ApplicationDBName               = "application"
	PlaceholderHash                 = "00"
	PlaceholderURL                  = "https://foo.bar:8080"
	PlaceholderServiceURL           = PlaceholderURL
	DefaultRemoteCLIURL             = "http://localhost"
	DefaultUserAgent                = ""
)

var (
	cdc *codec.Codec
	// the default fileseparator based on OS
	FS = string(fp.Separator)
	// app instance currently running
	PCA *PocketCoreApp
	// config
	GlobalConfig Config
	// HTTP CLIENT FOR TENDERMINT
	tmClient *client.HTTP
)

type Config struct {
	TendermintConfig con.Config   `json:"tendermint_config"`
	PocketConfig     PocketConfig `json:"pocket_config"`
}

type PocketConfig struct {
	DataDir                  string            `json:"data_dir"`
	GenesisName              string            `json:"genesis_file"`
	ChainsName               string            `json:"chains_name"`
	SessionDBType            dbm.DBBackendType `json:"session_db_type"`
	SessionDBName            string            `json:"session_db_name"`
	EvidenceDBType           dbm.DBBackendType `json:"evidence_db_type"`
	EvidenceDBName           string            `json:"evidence_db_name"`
	TendermintURI            string            `json:"tendermint_uri"`
	KeybaseName              string            `json:"keybase_name"`
	RPCPort                  string            `json:"rpc_port"`
	ClientBlockSyncAllowance int               `json:"client_block_sync_allowance"`
	MaxEvidenceCacheEntires  int               `json:"max_evidence_cache_entries"`
	MaxSessionCacheEntries   int               `json:"max_session_cache_entries"`
	JSONSortRelayResponses   bool              `json:"json_sort_relay_responses"`
	RemoteCLIURL             string            `json:"remote_cli_url"`
	UserAgent                string            `json:"user_agent"`
}

func DefaultConfig(dataDir string) Config {
	c := Config{
		TendermintConfig: *con.DefaultConfig(),
		PocketConfig: PocketConfig{
			DataDir:                  dataDir,
			RPCPort:                  DefaultRPCPort,
			GenesisName:              DefaultGenesisName,
			ChainsName:               DefaultChainsName,
			SessionDBType:            DefaultSessionDBType,
			SessionDBName:            DefaultSessionDBName,
			EvidenceDBType:           DefaultEvidenceDBType,
			EvidenceDBName:           DefaultEvidenceDBName,
			TendermintURI:            DefaultTMURI,
			KeybaseName:              DefaultKeybaseName,
			ClientBlockSyncAllowance: DefaultClientBlockSyncAllowance,
			MaxEvidenceCacheEntires:  DefaultMaxEvidenceCacheEntries,
			MaxSessionCacheEntries:   DefaultMaxSessionCacheEntries,
			JSONSortRelayResponses:   DefaultJSONSortRelayResponses,
			RemoteCLIURL:             DefaultRemoteCLIURL,
			UserAgent:                DefaultUserAgent,
		},
	}
	c.TendermintConfig.SetRoot(dataDir)
	c.TendermintConfig.NodeKey = DefaultNKName
	c.TendermintConfig.PrivValidatorKey = DefaultPVKName
	c.TendermintConfig.PrivValidatorState = DefaultPVSName
	c.TendermintConfig.P2P.AddrBookStrict = false
	c.TendermintConfig.Consensus.CreateEmptyBlocks = true
	c.TendermintConfig.Consensus.CreateEmptyBlocksInterval = time.Duration(1) * time.Minute
	c.TendermintConfig.Consensus.TimeoutCommit = time.Duration(1) * time.Minute
	c.TendermintConfig.P2P.MaxNumInboundPeers = 250
	c.TendermintConfig.P2P.MaxNumOutboundPeers = 250
	c.TendermintConfig.LogLevel = "*:info, *:error"
	c.TendermintConfig.TxIndex.Indexer = DefaultTxIndexer
	c.TendermintConfig.TxIndex.IndexTags = DefaultTxIndexTags
	c.TendermintConfig.DBBackend = DefaultDBBackend
	return c
}

func InitApp(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort string) *node.Node {
	// init config
	InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
	// init the keyfiles
	InitKeyfiles()
	// init cache
	InitPocketCoreConfig()
	// init genesis
	InitGenesis()
	// init the tendermint node
	return InitTendermint()
}

func InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort string) {
	// setup the codec
	MakeCodec()
	if datadir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log2.Fatal("could not get home directory for data dir creation: " + err.Error())
		}
		datadir = home + FS + DefaultDDName
	}
	c := DefaultConfig(datadir)
	// read from ccnfig file
	configFilepath := datadir + FS + ConfigDirName + FS + ConfigFileName
	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		// ensure directory path made
		err = os.MkdirAll(c.PocketConfig.DataDir+FS+ConfigDirName, os.ModePerm)
		if err != nil {
			log2.Fatal(err)
		}
	}
	var jsonFile *os.File
	defer jsonFile.Close()
	// if file exists open, else create and open
	if _, err := os.Stat(configFilepath); err == nil {
		jsonFile, err = os.OpenFile(configFilepath, os.O_RDWR, os.ModePerm)
		if err != nil {
			log2.Fatalf("cannot open config json file: " + err.Error())
		}
		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log2.Fatalf("cannot read config file: " + err.Error())
		}
		err = json.Unmarshal(b, &c)
		if err != nil {
			log2.Fatalf("cannot read config file into json: " + err.Error())
		}
	} else if os.IsNotExist(err) {
		// if does not exist create one
		jsonFile, err = os.OpenFile(configFilepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			log2.Fatalf("canot open/create config json file: " + err.Error())
		}
		b, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			log2.Fatalf("cannot marshal default config into json: " + err.Error())
		}
		// write to the file
		_, err = jsonFile.Write(b)
		if err != nil {
			log2.Fatalf("cannot write default config to json file: " + err.Error())
		}
	}
	// flags trump config file
	if tmNode != "" {
		c.PocketConfig.TendermintURI = tmNode
	}
	if persistentPeers != "" {
		c.TendermintConfig.P2P.PersistentPeers = persistentPeers
	}
	if seeds != "" {
		c.TendermintConfig.P2P.Seeds = seeds
	}
	if tmRPCPort != "" {
		c.TendermintConfig.RPC.ListenAddress = DefaultListenAddr + tmRPCPort
	}
	if tmPeersPort != "" {
		c.TendermintConfig.P2P.ListenAddress = DefaultListenAddr + tmPeersPort
	}
	GlobalConfig = c
}

func InitGenesis() {
	genFP := GlobalConfig.PocketConfig.DataDir + FS + ConfigDirName + FS + GlobalConfig.PocketConfig.GenesisName
	if _, err := os.Stat(genFP); os.IsNotExist(err) {
		// if file exists open, else create and open
		if _, err := os.Stat(genFP); err == nil {
			return
		} else if os.IsNotExist(err) {
			// if does not exist create one
			jsonFile, err := os.OpenFile(genFP, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				log2.Fatal(err)
			}
			// write to the file
			_, err = jsonFile.Write(newDefaultGenesisState())
			if err != nil {
				log2.Fatal(err)
			}
			// close the file
			err = jsonFile.Close()
			if err != nil {
				log2.Fatal(err)
			}
		}
	}
}

func InitTendermint() *node.Node {
	logger := log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			fmt.Printf("expected level key to be first, got %v", keyvals[0])
			log2.Fatal(1)
		}
		switch keyvals[1].(kitlevel.Value).String() {
		case "info":
			return term.FgBgColor{Fg: term.Green}
		case "debug":
			return term.FgBgColor{Fg: term.DarkBlue}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	})
	logger, err := flags.ParseLogLevel(GlobalConfig.TendermintConfig.LogLevel, logger, "info")
	if err != nil {
		log2.Fatal(err)
	}
	c := cfg.Config{
		TmConfig:    &GlobalConfig.TendermintConfig,
		Logger:      logger,
		TraceWriter: "",
	}
	tmNode, app, err := NewClient(config(c), func(logger log.Logger, db dbm.DB, _ io.Writer) *PocketCoreApp {
		return NewPocketCoreApp(nil, MustGetKeybase(), getTMClient(), NewHostedChains(), logger, db, baseapp.SetPruning(store.PruneNothing))
	})
	if err != nil {
		log2.Fatal(err)
	}
	if err := tmNode.Start(); err != nil {
		log2.Fatal(err)
	}
	app.SetTendermintNode(tmNode)
	PCA = app
	return tmNode
}

func InitKeyfiles() string {
	var password string
	datadir := GlobalConfig.PocketConfig.DataDir
	// Check if privvalkey file exist
	if _, err := os.Stat(datadir + FS + GlobalConfig.TendermintConfig.PrivValidatorKey); err != nil {
		// if not exist continue creating as other files may be missing
		if os.IsNotExist(err) {
			log2.Fatal(fmt.Errorf("validator address not set! please run set-validator"))
		} else {
			//panic on other errors
			log2.Fatal(err)
		}
	} else {
		// file exist so we can load pk from file.
		file, _ := loadPKFromFile(datadir + FS + GlobalConfig.TendermintConfig.PrivValidatorKey)
		types.InitPVKeyFile(file)
	}
	return password
}

func InitPocketCoreConfig() {
	types.InitConfig(GlobalConfig.PocketConfig.UserAgent, GlobalConfig.PocketConfig.DataDir, GlobalConfig.PocketConfig.DataDir, GlobalConfig.PocketConfig.SessionDBType, GlobalConfig.PocketConfig.EvidenceDBType, GlobalConfig.PocketConfig.MaxEvidenceCacheEntires, GlobalConfig.PocketConfig.MaxSessionCacheEntries, GlobalConfig.PocketConfig.EvidenceDBName, GlobalConfig.PocketConfig.SessionDBName)
	types.InitClientBlockAllowance(GlobalConfig.PocketConfig.ClientBlockSyncAllowance)
	types.InitJSONSorting(GlobalConfig.PocketConfig.JSONSortRelayResponses)
}

// get the global keybase
func MustGetKeybase() kb.Keybase {
	keys, err := GetKeybase()
	if err != nil {
		log2.Fatal(err)
	}
	return keys
}

// get the global keybase
func GetKeybase() (kb.Keybase, error) {
	keys := kb.New(GlobalConfig.PocketConfig.KeybaseName, GlobalConfig.PocketConfig.DataDir)
	kps, err := keys.List()
	if err != nil {
		return nil, err
	}
	if len(kps) == 0 {
		return nil, UninitializedKeybaseError
	}
	return keys, nil
}

func loadPKFromFile(path string) (privval.FilePVKey, string) {
	keyJSONBytes, err := ioutil.ReadFile(path)
	if err != nil {
		cmn.Exit(err.Error())
	}
	pvKey := privval.FilePVKey{}
	err = cdc.UnmarshalJSON(keyJSONBytes, &pvKey)
	if err != nil {
		cmn.Exit(fmt.Sprintf("Error reading PrivValidator key from %v: %v\n", path, err))
	}

	return pvKey, path
}

func privValKey(address sdk.Address, password string) string {
	keys := MustGetKeybase()
	if password == "" {
		fmt.Println("Initializing keyfiles: enter validator account passphrase")
		password = Credentials()
	}
	res, err := (keys).ExportPrivateKeyObject(address, password)
	if err != nil {
		log2.Fatal(err)
	}
	privValKey := privval.FilePVKey{
		Address: res.PubKey().Address(),
		PubKey:  res.PubKey(),
		PrivKey: res.PrivKey(),
	}
	pvkBz, err := cdc.MarshalJSONIndent(privValKey, "", "  ")
	if err != nil {
		log2.Fatal(err)
	}
	pvFile, err := os.OpenFile(GlobalConfig.PocketConfig.DataDir+FS+GlobalConfig.TendermintConfig.PrivValidatorKey, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log2.Fatal(err)
	}
	_, err = pvFile.Write(pvkBz)
	if err != nil {
		log2.Fatal(err)
	}
	types.InitPVKeyFile(privValKey)
	return password
}

func nodeKey(address sdk.Address, password string) {
	keys := MustGetKeybase()
	res, err := (keys).ExportPrivateKeyObject(address, password)
	if err != nil {
		log2.Fatal(err)
	}
	nodeKey := p2p.NodeKey{
		PrivKey: res.PrivKey(),
	}
	pvkBz, err := cdc.MarshalJSONIndent(nodeKey, "", "  ")
	if err != nil {
		log2.Fatal(err)
	}
	pvFile, err := os.OpenFile(GlobalConfig.PocketConfig.DataDir+FS+GlobalConfig.TendermintConfig.NodeKey, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log2.Fatal(err)
	}
	_, err = pvFile.Write(pvkBz)
	if err != nil {
		log2.Fatal(err)
	}
}

func privValState() {
	pvkBz, err := cdc.MarshalJSONIndent(privval.FilePVLastSignState{}, "", "  ")
	if err != nil {
		log2.Fatal(err)
	}
	pvFile, err := os.OpenFile(GlobalConfig.PocketConfig.DataDir+FS+GlobalConfig.TendermintConfig.PrivValidatorState, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log2.Fatal(err)
	}
	_, err = pvFile.Write(pvkBz)
	if err != nil {
		log2.Fatal(err)
	}
}

func getTMClient() client.Client {
	if tmClient == nil {
		if GlobalConfig.PocketConfig.TendermintURI == "" {
			tmClient = client.NewHTTP(DefaultTMURI, "/websocket")
		} else {
			tmClient = client.NewHTTP(GlobalConfig.PocketConfig.TendermintURI, "/websocket")
		}
	}
	return tmClient
}

// get the hosted chains variable
func NewHostedChains() *types.HostedBlockchains {
	// create the chains path
	var chainsPath = GlobalConfig.PocketConfig.DataDir + FS + ConfigDirName + FS + GlobalConfig.PocketConfig.ChainsName
	// if file exists open, else create and open
	var jsonFile *os.File
	var bz []byte
	if _, err := os.Stat(chainsPath); err != nil && os.IsNotExist(err) {
		log2.Println(fmt.Sprintf("no chains.json found @ %s, defaulting to empty chains", chainsPath))
		return &types.HostedBlockchains{} // default to empty object
	}
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(chainsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	bz, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	// unmarshal into the structure
	var hostedChainsSlice []types.HostedBlockchain
	err = json.Unmarshal(bz, &hostedChainsSlice)
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	m := make(map[string]types.HostedBlockchain)
	for _, chain := range hostedChainsSlice {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			log2.Fatal(fmt.Sprintf("invalid ID: %s in network identifier in %s file", chain.ID, GlobalConfig.PocketConfig.ChainsName))
		}
		m[chain.ID] = chain
	}
	// return the map
	return &types.HostedBlockchains{M: m}
}

func DeleteHostedChains() {
	// create the chains path
	var chainsPath = GlobalConfig.PocketConfig.DataDir + FS + ConfigDirName + FS + GlobalConfig.PocketConfig.ChainsName
	err := os.Remove(chainsPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not delete %s file: ", chainsPath) + err.Error())
		os.Exit(3)
	}
}

func Codec() *codec.Codec {
	if cdc == nil {
		MakeCodec()
	}
	return cdc
}

func MakeCodec() {
	// create a new codec
	cdc = codec.New()
	// register all of the app module types
	module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		gov.AppModuleBasic{},
		nodes.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).RegisterCodec(cdc)
	// register the sdk types
	sdk.RegisterCodec(cdc)
	// register the crypto types
	codec.RegisterCrypto(cdc)
}

func Credentials() string {
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(string(bytePassword))
}

func SetValidator(address sdk.Address, passphrase string) {
	resetFilePV(GlobalConfig.PocketConfig.DataDir+FS+GlobalConfig.TendermintConfig.PrivValidatorKey, GlobalConfig.PocketConfig.DataDir+FS+GlobalConfig.TendermintConfig.PrivValidatorState, log.NewNopLogger())
	privValKey(address, passphrase)
	privValState()
	nodeKey(address, passphrase)
}

func GetPrivValFile() (file privval.FilePVKey) {
	file, _ = loadPKFromFile(GlobalConfig.PocketConfig.DataDir + FS + GlobalConfig.TendermintConfig.PrivValidatorKey)
	return
}

func newDefaultGenesisState() []byte {
	keyb, err := GetKeybase()
	if err != nil {
		log2.Fatal(err)
	}
	cb, err := keyb.GetCoinbase()
	if err != nil {
		log2.Fatal(err)
	}
	pubKey := cb.PublicKey
	defaultGenesis := module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		gov.AppModuleBasic{},
		nodes.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).DefaultGenesis()
	// setup account genesis
	rawAuth := defaultGenesis[auth.ModuleName]
	var accountGenesis auth.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(rawAuth, &accountGenesis)
	accountGenesis.Accounts = append(accountGenesis.Accounts, &auth.BaseAccount{
		Address: cb.GetAddress(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000000))),
		PubKey:  pubKey,
	})
	res := Codec().MustMarshalJSON(accountGenesis)
	defaultGenesis[auth.ModuleName] = res
	// set default governance in genesis
	// setup pos genesis
	rawPOS := defaultGenesis[nodesTypes.ModuleName]
	var posGenesisState nodesTypes.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(rawPOS, &posGenesisState)
	posGenesisState.Validators = append(posGenesisState.Validators,
		nodesTypes.Validator{Address: sdk.Address(pubKey.Address()),
			PublicKey:    pubKey,
			Status:       sdk.Staked,
			Chains:       []string{PlaceholderHash},
			ServiceURL:   PlaceholderServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	res = types.ModuleCdc.MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res
	// set default governance in genesis
	var govGenesisState govTypes.GenesisState
	rawGov := defaultGenesis[govTypes.ModuleName]
	Codec().MustUnmarshalJSON(rawGov, &govGenesisState)
	mACL := createDummyACL(pubKey)
	govGenesisState.Params.ACL = mACL
	govGenesisState.Params.DAOOwner = sdk.Address(pubKey.Address())
	govGenesisState.Params.Upgrade = govTypes.NewUpgrade(0, "0")
	res4 := Codec().MustMarshalJSON(govGenesisState)
	defaultGenesis[govTypes.ModuleName] = res4
	// end genesis setup
	j, _ := types.ModuleCdc.MarshalJSONIndent(defaultGenesis, "", "    ")
	j, _ = types.ModuleCdc.MarshalJSONIndent(tmType.GenesisDoc{
		GenesisTime: time.Now(),
		ChainID:     "pocket-test",
		ConsensusParams: &tmType.ConsensusParams{
			Block: tmType.BlockParams{
				MaxBytes:   15000,
				MaxGas:     -1,
				TimeIotaMs: 1,
			},
			Evidence: tmType.EvidenceParams{
				MaxAge: 1000000,
			},
			Validator: tmType.ValidatorParams{
				PubKeyTypes: []string{"ed25519"},
			},
		},
		Validators: nil,
		AppHash:    nil,
		AppState:   j,
	}, "", "    ")
	return j
}

func createDummyACL(kp crypto.PublicKey) govTypes.ACL {
	addr := sdk.Address(kp.Address())
	acl := govTypes.ACL{}
	acl = make([]govTypes.ACLPair, 0)
	acl.SetOwner("auth/MaxMemoCharacters", addr)
	acl.SetOwner("auth/TxSigLimit", addr)
	acl.SetOwner("gov/daoOwner", addr)
	acl.SetOwner("gov/acl", addr)
	acl.SetOwner("pos/StakeDenom", addr)
	acl.SetOwner("pocketcore/SupportedBlockchains", addr)
	acl.SetOwner("pos/DowntimeJailDuration", addr)
	acl.SetOwner("pos/SlashFractionDoubleSign", addr)
	acl.SetOwner("pos/SlashFractionDowntime", addr)
	acl.SetOwner("application/ApplicationStakeMinimum", addr)
	acl.SetOwner("pocketcore/ClaimExpiration", addr)
	acl.SetOwner("pocketcore/SessionNodeCount", addr)
	acl.SetOwner("pocketcore/ReplayAttackBurnMultiplier", addr)
	acl.SetOwner("pos/MaxValidators", addr)
	acl.SetOwner("pos/ProposerPercentage", addr)
	acl.SetOwner("application/StabilityAdjustment", addr)
	acl.SetOwner("application/MaximumChains", addr)
	acl.SetOwner("application/AppUnstakingTime", addr)
	acl.SetOwner("application/ParticipationRateOn", addr)
	acl.SetOwner("pos/MaxEvidenceAge", addr)
	acl.SetOwner("pos/MinSignedPerWindow", addr)
	acl.SetOwner("pos/StakeMinimum", addr)
	acl.SetOwner("pos/UnstakingTime", addr)
	acl.SetOwner("application/BaseRelaysPerPOKT", addr)
	acl.SetOwner("pocketcore/ClaimSubmissionWindow", addr)
	acl.SetOwner("pos/DAOAllocation", addr)
	acl.SetOwner("pos/MaximumChains", addr)
	acl.SetOwner("pos/SignedBlocksWindow", addr)
	acl.SetOwner("pos/BlocksPerSession", addr)
	acl.SetOwner("application/MaxApplications", addr)
	acl.SetOwner("gov/daoOwner", addr)
	acl.SetOwner("pos/RelaysToTokensMultiplier", addr)
	acl.SetOwner("gov/upgrade", addr)
	return acl
}

// XXX: this is totally unsafe.
// it's only suitable for testnets.
func ResetWorldState(cmd *cobra.Command, args []string) {
	var datadir string
	if cmd.Flag("datadir").DefValue == cmd.Flag("datadir").Value.String() {
		home, err := os.UserHomeDir()
		if err != nil {
			log2.Fatal("could not get home directory for data dir creation: " + err.Error())
		}
		datadir = home + FS + DefaultDDName
	} else {
		datadir = cmd.Flag("datadir").Value.String()
	}
	c := DefaultConfig(datadir)
	GlobalConfig = c

	ResetAll(GlobalConfig.TendermintConfig.DBDir(),
		GlobalConfig.TendermintConfig.P2P.AddrBookFile(),
		GlobalConfig.PocketConfig.DataDir+FS+GlobalConfig.TendermintConfig.PrivValidatorKey,
		GlobalConfig.PocketConfig.DataDir+FS+GlobalConfig.TendermintConfig.PrivValidatorState,
		log.NewNopLogger())
}

// ResetAll removes address book files plus all data, and resets the privValidator data.
// Exported so other CLI tools can use it.
func ResetAll(dbDir, addrBookFile, privValKeyFile, privValStateFile string, logger log.Logger) {
	removeAddrBook(addrBookFile, logger)
	if err := os.RemoveAll(dbDir); err == nil {
		logger.Info("Removed all blockchain history", "dir", dbDir)
	} else {
		logger.Error("Error removing all blockchain history", "dir", dbDir, "err", err)
	}
	// recreate the dbDir since the privVal state needs to live there
	err := cmn.EnsureDir(dbDir, 0700)
	if err != nil {
		log2.Fatal(err)
	}
	resetFilePV(privValKeyFile, privValStateFile, logger)
}

func resetFilePV(privValKeyFile, privValStateFile string, logger log.Logger) {
	if _, err := os.Stat(privValKeyFile); err == nil {
		os.Remove(privValKeyFile)
		os.Remove(privValStateFile)
		os.Remove(GlobalConfig.PocketConfig.DataDir + FS + GlobalConfig.TendermintConfig.NodeKey)
	}
	logger.Info("Reset private validator file", "keyFile", privValKeyFile,
		"stateFile", privValStateFile)
}

func removeAddrBook(addrBookFile string, logger log.Logger) {
	if err := os.Remove(addrBookFile); err == nil {
		logger.Info("Removed existing address book", "file", addrBookFile)
	} else if !os.IsNotExist(err) {
		logger.Info("Error removing address book", "file", addrBookFile, "err", err)
	}
}
