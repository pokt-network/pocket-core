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
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/spf13/cobra"
	con "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/rpc/client"
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
	// passphrase needed for pocket core module
	passphrase string
	// the filepath to the genesis.json
	genesisFP string
	// the default fileseparator based on OS
	fs = string(fp.Separator)
)

func InitApp(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort string, blockTime int) *node.Node {
	pswrd := InitConfig(datadir)
	// setup coinbase password
	if pswrd == "" {
		fmt.Println("Pocket core needs your passphrase to start")
		pswrd = Credentials()
	}
	err := confirmCoinbasePassphrase(pswrd)
	if err != nil {
		panic("Coinbase Password could not be verified: " + err.Error())
	}
	setcoinbasePassphrase(pswrd)
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
			_, err = jsonFile.Write([]byte(testnetGenesis))
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

	c := cfg.Config{
		TmConfig:    newTMConfig,
		Logger:      logger,
		TraceWriter: "",
	}

	var err error
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

func setcoinbasePassphrase(pass string) {
	passphrase = pass
}

func getCoinbasePassphrase() string {
	return passphrase
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
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
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

func newDefaultGenesisState(pubKey crypto.PublicKey) []byte {
	defaultGenesis := module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).DefaultGenesis()
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
	res := types.ModuleCdc.MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res

	j, _ := types.ModuleCdc.MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
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

const testnetGenesis = `{
  "genesis_time": "2020-02-11T22:30:00.00Z",
  "chain_id": "pocket-testnet",
  "consensus_params": {
    "block": {
      "max_bytes": "4000000",
      "max_gas": "-1",
      "time_iota_ms": "1"
    },
    "evidence": {
      "max_age": "1000000"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    }
  },
  "app_hash": "",
  "app_state": {
    "bank": {
      "send_enabled": true
    },
    "params": null,
    "pos": {
      "params": {
        "unstaking_time": "1800000000000",
        "max_validators": "100000",
        "stake_denom": "upokt",
        "stake_minimum": "1000000",
        "session_block_frequency": "30",
        "dao_allocation": "10",
        "proposer_allocation": "1",
        "max_evidence_age": "120000000000",
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "600000000000",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "prevState_total_power": "0",
      "prevState_validator_powers": null,
      "validators": [
        {
          "address": "610cf8a6e8cefbaded845f1c1dc3b10a670be26b",
          "public_key": "1807948ea0041de2a9cd573d0edb073c1eaea60313c364c16c1bcd27629e305b",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node1.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "e6946760d9833f49da39aae9500537bef6f33a7a",
          "public_key": "4ac6202fca022b932be12a5bd51dc8375bfee843f4f90c412e83ad9af1069361",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node2.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "7674a47cc977326f1df6cb92c7b5a2ad36557ea2",
          "public_key": "257943d4255d60f9a042a2cd81ff64b711bedbf72db64d1f84b0e2455ce1dfd1",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node3.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "c7b7b7665d20a7172d0c0aa58237e425f333560a",
          "public_key": "d4448f629a19e4fb68a904a8d879fdd8b1b326d0fff39973f39af737a282be71",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node4.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "f6dc0b244c93232283cd1d8443363946d0a3d77a",
          "public_key": "1c03871c9f6d437a1856cc5141afa7beb1670e82ce692cb7d041d4bc90ab71ad",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node5.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "86209713befeca0807714bcdd5b79e81073faf8f",
          "public_key": "11ac3c35a531ec39f9c5f9164cdf13b19572181dc2048cba666ade0947df6a71",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node6.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "915a58ae437d2c2d6f35ac11b79f42972267700d",
          "public_key": "9290914ab72b4e1d377ac53350996b937ed466a01e3381a7de40282d11501f5b",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node7.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "b3d86cd8ab4aa0cb9861cb795d8d154e685a94cf",
          "public_key": "6b1072e5e5744f3cf9d0e318572f3feb7fdd20e46d67ac15d9607a3a1609bad0",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node8.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "17ca63e4ff7535a40512c550dd0267e519cafc1a",
          "public_key": "05f1b1bb09ddf26b7ba024e458c2712685ecad221beb83534aa8b7f9b19cee75",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node9.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        },
        {
          "address": "f99386c6d7cd42a486c63ccd80f5fbea68759cd7",
          "public_key": "2ea45fa7305a01d87d7471c3cc558f451a32d9dded958c002260069b9eb2249e",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://node10.testnet.pokt.network:8081",
          "chains": [
            "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
            "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
            "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
            "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
            "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
            "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
            "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
            "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
            "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
            "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
            "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
            "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
            "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
            "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
            "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
            "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
            "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
            "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
            "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
            "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
            "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
            "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
            "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
          ],
          "unstaking_time": "0001-01-01T00:00:00Z"
        }
      ],
      "exported": false,
      "dao": {
        "Tokens": "0"
      },
      "signing_infos": {},
      "missed_blocks": {},
      "previous_proposer": ""
    },
    "supply": {
      "supply": []
    },
    "pocketcore": {
      "params": {
        "session_node_count": "5",
        "proof_waiting_period": "3",
        "supported_blockchains": [
          "8ef9a7c67f6f8ad14f82c1f340963951245f912f037a7087f3f2d2f9f9ee38a8",
          "a969144c864bd87a92e974f11aca9d964fb84cf5fb67bcc6583fe91a407a9309",
          "0070eebec778ea95ef9c75551888971c27cce222e00b2f3f79168078b8a77ff9",
          "4ae7539e01ad2c42528b6a697f118a3535e404fe65999b2c6fee506465390367",
          "0de3141aec1e69aea9d45d9156269b81a3ab4ead314fbf45a8007063879e743b",
          "8cf7f8799c5b30d36c86d18f0f4ca041cf1803e0414ed9e9fd3a19ba2f0938ff",
          "10d1290eee169e3970afb106fe5417a11b81676ce1e2119a0292df29f0445d30",
          "d9d669583c2d2a88e54c0120be6f8195b2575192f178f925099813ff9095d139",
          "d9d77bce50d80e70026bd240fb0759f08aab7aee63d0a6d98c545f2b5ae0a0b8",
          "dcc98e38e1edb55a97265efca6c34f21e55f683abdded0aa71df3958a49c8b69",
          "26a2800156f76b66bcb5661f2988a9d09e76caaffd053fe17bf20d251b4cb823",
          "73d8dd1b7d8aa02254e75936b09780447c06729f3e55f7ae5eb94ab732c1ec05",
          "6cbb58da0b05d23022557dd2e479dd5cdf2441f20507b37383467d837ad40f5e",
          "54cb0d71117aa644e74bdea848d61bd2fd410d3d4a3ed92b46b0847769dc132e",
          "cb92cb81d6f72f55114140a7bbe5e0f63524d1200fe63250f58dfe5d907032bf",
          "e458822c5f4d927c29aa4240a34647e11aff75232ccb9ffb50af06dc4469a5fa",
          "0dfcabfb7f810f96cde01d65f775a565d3a60ad9e15575dfe3d188ff506c35a0",
          "866d7183a24fad1d0a32c399cf2a1101f3a3bdfdff999e142bd8f49b2ebc45d4",
          "4c0437dda63eff39f85c60d62ac936045da5e610aca97a3793771e271578c534",
          "773eda9368243afe027062d771b08cebddf22e03451e0eb5ed0ff4460288847e",
          "d5ddbb1ca49249438f552dccfd01918ee1fbdc6457997a142c8cfd144b40cd15",
          "4ecc78e62904c833ad5b727b9abf343a17d0d24fb27e9b5d2dd8c34361c23156",
          "d754973bdeab17eaed47729ee074ad87737c3ce51198263b8c4781568ea39e72"
        ],
        "claim_expiration": "100"
      },
      "proofs": null,
      "claims": null
    },
    "application": {
      "params": {
        "unstaking_time": "1800000000000",
        "max_applications": "18446744073709551615",
        "app_stake_minimum": "1000000",
        "base_relays_per_pokt": "100",
        "stability_adjustment": "0",
        "participation_rate_on": false
      },
      "applications": null,
      "exported": false
    },
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      },
      "Accounts": [
        {
          "type": "posmint/Account",
          "value": {
            "address": "610cf8a6e8cefbaded845f1c1dc3b10a670be26b",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "GAeUjqAEHeKpzVc9DtsHPB6upgMTw2TBbBvNJ2KeMFs="
            },
            "account_number": "0",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "e6946760d9833f49da39aae9500537bef6f33a7a",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "SsYgL8oCK5Mr4Spb1R3IN1v+6EP0+QxBLoOtmvEGk2E="
            },
            "account_number": "1",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "7674a47cc977326f1df6cb92c7b5a2ad36557ea2",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "JXlD1CVdYPmgQqLNgf9ktxG+2/cttk0fhLDiRVzh39E="
            },
            "account_number": "2",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "c7b7b7665d20a7172d0c0aa58237e425f333560a",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "1ESPYpoZ5PtoqQSo2Hn92LGzJtD/85lz85r3N6KCvnE="
            },
            "account_number": "3",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "f6dc0b244c93232283cd1d8443363946d0a3d77a",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "HAOHHJ9tQ3oYVsxRQa+nvrFnDoLOaSy30EHUvJCrca0="
            },
            "account_number": "4",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "86209713befeca0807714bcdd5b79e81073faf8f",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "Eaw8NaUx7Dn5xfkWTN8TsZVyGB3CBIy6ZmreCUffanE="
            },
            "account_number": "5",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "915a58ae437d2c2d6f35ac11b79f42972267700d",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "kpCRSrcrTh03esUzUJlrk37UZqAeM4Gn3kAoLRFQH1s="
            },
            "account_number": "6",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "b3d86cd8ab4aa0cb9861cb795d8d154e685a94cf",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "axBy5eV0Tzz50OMYVy8/63/dIORtZ6wV2WB6OhYJutA="
            },
            "account_number": "7",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "17ca63e4ff7535a40512c550dd0267e519cafc1a",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "BfGxuwnd8mt7oCTkWMJxJoXsrSIb64NTSqi3+bGc7nU="
            },
            "account_number": "8",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "f99386c6d7cd42a486c63ccd80f5fbea68759cd7",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "LqRfpzBaAdh9dHHDzFWPRRoy2d3tlYwAImAGm56yJJ4="
            },
            "account_number": "9",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "cad3b0b8f5b54f0750385c6ca17a5c745d9dba17",
            "coins": [
              {
                "denom": "upokt",
                "amount": "18446744073709551615"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "5gM3j0wP4cpX1UV0GoFQIxIYqj2eL2LAalAF37yjvz0="
            },
            "account_number": "10",
            "sequence": "0"
          }
        }
      ]
    }
  }
}`
