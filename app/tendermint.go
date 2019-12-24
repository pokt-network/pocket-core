package app

import (
	sdk "github.com/pokt-network/posmint/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"path/filepath"
)

type AppCreator func(log.Logger, dbm.DB, io.Writer) *pocketCoreApp

func NewClient(ctx config, creator AppCreator) (*node.Node, *pocketCoreApp, error) {
	// setup the database
	db, err := openDB(ctx.TmConfig.RootDir)
	if err != nil {
		return nil, nil, err
	}
	// open the tracewriter
	traceWriter, err := openTraceWriter(ctx.TraceWriter)
	if err != nil {
		return nil, nil, err
	}
	// create the instance
	app := creator(ctx.Logger, db, traceWriter)
	// load the node key
	nodeKey, err := p2p.LoadOrGenNodeKey(ctx.TmConfig.NodeKeyFile())
	if err != nil {
		return nil, nil, err
	}
	// upgrade the privVal file
	upgradePrivVal(ctx.TmConfig)
	// create & start tendermint node
	tmNode, err := node.NewNode(
		ctx.TmConfig,
		pvm.LoadOrGenFilePV(ctx.TmConfig.PrivValidatorKeyFile(), ctx.TmConfig.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(ctx.TmConfig),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(ctx.TmConfig.Instrumentation),
		ctx.Logger.With("module", "node"),
	)
	// setup the keybase and tendermint node for the proxy app
	app.SetNodeAndKeybase(tmNode, keybase)
	if err != nil {
		return nil, app, err
	}
	return tmNode, app, nil
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := sdk.NewLevelDB("application", dataDir)
	return db, err
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile != "" {
		w, err = os.OpenFile(
			traceWriterFile,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666,
		)
		return
	}
	return
}

// upgradePrivVal converts old priv_validator.json file (prior to Tendermint 0.28)
// to the new priv_validator_key.json and priv_validator_state.json files.
func upgradePrivVal(config *cfg.Config) {
	if _, err := os.Stat(config.OldPrivValidatorFile()); !os.IsNotExist(err) {
		if oldFilePV, err := pvm.LoadOldFilePV(config.OldPrivValidatorFile()); err == nil {
			oldFilePV.Upgrade(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
		}
	}
}

type config struct {
	TmConfig    *cfg.Config
	Logger      log.Logger
	TraceWriter string
}
