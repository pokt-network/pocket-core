package types

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
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
)

// "InitConfig" - Initializes the cache for sessions and evidence
func InitConfig(userAgent, evidenceDir, sessionDir string, sessionDBType, evidenceDBType db.DBBackendType, maxEvidenceEntries, maxSessionEntries int, evidenceDBName, sessionDBName string, chains HostedBlockchains, logger log.Logger, prometheusAddr string, maxOpenConn int, timeout int64) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceCache.Init(evidenceDir, evidenceDBName, evidenceDBType, maxEvidenceEntries)
		globalSessionCache.Init(sessionDir, sessionDBName, sessionDBType, maxSessionEntries)
		InitGlobalServiceMetric(chains, logger, prometheusAddr, maxOpenConn)
	})
	globalUserAgent = userAgent
	SetRPCTimeout(timeout)
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
