package mesh

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/akrylysov/pogreb"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	"github.com/xeipuuv/gojsonschema"
	"io"
	log2 "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type transport struct {
	http.RoundTripper
}

// RoundTrip - handle http requests before/after they run and hook to response handlers bases on path.
func (t *transport) RoundTrip(r *http.Request) (w *http.Response, err error) {
	w, err = t.RoundTripper.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if w.StatusCode != 200 {
		return w, nil
	}

	rr, err := newReusableReader(w.Body)
	if err != nil {
		return nil, err
	}

	w.Body = io.NopCloser(rr)

	return w, nil
}

type reusableReader struct {
	io.Reader
	readBuf *bytes.Buffer
	backBuf *bytes.Buffer
}

// Read - read the buffer and reset to allow multiple reads
func (r reusableReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == io.EOF {
		r.reset()
	}
	return n, err
}

// reset - reset buffer to allow other reads
func (r reusableReader) reset() {
	_, _ = io.Copy(r.readBuf, r.backBuf)
}

// newReusableReader - create new Reader that allow to be read multiple times.
func newReusableReader(r io.Reader) (io.Reader, error) {
	readBuf := bytes.Buffer{}
	_, err := readBuf.ReadFrom(r)
	if err != nil {
		return nil, err
	} // error handling ignored for brevity
	backBuf := bytes.Buffer{}

	return reusableReader{
		io.TeeReader(&readBuf, &backBuf),
		&readBuf,
		&backBuf,
	}, nil
}

// cors - set cors headers
func cors(w *http.ResponseWriter, r *http.Request) (isOptions bool) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	return (*r).Method == "OPTIONS"
}

// loadServicersFromFile return a sync.Map of nodes/servicers that could be used to start working or calculate a reload
func loadServicersFromFile() (nodes *xsync.MapOf[string, *fullNode], servicers *xsync.MapOf[string, *servicer]) {
	nodes = xsync.NewMapOf[*fullNode]()
	servicers = xsync.NewMapOf[*servicer]()

	path := getServicersFilePath()
	logger.Info("reading private key path=" + path)

	fallbackSchemaLoader := gojsonschema.NewSchemaLoader()
	fallbackSchemaStringLoader := gojsonschema.NewStringLoader(fallbackNodeFileSchema)
	fallbackSchema, fallbackSchemaError := fallbackSchemaLoader.Compile(fallbackSchemaStringLoader)
	if fallbackSchemaError != nil {
		log2.Fatal(fmt.Errorf("an error occurred loading fallback json schema: %s", fallbackSchemaError.Error()))
	}

	currentSchemaLoader := gojsonschema.NewSchemaLoader()
	currentSchemaStringLoader := gojsonschema.NewStringLoader(nodeFileSchema)
	currentSchema, currentSchemaError := currentSchemaLoader.Compile(currentSchemaStringLoader)
	if currentSchemaError != nil {
		log2.Fatal(fmt.Errorf("an error occurred loading json schema: %s", currentSchemaError.Error()))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log2.Fatal(fmt.Errorf("an error occurred attempting to read the servicer key file: %s", err.Error()))
	}

	strData := gojsonschema.NewStringLoader(string(data[:]))

	if r, e := fallbackSchema.Validate(strData); e != nil || len(r.Errors()) > 0 {
		if r2, e2 := currentSchema.Validate(strData); e2 != nil || len(r2.Errors()) > 0 {
			log2.Fatal(fmt.Errorf("unable to parse file %s to any of the supported key schemas", path))
		} else {
			var readServicers []nodeFileItem
			// load servicers with new format
			if err := json.Unmarshal(data, &readServicers); err != nil {
				log2.Fatal(fmt.Errorf("an error occurred attempting to parse the servicer key file: %s", err.Error()))
			}

			for _, n := range readServicers {
				var node *fullNode

				if v, ok := nodes.Load(n.URL); !ok {
					node = createNode(n.URL)
					nodes.Store(n.URL, node)
				} else {
					node = v
				}

				for index, pkStr := range n.Keys {
					pk, err := crypto.NewPrivateKey(pkStr)
					if err != nil {
						log2.Fatal(fmt.Errorf("error parsing private key on node=%s index=%d of the file %s", n.URL, index, path))
					}

					address, err := sdk.AddressFromHex(pk.PubKey().Address().String())
					if err != nil {
						log2.Fatal(fmt.Errorf("error getting address from private key on node=%s index=%d of the file %s", n.URL, index, path))
					}

					addressStr := address.String()

					if s, ok := servicers.Load(addressStr); ok {
						node.Servicers.Store(addressStr, s)
					} else {
						newServicer := &servicer{
							SessionCache: xsync.NewMapOf[*appSessionCache](),
							PrivateKey:   pk,
							Address:      address,
							Node:         node,
						}
						servicers.Store(addressStr, newServicer)
						node.Servicers.Store(addressStr, newServicer)
					}
				}
			}
		}
	} else {
		// load servicer with fallback one.
		var readServicers []fallbackNodeFileItem
		// load servicers with new format
		if err := json.Unmarshal(data, &readServicers); err != nil {
			log2.Fatal(fmt.Errorf("an error occurred attempting to parse the servicer key file: %s", err.Error()))
		}

		for index, n := range readServicers {
			var node *fullNode

			if v, ok := nodes.Load(n.ServicerUrl); !ok {
				node = createNode(n.ServicerUrl)
				nodes.Store(n.ServicerUrl, node)
			} else {
				node = v
			}

			pk, err := crypto.NewPrivateKey(n.PrivateKey)
			if err != nil {
				log2.Fatal(fmt.Errorf("error parsing private key on node=%s index=%d of the file %s", n.ServicerUrl, index, path))
			}

			address, err := sdk.AddressFromHex(pk.PubKey().Address().String())
			if err != nil {
				log2.Fatal(fmt.Errorf("error getting address from private key on node=%s index=%d of the file %s", n.ServicerUrl, index, path))
			}

			addressStr := address.String()

			if s, ok := servicers.Load(addressStr); ok {
				node.Servicers.Store(addressStr, s)
			} else {
				newServicer := &servicer{
					SessionCache: xsync.NewMapOf[*appSessionCache](),
					PrivateKey:   pk,
					Address:      address,
					Node:         node,
				}
				servicers.Store(addressStr, newServicer)
				node.Servicers.Store(addressStr, newServicer)
			}
		}
	}

	return
}

// loadServicerNodes - read servicer address and cast to sdk.Address
func loadServicerNodes() (totalNodes, totalServicers int) {
	nodes, servicers := loadServicersFromFile()

	nodes.Range(func(key string, value *fullNode) bool {
		nodesMap.Store(key, value)
		return true
	})

	servicers.Range(func(key string, value *servicer) bool {
		servicerMap.Store(key, value)
		return true
	})

	loadedServicerList := make([]string, 0)

	servicerMap.Range(func(key string, value *servicer) bool {
		loadedServicerList = append(loadedServicerList, value.Address.String())
		return true
	})

	totalNodes = nodes.Size()
	totalServicers = servicers.Size()
	mutex.Lock()
	servicerList = loadedServicerList
	mutex.Unlock()

	return
}

// retryRelaysPolicy - evaluate requests to understand if should or not retry depending on the servicer code response.
func retryRelaysPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error dispatching relay to servicer: %s",
				err.Error(),
			),
		)
		return true, nil
	}

	servicerAddress := resp.Request.Header.Get(ServicerHeader)

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 401 {
			// 401+ could be fixed between restart and reload of cache.
			// 5xx mean something go wrong on servicer node and after a restart could be fixed?
			return true, nil
		}

		if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
			return true, nil
		}

		result := meshRPCRelayResponse{}
		err = json.NewDecoder(resp.Body).Decode(&result)

		if err != nil {
			logger.Error(
				fmt.Sprintf(
					"error decoding servicer %s relay response: %s",
					servicerAddress,
					err.Error(),
				),
			)
			return true, err
		}

		ctxResult := ctx.Value("result").(*meshRPCRelayResponse)
		ctxResult.Success = result.Success
		ctxResult.Dispatch = result.Dispatch
		ctxResult.Error = result.Error

		if ctxResult.Error.Code == pocketTypes.CodeDuplicateProofError {
			return false, nil
		}

		return !isInvalidRelayCode(result.Error.Code), nil
	}

	return false, nil
}

// serveReverseProxy - forward request to ServicerURL
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	u, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &transport{http.DefaultTransport}

	// Update the headers to allow for SSL redirection
	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = u.Host

	// Note that ServeHttp is non-blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// proxyRequest - proxy request to ServicerURL
func proxyRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	serveReverseProxy(getRandomNode().URL, w, r)
}

// reuseBody - transform request body in a reusable reader to allow multiple source read it.
func reuseBody(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		rr, err := newReusableReader(r.Body)
		if err != nil {
			rpc.WriteErrorResponse(w, 500, fmt.Sprintf("error in RPC Handler WriteErrorResponse: %v", err))
		} else {
			r.Body = io.NopCloser(rr)
			handler(w, r, ps)
		}
	}
}

// getMeshRoutes - return routes that will be handled/proxied by mesh rpc server
func getMeshRoutes(simulation bool) rpc.Routes {
	routes := rpc.Routes{
		// Proxy
		rpc.Route{Name: "AppVersion", Method: "GET", Path: "/v1", HandlerFunc: proxyRequest},
		rpc.Route{Name: "Health", Method: "GET", Path: "/v1/health", HandlerFunc: proxyRequest},
		rpc.Route{Name: "Challenge", Method: "POST", Path: "/v1/client/challenge", HandlerFunc: proxyRequest},
		rpc.Route{Name: "ChallengeCORS", Method: "OPTIONS", Path: "/v1/client/challenge", HandlerFunc: proxyRequest},
		rpc.Route{Name: "HandleDispatch", Method: "POST", Path: "/v1/client/dispatch", HandlerFunc: proxyRequest},
		rpc.Route{Name: "HandleDispatchCORS", Method: "OPTIONS", Path: "/v1/client/dispatch", HandlerFunc: proxyRequest},
		rpc.Route{Name: "SendRawTx", Method: "POST", Path: "/v1/client/rawtx", HandlerFunc: proxyRequest},
		rpc.Route{Name: "Stop", Method: "POST", Path: "/v1/private/stop", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryChains", Method: "POST", Path: "/v1/private/chains", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryAccount", Method: "POST", Path: "/v1/query/account", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryAccounts", Method: "POST", Path: "/v1/query/accounts", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryAccountTxs", Method: "POST", Path: "/v1/query/accounttxs", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryACL", Method: "POST", Path: "/v1/query/acl", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryAllParams", Method: "POST", Path: "/v1/query/allparams", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryApp", Method: "POST", Path: "/v1/query/app", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryAppParams", Method: "POST", Path: "/v1/query/appparams", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryApps", Method: "POST", Path: "/v1/query/apps", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryBalance", Method: "POST", Path: "/v1/query/balance", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryBlock", Method: "POST", Path: "/v1/query/block", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryBlockTxs", Method: "POST", Path: "/v1/query/blocktxs", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryDAOOwner", Method: "POST", Path: "/v1/query/daoowner", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryHeight", Method: "POST", Path: "/v1/query/height", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryNode", Method: "POST", Path: "/v1/query/node", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryNodeClaim", Method: "POST", Path: "/v1/query/nodeclaim", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryNodeClaims", Method: "POST", Path: "/v1/query/nodeclaims", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryNodeParams", Method: "POST", Path: "/v1/query/nodeparams", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryNodes", Method: "POST", Path: "/v1/query/nodes", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryParam", Method: "POST", Path: "/v1/query/param", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryPocketParams", Method: "POST", Path: "/v1/query/pocketparams", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryState", Method: "POST", Path: "/v1/query/state", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QuerySupply", Method: "POST", Path: "/v1/query/supply", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QuerySupportedChains", Method: "POST", Path: "/v1/query/supportedchains", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryTX", Method: "POST", Path: "/v1/query/tx", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryUpgrade", Method: "POST", Path: "/v1/query/upgrade", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QuerySigningInfo", Method: "POST", Path: "/v1/query/signinginfo", HandlerFunc: proxyRequest},
		rpc.Route{Name: "LocalNodes", Method: "POST", Path: "/v1/private/nodes", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryUnconfirmedTxs", Method: "POST", Path: "/v1/query/unconfirmedtxs", HandlerFunc: proxyRequest},
		rpc.Route{Name: "QueryUnconfirmedTx", Method: "POST", Path: "/v1/query/unconfirmedtx", HandlerFunc: proxyRequest},
		//
		rpc.Route{Name: "MeshService", Method: "POST", Path: "/v1/client/relay", HandlerFunc: reuseBody(meshNodeRelay)},
		// mesh private routes
		rpc.Route{Name: "MeshHealth", Method: "GET", Path: "/v1/private/mesh/health", HandlerFunc: meshHealth},
		rpc.Route{Name: "QueryMeshNodeChains", Method: "POST", Path: "/v1/private/mesh/chains", HandlerFunc: meshChains},
		rpc.Route{Name: "MeshNodeServicer", Method: "POST", Path: "/v1/private/mesh/servicers", HandlerFunc: meshServicerNode},
		rpc.Route{Name: "UpdateMeshNodeChains", Method: "POST", Path: "/v1/private/mesh/updatechains", HandlerFunc: meshUpdateChains},
		rpc.Route{Name: "StopMeshNode", Method: "POST", Path: "/v1/private/mesh/stop", HandlerFunc: meshStop},
	}

	// check if simulation is turn on
	if simulation {
		simRoute := rpc.Route{Name: "SimulateRequest", Method: "POST", Path: "/v1/client/sim", HandlerFunc: meshSimulateRelay}
		routes = append(routes, simRoute)
	}

	return routes
}

// prepareHttpClients - prepare http clients & transports
func prepareHttpClients() {
	logger.Info("initializing http clients")
	chainsTransport := http.DefaultTransport.(*http.Transport).Clone()
	chainsTransport.MaxIdleConns = 1000
	chainsTransport.MaxConnsPerHost = 1000
	chainsTransport.MaxIdleConnsPerHost = 1000

	servicerTransport := http.DefaultTransport.(*http.Transport).Clone()
	servicerTransport.MaxIdleConns = 50
	servicerTransport.MaxConnsPerHost = 50
	servicerTransport.MaxIdleConnsPerHost = 50

	chainsClient = &http.Client{
		Timeout:   time.Duration(app.GlobalMeshConfig.RPCTimeout) * time.Millisecond,
		Transport: chainsTransport,
	}
	servicerClient = &http.Client{
		Timeout:   time.Duration(app.GlobalMeshConfig.RPCTimeout) * time.Millisecond,
		Transport: servicerTransport,
	}

	relaysClient = retryablehttp.NewClient()
	relaysClient.RetryMax = app.GlobalMeshConfig.ServicerRetryMaxTimes
	relaysClient.HTTPClient = servicerClient
	relaysClient.Logger = &LevelHTTPLogger{}
	relaysClient.RetryWaitMin = time.Duration(app.GlobalMeshConfig.ServicerRetryWaitMin) * time.Millisecond
	relaysClient.RetryWaitMax = time.Duration(app.GlobalMeshConfig.ServicerRetryWaitMax) * time.Millisecond
	relaysClient.CheckRetry = retryRelaysPolicy
}

// catchSignal - catch system signals and process them
func catchSignal() {
	terminateSignals := make(chan os.Signal, 1)
	reloadSignals := make(chan os.Signal, 1)

	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, os.Kill, os.Interrupt) //NOTE:: syscall.SIGKILL we cannot catch kill -9 as its force kill signal.

	signal.Notify(reloadSignals, syscall.SIGUSR1)

	for { // We are looping here because config reload can happen multiple times.
		select {
		case s := <-terminateSignals:
			logger.Info("shutting down server gracefully, SIGNAL NAME:", s)
			StopMeshRPC()
			finish()
			break //break is not necessary to add here as if server is closed our main function will end.
		case s := <-reloadSignals:
			logger.Debug("reloading, SIGNAL NAME:", s)
			// todo: reload config? reload chains? reload auth/key? is really this needed?
		}
	}
}

// initCache - initialize cache
func initCache() {
	var err error

	logger.Info("initializing relays cache")
	relaysCacheFilePath := app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.RelayCacheFile
	relaysCacheDb, err = pogreb.Open(relaysCacheFilePath, &pogreb.Options{
		// BackgroundSyncInterval sets the amount of time between background Sync() calls.
		//
		// Setting the value to 0 disables the automatic background synchronization.
		// Setting the value to -1 makes the DB call Sync() after every write operation.
		BackgroundSyncInterval: time.Duration(app.GlobalMeshConfig.RelayCacheBackgroundSyncInterval) * time.Millisecond,
		// BackgroundCompactionInterval sets the amount of time between background Compact() calls.
		//
		// Setting the value to 0 disables the automatic background compaction.
		BackgroundCompactionInterval: time.Duration(app.GlobalMeshConfig.RelayCacheBackgroundCompactionInterval) * time.Millisecond,
	})
	if err != nil {
		log2.Fatal(err)
		return
	}

	logger.Info(fmt.Sprintf("resuming %d relays from cache", relaysCacheDb.Count()))
	it := relaysCacheDb.Items()
	for {
		key, val, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log2.Fatal(err)
		}

		logger.Debug("loading relay hash=%s", hex.EncodeToString(key))
		relay := decodeCacheRelay(val)

		if relay != nil {
			servicerAddress, err := getAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)
			if err != nil {
				logger.Debug(
					fmt.Sprintf(
						"removing relay hash=%s from cache because was unable decode pk from pk file",
						relay.RequestHashString(),
					),
				)
				deleteCacheRelay(relay)
				continue
			}

			servicerNode, ok := servicerMap.Load(servicerAddress)
			if !ok {
				logger.Debug(
					fmt.Sprintf(
						"removing relay hash=%s from cache because was unable to load servicer %s from pk file",
						relay.RequestHashString(),
						hex.EncodeToString(key),
					),
				)
				deleteCacheRelay(relay)
				continue
			}

			if !ok {
				logger.Debug(
					fmt.Sprintf(
						"removing relay hash=%s from cache because was unable to cast *servicer instance for %s",
						relay.RequestHashString(),
						hex.EncodeToString(key),
					),
				)
				deleteCacheRelay(relay)
				continue
			}

			// the node worker pool is dynamic so if the keys are reloaded and the current worker is reloaded + modified
			// it will need to be resized, for that the current way is stop the current worker and create a new one
			// so at that moment the node will have this flag on "true" until it get done.
			if servicerNode.Node.ResizingWorker {
				go func(s *servicer, relay *pocketTypes.Relay) {
					for {
						time.Sleep(10 * time.Millisecond)
						if !s.Node.ResizingWorker {
							s.Node.Worker.Submit(func() {
								notifyServicer(relay)
							})
						}
					}
				}(servicerNode, relay)
			} else {
				servicerNode.Node.Worker.Submit(func() {
					notifyServicer(relay)
				})
			}
		}
	}
}

// initCrons - initialize in memory cron jobs
func initCrons() {
	// start cron for height pooling
	cronJobs = cron.New()

	logger.Info("initializing session cache clean up")
	// schedule clean old session job
	cleanOldSessions(cronJobs)

	// start all the cron jobs
	cronJobs.Start()
}
