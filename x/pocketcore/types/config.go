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
	globalUserAgent  string
	globalRPCTimeout time.Duration
	globalClientBlockAllowance int
	globalSortJSONResponses    bool
)

// "InitConfig" - Initializes the cache for sessions and evidence
func InitConfig(chains *HostedBlockchains, logger log.Logger, c types.PocketConfig) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceCache.Init(c.DataDir, c.EvidenceDBName, c.EvidenceDBType, c.MaxEvidenceCacheEntires)
		globalSessionCache.Init(c.DataDir, c.SessionDBName, c.SessionDBType, c.MaxSessionCacheEntries)
		InitGlobalServiceMetric(*chains, logger, c.PrometheusAddr, c.PrometheusMaxOpenfiles)
	})
	globalUserAgent = c.UserAgent
	globalClientBlockAllowance = c.ClientBlockSyncAllowance
	globalSortJSONResponses = c.JSONSortRelayResponses
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
