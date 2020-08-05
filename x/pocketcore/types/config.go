package types

import (
	"fmt"
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
func InitConfig(userAgent, evidenceDir, sessionDir string, sessionDBType, evidenceDBType db.DBBackendType, maxEvidenceEntries, maxSessionEntries int, evidenceDBName, sessionDBName string, timeout int64) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceCache.Init(evidenceDir, evidenceDBName, evidenceDBType, maxEvidenceEntries)
		globalSessionCache.Init(sessionDir, sessionDBName, sessionDBType, maxSessionEntries)
	})
	globalUserAgent = userAgent
	SetRPCTimeout(timeout)
}

func FlushCache() {
	err := globalSessionCache.FlushToDB()
	if err != nil {
		fmt.Printf("unable to flush sessions to the database before shutdown!! %s\n", err.Error())
	}
	err = globalEvidenceCache.FlushToDB()
	if err != nil {
		fmt.Printf("unable to flush evidence to the database before shutdown!! %s\n", err.Error())
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
