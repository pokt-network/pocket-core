package types

import (
	"github.com/tendermint/tendermint/config"
	db "github.com/tendermint/tm-db"
	"sync"
)

// TmConfig is the structure that holds the SDK configuration parameters.
// This could be used to initialize certain configuration parameters for the SDK.
type SDKConfig struct {
	mtx             sync.RWMutex
	sealed          bool
	txEncoder       TxEncoder
	addressVerifier func([]byte) error
}

type PocketConfig struct {
	DataDir                  string           `json:"data_dir"`
	GenesisName              string           `json:"genesis_file"`
	ChainsName               string           `json:"chains_name"`
	SessionDBType            db.DBBackendType `json:"session_db_type"`
	SessionDBName            string           `json:"session_db_name"`
	EvidenceDBType           db.DBBackendType `json:"evidence_db_type"`
	EvidenceDBName           string           `json:"evidence_db_name"`
	TendermintURI            string           `json:"tendermint_uri"`
	KeybaseName              string           `json:"keybase_name"`
	RPCPort                  string           `json:"rpc_port"`
	ClientBlockSyncAllowance int              `json:"client_block_sync_allowance"`
	MaxEvidenceCacheEntires  int              `json:"max_evidence_cache_entries"`
	MaxSessionCacheEntries   int              `json:"max_session_cache_entries"`
	JSONSortRelayResponses   bool             `json:"json_sort_relay_responses"`
	RemoteCLIURL             string           `json:"remote_cli_url"`
	UserAgent                string           `json:"user_agent"`
	ValidatorCacheSize       int64            `json:"validator_cache_size"`
	ApplicationCacheSize     int64            `json:"application_cache_size"`
	RPCTimeout               int64            `json:"rpc_timeout"`
	PrometheusAddr           string           `json:"pocket_prometheus_port"`
	PrometheusMaxOpenfiles   int              `json:"prometheus_max_open_files"`
	MaxClaimAgeForProofRetry int              `json:"max_claim_age_for_proof_retry"`
	ProofPrevalidation       bool             `json:"proof_prevalidation"`
}

type Config struct {
	TendermintConfig config.Config `json:"tendermint_config"`
	PocketConfig     PocketConfig  `json:"pocket_config"`
}

const (
	DefaultDDName                     = ".pocket"
	DefaultKeybaseName                = "pocket-keybase"
	DefaultPVKName                    = "priv_val_key.json"
	DefaultPVSName                    = "priv_val_state.json"
	DefaultNKName                     = "node_key.json"
	DefaultChainsName                 = "chains.json"
	DefaultGenesisName                = "genesis.json"
	DefaultRPCPort                    = "8081"
	DefaultSessionDBType              = db.CLevelDBBackend
	DefaultEvidenceDBType             = db.CLevelDBBackend
	DefaultSessionDBName              = "session"
	DefaultEvidenceDBName             = "pocket_evidence"
	DefaultTMURI                      = "tcp://localhost:26657"
	DefaultMaxSessionCacheEntries     = 500
	DefaultMaxEvidenceCacheEntries    = 500
	DefaultListenAddr                 = "tcp://0.0.0.0:"
	DefaultClientBlockSyncAllowance   = 10
	DefaultJSONSortRelayResponses     = true
	DefaultDBBackend                  = string(db.CLevelDBBackend)
	DefaultTxIndexer                  = "kv"
	DefaultTxIndexTags                = "tx.hash,tx.height,message.sender,transfer.recipient"
	ConfigDirName                     = "config"
	ConfigFileName                    = "config.json"
	ApplicationDBName                 = "application"
	PlaceholderHash                   = "0001"
	PlaceholderURL                    = "http://127.0.0.1:8081"
	PlaceholderServiceURL             = PlaceholderURL
	DefaultRemoteCLIURL               = "http://localhost:8081"
	DefaultUserAgent                  = ""
	DefaultValidatorCacheSize         = 500
	DefaultApplicationCacheSize       = DefaultValidatorCacheSize
	DefaultPocketPrometheusListenAddr = "8083"
	DefaultPrometheusMaxOpenFile      = 3
	DefaultRPCTimeout                 = 3000
	DefaultMaxClaimProofRetryAge      = 32
	DefaultProofPrevalidation         = false
)

func DefaultConfig(dataDir string) Config {
	c := Config{
		TendermintConfig: *config.DefaultConfig(),
		PocketConfig: PocketConfig{
			DataDir:                  dataDir,
			GenesisName:              DefaultGenesisName,
			ChainsName:               DefaultChainsName,
			SessionDBType:            DefaultSessionDBType,
			SessionDBName:            DefaultSessionDBName,
			EvidenceDBType:           DefaultEvidenceDBType,
			EvidenceDBName:           DefaultEvidenceDBName,
			TendermintURI:            DefaultTMURI,
			KeybaseName:              DefaultKeybaseName,
			RPCPort:                  DefaultRPCPort,
			ClientBlockSyncAllowance: DefaultClientBlockSyncAllowance,
			MaxEvidenceCacheEntires:  DefaultMaxEvidenceCacheEntries,
			MaxSessionCacheEntries:   DefaultMaxSessionCacheEntries,
			JSONSortRelayResponses:   DefaultJSONSortRelayResponses,
			RemoteCLIURL:             DefaultRemoteCLIURL,
			UserAgent:                DefaultUserAgent,
			ValidatorCacheSize:       DefaultValidatorCacheSize,
			ApplicationCacheSize:     DefaultApplicationCacheSize,
			RPCTimeout:               DefaultRPCTimeout,
			PrometheusAddr:           DefaultPocketPrometheusListenAddr,
			PrometheusMaxOpenfiles:   DefaultPrometheusMaxOpenFile,
			MaxClaimAgeForProofRetry: DefaultMaxClaimProofRetryAge,
			ProofPrevalidation:       DefaultProofPrevalidation,
		},
	}
	c.TendermintConfig.SetRoot(dataDir)
	c.TendermintConfig.NodeKey = DefaultNKName
	c.TendermintConfig.PrivValidatorKey = DefaultPVKName
	c.TendermintConfig.PrivValidatorState = DefaultPVSName
	c.TendermintConfig.P2P.AddrBookStrict = false
	c.TendermintConfig.P2P.MaxNumInboundPeers = 250
	c.TendermintConfig.P2P.MaxNumOutboundPeers = 250
	c.TendermintConfig.LogLevel = "*:info, *:error"
	c.TendermintConfig.TxIndex.Indexer = DefaultTxIndexer
	c.TendermintConfig.TxIndex.IndexTags = DefaultTxIndexTags
	c.TendermintConfig.DBBackend = DefaultDBBackend
	c.TendermintConfig.RPC.GRPCMaxOpenConnections = 2500
	c.TendermintConfig.RPC.MaxOpenConnections = 2500
	c.TendermintConfig.Mempool.Size = 9000
	c.TendermintConfig.Mempool.CacheSize = 9000
	c.TendermintConfig.Consensus.TimeoutPropose = 60000000000
	c.TendermintConfig.Consensus.TimeoutProposeDelta = 10000000000
	c.TendermintConfig.Consensus.TimeoutPrevote = 60000000000
	c.TendermintConfig.Consensus.TimeoutPrevoteDelta = 10000000000
	c.TendermintConfig.Consensus.TimeoutPrecommit = 60000000000
	c.TendermintConfig.Consensus.TimeoutPrecommitDelta = 10000000000
	c.TendermintConfig.Consensus.TimeoutCommit = 900000000000
	c.TendermintConfig.Consensus.SkipTimeoutCommit = false
	c.TendermintConfig.Consensus.CreateEmptyBlocks = true
	c.TendermintConfig.Consensus.CreateEmptyBlocksInterval = 900000000000
	c.TendermintConfig.Consensus.PeerGossipSleepDuration = 100000000
	c.TendermintConfig.Consensus.PeerQueryMaj23SleepDuration = 2000000000
	c.TendermintConfig.P2P.AllowDuplicateIP = true
	return c
}

func DefaultTestingPocketConfig() PocketConfig {
	c := DefaultConfig("data")
	c.PocketConfig.EvidenceDBType = db.MemDBBackend
	c.PocketConfig.SessionDBType = db.MemDBBackend
	return c.PocketConfig
}

var (
	// Initializing an instance of TmConfig
	sdkConfig = &SDKConfig{
		sealed:    false,
		txEncoder: nil,
	}
)

// GetConfig returns the config instance for the SDK.
func GetConfig() *SDKConfig {
	return sdkConfig
}

func (config *SDKConfig) assertNotSealed() {
	config.mtx.Lock()
	defer config.mtx.Unlock()

	if config.sealed {
		panic("TmConfig is sealed")
	}
}

// SetTxEncoder builds the TmConfig with TxEncoder used to marshal StdTx to bytes
func (config *SDKConfig) SetTxEncoder(encoder TxEncoder) {
	config.assertNotSealed()
	config.txEncoder = encoder
}

// SetAddressVerifier builds the TmConfig with the provided function for verifying that addresses
// have the correct format
func (config *SDKConfig) SetAddressVerifier(addressVerifier func([]byte) error) {
	config.assertNotSealed()
	config.addressVerifier = addressVerifier
}

// Set the BIP-0044 CoinType code on the config
func (config *SDKConfig) SetCoinType(coinType uint32) {
	config.assertNotSealed()
}

// Seal seals the config such that the config state could not be modified further
func (config *SDKConfig) Seal() *SDKConfig {
	config.mtx.Lock()
	defer config.mtx.Unlock()

	config.sealed = true
	return config
}

// GetTxEncoder return function to encode transactions
func (config *SDKConfig) GetTxEncoder() TxEncoder {
	return config.txEncoder
}

// GetAddressVerifier returns the function to verify that addresses have the correct format
func (config *SDKConfig) GetAddressVerifier() func([]byte) error {
	return config.addressVerifier
}
