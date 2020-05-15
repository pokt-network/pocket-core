package types

import db "github.com/tendermint/tm-db"

var (
	globalUserAgent string
)

// "InitConfig" - Initializes the cache for sessions and evidence
func InitConfig(userAgent, evidenceDir, sessionDir string, sessionDBType, evidenceDBType db.DBBackendType, maxEvidenceEntries, maxSessionEntries int, evidenceDBName, sessionDBName string) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceCache.Init(evidenceDir, evidenceDBName, evidenceDBType, maxEvidenceEntries)
		globalSessionCache.Init(sessionDir, sessionDBName, sessionDBType, maxSessionEntries)
	})
	globalUserAgent = userAgent
}
