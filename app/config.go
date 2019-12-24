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
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/config"
	kb "github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	tmTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	filepath2 "path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	keybaseName     = "pocket-keybase"
	kbDirName       = "keybase"
	chainsName      = "chains.json"
	dummyChainsHash = "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"
	dummyChainsURL  = "https://foo.bar:8080"
	dummyServiceURL = "https://myPocketNode:8080"
)

var (
	// global pocket core instance
	pcInstance *pocketCoreApp
	// the tendermint node running in process
	tmNode *node.Node
	// global keybase instance
	keybase *kb.Keybase
	// global hosted blockchains instance
	hostedBlockchains *types.HostedBlockchains
	// expose coded to pcInstance module
	Cdc *codec.Codec
	// passphrase needed for pocket core module
	passphrase string
	// the filepath to the genesis.json
	genesisFP string
	// the default fileseparator based on OS
	fs = string(filepath2.Separator)
)

// keybase creation
func CreateKeybase(datadir, passphrase string) error {
	keys := kb.New(keybaseName, datadir+fs+kbDirName)
	_, err := keys.Create(passphrase)
	if err != nil {
		return err
	}
	SetKeybase(&keys)
	return nil
}

// keybase initialization
func InitKeybase(datadir string) error {
	keys := kb.New(keybaseName, datadir+fs+kbDirName)
	kps, err := keys.List()
	if err != nil {
		return err
	}
	if len(kps) == 0 {
		return UninitializedKeybaseError
	}
	SetKeybase(&keys)
	return nil
}

// get the global keybase
func GetKeybase() *kb.Keybase {
	if keybase == nil {
		panic(UninitializedKeybaseError)
	}
	return keybase
}

// set the global keybase
func SetKeybase(k *kb.Keybase) { keybase = k }

// initialize the hosted chains
func InitHostedChains(filepath string) {
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
			panic(err)
		}
		// create dummy input for the file
		res, err := json.MarshalIndent(map[string]types.HostedBlockchain{dummyChainsHash: {Hash: dummyChainsHash, URL: dummyChainsURL}}, "", "  ")
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
	}
	// reopen the file to read into the variable
	jsonFile, err = os.OpenFile(chainsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	bz, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	// return the map
	hostedBlockchains = &types.HostedBlockchains{M: m}
}

// get the hosted chains variable
func GetHostedChains() types.HostedBlockchains {
	if hostedBlockchains == nil || len(hostedBlockchains.M) == 0 {
		panic(InvalidChainsError)
	}
	return *hostedBlockchains
}

// init the tendermint node with a logger and configruation
func InitTendermintNode(rootDir, dataDir, nodeKey, privValKey, privValState, persistentPeers, seeds, listenAddr string) *node.Node {
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
	cfg := *config.NewConfig(rootDir, dataDir, nodeKey, privValKey, privValState, persistentPeers, seeds, listenAddr, true, 0, 40, 10, logger, "")
	var err error
	tmNode, pcInstance, err = NewClient(Config(cfg), newApp)
	if err != nil {
		panic(err)
	}
	return tmNode
}

func newDefaultGenesisState(pubKey crypto.PubKey) []byte {
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
		nodesTypes.Validator{Address: sdk.ValAddress(pubKey.Address()),
			ConsPubKey:   pubKey,
			Status:       sdk.Bonded,
			Chains:       []string{dummyChainsHash},
			ServiceURL:   dummyServiceURL,
			StakedTokens: sdk.NewInt(10000000)})
	res := types.ModuleCdc.MustMarshalJSON(posGenesisState)
	defaultGenesis[nodesTypes.ModuleName] = res
	j, _ := types.ModuleCdc.MarshalJSONIndent(defaultGenesis, "", "    ")
	return j
}

func InitGenesis(filepath string) {
	SetGenesisFilepath(filepath + fs + "genesis.json")
	if _, err := os.Stat(GetGenesisFilePath()); os.IsNotExist(err) {
		keys := GetKeybase()
		kps, err := (*keys).List()
		if err != nil {
			panic(err)
		}
		publicKey := kps[0].PubKey
		// ensure directory path made
		err = os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defaultGenesis := tmTypes.GenesisDoc{
			GenesisTime: time.Time{},
			ChainID:     "pocket-test",
			ConsensusParams: &tmTypes.ConsensusParams{
				Block: tmTypes.BlockParams{
					MaxBytes:   15000,
					MaxGas:     -1,
					TimeIotaMs: 1,
				},
				Evidence: tmTypes.EvidenceParams{
					MaxAge: 1000000,
				},
				Validator: tmTypes.ValidatorParams{
					PubKeyTypes: []string{"ed25519"},
				},
			},
			Validators: nil,
			AppHash:    nil,
			AppState:   newDefaultGenesisState(publicKey),
		}
		// create the genesis path
		var genFP = GetGenesisFilePath()
		// if file exists open, else create and open
		if _, err := os.Stat(genFP); err == nil {
			return
		} else if os.IsNotExist(err) {
			// if does not exist create one
			jsonFile, err := os.OpenFile(genFP, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			bz, err := Cdc.MarshalJSONIndent(defaultGenesis, "", "    ")
			if err != nil {
				panic(err)
			}
			// write to the file
			_, err = jsonFile.Write(bz)
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

func ConfirmCoinbasePassword(pswrd string) error {
	keys := GetKeybase()
	kps, err := (*keys).List()
	if err != nil {
		panic(err)
	}
	kp := kps[0]
	err = (*keys).Update(kp.GetAddress(), pswrd, pswrd)
	if err != nil {
		return err
	}
	return nil
}

// get the in process tendermint node
func GetTendermintNode() (*node.Node, error) {
	if tmNode == nil {
		return nil, UninitializedTendermintError
	}
	return tmNode, nil
}

func SetCoinbasePassphrase(pass string) {
	passphrase = pass
}

func GetCoinbasePassphrase() string {
	return passphrase
}

func SetGenesisFilepath(filepath string) {
	genesisFP = filepath
}

func GetGenesisFilePath() string {
	return genesisFP
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) *pocketCoreApp {
	return NewPocketCoreApp(logger, db)
}

func InitTendermint(datadir, persistentPeers, seeds string) {
	tmNode := InitTendermintNode(datadir, "", "node_key.json", "priv_val_key.json",
		"priv_val_state.json", persistentPeers, seeds, "0.0.0.0:46656")
	if err := tmNode.Start(); err != nil {
		panic(err)
	}
	// We trap kill signals (2,3,15,9)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Kill,
		os.Interrupt)

	defer func() {
		sig := <-signalChannel
		tmNode, _ := GetTendermintNode()
		err := tmNode.Stop()
		if err != nil {
			panic(err)
		}
		message := fmt.Sprintf("Exit signal %s received\n", sig)
		fmt.Println(message)
		os.Exit(3)
	}()
}

func InitDataDirectory(datadir string) string {
	// check for empty data_dir
	if datadir == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// set the default data directory
		datadir = home + fs + ".pocket"
	}
	// create the folder if not already created
	err := os.MkdirAll(datadir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return datadir
}

func InitKeyfiles(datadir string) string {
	var password string
	if err := InitKeybase(datadir); err != nil {
		fmt.Println("Initializing keybase: enter coinbase passphrase")
		password = Credentials()
		err := CreateKeybase(datadir, password)
		if err != nil {
			panic(err)
		}
		kb := GetKeybase()
		keypairs, err := (*kb).List()
		if err != nil {
			panic(err)
		}
		coinbaseKeypair := keypairs[0]
		res, err := (*kb).ExportPrivateKeyObject(coinbaseKeypair.GetAddress(), password)
		if err != nil {
			panic(err)
		}
		privValKey := privval.FilePVKey{
			Address: res.PubKey().Address(),
			PubKey:  res.PubKey(),
			PrivKey: res,
		}
		privValState := privval.FilePVLastSignState{}
		nodeKey := p2p.NodeKey{
			PrivKey: res,
		}
		pvkBz, err := Cdc.MarshalJSONIndent(privValKey, "", "  ")
		if err != nil {
			panic(err)
		}
		nkBz, err := Cdc.MarshalJSONIndent(nodeKey, "", "  ")
		if err != nil {
			panic(err)
		}
		pvsBz, err := Cdc.MarshalJSONIndent(privValState, "", "  ")
		if err != nil {
			panic(err)
		}
		pvFile, err := os.OpenFile(datadir+fs+"priv_val_key.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		_, err = pvFile.Write(pvkBz)
		if err != nil {
			panic(err)
		}
		pvStateFile, err := os.OpenFile(datadir+fs+"priv_val_state.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		_, err = pvStateFile.Write(pvsBz)
		if err != nil {
			panic(err)
		}
		nkFile, err := os.OpenFile(datadir+fs+"node_key.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		_, err = nkFile.Write(nkBz)
		if err != nil {
			panic(err)
		}
	}
	return password
}

func SetCodec() {
	// create a new codec
	Cdc = codec.New()
	// register all of the app module types
	module.NewBasicManager(
		apps.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nodes.AppModuleBasic{},
		supply.AppModuleBasic{},
		pocket.AppModuleBasic{},
	).RegisterCodec(Cdc)
	// register the sdk types
	sdk.RegisterCodec(Cdc)
	// register the crypto types
	codec.RegisterCrypto(Cdc)
}

func Credentials() string {
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(bytePassword))
}
