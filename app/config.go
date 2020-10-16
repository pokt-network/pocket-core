package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	log2 "log"
	"os"
	fp "path/filepath"
	"strings"
	"syscall"

	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/pokt-network/pocket-core/baseapp"
	"github.com/pokt-network/pocket-core/codec"
	cfg "github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/crypto"
	kb "github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/store"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/gov"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/spf13/cobra"
	con "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	cmn "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/client/local"
	dbm "github.com/tendermint/tm-db"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	DefaultDDName                     = ".pocket"
	DefaultKeybaseName                = "pocket-keybase"
	DefaultPVKName                    = "priv_val_key.json"
	DefaultPVSName                    = "priv_val_state.json"
	DefaultNKName                     = "node_key.json"
	DefaultChainsName                 = "chains.json"
	DefaultGenesisName                = "genesis.json"
	DefaultRPCPort                    = "8081"
	DefaultSessionDBType              = dbm.CLevelDBBackend
	DefaultEvidenceDBType             = dbm.CLevelDBBackend
	DefaultSessionDBName              = "session"
	DefaultEvidenceDBName             = "pocket_evidence"
	DefaultTMURI                      = "tcp://localhost:26657"
	DefaultMaxSessionCacheEntries     = 500
	DefaultMaxEvidenceCacheEntries    = 500
	DefaultListenAddr                 = "tcp://0.0.0.0:"
	DefaultClientBlockSyncAllowance   = 10
	DefaultJSONSortRelayResponses     = true
	DefaultDBBackend                  = string(dbm.CLevelDBBackend)
	DefaultTxIndexer                  = "kv"
	DefaultTxIndexTags                = "tx.hash,tx.height,message.sender,transfer.recipient"
	ConfigDirName                     = "config"
	ConfigFileName                    = "config.json"
	ApplicationDBName                 = "application"
	PlaceholderHash                   = "0001"
	PlaceholderURL                    = "http://127.0.0.1:8081"
	PlaceholderServiceURL             = PlaceholderURL
	DefaultRemoteCLIURL               = "http://localhost:8081"
	DefaultUserAgent                  = ""
	DefaultValidatorCacheSize         = 500
	DefaultApplicationCacheSize       = DefaultValidatorCacheSize
	DefaultPocketPrometheusListenAddr = "8083"
	DefaultPrometheusMaxOpenFile      = 3
	DefaultRPCTimeout                 = 3000
	DefaultCtxCacheSize               = 20
)

var (
	cdc *codec.Codec
	// the default fileseparator based on OS
	FS = string(fp.Separator)
	// app instance currently running
	PCA *PocketCoreApp
	// config
	GlobalConfig sdk.Config
	// HTTP CLIENT FOR TENDERMINT
	tmClient *http.HTTP
	// global genesis type
	GlobalGenesisType GenesisType
)

type Config struct {
	TendermintConfig con.Config   `json:"tendermint_config"`
	PocketConfig     PocketConfig `json:"pocket_config"`
}

type PocketConfig struct {
	DataDir                  string          `json:"data_dir"`
	GenesisName              string          `json:"genesis_file"`
	ChainsName               string          `json:"chains_name"`
	SessionDBType            dbm.BackendType `json:"session_db_type"`
	SessionDBName            string          `json:"session_db_name"`
	EvidenceDBType           dbm.BackendType `json:"evidence_db_type"`
	EvidenceDBName           string          `json:"evidence_db_name"`
	TendermintURI            string          `json:"tendermint_uri"`
	KeybaseName              string          `json:"keybase_name"`
	RPCPort                  string          `json:"rpc_port"`
	ClientBlockSyncAllowance int             `json:"client_block_sync_allowance"`
	MaxEvidenceCacheEntires  int             `json:"max_evidence_cache_entries"`
	MaxSessionCacheEntries   int             `json:"max_session_cache_entries"`
	JSONSortRelayResponses   bool            `json:"json_sort_relay_responses"`
	RemoteCLIURL             string          `json:"remote_cli_url"`
	UserAgent                string          `json:"user_agent"`
	ValidatorCacheSize       int64           `json:"validator_cache_size"`
	ApplicationCacheSize     int64           `json:"application_cache_size"`
	RPCTimeout               int64           `json:"rpc_timeout"`
	PrometheusAddr           string          `json:"pocket_prometheus_port"`
	PrometheusMaxOpenfiles   int             `json:"prometheus_max_open_files"`
	CtxCacheSize             int             `json:"ctx_cache_size"`
}

type GenesisType int

const (
	MainnetGenesisType GenesisType = iota + 1
	TestnetGenesisType
	DefaultGenesisType
)

func DefaultConfig(dataDir string) Config {
	c := Config{
		TendermintConfig: *con.DefaultConfig(),
		PocketConfig: PocketConfig{
			DataDir:                  dataDir,
			GenesisName:              DefaultGenesisName,
			ChainsName:               DefaultChainsName,
			SessionDBType:            DefaultSessionDBType,
			SessionDBName:            DefaultSessionDBName,
			EvidenceDBType:           DefaultEvidenceDBType,
			EvidenceDBName:           DefaultEvidenceDBName,
			TendermintURI:            DefaultTMURI,
			KeybaseName:              DefaultKeybaseName,
			RPCPort:                  DefaultRPCPort,
			ClientBlockSyncAllowance: DefaultClientBlockSyncAllowance,
			MaxEvidenceCacheEntires:  DefaultMaxEvidenceCacheEntries,
			MaxSessionCacheEntries:   DefaultMaxSessionCacheEntries,
			JSONSortRelayResponses:   DefaultJSONSortRelayResponses,
			RemoteCLIURL:             DefaultRemoteCLIURL,
			UserAgent:                DefaultUserAgent,
			ValidatorCacheSize:       DefaultValidatorCacheSize,
			ApplicationCacheSize:     DefaultApplicationCacheSize,
			RPCTimeout:               DefaultRPCTimeout,
			PrometheusAddr:           DefaultPocketPrometheusListenAddr,
			PrometheusMaxOpenfiles:   DefaultPrometheusMaxOpenFile,
			CtxCacheSize:             DefaultCtxCacheSize,
		},
	}
	c.TendermintConfig.SetRoot(dataDir)
	c.TendermintConfig.NodeKey = DefaultNKName
	c.TendermintConfig.PrivValidatorKey = DefaultPVKName
	c.TendermintConfig.PrivValidatorState = DefaultPVSName
	c.TendermintConfig.P2P.AddrBookStrict = false
	c.TendermintConfig.P2P.MaxNumInboundPeers = 250
	c.TendermintConfig.P2P.MaxNumOutboundPeers = 250
	c.TendermintConfig.LogLevel = "*:info, *:error"
	c.TendermintConfig.TxIndex.Indexer = DefaultTxIndexer
	c.TendermintConfig.TxIndex.IndexKeys = DefaultTxIndexTags
	c.TendermintConfig.DBBackend = DefaultDBBackend
	c.TendermintConfig.RPC.GRPCMaxOpenConnections = 2500
	c.TendermintConfig.RPC.MaxOpenConnections = 2500
	c.TendermintConfig.Mempool.Size = 9000
	c.TendermintConfig.Mempool.CacheSize = 9000
	c.TendermintConfig.Consensus.TimeoutPropose = 60000000000
	c.TendermintConfig.Consensus.TimeoutProposeDelta = 10000000000
	c.TendermintConfig.Consensus.TimeoutPrevote = 60000000000
	c.TendermintConfig.Consensus.TimeoutPrevoteDelta = 10000000000
	c.TendermintConfig.Consensus.TimeoutPrecommit = 60000000000
	c.TendermintConfig.Consensus.TimeoutPrecommitDelta = 10000000000
	c.TendermintConfig.Consensus.TimeoutCommit = 900000000000
	c.TendermintConfig.Consensus.SkipTimeoutCommit = false
	c.TendermintConfig.Consensus.CreateEmptyBlocks = true
	c.TendermintConfig.Consensus.CreateEmptyBlocksInterval = 900000000000
	c.TendermintConfig.Consensus.PeerGossipSleepDuration = 100000000
	c.TendermintConfig.Consensus.PeerQueryMaj23SleepDuration = 2000000000
	c.TendermintConfig.P2P.AllowDuplicateIP = true
	sdk.InitCtxCache(c.PocketConfig.CtxCacheSize)
	return c
}

func InitApp(datadir, tmNode, persistentPeers, seeds, remoteCLIURL string, keybase bool, genesisType GenesisType) *node.Node {
	// init config
	InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
	// init the keyfiles
	InitKeyfiles()
	// get hosted blockchains
	chains := NewHostedChains(false)
	// create logger
	logger := InitLogger()
	// init cache
	InitPocketCoreConfig(chains, logger)
	// init genesis
	InitGenesis(genesisType)
	// init the tendermint node
	return InitTendermint(keybase, chains, logger)
}

func InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL string) {
	// setup the codec
	MakeCodec()
	if datadir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log2.Fatal("could not get home directory for data dir creation: " + err.Error())
		}
		datadir = home + FS + sdk.DefaultDDName
	}
	c := sdk.DefaultConfig(datadir)
	// read from ccnfig file
	configFilepath := datadir + FS + sdk.ConfigDirName + FS + sdk.ConfigFileName
	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		// ensure directory path made
		err = os.MkdirAll(c.PocketConfig.DataDir+FS+sdk.ConfigDirName, os.ModePerm)
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
	if remoteCLIURL != "" {
		c.PocketConfig.RemoteCLIURL = strings.TrimRight(remoteCLIURL, "/")
	}
	GlobalConfig = c
}

func InitGenesis(genesisType GenesisType) {
	// set global variable for init
	GlobalGenesisType = genesisType
	// setup file if not exists
	genFP := GlobalConfig.PocketConfig.DataDir + FS + sdk.ConfigDirName + FS + GlobalConfig.PocketConfig.GenesisName
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
			if genesisType == MainnetGenesisType {
				_, err = jsonFile.Write([]byte(mainnetGenesis))
				if err != nil {
					log2.Fatal(err)
				}
			} else if genesisType == TestnetGenesisType {
				_, err = jsonFile.Write([]byte(testnetGenesis))
				if err != nil {
					log2.Fatal(err)
				}
			} else {
				_, err = jsonFile.Write(newDefaultGenesisState())
				if err != nil {
					log2.Fatal(err)
				}
			}
			// close the file
			err = jsonFile.Close()
			if err != nil {
				log2.Fatal(err)
			}
		}
	}
}

func InitTendermint(keybase bool, chains *types.HostedBlockchains, logger log.Logger) *node.Node {
	c := cfg.Config{
		TmConfig:    &GlobalConfig.TendermintConfig,
		Logger:      logger,
		TraceWriter: "",
	}
	var keys kb.Keybase
	switch keybase {
	case false:
		keys, _ = GetKeybase()
	default:
		keys = MustGetKeybase()
	}
	appCreatorFunc := func(logger log.Logger, db dbm.DB, _ io.Writer) *PocketCoreApp {
		return NewPocketCoreApp(nil, keys, getTMClient(), chains, logger, db, baseapp.SetPruning(store.PruneNothing))
	}
	tmNode, app, err := NewClient(config(c), appCreatorFunc)
	if err != nil {
		log2.Fatal(err)
	}
	if err := tmNode.Start(); err != nil {
		log2.Fatal(err)
	}
	app.SetTendermintNode(tmNode)
	app.pocketKeeper.TmNode = local.New(tmNode)
	PCA = app
	return tmNode
}

func InitKeyfiles() {
	datadir := GlobalConfig.PocketConfig.DataDir
	// Check if privvalkey file exist
	if _, err := os.Stat(datadir + FS + GlobalConfig.TendermintConfig.PrivValidatorKey); err != nil {
		// if not exist continue creating as other files may be missing
		if os.IsNotExist(err) {
			// generate random key for easy orchestration
			randomKey := crypto.GenerateEd25519PrivKey()
			privValKey(randomKey)
			privValState()
			nodeKey(randomKey)
			log2.Printf("No Validator Set! Creating Random Key: %s", randomKey.PublicKey().RawString())
			return
		} else {
			//panic on other errors
			log2.Fatal(err)
		}
	} else {
		// file exist so we can load pk from file.
		file, _ := loadPKFromFile(datadir + FS + GlobalConfig.TendermintConfig.PrivValidatorKey)
		types.InitPVKeyFile(file)
	}
}

func InitLogger() (logger log.Logger) {
	logger = log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), func(keyvals ...interface{}) term.FgBgColor {
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
	return
}

func InitPocketCoreConfig(chains *types.HostedBlockchains, logger log.Logger) {
	types.InitConfig(chains, logger, GlobalConfig.PocketConfig)
	nodesTypes.InitConfig(GlobalConfig.PocketConfig.ValidatorCacheSize)
	appsTypes.InitConfig(GlobalConfig.PocketConfig.ApplicationCacheSize)
}

func ShutdownPocketCore() {
	types.FlushSessionCache()
	types.StopServiceMetrics()
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

func privValKey(res crypto.PrivateKey) {
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
}

func nodeKey(res crypto.PrivateKey) {
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
			tmClient, _ = http.New(sdk.DefaultTMURI, "/websocket")
		} else {
			tmClient, _ = http.New(GlobalConfig.PocketConfig.TendermintURI, "/websocket")
		}
	}
	return tmClient
}

// get the hosted chains variable
func NewHostedChains(generate bool) *types.HostedBlockchains {
	// create the chains path
	var chainsPath = GlobalConfig.PocketConfig.DataDir + FS + sdk.ConfigDirName + FS + GlobalConfig.PocketConfig.ChainsName
	// if file exists open, else create and open
	var jsonFile *os.File
	var bz []byte
	if _, err := os.Stat(chainsPath); err != nil && os.IsNotExist(err) {
		if !generate {
			log2.Println(fmt.Sprintf("no chains.json found @ %s, defaulting to empty chains", chainsPath))
			return &types.HostedBlockchains{} // default to empty object
		}
		return generateChainsJson(chainsPath)
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

func generateChainsJson(chainsPath string) *types.HostedBlockchains {
	var jsonFile *os.File
	// if does not exist create one
	jsonFile, err := os.OpenFile(chainsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return &types.HostedBlockchains{} // default to empty object
	}
	// generate hosted chains from user input
	c := GenerateHostedChains()
	// create dummy input for the file
	res, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	// write to the file
	_, err = jsonFile.Write(res)
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	m := make(map[string]types.HostedBlockchain)
	for _, chain := range c {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			log2.Fatal(fmt.Sprintf("invalid ID: %s in network identifier in %s file", chain.ID, GlobalConfig.PocketConfig.ChainsName))
		}
		m[chain.ID] = chain
	}
	// return the map
	return &types.HostedBlockchains{M: m}
}

const (
	enterIDPrompt     = `Enter the ID of the network identifier:`
	enterURLPrompt    = `Enter the URL of the network identifier:`
	addNewChainPrompt = `Would you like to enter another network identifier? (y/n)`
	ReadInError       = `An error occurred reading in the information: `
)

func GenerateHostedChains() (chains []types.HostedBlockchain) {
	for {
		var ID, URL, again string
		fmt.Println(enterIDPrompt)
		reader := bufio.NewReader(os.Stdin)
		ID, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(ReadInError + err.Error())
			os.Exit(3)
		}
		ID = strings.Trim(strings.TrimSpace(ID), "\n")
		if err := nodesTypes.ValidateNetworkIdentifier(ID); err != nil {
			fmt.Println(err)
			fmt.Println("please try again")
			continue
		}
		fmt.Println(enterURLPrompt)
		URL, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(ReadInError + err.Error())
			os.Exit(3)
		}
		URL = strings.Trim(strings.TrimSpace(URL), "\n")
		chains = append(chains, types.HostedBlockchain{
			ID:  ID,
			URL: URL,
		})
		fmt.Println(addNewChainPrompt)
		for {
			again, err = reader.ReadString('\n')
			if err != nil {
				log2.Fatal(ReadInError + err.Error())
			}
			switch strings.TrimSpace(strings.ToLower(again)) {
			case "y":
				// break out of switch
				break
			case "n":
				// return chains
				return
			default:
				fmt.Println("unrecognized response, please try again")
				// try switch again
				continue
			}
			// break out of for loop
			break
		}
	}
}

func DeleteHostedChains() {
	// create the chains path
	var chainsPath = GlobalConfig.PocketConfig.DataDir + FS + sdk.ConfigDirName + FS + GlobalConfig.PocketConfig.ChainsName
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
	keys := MustGetKeybase()
	res, err := (keys).ExportPrivateKeyObject(address, passphrase)
	if err != nil {
		log2.Fatal(err)
	}
	privValKey(res)
	privValState()
	nodeKey(res)
}

func GetPrivValFile() (file privval.FilePVKey) {
	file, _ = loadPKFromFile(GlobalConfig.PocketConfig.DataDir + FS + GlobalConfig.TendermintConfig.PrivValidatorKey)
	return
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
		datadir = home + FS + sdk.DefaultDDName
	} else {
		datadir = cmd.Flag("datadir").Value.String()
	}
	c := sdk.DefaultConfig(datadir)
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
		_ = os.Remove(privValKeyFile)
		_ = os.Remove(privValStateFile)
		_ = os.Remove(GlobalConfig.PocketConfig.DataDir + FS + GlobalConfig.TendermintConfig.NodeKey)
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
