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
func InitConfig(chains *HostedBlockchains, logger log.Logger, c types.PocketConfig) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceSealedMap = make(map[string]struct{})
		globalEvidenceCache.Init(c.DataDir, c.EvidenceDBName, c.LevelDBOptions, c.MaxEvidenceCacheEntires)
		globalSessionCache.Init(c.DataDir, c.SessionDBName, c.LevelDBOptions, c.MaxSessionCacheEntries)
		InitGlobalServiceMetric(*chains, logger, c.PrometheusAddr, c.PrometheusMaxOpenfiles)
	})
	GlobalPocketConfig = c
	SetRPCTimeout(c.RPCTimeout)
}

// NOTE: evidence cache is flushed every time db iterator is created (every claim/proof submission)
func FlushSessionCache() {
	err := globalSessionCache.FlushToDB()
	if err != nil {
		fmt.Printf("unable to flush sessions to the database before shutdown!! %s\n", err.Error())
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
