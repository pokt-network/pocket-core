package types

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/libs/log"
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
		globalEvidenceSealedMap = sync.Map{}
		globalEvidenceCache.Init(c.PocketConfig.DataDir, c.PocketConfig.EvidenceDBName, c.TendermintConfig.LevelDBOptions, c.PocketConfig.MaxEvidenceCacheEntires, false)
		globalSessionCache.Init(c.PocketConfig.DataDir, "", c.TendermintConfig.LevelDBOptions, c.PocketConfig.MaxSessionCacheEntries, true)
		InitGlobalServiceMetric(chains, logger, c.PocketConfig.PrometheusAddr, c.PocketConfig.PrometheusMaxOpenfiles)
	})
	GlobalPocketConfig = c.PocketConfig
	SetRPCTimeout(c.PocketConfig.RPCTimeout)
}

func ConvertEvidenceToProto(config types.Config) error {
	InitConfig(nil, log.NewNopLogger(), config)
	gec := globalEvidenceCache
	it, err := gec.Iterator()
	if err != nil {
		return fmt.Errorf("error creating evidence iterator: %s", err.Error())
	}
	defer it.Close()
	for ; it.Valid(); it.Next() {
		ev, err := Evidence{}.LegacyAminoUnmarshal(it.Value())
		if err != nil {
			return fmt.Errorf("error amino unmarshalling evidence: %s", err.Error())
		}
		k, err := ev.Key()
		if err != nil {
			return fmt.Errorf("error creating key from evidence object: %s", err.Error())
		}
		gec.SetWithoutLockAndSealCheck(hex.EncodeToString(k), ev)
	}
	err = gec.FlushToDBWithoutLock()
	if err != nil {
		return fmt.Errorf("error flushing evidence objects to the database: %s", err.Error())
	}
	return nil
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
