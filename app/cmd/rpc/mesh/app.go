package mesh

import (
	"context"
	"fmt"
	"github.com/akrylysov/pogreb"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	"github.com/tendermint/tendermint/libs/log"
	log2 "log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	ModuleName              = "pocketcore"
	ServicerHeader          = "X-Servicer"
	ServicerRelayEndpoint   = "/v1/private/mesh/relay"
	ServicerSessionEndpoint = "/v1/private/mesh/session"
	ServicerCheckEndpoint   = "/v1/private/mesh/check"
	AppVersion              = "RC-0.3.0"
)

var (
	srv               *http.Server
	finish            context.CancelFunc
	logger            log.Logger
	chainsClient      *http.Client
	servicerClient    *http.Client
	relaysClient      *retryablehttp.Client
	relaysCacheDb     *pogreb.DB
	servicerMap       = xsync.NewMapOf[*servicer]()
	nodesMap          = xsync.NewMapOf[*fullNode]()
	servicerList      []string
	chains            *pocketTypes.HostedBlockchains
	meshAuthToken     sdk.AuthToken
	servicerAuthToken sdk.AuthToken
	cronJobs          *cron.Cron
	mutex             = sync.Mutex{}
)

// validate payload
//	modulename: pocketcore CodeEmptyPayloadDataError = 25
// ensures the block height is within the acceptable range
//	modulename: pocketcore CodeOutOfSyncRequestError            = 75
// validate the relay merkleHash = request merkleHash
// 	modulename: pocketcore CodeRequestHash                      = 74
// ensure the blockchain is supported locally
// 	CodeUnsupportedBlockchainNodeError   = 26
// ensure session block height == one in the relay proof
// 	CodeInvalidBlockHeightError          = 60
// get the session context
// 	CodeInternal              CodeType = 1
// get the application that staked on behalf of the client
// 	CodeAppNotFoundError                 = 45
// validate unique relay
// 	CodeEvidenceSealed                   = 90
// get evidence key by proof
// 	CodeDuplicateProofError              = 37
// validate not over service
// 	CodeOverServiceError                 = 71
// "ValidateLocal" - Validates the proof object, where the owner of the proof is the local node
// 	CodeInvalidBlockHeightError          = 60
// 	CodePublKeyDecodeError               = 6
// 	CodePubKeySizeError                  = 42
// 	CodeNewHexDecodeError                = 52
// 	CodeEmptyBlockHashError              = 23
// 	CodeInvalidHashLengthError           = 62
// 	CodeInvalidEntropyError              = 29
// 	CodeInvalidTokenError                = 4
// 	CodeSigDecodeError                   = 39
// 	CodeInvalidSignatureSizeError        = 38
// 	CodePublKeyDecodeError               = 6
// 	CodeMsgDecodeError                   = 40
// 	CodeInvalidSigError                  = 41
// 	CodeInvalidEntropyError              = 29
// 	CodeInvalidNodePubKeyError           = 34
// 	CodeUnsupportedBlockchainAppError    = 13
var invalidCodes = []sdk.CodeType{
	pocketTypes.CodeRequestHash,
	pocketTypes.CodeAppNotFoundError,
	pocketTypes.CodeEvidenceSealed,
	pocketTypes.CodeOverServiceError,
	pocketTypes.CodeOutOfSyncRequestError,
	pocketTypes.CodeInvalidBlockHeightError,
}

// StopRPC - stop http server
func StopRPC() {
	// stop receiving new requests
	logger.Info("stopping http server...")
	if srv != nil {
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("http server shutdown error: %s", err.Error()))
		}
	}
	logger.Info("http server stopped!")

	// close relays cache db
	logger.Info("stopping relays cache database...")
	if err := relaysCacheDb.Close(); err != nil {
		logger.Error(fmt.Sprintf("relays cache db shutdown error: %s", err.Error()))
	}
	logger.Info("relays cache database stopped!")

	// stop accepting new tasks and signal all workers to stop processing new tasks. Tasks being processed by workers
	// will continue until completion unless the process is terminated.
	logger.Info("stopping worker pools...")
	nodesMap.Range(func(key string, node *fullNode) bool {
		node.stop()
		return true
	})
	logger.Info("worker pools stopped!")

	logger.Info("stopping clean session cron job")
	cronJobs.Stop()
	logger.Info("clean session job stopped!")

	// Stop prometheus server
	StopPrometheusServer()
}

// StartRPC - Start mesh rpc server
func StartRPC(router *httprouter.Router) {
	ctx, cancel := context.WithCancel(context.Background())
	finish = cancel
	defer cancel()
	// initialize logger
	logger = initLogger()
	// initialize pseudo random to choose servicer url
	rand.Seed(time.Now().Unix())
	// load auth token files (servicer and mesh node)
	loadAuthTokens()
	// instantiate all the http clients used to call Chains and Servicer
	prepareHttpClients()
	// retrieve the nonNative blockchains your node is hosting
	chains = loadHostedChains()
	// load chain name map use on metrics. this will not raise or throw an error.
	loadChainsNameMap()
	// turn on chains hot reload
	go initKeysHotReload()
	go initChainsHotReload()
	// initialize prometheus metrics
	StartPrometheusServer()
	// read servicer
	totalNodes, totalServicers := loadServicerNodes()
	// check servicers are reachable at required endpoints
	connectivityChecks(mapset.NewSet[string]())
	// initialize crons
	initCrons()
	// bootstrap cache
	initCache()

	srv = &http.Server{
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      60 * time.Second,
		Addr:              ":" + app.GlobalMeshConfig.RPCPort,
		Handler: http.TimeoutHandler(
			router,
			time.Duration(app.GlobalMeshConfig.ClientRPCTimeout)*time.Millisecond,
			"server Timeout Handling Request",
		),
	}

	go catchSignal()

	logger.Info(
		fmt.Sprintf(
			"start serving relay as mesh node on http://0.0.0.0:%s for %d servicer in %d nodes",
			app.GlobalMeshConfig.RPCPort,
			totalServicers,
			totalNodes,
		),
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log2.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		// Shutdown the server when the context is canceled
		logger.Info("bye bye! bip bop!")
	}
}
