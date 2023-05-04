package types

import (
	"encoding/hex"
	"fmt"
	"github.com/alitto/pond"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"github.com/puzpuzpuz/xsync"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	log2 "log"
	"time"
)

const (
	DefaultRPCTimeout = 3000
	MaxRPCTimeout     = 1000000
	MinRPCTimeout     = 1
)

var (
	globalRPCTimeout        time.Duration
	GlobalPocketConfig      types.PocketConfig
	GlobalTenderMintConfig  config.Config
	GlobalEvidenceWorkerMap = xsync.NewMapOf[*pond.WorkerPool]()
)

func InitConfig(chains *HostedBlockchains, logger log.Logger, c types.Config) {
	ConfigOnce.Do(func() {
		InitGlobalServiceMetric(chains, logger, c.PocketConfig.PrometheusAddr, c.PocketConfig.PrometheusMaxOpenfiles)
	})
	InitHttpClient(c.PocketConfig.RPCMaxIdleConns, c.PocketConfig.RPCMaxConnsPerHost, c.PocketConfig.RPCMaxIdleConnsPerHost)
	InitPocketNodeCaches(c, logger)
	InitEvidenceWorker(c, logger)
	GlobalPocketConfig = c.PocketConfig
	GlobalTenderMintConfig = c.TendermintConfig
	if GlobalPocketConfig.LeanPocket {
		GlobalTenderMintConfig.PrivValidatorState = types.DefaultPVSNameLean
		GlobalTenderMintConfig.PrivValidatorKey = types.DefaultPVKNameLean
		GlobalTenderMintConfig.NodeKey = types.DefaultPVSNameLean
	}
	SetRPCTimeout(c.PocketConfig.RPCTimeout)
}

// NewWorkerPool - create pond.WorkerPool instance with the right params in place.
func NewWorkerPool(address string, c types.Config, logger log.Logger) *pond.WorkerPool {
	panicHandler := func(p interface{}) {
		logger.Error(fmt.Sprintf("evidence worker of %s panic a task with: %v", address, p))
	}

	var strategy pond.ResizingStrategy

	switch c.PocketConfig.EvidenceWorker.Strategy {
	case "lazy":
		strategy = pond.Lazy()
		break
	case "eager":
		strategy = pond.Eager()
		break
	case "balanced":
		strategy = pond.Balanced()
		break
	default:
		log2.Fatal(
			fmt.Sprintf(
				"evidence_worker.strategy %s is not a valid option; allowed values are: lazy|eager|balanced",
				c.PocketConfig.EvidenceWorker.Strategy,
			),
		)
	}

	return pond.New(
		// avoid race condition writing evidence in parallel
		1, c.PocketConfig.EvidenceWorker.MaxCapacity,
		pond.IdleTimeout(time.Duration(c.PocketConfig.EvidenceWorker.IdleTimeout)*time.Millisecond),
		pond.PanicHandler(panicHandler),
		pond.Strategy(strategy),
	)
}

func InitEvidenceWorker(c types.Config, logger log.Logger) {
	for address, node := range GlobalPocketNodes {
		if node == nil {
			continue
		}

		logger.Debug(
			fmt.Sprintf(
				"starting worker for: Address=%s Strategy=%s MaxCapacity=%d IdleTimeout=%d",
				address, c.PocketConfig.EvidenceWorker.Strategy,
				c.PocketConfig.EvidenceWorker.MaxCapacity,
				c.PocketConfig.EvidenceWorker.IdleTimeout,
			),
		)

		GlobalEvidenceWorkerMap.Store(address, NewWorkerPool(address, c, logger))
	}
}

func ConvertEvidenceToProto(config types.Config) error {
	// we have to add a random pocket node so that way lean pokt can still support getting the global evidence cache
	node := AddPocketNode(crypto.GenerateEd25519PrivKey().GenPrivateKey(), log.NewNopLogger())

	InitConfig(nil, log.NewNopLogger(), config)

	gec := node.EvidenceStore
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

func StopEvidenceWorker() {
	GlobalEvidenceWorkerMap.Range(func(address string, worker *pond.WorkerPool) bool {
		if !worker.Stopped() {
			worker.StopAndWait()
		}
		return true
	})
}

func FlushSessionCache() {
	for _, k := range GlobalPocketNodes {
		if k.SessionStore != nil {
			err := k.SessionStore.FlushToDB()
			if err != nil {
				fmt.Printf("unable to flush sessions to the database before shutdown!! %s\n", err.Error())
			}
		}
		if k.EvidenceStore != nil {
			err := k.EvidenceStore.FlushToDB()
			if err != nil {
				fmt.Printf("unable to flush GOBEvidence to the database before shutdown!! %s\n", err.Error())
			}
		}
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
