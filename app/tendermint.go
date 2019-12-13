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
	"time"
)

func NewClient(ctx Config, creator AppCreator) (*node.Node, *pocketCoreApp, error) {
	config := ctx.TmConfig
	home := config.RootDir
	db, err := openDB(home)
	if err != nil {
		return nil, nil, err
	}

	traceWriter, err := openTraceWriter(ctx.TraceWriter)
	if err != nil {
		return nil, nil, err
	}

	app := creator(ctx.Logger, db, traceWriter)

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, nil, err
	}

	UpgradeOldPrivValFile(config)

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
	)
	app.SetTendermintNode(tmNode) // todo
	if err != nil {
		return nil, app, err
	}

	//if err := tmNode.Start(); err != nil {
	//	return nil, err
	//}
	return tmNode, app, nil
}


type (
	AppCreator func(log.Logger, dbm.DB, io.Writer) *pocketCoreApp
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


type Config struct {
	TmConfig    *cfg.Config
	Logger      log.Logger
	TraceWriter string
}

func NewDefaultConfig() *Config {
	return &Config{
		cfg.DefaultConfig(),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		"",
	}
}

func NewConfig(rootDir, datadir, nodekey, privValKey, privValState, persistentPeers, seeds, listenAddr string, createEmptyBlocks bool, createEmptyBlocksInterval time.Duration,
	MaxNumberInboundPeers, MaxNumberOutboundPeers int, logger log.Logger, traceWriterPath string) *Config {
	// setup tendermint node config
	newTMConfig := cfg.DefaultConfig()
	newTMConfig.RootDir = rootDir
	newTMConfig.DBPath = datadir
	newTMConfig.NodeKey = nodekey
	newTMConfig.PrivValidatorKey = privValKey
	newTMConfig.PrivValidatorState = privValState
	newTMConfig.P2P.ListenAddress = listenAddr                  // node listen address. (0.0.0.0:0 means any interface, any port)
	newTMConfig.P2P.PersistentPeers = persistentPeers           // Comma-delimited ID@host:port persistent peers
	newTMConfig.P2P.Seeds = seeds                               // Comma-delimited ID@host:port seed nodes
	newTMConfig.Consensus.CreateEmptyBlocks = createEmptyBlocks // Set this to false to only produce blocks when there are txs or when the AppHash changes
	newTMConfig.Consensus.CreateEmptyBlocksInterval = createEmptyBlocksInterval
	newTMConfig.P2P.MaxNumInboundPeers = MaxNumberInboundPeers
	newTMConfig.P2P.MaxNumOutboundPeers = MaxNumberOutboundPeers

	return &Config{newTMConfig, logger, traceWriterPath}
}
