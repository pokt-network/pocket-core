package app

import (
	"encoding/json"
	"fmt"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/config"
	"github.com/pokt-network/posmint/crypto/keys"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	dbm "github.com/tendermint/tm-db"
	"io"
	"io/ioutil"
	"os"
)

func GetKeybase(keybaseName, keybaseDirectory string) keys.Keybase {
	return keys.New(keybaseName, keybaseDirectory)
}

func GetHostedChains(filepath string) (chains types.HostedBlockchains) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	bz, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bz, &chains)
	if err != nil {
		panic(err)
	}
	return
}

func TendermintNode(persistentPeers, seeds, listenAddr string) *node.Node {
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
	cfg := *config.NewConfig(persistentPeers, seeds, listenAddr, false, 0, 40, 10, logger, "")
	n, err := config.NewClient(cfg, newApp)
	if err != nil {
		panic(err)
	}
	return n
}

func CoinbasePassphrase(passphrase string) string {
	return passphrase
}

func GenesisFile(filepath string) string {
	return filepath
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return NewPocketCoreApp(logger, db)
}
