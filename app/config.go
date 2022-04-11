package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/pokt-network/pocket-core/baseapp"
	"github.com/pokt-network/pocket-core/codec"
	types2 "github.com/pokt-network/pocket-core/codec/types"
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
	config2 "github.com/tendermint/tendermint/config"
	cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	cmn "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/client/local"
	dbm "github.com/tendermint/tm-db"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	log2 "log"
	"os"
	fp "path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
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
	// current authToken for secured rpc calls
	AuthToken sdk.AuthToken
)

type GenesisType int

const (
	MainnetGenesisType GenesisType = iota + 1
	TestnetGenesisType
	DefaultGenesisType
)

func InitApp(datadir, tmNode, persistentPeers, seeds, remoteCLIURL string, keybase bool, genesisType GenesisType, useCache bool) *node.Node {
	// init config
	InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
	GlobalConfig.PocketConfig.Cache = useCache
	// init AuthToken
	InitAuthToken()
	// init the keyfiles
	InitKeyfiles()
	// get hosted blockchains
	chains := NewHostedChains(false)
	if GlobalConfig.PocketConfig.ChainsHotReload {
		// hot reload chains
		HotReloadChains(chains)
	}
	// create logger
	logger := InitLogger()
	// init cache
	InitPocketCoreConfig(chains, logger)
	// init genesis
	InitGenesis(genesisType, logger)
	// log the config and chains
	logger.Debug(fmt.Sprintf("Pocket Config: \n%v", GlobalConfig))
	// init the tendermint node
	return InitTendermint(keybase, chains, logger)
}

func InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL string) {
	log2.Println("Initializing Pocket Datadir")
	// setup the codec
	MakeCodec()
	if datadir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log2.Fatal("could not get home directory for data dir creation: " + err.Error())
		}
		datadir = home + FS + sdk.DefaultDDName
		log2.Println("datadir = " + datadir)
	}
	c := sdk.DefaultConfig(datadir)
	// read from ccnfig file
	configFilepath := datadir + FS + sdk.ConfigDirName + FS + sdk.ConfigFileName
	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		log2.Println("no config file found... creating the datadir @ "+c.PocketConfig.DataDir+FS+sdk.ConfigDirName, os.ModePerm)
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
	//Always Allow Duplicate IP
	c.TendermintConfig.P2P.AllowDuplicateIP = true

	GlobalConfig = c
}

func UpdateConfig(datadir string) {
	//Check DataDir is present or use Default home dir.

	if datadir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log2.Fatal("could not get home directory for data dir creation: " + err.Error())
		}
		datadir = home + FS + sdk.DefaultDDName
	}

	//Write Defaults on GlobalConfig.
	GlobalConfig.TendermintConfig.LevelDBOptions = config2.DefaultLevelDBOpts()
	sdk.DefaultPocketConsensusConfig(GlobalConfig.TendermintConfig.Consensus)
	GlobalConfig.TendermintConfig.P2P.AllowDuplicateIP = true
	GlobalConfig.TendermintConfig.P2P.AddrBookStrict = false
	GlobalConfig.TendermintConfig.P2P.MaxNumInboundPeers = 14
	GlobalConfig.TendermintConfig.P2P.MaxNumOutboundPeers = 7
	GlobalConfig.TendermintConfig.RPC.GRPCMaxOpenConnections = 2500
	GlobalConfig.TendermintConfig.RPC.MaxOpenConnections = 2500
	GlobalConfig.TendermintConfig.Mempool.Size = 9000
	GlobalConfig.TendermintConfig.Mempool.CacheSize = 9000
	GlobalConfig.TendermintConfig.FastSync = &config2.FastSyncConfig{
		Version: "v1",
	}
	GlobalConfig.PocketConfig.ValidatorCacheSize = sdk.DefaultValidatorCacheSize
	GlobalConfig.PocketConfig.ApplicationCacheSize = sdk.DefaultApplicationCacheSize
	GlobalConfig.PocketConfig.CtxCacheSize = sdk.DefaultCtxCacheSize
	GlobalConfig.PocketConfig.RPCTimeout = sdk.DefaultRPCTimeout
	GlobalConfig.PocketConfig.IavlCacheSize = sdk.DefaultIavlCacheSize

	// Backup and Save the File
	var jsonFile *os.File
	defer jsonFile.Close()

	configFilepath := datadir + FS + sdk.ConfigDirName + FS + sdk.ConfigFileName
	configFileBackupPath := configFilepath + ".bk"

	backupConfigFile(configFilepath, configFileBackupPath)

	writeConfigFile(configFilepath, jsonFile)

}

func writeConfigFile(configFilepath string, jsonFile *os.File) {
	if _, err := os.Stat(configFilepath); err == nil {
		jsonFile, err = os.OpenFile(configFilepath, os.O_RDWR, os.ModePerm)
		if err != nil {
			log2.Fatalf("cannot open config json file: " + err.Error())
		}
		err = jsonFile.Truncate(0)
		if err != nil {
			log2.Fatalf("cannot truncate config json file: " + err.Error())
		}
		b, err := json.MarshalIndent(GlobalConfig, "", "    ")
		if err != nil {
			log2.Fatalf("cannot marshal default config into json: " + err.Error())
		}
		// write to the file
		_, err = jsonFile.Write(b)
		if err != nil {
			log2.Fatalf("cannot write default config to json file: " + err.Error())
		}
	}
}

func backupConfigFile(filepath string, filepath2 string) {
	var jsonFile *os.File
	defer jsonFile.Close()

	if _, err := os.Stat(filepath); err == nil {
		jsonFile, err = os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log2.Fatalf("cannot open config json file: " + err.Error())
		}
		destination, err := os.Create(filepath2)
		if err != nil {
			log2.Fatalf("cannot create backup config json file: " + err.Error())
		}

		_, err = io.Copy(destination, jsonFile)
		if err != nil {
			log2.Fatalf("cannot create backup config json file: " + err.Error())
		}
		_ = destination.Close()
	}
}

func InitGenesis(genesisType GenesisType, logger log.Logger) {
	logger.Info("Initializing genesis file")
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

type Config struct {
	TmConfig    *config2.Config
	Logger      log.Logger
	TraceWriter string
}

func InitTendermint(keybase bool, chains *types.HostedBlockchains, logger log.Logger) *node.Node {
	logger.Info("Initializing Tendermint")
	c := Config{
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
		return NewPocketCoreApp(nil, keys, getTMClient(), chains, logger, db, GlobalConfig.PocketConfig.Cache, GlobalConfig.PocketConfig.IavlCacheSize, baseapp.SetPruning(store.PruneNothing))
	}
	tmNode, app, err := NewClient(config(c), appCreatorFunc)
	if err != nil {
		log2.Fatal(err)
	}
	app.pocketKeeper.TmNode = local.New(tmNode)
	if err := tmNode.Start(); err != nil {
		log2.Fatal(err)
	}
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
	logger.Info("Initializing pocket core config")
	types.InitConfig(chains, logger, GlobalConfig)
	logger.Info("Initializing ctx cache")
	sdk.InitCtxCache(GlobalConfig.PocketConfig.CtxCacheSize)
	logger.Info("Initializing pos config")
	nodesTypes.InitConfig(GlobalConfig.PocketConfig.ValidatorCacheSize)
	logger.Info("Initializing app config")
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

func HotReloadChains(chains *types.HostedBlockchains) {
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			// create the chains path
			var chainsPath = GlobalConfig.PocketConfig.DataDir + FS + sdk.ConfigDirName + FS + GlobalConfig.PocketConfig.ChainsName
			// if file exists open, else create and open
			var jsonFile *os.File
			var bz []byte
			if _, err := os.Stat(chainsPath); err != nil && os.IsNotExist(err) {
				log2.Println(fmt.Sprintf("no chains.json found @ %s, defaulting to empty chains", chainsPath))
				return
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
			chains.L.Lock()
			chains.M = m
			chains.L.Unlock()
		}
	}()
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
			log2.Println(fmt.Sprintf("no chains.json found @ %s, defaulting to sample chains", chainsPath))
			// added for hot reload compatibility chain.json should exist even if empty
			createMissingChainsJson(chainsPath)
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
	return &types.HostedBlockchains{
		M: m,
		L: sync.Mutex{},
	}
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
	return &types.HostedBlockchains{M: m, L: sync.Mutex{}}
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
	cdc = codec.NewCodec(types2.NewInterfaceRegistry())
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
	crypto.RegisterAmino(cdc.AminoCodec().Amino)
	cryptoamino.RegisterAmino(cdc.AminoCodec().Amino)
	codec.RegisterEvidences(cdc.AminoCodec(), cdc.ProtoCodec())
}

func Credentials(pwd string) string {
	if pwd != "" && strings.TrimSpace(pwd) != "" {
		return strings.TrimSpace(pwd)
	} else {
		bytePassword, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println(err)
		}
		return strings.TrimSpace(string(bytePassword))
	}
}

func Confirmation(pwd string) bool {

	if pwd != "" && strings.TrimSpace(pwd) != "" {
		return true
	} else {
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Println("yes | no")
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading string: ", err.Error())
				return false
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response == "y" || response == "yes" {
				return true
			} else if response == "n" || response == "no" {
				return false
			}
		}
	}

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

func GetDefaultConfig(datadir string) string {

	if datadir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log2.Fatal("could not get home directory for data dir creation: " + err.Error())
		}
		datadir = home + FS + sdk.DefaultDDName
	}
	c := sdk.DefaultConfig(datadir)

	jsonbytes, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return ""
	}

	return string(jsonbytes)
}

func InitAuthToken() {
	var t = sdk.AuthToken{
		Value:  rand.Str(25),
		Issued: time.Now(),
	}
	datadir := GlobalConfig.PocketConfig.DataDir
	configFilepath := datadir + FS + sdk.ConfigDirName + FS + sdk.AuthFileName

	var jsonFile *os.File
	defer jsonFile.Close()

	jsonFile, err := os.OpenFile(configFilepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log2.Fatalf("canot open/create auth token json file: " + err.Error())
	}
	err = jsonFile.Truncate(0)

	b, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		log2.Fatalf("cannot marshal auth token into json: " + err.Error())
	}
	// write to the file
	_, err = jsonFile.Write(b)
	if err != nil {
		log2.Fatalf("cannot write auth token to json file: " + err.Error())
	}

	AuthToken = t
}

func GetAuthTokenFromFile() sdk.AuthToken {
	t := sdk.AuthToken{}
	datadir := GlobalConfig.PocketConfig.DataDir
	configFilepath := datadir + FS + sdk.ConfigDirName + FS + sdk.AuthFileName

	var jsonFile *os.File
	defer jsonFile.Close()

	if _, err := os.Stat(configFilepath); err == nil {
		jsonFile, err = os.OpenFile(configFilepath, os.O_RDWR, os.ModePerm)
		if err != nil {
			log2.Fatalf("cannot open config json file: " + err.Error())
		}
		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log2.Fatalf("cannot read config file: " + err.Error())
		}
		err = json.Unmarshal(b, &t)
		if err != nil {
			log2.Fatalf("cannot read config file into json: " + err.Error())
		}
	}

	return t
}

func createMissingChainsJson(chainsPath string) {
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(chainsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}

	var hostedChainsSlice []types.HostedBlockchain

	hostedChainsSlice = append(hostedChainsSlice, types.HostedBlockchain{
		ID:  "0001",
		URL: "http://localhost:8081/",
	})

	// write to the file
	res, err := json.MarshalIndent(hostedChainsSlice, "", "  ")
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	_, err = jsonFile.Write(res)
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		log2.Fatal(NewInvalidChainsError(err))
	}
}
