package app

import (
	"errors"
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

// loadFilePVWithConfig returns an array of pvkeys & last sign state for leanpokt or constructs an array of pv keys & lastsignstate if using pre leanpokt to maintain backwards compability
func loadFilePVWithConfig(c config) *pvm.FilePVLean {
	privValPath := c.TmConfig.PrivValidatorKeyFile()
	privStatePath := c.TmConfig.PrivValidatorStateFile()
	if GlobalConfig.PocketConfig.LeanPocket {
		return pvm.LoadOrGenFilePVLean(privValPath, privStatePath)
	}
	globalFilePV := pvm.LoadOrGenFilePV(privValPath, privStatePath)
	return &pvm.FilePVLean{
		Keys:           []pvm.FilePVKey{globalFilePV.Key},
		LastSignStates: []pvm.FilePVLastSignState{globalFilePV.LastSignState},
		KeyFilepath:    privValPath,
		StateFilepath:  privStatePath,
	}
}

func ReloadValidatorKeys(c config, tmNode *node.Node) error {

	keys, err := ReadValidatorPrivateKeyFileLean(GlobalConfig.PocketConfig.GetLeanPocketUserKeyFilePath())
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return errors.New("user key file contained zero validator keys")
	}

	err = SetValidatorsFilesLean(keys)
	if err != nil {
		return err
	}

	validators := loadFilePVWithConfig(c)
	tmNode.ConsensusState().SetPrivValidators(validators) // set new lean nodes

	err = InitNodesLean(c.Logger) // initialize lean nodes
	if err != nil {
		return err
	}

	return nil
}

// hotReloadValidatorsLean - spins off a goroutine that reads from validator files
// TODO: Flesh out hot reloading (removing/adding) lean nodes
func hotReloadValidatorsLean(c config, tmNode *node.Node) {
	userKeysPath := GlobalConfig.PocketConfig.GetLeanPocketUserKeyFilePath()
	stat, err := os.Stat(userKeysPath)
	if err != nil {
		c.Logger.Error("Cannot find user provided key file to hot reload")
		return
	}
	for {
		time.Sleep(time.Second * 5)
		c.Logger.Info("Checking for hot reload")
		newStat, err := os.Stat(userKeysPath)
		if err != nil {
			continue
		}
		if newStat.Size() != stat.Size() || stat.ModTime() != newStat.ModTime() {
			c.Logger.Info("Detected change in files, hot reloading validators")
			err := ReloadValidatorKeys(c, tmNode)
			if err != nil {
				c.Logger.Error("Failed to hot reload validators")
				continue
			}
			c.Logger.Info("Successfully hot reloaded validators")
			stat = newStat
		}
	}
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

	// TODO: Flesh out hotreloading(removing/adding) lean nodes
	//if GlobalConfig.PocketConfig.LeanPocket {
	//	go hotReloadValidatorsLean(c, tmNode)
	//}

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
//	filePv := pvm.LoadOrGenFilePVLean(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
//	filePv.LastSignState.Height = rollbackHeight
//	filePv.LastSignState.Round = 0
//	filePv.LastSignState.Step = 0
//	filePv.LastSignState.Signature = sig
//	filePv.LastSignState.SignBytes = nil
//	filePv.Save()
//	return nil
//}
