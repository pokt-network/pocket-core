package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

const (
	DefaultRPCTimeout = 3000
	MaxRPCTimeout     = 1000000
	MinRPCTimeout     = 1
)

var (
	globalRPCTimeout   time.Duration
	GlobalPocketConfig types.PocketConfig
)

// "InitConfig" - Initializes the cache for sessions and evidence
func InitConfig(chains *HostedBlockchains, logger log.Logger, c types.Config) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceSealedMap = make(map[string]struct{})
		globalEvidenceCache.Init(c.PocketConfig.DataDir, c.PocketConfig.EvidenceDBName, c.TendermintConfig.LevelDBOptions, c.PocketConfig.MaxEvidenceCacheEntires)
		globalSessionCache.Init(c.PocketConfig.DataDir, c.PocketConfig.SessionDBName, c.TendermintConfig.LevelDBOptions, c.PocketConfig.MaxSessionCacheEntries)
		InitGlobalServiceMetric(*chains, logger, c.PocketConfig.PrometheusAddr, c.PocketConfig.PrometheusMaxOpenfiles)
	})
	GlobalPocketConfig = c.PocketConfig
	SetRPCTimeout(c.PocketConfig.RPCTimeout)
}

// NOTE: evidence cache is flushed every time db iterator is created (every claim/proof submission)
func FlushSessionCache() {
	err := globalSessionCache.FlushToDB()
	if err != nil {
		fmt.Printf("unable to flush sessions to the database before shutdown!! %s\n", err.Error())
	}
	err = globalEvidenceCache.FlushToDB()
	if err != nil {
		fmt.Printf("unable to flush GOBEvidence to the database before shutdown!! %s\n", err.Error())
	}
}

func GetRPCTimeout() time.Duration {
	return globalRPCTimeout
}

func SetRPCTimeout(timeout int64) {
	if timeout < MinRPCTimeout || timeout > MaxRPCTimeout {
		timeout = DefaultRPCTimeout
	}

	globalRPCTimeout = time.Duration(timeout)
}
