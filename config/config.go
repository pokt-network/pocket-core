package config

import (
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"time"
)

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
	newTMConfig.SetRoot(rootDir)
	newTMConfig.DBPath = datadir
	newTMConfig.NodeKey = nodekey
	newTMConfig.PrivValidatorKey = privValKey
	newTMConfig.PrivValidatorState = privValState
	newTMConfig.P2P.ListenAddress = listenAddr                  // Node listen address. (0.0.0.0:0 means any interface, any port)
	newTMConfig.P2P.PersistentPeers = persistentPeers           // Comma-delimited ID@host:port persistent peers
	newTMConfig.P2P.Seeds = seeds                               // Comma-delimited ID@host:port seed nodes
	newTMConfig.Consensus.CreateEmptyBlocks = createEmptyBlocks // Set this to false to only produce blocks when there are txs or when the AppHash changes
	newTMConfig.Consensus.CreateEmptyBlocksInterval = createEmptyBlocksInterval
	newTMConfig.P2P.MaxNumInboundPeers = MaxNumberInboundPeers
	newTMConfig.P2P.MaxNumOutboundPeers = MaxNumberOutboundPeers

	return &Config{newTMConfig, logger, traceWriterPath}
}
