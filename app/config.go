package app

import (
	"encoding/json"
	"fmt"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/mitchellh/go-homedir"
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
	"io"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	KeybaseName       = "pocket-keybase"
	privValKeyName    = "priv_val_key.json"
	privValStateName  = "priv_val_state.json"
	nodeKeyName       = "node_key.json"
	KBDirectoryName   = "keybase"
	chainsName        = "chains.json"
	dummyChainsHash   = "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"
	dummyChainsURL    = "https://foo.bar:8080"
	dummyServiceURL   = "0.0.0.0:8081"
	defaultTMURI      = "tcp://localhost:26657"
	defaultNodeKey    = "node_key.json"
	defaultValKey     = "priv_val_key.json"
	defaultValState   = "priv_val_state.json"
	defaultListenAddr = "tcp://0.0.0.0:"
)

var (
	datadir string
	// expose coded to pcInstance module
	cdc *codec.Codec
	// tendermint node uri
	tmNodeURI string
	// privateKey needed for pocket core module
	privateKeyPath string
	// the filepath to the genesis.json
	genesisFP string
	// the default fileseparator based on OS
	fs = string(fp.Separator)
	// app instance currently running
	pca *pocketCoreApp
)

func InitApp(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort string, blockTime int) *node.Node {
	InitConfig(datadir)
	// set tendermint node
	SetTMNode(tmNode)
	// init the tendermint node
	return InitTendermint(persistentPeers, seeds, tmRPCPort, tmPeersPort, blockTime)
}

func InitConfig(datadir string) string {
	// setup the codec
	MakeCodec()
	// setup data directory
	InitDataDirectory(datadir)
	// init the keyfiles
	pswrd := InitKeyfiles()
	// init cache
	InitPocketCoreCache(datadir)
	// init genesis
	InitGenesis()
	return pswrd
}

func InitGenesis() {
	setGenesisPath(getDataDir() + fs + "config" + fs + "genesis.json")
	if _, err := os.Stat(genesisPath()); os.IsNotExist(err) {
		// ensure directory path made
		err = os.MkdirAll(datadir+fs+"config", os.ModePerm)
		if err != nil {
			panic(err)
		}
		// create the genesis path
		var genFP = genesisPath()
		// if file exists open, else create and open
		if _, err := os.Stat(genFP); err == nil {
			return
		} else if os.IsNotExist(err) {
			// if does not exist create one
			jsonFile, err := os.OpenFile(genFP, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			// write to the file
			_, err = jsonFile.Write(newDefaultGenesisState())
			if err != nil {
				panic(err)
			}
			// close the file
			err = jsonFile.Close()
			if err != nil {
				panic(err)
			}
		}
	}
}

func InitTendermint(persistentPeers, seeds, tmRPCPort, tmPeersPort string, blockTime int) *node.Node {
	datadir := getDataDir()
	// setup the logger
	logger := log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			panic(fmt.Sprintf("expected level key to be first, got %v", keyvals[0]))
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

	// setup tendermint node config
	newTMConfig := con.DefaultConfig()
	newTMConfig.SetRoot(datadir)
	//newTMConfig.DBPath = datadir
	newTMConfig.NodeKey = defaultNodeKey
	newTMConfig.PrivValidatorKey = defaultValKey
	newTMConfig.PrivValidatorState = defaultValState
	newTMConfig.P2P.AddrBookStrict = false
	newTMConfig.RPC.ListenAddress = defaultListenAddr + tmRPCPort
	newTMConfig.P2P.ListenAddress = defaultListenAddr + tmPeersPort // Node listen address. (0.0.0.0:0 means any interface, any port)
	newTMConfig.P2P.PersistentPeers = persistentPeers               // Comma-delimited ID@host:port persistent peers
	newTMConfig.P2P.Seeds = seeds                                   // Comma-delimited ID@host:port seed nodes
	newTMConfig.Consensus.CreateEmptyBlocks = true                  // Set this to false to only produce blocks when there are txs or when the AppHash changes
	newTMConfig.Consensus.CreateEmptyBlocksInterval = time.Duration(blockTime) * time.Minute
	newTMConfig.Consensus.TimeoutCommit = time.Duration(blockTime) * time.Minute
	newTMConfig.P2P.MaxNumInboundPeers = 40
	newTMConfig.P2P.MaxNumOutboundPeers = 10
	newTMConfig.LogLevel = "*:info, *:error"
	logger, err := flags.ParseLogLevel(newTMConfig.LogLevel, logger, "info")
	if err != nil {
		panic(err)
	}
	c := cfg.Config{
		TmConfig:    newTMConfig,
		Logger:      logger,
		TraceWriter: "",
	}
	tmNode, app, err := NewClient(config(c), func(logger log.Logger, db dbm.DB, _ io.Writer) *pocketCoreApp {
		return NewPocketCoreApp(logger, db, baseapp.SetPruning(store.PruneNothing))
	})
	if err != nil {
		panic(err)
	}
	if err := tmNode.Start(); err != nil {
		panic(err)
	}
	app.SetTendermintNode(tmNode)
	pca = app
	return tmNode
}

func InitDataDirectory(d string) string {
	// check for empty data_dir
	if d == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// set the default data directory
		d = home + fs + ".pocket"
	}
	// create the folder if not already created
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		panic(err)
	}
	datadir = d
	return d
}

func InitKeyfiles() string {
	var password string
	datadir := getDataDir()

	//Check if privvalkey file exist
	if _, err := os.Stat(datadir + fs + privValKeyName); err != nil {

		//if not exist continue creating as other files may be missing
		if os.IsNotExist(err) {

			password = initFiles(password, datadir)

		} else {
			//panic on other errors
			panic(err)
		}

	} else {
		// file exist so we can load pk from file.
		file, _ := loadPKFromFile(getDataDir() + fs + privValKeyName)
		types.InitPvKeyFile(file)
	}

	return password
}

func initFiles(password string, datadir string) string {
	if _, err := GetKeybase(); err != nil {
		fmt.Println("Initializing keybase: enter coinbase passphrase")
		password = Credentials()
		if password == "" {
			panic("you must have a coinbase password")
		}
		err := newKeybase(password)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(datadir + fs + privValKeyName); err != nil {
		if os.IsNotExist(err) {
			password = privValKey(password)
		} else {
			panic(err)
		}
	}
	if _, err := os.Stat(datadir + fs + privValStateName); err != nil {
		if os.IsNotExist(err) {
			privValState()
		} else {
			panic(err)
		}
	}
	if _, err := os.Stat(datadir + fs + nodeKeyName); err != nil {
		if os.IsNotExist(err) {
			nodeKey(password)
		} else {
			panic(err)
		}
	}
	return password
}

func InitPocketCoreCache(dataDir string) {
	types.InitCache(dataDir, dataDir, dbm.GoLevelDBBackend, dbm.GoLevelDBBackend, 100, 100)
}

// get the global keybase
func MustGetKeybase() kb.Keybase {
	keys, err := GetKeybase()
	if err != nil {
		panic(err)
	}
	return keys
}

// get the global keybase
func GetKeybase() (kb.Keybase, error) {
	keys := kb.New(KeybaseName, getDataDir()+fs+KBDirectoryName)
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

func privValKey(password string) string {
	keys := MustGetKeybase()
	coinbaseKeypair, err := keys.GetCoinbase()
	if err != nil {
		panic(err)
	}
	if password == "" {
		fmt.Println("Initializing keyfiles: enter coinbase passphrase")
		password = Credentials()
	}
	res, err := (keys).ExportPrivateKeyObject(coinbaseKeypair.GetAddress(), password)
	if err != nil {
		panic(err)
	}
	privValKey := privval.FilePVKey{
		Address: res.PubKey().Address(),
		PubKey:  res.PubKey(),
		PrivKey: res.PrivKey(),
	}
	pvkBz, err := cdc.MarshalJSONIndent(privValKey, "", "  ")
	if err != nil {
		panic(err)
	}
	pvFile, err := os.OpenFile(datadir+fs+privValKeyName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	_, err = pvFile.Write(pvkBz)
	if err != nil {
		panic(err)
	}
	types.InitPvKeyFile(privValKey)

	return password
}

func nodeKey(password string) {
	keys := MustGetKeybase()
	coinbaseKeypair, err := keys.GetCoinbase()
	if err != nil {
		panic(err)
	}
	res, err := (keys).ExportPrivateKeyObject(coinbaseKeypair.GetAddress(), password)
	if err != nil {
		panic(err)
	}
	nodeKey := p2p.NodeKey{
		PrivKey: res.PrivKey(),
	}
	pvkBz, err := cdc.MarshalJSONIndent(nodeKey, "", "  ")
	if err != nil {
		panic(err)
	}
	pvFile, err := os.OpenFile(datadir+fs+nodeKeyName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	_, err = pvFile.Write(pvkBz)
	if err != nil {
		panic(err)
	}
}

func privValState() {
	pvkBz, err := cdc.MarshalJSONIndent(privval.FilePVLastSignState{}, "", "  ")
	if err != nil {
		panic(err)
	}
	pvFile, err := os.OpenFile(datadir+fs+privValStateName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	_, err = pvFile.Write(pvkBz)
	if err != nil {
		panic(err)
	}
}

func getTMClient() client.Client {
	if tmNodeURI == "" {
		return client.NewHTTP(defaultTMURI, "/websocket")
	}
	return client.NewHTTP(tmNodeURI, "/websocket")
}

// get the hosted chains variable
func getHostedChains() types.HostedBlockchains {
	filepath := getDataDir() + fs + "config"
	// create the chains path
	var chainsPath = filepath + fs + chainsName
	// ensure directory path made
	err := os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	// if file exists open, else create and open
	var jsonFile *os.File
	var bz []byte
	if _, err := os.Stat(chainsPath); err == nil {
		// if file exists
	} else if os.IsNotExist(err) {
		// if does not exist create one
		jsonFile, err = os.OpenFile(chainsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(NewInvalidChainsError(err))
		}
		// create dummy input for the file
		res, err := json.MarshalIndent(map[string]types.HostedBlockchain{dummyChainsHash: {Hash: dummyChainsHash, URL: dummyChainsURL}}, "", "  ")
		if err != nil {
			panic(NewInvalidChainsError(err))
		}
		// write to the file
		_, err = jsonFile.Write(res)
		if err != nil {
			panic(NewInvalidChainsError(err))
		}
		// close the file
		err = jsonFile.Close()
		if err != nil {
			panic(NewInvalidChainsError(err))
		}
	}
	// reopen the file to read into the variable
	jsonFile, err = os.OpenFile(chainsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(NewInvalidChainsError(err))
	}
	bz, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(NewInvalidChainsError(err))
	}
	// unmarshal into the structure
	m := map[string]types.HostedBlockchain{}
	err = json.Unmarshal(bz, &m)
	if err != nil {
		panic(NewInvalidChainsError(err))
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		panic(NewInvalidChainsError(err))
	}
	// return the map
	return types.HostedBlockchains{M: m}
}

func confirmCoinbasePassphrase(pswrd string) error {
	keys := MustGetKeybase()
	kp, err := keys.GetCoinbase()
	if err != nil {
		return err
	}
	err = (keys).Update(kp.GetAddress(), pswrd, pswrd)
	if err != nil {
		return err
	}
	return nil
}

func SetTMNode(n string) {
	tmNodeURI = n
}

func setGenesisPath(filepath string) {
	genesisFP = filepath
}

func genesisPath() string {
	return genesisFP
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
		panic(err)
	}
	return strings.TrimSpace(string(bytePassword))
}

func getDataDir() string {
	InitDataDirectory(datadir)
	return datadir
}

// keybase creation
func newKeybase(passphrase string) error {
	keys := kb.New(KeybaseName, getDataDir()+fs+KBDirectoryName)
	_, err := keys.Create(passphrase)
	if err != nil {
		return err
	}
	return nil
}

func newDefaultGenesisState() []byte {
	keyb, err := GetKeybase()
	if err != nil {
		panic(err)
	}
	cb, err := keyb.GetCoinbase()
	if err != nil {
		panic(err)
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
		PubKey:  nil,
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
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
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
	acl.SetOwner("pos/MaxValidators", addr)
	acl.SetOwner("pos/ProposerPercentage", addr)
	acl.SetOwner("application/StabilityAdjustment", addr)
	acl.SetOwner("application/AppUnstakingTime", addr)
	acl.SetOwner("application/ParticipationRateOn", addr)
	acl.SetOwner("pos/MaxEvidenceAge", addr)
	acl.SetOwner("pos/MinSignedPerWindow", addr)
	acl.SetOwner("pos/StakeMinimum", addr)
	acl.SetOwner("pos/UnstakingTime", addr)
	acl.SetOwner("application/BaseRelaysPerPOKT", addr)
	acl.SetOwner("pocketcore/ClaimSubmissionWindow", addr)
	acl.SetOwner("pos/DAOAllocation", addr)
	acl.SetOwner("pos/SignedBlocksWindow", addr)
	acl.SetOwner("pos/SessionBlockFrequency", addr)
	acl.SetOwner("application/MaxApplications", addr)
	acl.SetOwner("gov/daoOwner", addr)
	acl.SetOwner("gov/upgrade", addr)
	return acl
}

// XXX: this is totally unsafe.
// it's only suitable for testnets.
func ResetWorldState(cmd *cobra.Command, args []string) {
	// setup the logger
	logger := log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			panic(fmt.Sprintf("expected level key to be first, got %v", keyvals[0]))
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

	// setup tendermint node config
	newTMConfig := con.DefaultConfig()
	newTMConfig.SetRoot(getDataDir())

	ResetAll(newTMConfig.DBDir(), newTMConfig.P2P.AddrBookFile(), getDataDir()+fs+privValKeyName,
		getDataDir()+fs+privValStateName, logger)
}

var keepAddrBook = false

// ResetAll removes address book files plus all data, and resets the privValdiator data.
// Exported so other CLI tools can use it.
func ResetAll(dbDir, addrBookFile, privValKeyFile, privValStateFile string, logger log.Logger) {
	if keepAddrBook {
		logger.Info("The address book remains intact")
	} else {
		removeAddrBook(addrBookFile, logger)
	}
	if err := os.RemoveAll(dbDir); err == nil {
		logger.Info("Removed all blockchain history", "dir", dbDir)
	} else {
		logger.Error("Error removing all blockchain history", "dir", dbDir, "err", err)
	}
	// recreate the dbDir since the privVal state needs to live there
	cmn.EnsureDir(dbDir, 0700)
	resetFilePV(privValKeyFile, privValStateFile, logger)
}

func resetFilePV(privValKeyFile, privValStateFile string, logger log.Logger) {

	if _, err := os.Stat(privValKeyFile); err == nil {
		os.Remove(privValKeyFile)
		os.Remove(privValStateFile)
		os.Remove(getDataDir() + fs + nodeKeyName)
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
