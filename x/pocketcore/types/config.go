package types

import (
	"fmt"
	db "github.com/tendermint/tm-db"
)

var (
	globalUserAgent string
	GlobalAltruism  bool
)

// "InitConfig" - Initializes the cache for sessions and evidence
func InitConfig(userAgent, evidenceDir, sessionDir string, sessionDBType, evidenceDBType db.DBBackendType, maxEvidenceEntries, maxSessionEntries int, evidenceDBName, sessionDBName string, altruism bool) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceCache.Init(evidenceDir, evidenceDBName, evidenceDBType, maxEvidenceEntries)
		globalSessionCache.Init(sessionDir, sessionDBName, sessionDBType, maxSessionEntries)
	})
	globalUserAgent = userAgent
	GlobalAltruism = altruism
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
