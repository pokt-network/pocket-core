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
  "genesis_time": "2020-02-11T08:30:00.000000Z",
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
            "address": "808053795c7b302218a26af6c40f8c39565ebe02",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "uC2Etat3V0/AlbeFJqTRfEBWSo1kFWYVrYCBiILf8r4="
            },
            "account_number": "0",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "e9ee23ea88967a3493c11d783d69b14a8a448f36",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "uomLLP25BiEJqhSWNJ/syzfze61TSn6bNRGgmxkbjX0="
            },
            "account_number": "1",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "99438bc8937b3c5711886ca5c4ed657e17174657",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "ySKNV8+8iBBZ4+B95rvUb9XVokkSYwkKvYXwe1rnFBc="
            },
            "account_number": "2",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "deb8f5b8be1fab076db014ac9ecf92068e616d93",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "8PupreqsBPoFq8C6lA57wm/l0k8bGhMyJPc6oHaxLbo="
            },
            "account_number": "3",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "f31f77c8a882504ef8525e6557351295107f1680",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "2cFggupl7wCwdkOlVsNq8Yp67I3QAz/JhfoB72Dkw5Q="
            },
            "account_number": "4",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "cf37ae72a13de919705990b094b765eac5d8a04c",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "ioq5A4OxGsusKC9UngGvwqCQACrPeb07XKJfq1HoqkE="
            },
            "account_number": "5",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "24486271c197bdd1b58af70d69697aaedf01a569",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "5rbVITCATgOsJoJGJkVrkvTk+Aprav4A3o15PHdXbR4="
            },
            "account_number": "6",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "1b33a9a2f9108c7bf9a60c12e2cf92326175fbd2",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "FaJHHhJufXztMU5C4O5vCp3BenyRXOHZocVIgwRWUSc="
            },
            "account_number": "7",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "75e65c143a9ec7f9c400c266116e93f9d7a3a3bb",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "yY/D3Spj4ej6dfHb24CLOhkSifxQS766ixlWZSxKORI="
            },
            "account_number": "8",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "bb466f281ad6eb17b2dfb6a0995a021d1a957253",
            "coins": [
              {
                "denom": "upokt",
                "amount": "1000000000"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "92rjfwXy3rmF+p6kUogUKg9FgaVBPxfxefthL+8aemY="
            },
            "account_number": "9",
            "sequence": "0"
          }
        },
        {
          "type": "posmint/Account",
          "value": {
            "address": "61698cd28ef417be540db38eb37d559742dee41e",
            "coins": [
              {
                "denom": "upokt",
                "amount": "18446744073709551615"
              }
            ],
            "public_key": {
              "type": "crypto/ed25519_public_key",
              "value": "0ELLVHbMzt2nZTkbXlm8OVy/qlzR4hCo5rfJQnj8bjo="
            },
            "account_number": "10",
            "sequence": "0"
          }
        }
      ]
    },
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
        "session_block_frequency": "60",
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
          "address": "808053795c7b302218a26af6c40f8c39565ebe02",
          "public_key": "b82d84b5ab77574fc095b78526a4d17c40564a8d64156615ad80818882dff2be",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet0:8081",
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
          "address": "e9ee23ea88967a3493c11d783d69b14a8a448f36",
          "public_key": "ba898b2cfdb9062109aa1496349feccb37f37bad534a7e9b3511a09b191b8d7d",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet1:8081",
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
          "address": "99438bc8937b3c5711886ca5c4ed657e17174657",
          "public_key": "c9228d57cfbc881059e3e07de6bbd46fd5d5a2491263090abd85f07b5ae71417",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet2:8081",
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
          "address": "deb8f5b8be1fab076db014ac9ecf92068e616d93",
          "public_key": "f0fba9adeaac04fa05abc0ba940e7bc26fe5d24f1b1a133224f73aa076b12dba",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet3:8081",
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
          "address": "f31f77c8a882504ef8525e6557351295107f1680",
          "public_key": "d9c16082ea65ef00b07643a556c36af18a7aec8dd0033fc985fa01ef60e4c394",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet4:8081",
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
          "address": "cf37ae72a13de919705990b094b765eac5d8a04c",
          "public_key": "8a8ab90383b11acbac282f549e01afc2a090002acf79bd3b5ca25fab51e8aa41",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet5:8081",
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
          "address": "24486271c197bdd1b58af70d69697aaedf01a569",
          "public_key": "e6b6d52130804e03ac26824626456b92f4e4f80a6b6afe00de8d793c77576d1e",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet6:8081",
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
          "address": "1b33a9a2f9108c7bf9a60c12e2cf92326175fbd2",
          "public_key": "15a2471e126e7d7ced314e42e0ee6f0a9dc17a7c915ce1d9a1c5488304565127",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet7:8081",
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
          "address": "75e65c143a9ec7f9c400c266116e93f9d7a3a3bb",
          "public_key": "c98fc3dd2a63e1e8fa75f1dbdb808b3a191289fc504bbeba8b1956652c4a3912",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet8:8081",
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
          "address": "bb466f281ad6eb17b2dfb6a0995a021d1a957253",
          "public_key": "f76ae37f05f2deb985fa9ea45288142a0f4581a5413f17f179fb612fef1a7a66",
          "jailed": false,
          "status": 2,
          "tokens": "1000000000",
          "service_url": "http://www.pocket-core-testnet9:8081",
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
    }
  }
}`
