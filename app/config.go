package app

import (
	"encoding/json"
	"errors"
	"fmt"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/config"
	kb "github.com/pokt-network/posmint/crypto/keys"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	dbm "github.com/tendermint/tm-db"
	"io"
	"io/ioutil"
	"os"
	filepath2 "path/filepath"
)

var (
	app *pocketCoreApp
	// config variables
	tmNode            *node.Node
	keybase           *kb.Keybase
	hostedBlockchains *types.HostedBlockchains
	// expose coded to app module
	Cdc *codec.Codec
	// expose nodes, apps, and core module for tx
	passphrase string
	genesisFP  string
)

var (
	fs = string(filepath2.Separator)
)

func GetKeybase() *kb.Keybase {
	return keybase
}

func SetKeybase(k *kb.Keybase) {
	keybase = k
}

func InitKeybase(datadir string) error {
	keys := kb.New("pocket-keybase", datadir+fs+"keybase")
	kps, err := keys.List()
	if err != nil {
		return err
	}
	if len(kps) == 0 {
		return errors.New("uninitialized keybase")
	}
	SetKeybase(&keys)
	return nil
}

func CreateKeybase(datadir, passphrase string) error {
	keys := kb.New("pocket-keybase", datadir+fs+"keybase")
	_, err := keys.Create(passphrase)
	if err != nil {
		return err
	}
	SetKeybase(&keys)
	return nil
}

func GetHostedChains() (types.HostedBlockchains, error) {
	if hostedBlockchains == nil || len(hostedBlockchains.M) == 0 {
		return types.HostedBlockchains{}, errors.New("nil or empty hosted chains object")
	}
	return *hostedBlockchains, nil
}

func InitHostedChains(filepath string) (chains types.HostedBlockchains) {
	var chainsPath = filepath + fs + "chains.json"
	err := os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	var jsonFile *os.File
	if _, err := os.Stat(chainsPath); err == nil {
		// if file exists
		jsonFile, err = os.OpenFile(filepath+fs+"chains.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		// if does not exist create one
		jsonFile, err = os.OpenFile(filepath+fs+"chains.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		res, err := json.MarshalIndent(map[string]types.HostedBlockchain{"36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80": {
			Hash: "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80",
			URL:  "ttps://foo.bar:8080/lolo?token=lolo",
		},
		}, "", "  ")
		if err != nil {
			panic(err)
		}
		_, err = jsonFile.Write(res)
		if err != nil {
			panic(err)
		}
		err = jsonFile.Close()
		if err != nil {
			panic(err)
		}
		jsonFile, err = os.OpenFile(filepath+fs+"chains.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	bz, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	m := map[string]types.HostedBlockchain{}
	err = json.Unmarshal(bz, &m)
	if err != nil {
		panic("invalid chains.json: " + err.Error())
	}
	hostedBlockchains = &types.HostedBlockchains{
		M: m,
	}
	jsonFile.Close()
	return *hostedBlockchains
}

func InitTendermintNode(rootDir, dataDir, nodeKey, privValKey, privValState, persistentPeers, seeds, listenAddr string) *node.Node {
	// Color by level value
	colorFn := func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			panic(fmt.Sprintf("expected level key to be first, got %v", keyvals[0]))
		}
		switch keyvals[1].(kitlevel.Value).String() {
		case "info":
			return term.FgBgColor{Fg: term.Blue}
		case "debug":
			return term.FgBgColor{Fg: term.DarkMagenta}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	}
	logger := log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)
	cfg := *config.NewConfig(rootDir, dataDir, nodeKey, privValKey, privValState, persistentPeers, seeds, listenAddr, true, 0, 40, 10, logger, "")
	var err error
	tmNode, app, err = NewClient(Config(cfg), newApp)
	if err != nil {
		panic(err)
	}
	return tmNode
}

func GetTendermintNode() (*node.Node, error) {
	if tmNode == nil {
		return nil, errors.New("uninitialized tendermint node")
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

func GetApp() *pocketCoreApp {
	return app
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) *pocketCoreApp {
	return NewPocketCoreApp(logger, db)
}

func SetCodec() {
	Cdc = MakeCodec()
}
