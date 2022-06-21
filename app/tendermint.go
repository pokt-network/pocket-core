package app

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
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

type AppCreator func(log.Logger, dbm.DB, io.Writer) *PocketCoreApp

func loadFilePVWithConfig(c config) *pvm.FilePVLite {
	privValPath := c.TmConfig.PrivValidatorKeyFile()
	privStatePath := c.TmConfig.PrivValidatorStateFile()
	if GlobalConfig.PocketConfig.LeanPocket {
		return pvm.LoadOrGenFilePV(privValPath, privStatePath)
	}
	legacyFilePV := pvm.LoadOrGenFilePVLegacy(privValPath, privStatePath)
	return &pvm.FilePVLite{
		Key:           []pvm.FilePVKey{legacyFilePV.Key},
		LastSignState: []pvm.FilePVLastSignState{legacyFilePV.LastSignState},
		KeyFilepath:   privValPath,
		StateFilepath: privStatePath,
	}
}


func loadValidatorsLean(c config, tmNode *node.Node) error {
	// add a read to nodes.json -> convert with setValidators()
	validators := loadFilePVWithConfig(c)
	err := InitNodesLean() // reinitialize lean nodes
	tmNode.ConsensusState().SetPrivValidators(validators) // set new lean nodes
	if err != nil {
		return err
	}
	return nil
}

// hotReloadValidatorsLean - spins off a goroutine that reads from validator files  TODO: add load file with error
func hotReloadValidatorsLean(c config, tmNode *node.Node) {
	go func() {
		for {
			loadValidatorsLean(c, tmNode)
			// init light node, but add removal code for validators that aren't part of it.
			time.Sleep(time.Minute * 1)
		}
	}()
}

func NewClient(c config, creator AppCreator) (*node.Node, *PocketCoreApp, error) {
	// setup the database
	appDB, err := OpenApplicationDB(GlobalConfig)
	if err != nil {
		return nil, nil, err
	}
	// setup the transaction indexer
	txDB, err := OpenTxIndexerDB(GlobalConfig)
	if err != nil {
		return nil, nil, err
	}
	transactionIndexer := sdk.NewTransactionIndexer(txDB)
	// open the tracewriter
	traceWriter, err := openTraceWriter(c.TraceWriter)
	if err != nil {
		return nil, nil, err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(c.TmConfig.NodeKeyFile())
	if err != nil {
		return nil, nil, err
	}

	// upgrade the privVal file

	app := creator(c.Logger, appDB, traceWriter)
	PCA = app
	// create & start tendermint node
	tmNode, err := node.NewNode(app,
		c.TmConfig,
		codec.GetCodecUpgradeHeight(),
		loadFilePVWithConfig(c),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		transactionIndexer,
		node.DefaultGenesisDocProviderFunc(c.TmConfig),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(c.TmConfig.Instrumentation),
		c.Logger.With("module", "node"),
	)



	if err != nil {
		return nil, nil, err
	}

	if GlobalConfig.PocketConfig.LeanPocket {
		go hotReloadValidatorsLean(c, tmNode)
	}

	return tmNode, app, nil
}

func OpenApplicationDB(config sdk.Config) (dbm.DB, error) {
	dataDir := filepath.Join(config.TendermintConfig.RootDir, GlobalConfig.TendermintConfig.DBPath)
	return sdk.NewLevelDB(sdk.ApplicationDBName, dataDir, config.TendermintConfig.LevelDBOptions.ToGoLevelDBOpts())
}

func OpenTxIndexerDB(config sdk.Config) (dbm.DB, error) {
	dataDir := filepath.Join(config.TendermintConfig.RootDir, GlobalConfig.TendermintConfig.DBPath)
	return sdk.NewLevelDB(sdk.TransactionIndexerDBName, dataDir, config.TendermintConfig.LevelDBOptions.ToGoLevelDBOpts())
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

//// upgradePrivVal converts old priv_validator.json file (prior to Tendermint 0.28)
//// to the new priv_validator_key.json and priv_validator_state.json files.
//func upgradePrivVal(config *cfg.Config) {
//	if _, err := os.Stat(config.OldPrivValidatorFile()); !os.IsNotExist(err) {
//		if oldFilePV, err := pvm.LoadOldFilePV(config.OldPrivValidatorFile()); err == nil {
//			oldFilePV.Upgrade(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
//		}
//	}
//}

type config struct {
	TmConfig    *cfg.Config
	Logger      log.Logger
	TraceWriter string
}

//func modifyPrivValidatorsFile(config *cfg.Config, rollbackHeight int64) error {
//	var sig []byte
//	filePv := pvm.LoadOrGenFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
//	filePv.LastSignState.Height = rollbackHeight
//	filePv.LastSignState.Round = 0
//	filePv.LastSignState.Step = 0
//	filePv.LastSignState.Signature = sig
//	filePv.LastSignState.SignBytes = nil
//	filePv.Save()
//	return nil
//}
