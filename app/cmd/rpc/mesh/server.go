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

	rr, err := NewReusableReader(w.Body)
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

// NewReusableReader - create new Reader that allow to be read multiple times.
func NewReusableReader(r io.Reader) (io.Reader, error) {
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
					node = createNode(n.URL, n.Name)
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
							SessionCache: xsync.NewMapOf[*AppSessionCache](),
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
				node = createNode(n.ServicerUrl, "")
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
					SessionCache: xsync.NewMapOf[*AppSessionCache](),
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

		result := RPCRelayResponse{}
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

		ctxResult := ctx.Value("result").(*RPCRelayResponse)
		ctxResult.Success = result.Success
		ctxResult.Dispatch = result.Dispatch
		ctxResult.Error = result.Error

		if ctxResult.Error.Code == pocketTypes.CodeDuplicateProofError {
			return false, nil
		}

		return !IsInvalidRelayCode(result.Error.Code), nil
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

// ProxyRequest - proxy request to ServicerURL
func ProxyRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	serveReverseProxy(GetRandomNode().URL, w, r)
}

// prepareHttpClients - prepare http clients & transports
func prepareHttpClients() {
	logger.Info("initializing http clients")
	chainsTransport := http.DefaultTransport.(*http.Transport).Clone()
	chainsTransport.MaxIdleConns = app.GlobalMeshConfig.ChainRPCMaxIdleConnections
	chainsTransport.MaxConnsPerHost = app.GlobalMeshConfig.ChainRPCMaxConnsPerHost
	chainsTransport.MaxIdleConnsPerHost = app.GlobalMeshConfig.ChainRPCMaxIdleConnsPerHost

	servicerTransport := http.DefaultTransport.(*http.Transport).Clone()
	servicerTransport.MaxIdleConns = app.GlobalMeshConfig.ServicerRPCMaxIdleConnections
	servicerTransport.MaxConnsPerHost = app.GlobalMeshConfig.ServicerRPCMaxConnsPerHost
	servicerTransport.MaxIdleConnsPerHost = app.GlobalMeshConfig.ServicerRPCMaxIdleConnsPerHost

	chainsClient = &http.Client{
		Timeout:   time.Duration(app.GlobalMeshConfig.ChainRPCTimeout) * time.Millisecond,
		Transport: chainsTransport,
	}
	servicerClient = &http.Client{
		Timeout:   time.Duration(app.GlobalMeshConfig.ServicerRPCTimeout) * time.Millisecond,
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
			StopRPC()
			finish()
			break // break is not necessary to add here as if server is closed our main function will end.
		case s := <-reloadSignals:
			logger.Debug("reloading SIGNAL received:", s)
			reloadChains()
			reloadServicers()
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
			servicerAddress, err := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)
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

			servicerNode.Node.Worker.Submit(func() {
				notifyServicer(relay)
			})
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
