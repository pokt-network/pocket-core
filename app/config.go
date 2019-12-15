package app

import (
	"encoding/json"
	"errors"
	"fmt"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/config"
	"github.com/pokt-network/posmint/crypto/keys"
	kb "github.com/pokt-network/posmint/crypto/keys"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	dbm "github.com/tendermint/tm-db"
	"io"
	"io/ioutil"
	"os"
)

var (
	pocketCoreInstance *pocketCoreApp
	// config variables
	tmNode            *node.Node
	keybase           *kb.Keybase
	hostedBlockchains *types.HostedBlockchains
	passphrase        string
	genesisFP         string
)

func GetKeybase() (keys.Keybase, error) {
	kp, err := (*keybase).List()
	if err != nil {
		return nil, err
	}
	if len(kp) < 1 {
		return nil, errors.New("uninitialized keybase")
	}
	return *keybase, nil
}

func InitKeybase(dataDirectoryPath, coinbasePassphrase string) keys.Keybase {
	k := keys.New("pocket-keybase", dataDirectoryPath+string(os.PathSeparator)+"keybase"+string(os.PathSeparator))
	_, err := k.Create(coinbasePassphrase)
	if err != nil {
		panic(err)
	}
	keybase = &k
	return k
}

func GetHostedChains() (types.HostedBlockchains, error) {
	if hostedBlockchains == nil || len(hostedBlockchains.M) == 0 {
		return types.HostedBlockchains{}, errors.New("nil or empty hosted chains object")
	}
	return *hostedBlockchains, nil
}

func InitHostedChains(filepath string) (chains types.HostedBlockchains) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	bz, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	m := map[string]types.HostedBlockchain{}
	err = json.Unmarshal(bz, &m)
	if err != nil {
		panic(err)
	}
	hostedBlockchains = &types.HostedBlockchains{
		M: m,
	}
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
	cfg := *config.NewConfig(rootDir, dataDir, nodeKey, privValKey, privValState, persistentPeers, seeds, listenAddr, false, 0, 40, 10, logger, "")
	var err error
	tmNode, pocketCoreInstance, err = NewClient(Config(cfg), newApp)
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
	return pocketCoreInstance
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) *pocketCoreApp {
	return NewPocketCoreApp(logger, db)
}
