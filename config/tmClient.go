package config

import (
	sdk "github.com/pokt-network/pocket-core/types"
	abci "github.com/tendermint/tendermint/abci/types"
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

func NewClient(ctx Config, appCreator AppCreator) (*node.Node, error) {
	config := ctx.TmConfig
	home := config.RootDir
	db, err := openDB(home)
	if err != nil {
		return nil, err
	}

	traceWriter, err := openTraceWriter(ctx.TraceWriter)
	if err != nil {
		return nil, err
	}
	app := appCreator(ctx.Logger, db, traceWriter)

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	UpgradeOldPrivValFile(config)

	txIndexer, err := node.CreateTxIndexer(config, node.DefaultDBProvider)
	if err != nil {
		return nil, err
	}
	blockStore, stateDB, err := node.InitDBs(config, node.DefaultDBProvider)
	if err != nil {
		return nil, err
	}
	// create & start tendermint node
	tmNode, err := node.NewNode(
		config,
		pvm.LoadOrGenFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(config.Instrumentation),
		ctx.Logger.With("module", "node"),
		txIndexer,
		blockStore,
		stateDB,
	)
	if err != nil {
		return nil, err
	}

	if err := tmNode.Start(); err != nil {
		return nil, err
	}
	return tmNode, nil
}

type (
	// AppCreator is a function that allows us to lazily initialize an
	// application using various configurations.
	AppCreator func(log.Logger, dbm.DB, io.Writer) abci.Application
)

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

// UpgradeOldPrivValFile converts old priv_validator.json file (prior to Tendermint 0.28)
// to the new priv_validator_key.json and priv_validator_state.json files.
func UpgradeOldPrivValFile(config *cfg.Config) {
	if _, err := os.Stat(config.OldPrivValidatorFile()); !os.IsNotExist(err) {
		if oldFilePV, err := pvm.LoadOldFilePV(config.OldPrivValidatorFile()); err == nil {
			oldFilePV.Upgrade(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
		}
	}
}
