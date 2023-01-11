package rpc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akrylysov/pogreb"
	"github.com/alitto/pond"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	types4 "github.com/pokt-network/pocket-core/app/cmd/rpc/types"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/robfig/cron/v3"
	"github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	"io"
	"io/ioutil"
	log2 "log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	ModuleName              = "pocketcore"
	ServicerHeader          = "X-Servicer"
	ServicerRelayEndpoint   = "/v1/private/mesh/relay"
	ServicerSessionEndpoint = "/v1/private/mesh/session"
	AppVersion              = "ALPHA-0.2.6"
)

type appCache struct {
	PublicKey       string
	Chain           string
	Dispatch        *dispatchResponse
	RemainingRelays int64
	IsValid         bool
	Error           *sdkErrorResponse
}

type servicerFile struct {
	PrivateKey  string `json:"priv_key"`
	ServicerUrl string `json:"servicer_url"`
}

type meshHealthResponse struct {
	Version string `json:"version"`
}

type meshRPCRelayResult struct {
	Success  bool                          `json:"signature"`
	Error    error                         `json:"error"`
	Dispatch *pocketTypes.DispatchResponse `json:"dispatch"`
}

type meshRPCSessionResult struct {
	Success         bool              `json:"success"`
	Error           *sdkErrorResponse `json:"error"`
	Dispatch        *dispatchResponse `json:"dispatch"`
	RemainingRelays json.Number       `json:"remaining_relays"`
}

type meshRPCRelayResponse struct {
	Success  bool              `json:"signature"`
	Error    *sdkErrorResponse `json:"error"`
	Dispatch *dispatchResponse `json:"dispatch"`
}

type sdkErrorResponse struct {
	Code      sdk.CodeType      `json:"code"`
	Codespace sdk.CodespaceType `json:"codespace"`
	Error     string            `json:"message"`
}

type transport struct {
	http.RoundTripper
}

type reusableReader struct {
	io.Reader
	readBuf *bytes.Buffer
	backBuf *bytes.Buffer
}

type dispatchSessionNode struct {
	Address       string          `json:"address"`
	Chains        []string        `json:"chains"`
	Jailed        bool            `json:"jailed"`
	OutputAddress string          `json:"output_address"`
	PublicKey     string          `json:"public_key"`
	ServiceUrl    string          `json:"service_url"`
	Status        sdk.StakeStatus `json:"status"`
	Tokens        string          `json:"tokens"`
	UnstakingTime time.Time       `json:"unstaking_time"`
}

type dispatchSession struct {
	Header pocketTypes.SessionHeader `json:"header"`
	Key    string                    `json:"key"`
	Nodes  []dispatchSessionNode     `json:"nodes"`
}

// handle /v1/client/dispatch response due to was unable to inflate it using pocket core struct
// it was throwing an error about Nodes unmarshalling
type dispatchResponse struct {
	BlockHeight int64           `json:"block_height"`
	Session     dispatchSession `json:"session"`
}

// servicer - represents a node load from servicer_private_key_file
// also handle status and dedicated worker group.
type servicer struct {
	PrivateKey   crypto.PrivateKey
	Address      sdk.Address
	ServicerURL  string
	Status       *app.HealthResponse
	Crons        *cron.Cron
	Worker       *pond.WorkerPool
	SessionCache sync.Map
}

type LevelHTTPLogger struct {
	retryablehttp.LeveledLogger
}

var (
	srv               *http.Server
	finish            context.CancelFunc
	logger            log.Logger
	chainsClient      *http.Client
	servicerClient    *http.Client
	relaysClient      *retryablehttp.Client
	relaysCacheDb     *pogreb.DB
	servicerPkMap     sync.Map
	servicerList      []string
	chains            *pocketTypes.HostedBlockchains
	meshAuthToken     sdk.AuthToken
	servicerAuthToken sdk.AuthToken
	cronJobs          *cron.Cron
	interceptorPool   *pond.WorkerPool
	mutex             = sync.Mutex{}
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
	invalidCodes = []sdk.CodeType{
		pocketTypes.CodeRequestHash,
		pocketTypes.CodeAppNotFoundError,
		pocketTypes.CodeEvidenceSealed,
		pocketTypes.CodeOverServiceError,
		pocketTypes.CodeOutOfSyncRequestError,
		pocketTypes.CodeInvalidBlockHeightError,
	}
)

// fields - mutate interface to key/value object to be print on stdout
func (l *LevelHTTPLogger) fields(keysAndValues ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	for i := 0; i < len(keysAndValues)-1; i += 2 {
		fields[keysAndValues[i].(string)] = keysAndValues[i+1]
	}

	return fields
}

// Error - log to stdout as error level
func (l *LevelHTTPLogger) Error(msg string, keysAndValues ...interface{}) {
	fields := l.fields(keysAndValues...)
	err := fields["error"].(error)
	_url := fields["url"]
	if _url != nil {
		_url2, ok := _url.(*url.URL)
		if !ok {
			logger.Error("request error", "error", _url)
			return
		}

		logger.Error(
			fmt.Sprintf(
				"%s at %s %s://%s%s\n",
				msg,
				fields["method"].(string),
				_url2.Scheme,
				_url2.Host,
				_url2.Path,
			),
		)
		return
	}
	logger.Error(msg, err, fields)
}

// Info - log to stdout as info level
func (l *LevelHTTPLogger) Info(msg string, keysAndValues ...interface{}) {
	logger.Info(msg, l.fields(keysAndValues...))
}

// Debug - log to stdout as debug level
func (l *LevelHTTPLogger) Debug(msg string, keysAndValues ...interface{}) {
	fields := l.fields(keysAndValues...)
	_url := fields["url"]
	if _url != nil {
		_url2, ok := _url.(*url.URL)
		if !ok {
			logger.Error(fmt.Sprintf("unable to cast to url.URL %v", _url))
			return
		}
		logger.Debug(
			fmt.Sprintf(
				"%s:\nURL=%s://%s%s?%s\nMETHOD=%s",
				msg,
				_url2.Scheme, _url2.Host, _url2.Path, _url2.RawQuery,
				fields["method"].(string),
			),
		)
		return
	}
	logger.Debug(msg, fields)
}

// Warn - log to stdout as warning level
func (l *LevelHTTPLogger) Warn(msg string, keysAndValues ...interface{}) {
	logger.Debug(msg, l.fields(keysAndValues...))
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

// Contains - evaluate if the dispatch response contains passed address in their node list
func (sn dispatchResponse) Contains(addr sdk.Address) bool {
	// if nil return
	if addr == nil {
		return false
	}
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		// There is reference to node address so that way we don't have to recreate address twice for pre-leanpokt
		address, err := sdk.AddressFromHex(node.Address)
		if err != nil {
			log2.Fatal(err)
		}
		if _, ok := servicerPkMap.Load(address); ok {
			return true
		}
	}
	return false
}

// ShouldKeep - evaluate if this dispatch response is one that we need to keep for the running mesh node.
func (sn dispatchResponse) ShouldKeep() bool {
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		if _, ok := servicerPkMap.Load(node.Address); ok {
			return true
		}
	}
	// if hit here, no one of in the map match the dispatch response nodes.
	return false
}

// GetSupportedNodes - return a list of the supported nodes of running mesh node from the dispatchResponse payload.
func (sn dispatchResponse) GetSupportedNodes() []string {
	nodes := make([]string, 0)
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		// There is reference to node address so that way we don't have to recreate address twice for pre-leanpokt
		if _, ok := servicerPkMap.Load(node.Address); ok {
			nodes = append(nodes, node.Address)
		}
	}
	// if hit here, no one of in the map match the dispatch response nodes.
	return nodes
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

// isInvalidRelayCode - check if the error code is someone that block incoming relays for current session.
func isInvalidRelayCode(code sdk.CodeType) bool {
	for _, c := range invalidCodes {
		if c == code {
			return true
		}
	}

	return false
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

// reuseBody - transform request body in a reusable reader to allow multiple source read it.
func reuseBody(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		rr, err := newReusableReader(r.Body)
		if err != nil {
			WriteErrorResponse(w, 500, fmt.Sprintf("error in RPC Handler WriteErrorResponse: %v", err))
		} else {
			r.Body = io.NopCloser(rr)
			handler(w, r, ps)
		}
	}
}

// getAppSession - call ServicerURL to get an application session using retrieve header
func getAppSession(relay *pocketTypes.Relay, model interface{}) *sdkErrorResponse {
	servicerNode := getRandomServicer()
	payload := pocketTypes.MeshSession{
		SessionHeader: pocketTypes.SessionHeader{
			ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
			Chain:              relay.Proof.Blockchain,
			SessionBlockHeight: relay.Proof.SessionBlockHeight,
		},
		Meta:               relay.Meta,
		ServicerPubKey:     relay.Proof.ServicerPubKey,
		Blockchain:         relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}
	logger.Debug("reading session from servicer")
	jsonData, e := json.Marshal(payload)
	if e != nil {
		return newSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
	}

	requestURL := fmt.Sprintf(
		"%s%s?authtoken=%s",
		servicerNode.ServicerURL,
		ServicerSessionEndpoint,
		servicerAuthToken.Value,
	)
	req, e := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, e := servicerClient.Do(req)

	if e != nil {
		return newSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return // add log here
		}
	}(resp.Body)

	if resp.StatusCode == 401 {
		return newSdkErrorFromPocketSdkError(
			sdk.ErrUnauthorized(
				fmt.Sprintf("wrong auth form %s", ServicerSessionEndpoint),
			),
		)
	}

	isSuccess := resp.StatusCode == 200

	if !isSuccess {
		result := meshRPCSessionResult{}
		e = json.NewDecoder(resp.Body).Decode(&result)
		if e != nil {
			return newSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
		}
		return nil
	} else {
		e = json.NewDecoder(resp.Body).Decode(model)
		if e != nil {
			return newSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
		}
		return nil
	}
}

// getRandomServicer - return a random servicer object from the list load at the start
func getRandomServicer() *servicer {
	mutex.Lock()
	address := servicerList[rand.Intn(len(servicerList))]
	mutex.Unlock()
	s, ok := servicerPkMap.Load(address)
	if !ok {
		return nil
	}
	return s.(*servicer)
}

// getServicerAddressFromPubKey - return an address as string from a public key string
func getServicerAddressFromPubKeyAsString(pubKey string) (string, error) {
	key, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	return sdk.GetAddress(key).String(), nil
}

// proxyRequest - proxy request to ServicerURL
func proxyRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	servicerNode := getRandomServicer()
	serveReverseProxy(servicerNode.ServicerURL, w, r)
}

// updateChains - update chainName file with the retrieve chains value.
func updateChains(chains []pocketTypes.HostedBlockchain) {
	var chainsPath = app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ChainsName
	var jsonFile *os.File
	if _, err := os.Stat(chainsPath); err != nil && os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("no chains.json found @ %s", chainsPath))
		return
	}
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(chainsPath, os.O_WRONLY, os.ModePerm)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// create dummy input for the file
	res, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// write to the file
	_, err = jsonFile.Write(res)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
}

// reloadChains - reload chainsName file
func reloadChains(chainsPath string) {
	// if file exists open, else create and open
	var jsonFile *os.File
	var bz []byte
	if !fileExist(chainsPath) {
		log2.Println(fmt.Sprintf("chains file no found at %s; ignoring reload", chainsPath))
		return
	}
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(chainsPath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	bz, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// unmarshal into the structure
	var hostedChainsSlice []pocketTypes.HostedBlockchain
	err = json.Unmarshal(bz, &hostedChainsSlice)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	m := make(map[string]pocketTypes.HostedBlockchain)
	for _, chain := range hostedChainsSlice {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			log2.Fatal(fmt.Sprintf("invalid ID: %s in network identifier in %s file", chain.ID, app.GlobalMeshConfig.ChainsName))
		}
		m[chain.ID] = chain
	}
	chains.L.Lock()
	chains.M = m
	chains.L.Unlock()
}

// reloadServicers - reload servicersName file
func reloadServicers(servicersPath string) {
	if !fileExist(servicersPath) {
		log2.Println(fmt.Sprintf("servicers file not found at %s; ignoring reload", servicersPath))
		return
	}

	newServicers := getServicersFromFile()
	newServicersMap := map[string]bool{}

	for i := range newServicers {
		ns := &newServicers[i]

		pk, err := crypto.NewPrivateKey(ns.PrivateKey)
		if err != nil {
			log2.Fatal(fmt.Errorf("error parsing private key at index=%d of the file %s", i, servicersPath))
		}

		address, err := sdk.AddressFromHex(pk.PubKey().Address().String())
		if err != nil {
			log2.Fatal(fmt.Errorf("error getting address from private key at index=%d of the file %s", i, servicersPath))
		}

		addressStr := address.String()

		newServicersMap[addressStr] = true
	}

	// looking for removed pk if any
	mutex.Lock()
	removedAddresses := make([]string, 0)
	for _, address := range servicerList {
		if _, ok := newServicersMap[address]; ok {
			// still there
			continue
		}

		removedAddresses = append(removedAddresses, address)
	}
	mutex.Unlock()

	// remove servicers
	removeServicers(removedAddresses)

	// reload servicer keys
	loadServicerNodes()
}

// initHotReload - initialize keys and chains file change detection
func initHotReload() {
	chainsPath := getChainsFilePath()
	servicersPath := getServicersFilePath()

	if app.GlobalMeshConfig.HotReloadInterval <= 0 {
		logger.Info("skipping hot reload due to hot_reload_interval is less or equal to 0")
		return
	}

	for {
		time.Sleep(time.Duration(app.GlobalMeshConfig.HotReloadInterval) * time.Millisecond)
		reloadChains(chainsPath)
		reloadServicers(servicersPath)
	}
}

// fileExist - check if file exists or not.
func fileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// loadHostedChains - load chainName file and read the content of it.
func loadHostedChains() *pocketTypes.HostedBlockchains {
	// create the chains path
	var chainsPath = getChainsFilePath()
	logger.Info("reading chains from path=" + chainsPath)
	// if file exists open, else create and open
	var jsonFile *os.File
	var bz []byte
	if _, err := os.Stat(chainsPath); err != nil && os.IsNotExist(err) {
		log2.Fatal(fmt.Sprintf("no chains.json found @ %s", chainsPath))
	}
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(chainsPath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	bz, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// unmarshal into the structure
	var hostedChainsSlice []pocketTypes.HostedBlockchain
	err = json.Unmarshal(bz, &hostedChainsSlice)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	m := make(map[string]pocketTypes.HostedBlockchain)
	for _, chain := range hostedChainsSlice {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			log2.Fatal(fmt.Sprintf("invalid ID: %s in network identifier in %s file", chain.ID, app.GlobalMeshConfig.ChainsName))
		}
		m[chain.ID] = chain
	}
	// return the map
	return &pocketTypes.HostedBlockchains{
		M: m,
		L: sync.RWMutex{},
	}
}

// storeRelay - persist relay to disk
func storeRelay(relay *pocketTypes.Relay) {
	hash := relay.RequestHash()
	logger.Debug(fmt.Sprintf("storing relay %s", relay.RequestHashString()))
	rb, err := json.Marshal(relay)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = relaysCacheDb.Put(hash, rb)
	if err != nil {
		logger.Error(fmt.Sprintf("error adding relay %s to cache. %s", relay.RequestHashString(), err.Error()))
	}

	return
}

// decodeCacheRelay - decode []byte relay from cache to pocketTypes.Relay
func decodeCacheRelay(body []byte) (relay *pocketTypes.Relay) {
	if err := json.Unmarshal(body, &relay); err != nil {
		logger.Error("error decoding cache relay")
		// todo: delete key from cache?
		return nil
	}
	return
}

// deleteCacheRelay - delete a key from relay cache
func deleteCacheRelay(relay *pocketTypes.Relay) {
	hash := relay.RequestHash()
	err := relaysCacheDb.Delete(hash)
	if err != nil {
		logger.Error("error deleting relay from cache %s", hex.EncodeToString(hash))
		return
	}

	return
}

// LoadAppSession - retrieve from cache (memory or persistent) an app session cache
func (s *servicer) LoadAppSession(hash []byte) (*appCache, bool) {
	sHash := hex.EncodeToString(hash)
	if v, ok := s.SessionCache.Load(sHash); ok {
		return v.(*appCache), ok
	}

	return nil, false
}

// decodeAppSession - decode []byte app session cache to appCache
//func decodeAppSession(body []byte) (appSession *appCache) {
//	if err := json.Unmarshal(body, &appSession); err != nil {
//		logger.Error("error decoding app session from cache")
//		return nil
//	}
//	return appSession
//}

// StoreAppSession - store in cache (memory and persistent) an appCache
func (s *servicer) StoreAppSession(hash []byte, appSession *appCache) {
	hashString := hex.EncodeToString(hash)
	s.SessionCache.Store(hashString, appSession)

	return
}

// DeleteAppSession - delete an app session from cache (memory and persistent)
func (s *servicer) DeleteAppSession(hash []byte) {
	sHash := hex.EncodeToString(hash)
	s.SessionCache.Delete(sHash)
}

// evaluateServicerError - this will change internalCache[hash].IsValid bool depending on the result of the evaluation
func evaluateServicerError(r *pocketTypes.Relay, err *sdkErrorResponse) (isSessionStillValid bool) {
	hash := getSessionHashFromRelay(r)

	isSessionStillValid = !isInvalidRelayCode(err.Code) // we should not retry if is invalid

	if isSessionStillValid {
		return isSessionStillValid
	}

	servicerNode := getServicerFromPubKey(r.Proof.ServicerPubKey)

	if appSession, ok := servicerNode.LoadAppSession(hash); ok {
		appSession.IsValid = isSessionStillValid
		appSession.Error = err
		servicerNode.StoreAppSession(hash, appSession)
	} else {
		logger.Error(
			fmt.Sprintf(
				"missing session hash=%s from cache; it should be there but if u see this after a restart it's ok.",
				hex.EncodeToString(hash),
			),
		)
	}

	return
}

// getSessionHashFromRelay - calculate the session header and late the hash of it
func getSessionHashFromRelay(r *pocketTypes.Relay) []byte {
	header := pocketTypes.SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: r.Proof.SessionBlockHeight,
	}

	return header.Hash()
}

func getServicerFromPubKey(pubKey string) *servicer {
	servicerAddress, err := getServicerAddressFromPubKeyAsString(pubKey)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"unable to decode servicer public key %s",
				pubKey,
			),
		)
		return nil
	}

	s, ok := servicerPkMap.Load(servicerAddress)

	if !ok {
		logger.Error(
			fmt.Sprintf(
				"unable to find servicer with address=%s",
				servicerAddress,
			),
		)
		return nil
	}

	return s.(*servicer)
}

// notifyServicer - call servicer to ack about the processed relay.
func notifyServicer(r *pocketTypes.Relay) {
	// discard this relay at the end of this function, to end this function the servicer will be retried N times
	defer deleteCacheRelay(r)

	result := meshRPCRelayResponse{}
	ctx := context.WithValue(context.Background(), "result", &result)
	jsonData, err := json.Marshal(r)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error encoding relay %s for servicer: %s",
				r.RequestHashString(),
				err.Error(),
			),
		)

		return
	}

	servicerAddress, err := getServicerAddressFromPubKeyAsString(r.Proof.ServicerPubKey)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"unable to decode service public key from relay %s to address",
				r.RequestHashString(),
			),
		)
		return
	}

	logger.Debug(
		fmt.Sprintf(
			"delivery relay %s notification to servicer %s",
			r.RequestHashString(),
			servicerAddress,
		),
	)

	s, ok := servicerPkMap.Load(servicerAddress)

	if !ok {
		logger.Error(
			fmt.Sprintf(
				"unable to find servicer with address=%s to notify relay %s",
				servicerAddress,
				r.RequestHashString(),
			),
		)
		return
	}

	servicerNode := s.(*servicer)

	requestURL := fmt.Sprintf(
		"%s%s?authtoken=%s&chain=%s&app=%s",
		servicerNode.ServicerURL,
		ServicerRelayEndpoint,
		servicerAuthToken.Value,
		r.Proof.Blockchain,
		r.Proof.Token.ApplicationPublicKey,
	)
	req, err := retryablehttp.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(fmt.Sprintf("error formatting Servicer URL: %s", err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ServicerHeader, servicerAddress)
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := relaysClient.Do(req)

	if err != nil {
		logger.Error(fmt.Sprintf("error dispatching relay to Servicer: %s", err.Error()))
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	isSuccess := resp.StatusCode == 200

	if result.Dispatch != nil && result.Dispatch.BlockHeight > servicerNode.Status.Height {
		servicerNode.Status.Height = result.Dispatch.BlockHeight
	}

	if !isSuccess {
		logger.Debug(
			fmt.Sprintf(
				"servicer %s reject relay %s\n: CODE=%d\nCODESPACE=%s\nMESSAGE=%s",
				servicerAddress, r.RequestHashString(),
				result.Error.Code, result.Error.Codespace, result.Error.Error,
			),
		)

		evaluateServicerError(r, result.Error)
	} else {
		logger.Debug(fmt.Sprintf("servicer processed relay %s successfully", r.RequestHashString()))

		header := pocketTypes.SessionHeader{
			ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
			Chain:              r.Proof.Blockchain,
			SessionBlockHeight: r.Proof.SessionBlockHeight,
		}

		hash := header.Hash()
		if appSession, ok := servicerNode.LoadAppSession(hash); ok {
			appSession.RemainingRelays -= 1
			logger.Debug(
				fmt.Sprintf(
					"servicer %s has %d remaining relays to process for app %s at blockchain %s",
					servicerNode.Address.String(),
					appSession.RemainingRelays,
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
				),
			)
			if appSession.RemainingRelays <= 0 {
				logger.Debug(
					fmt.Sprintf(
						"servicer %s exhaust relays for app %s at blockchain %s",
						servicerNode.Address.String(),
						r.Proof.Token.ApplicationPublicKey,
						r.Proof.Blockchain,
					),
				)
				appSession.IsValid = false
				appSession.Error = newSdkErrorFromPocketSdkError(pocketTypes.NewOverServiceError(ModuleName))
			}
			servicerNode.StoreAppSession(hash, appSession)
		}
	}

	return
}

// validate - evaluate relay to understand if should or not processed.
func validate(r *pocketTypes.Relay) sdk.Error {
	logger.Debug(fmt.Sprintf("validating relay %s", r.RequestHashString()))
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return pocketTypes.NewEmptyPayloadDataError(ModuleName)
	}
	// validate appPubKey
	if err := pocketTypes.PubKeyVerification(r.Proof.Token.ApplicationPublicKey); err != nil {
		return err
	}
	// validate chain
	if err := pocketTypes.NetworkIdentifierVerification(r.Proof.Blockchain); err != nil {
		return pocketTypes.NewEmptyChainError(ModuleName)
	}
	// validate the relay merkleHash = request merkleHash
	if r.Proof.RequestHash != r.RequestHashString() {
		return pocketTypes.NewRequestHashError(ModuleName)
	}
	// validate servicer public key
	servicerAddress, e := getServicerAddressFromPubKeyAsString(r.Proof.ServicerPubKey)
	if e != nil {
		return sdk.ErrInternal("could not convert servicer hex to public key")
	}
	// load servicer from servicer map, if not there maybe the servicer is pk is not loaded
	if _, ok := servicerPkMap.Load(servicerAddress); !ok {
		return pocketTypes.NewInvalidSessionError(ModuleName)
	}

	header := pocketTypes.SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: r.Proof.SessionBlockHeight,
	}

	hash := header.Hash()

	servicerNode := getServicerFromPubKey(r.Proof.ServicerPubKey)

	if appSession, ok := servicerNode.LoadAppSession(hash); !ok {
		result := &meshRPCSessionResult{}
		e2 := getAppSession(r, result)

		if e2 != nil {
			return newPocketSdkErrorFromSdkError(e2)
		}

		remainingRelays, _ := result.RemainingRelays.Int64()

		isValid := result.Success && remainingRelays > 0 && result.Error == nil

		servicerNode.StoreAppSession(header.Hash(), &appCache{
			PublicKey:       header.ApplicationPubKey,
			Chain:           header.Chain,
			Dispatch:        result.Dispatch,
			RemainingRelays: remainingRelays,
			IsValid:         isValid,
			Error:           result.Error,
		})
	} else {
		if !appSession.IsValid {
			if appSession.Error != nil {
				return newPocketSdkErrorFromSdkError(appSession.Error)
			} else {
				return sdk.ErrInternal("invalid session")
			}
		}
	}

	// is needed we call the node and validate if there is not a validation already in place get done by the cron?
	return nil
}

// addServiceMetricErrorFor - add to prometheus metrics an error for a servicer
func addServiceMetricErrorFor(blockchain string, address *sdk.Address) {
	pocketTypes.GlobalServiceMetric().AddErrorFor(blockchain, address)
}

// executeMeshHTTPRequest - run the non-native blockchain http request reusing chains http client.
func executeMeshHTTPRequest(payload, url, userAgent string, basicAuth pocketTypes.BasicAuth, method string, headers map[string]string) (string, error) {
	var m string
	if method == "" {
		m = pocketTypes.DEFAULTHTTPMETHOD
	} else {
		m = method
	}
	// generate an http request
	req, err := http.NewRequest(m, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	if basicAuth.Username != "" {
		req.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	// add headers if needed
	if len(headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	// execute the request
	resp, err := chainsClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// read all bz
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if app.GlobalMeshConfig.JSONSortRelayResponses {
		body = []byte(sortJSONResponse(string(body)))
	}
	logger.Debug(fmt.Sprintf("executing blockchain request:\nURL=%s\nMETHOD=%s\nREQ=%s\nSTATUS=%d\nRES=%s", url, m, payload, resp.StatusCode, string(body)))
	// return
	return string(body), nil
}

// sortJSONResponse - sorts json from a relay response
func sortJSONResponse(response string) string {
	var rawJSON map[string]interface{}
	// unmarshal into json
	if err := json.Unmarshal([]byte(response), &rawJSON); err != nil {
		return response
	}
	// marshal into json
	bz, err := json.Marshal(rawJSON)
	if err != nil {
		return response
	}
	return string(bz)
}

// execute - Attempts to do a request on the non-native blockchain specified
func execute(r *pocketTypes.Relay, hostedBlockchains *pocketTypes.HostedBlockchains, address *sdk.Address) (string, sdk.Error) {
	// retrieve the hosted blockchain url requested
	chain, err := hostedBlockchains.GetChain(r.Proof.Blockchain)
	if err != nil {
		// metric track
		addServiceMetricErrorFor(r.Proof.Blockchain, address)
		return "", err
	}
	_url := strings.Trim(chain.URL, `/`)
	if len(r.Payload.Path) > 0 {
		_url = _url + "/" + strings.Trim(r.Payload.Path, `/`)
	}

	// do basic http request on the relay
	res, er := executeMeshHTTPRequest(
		r.Payload.Data, _url,
		app.GlobalMeshConfig.UserAgent, chain.BasicAuth,
		r.Payload.Method, r.Payload.Headers,
	)
	if er != nil {
		// metric track
		addServiceMetricErrorFor(r.Proof.Blockchain, address)
		return res, pocketTypes.NewHTTPExecutionError(ModuleName, er)
	}
	return res, nil
}

// processRelay - call execute and create RelayResponse or Error in case. Also trigger relay metrics.
func processRelay(relay *pocketTypes.Relay) (*pocketTypes.RelayResponse, sdk.Error) {
	relayTimeStart := time.Now()
	logger.Debug(fmt.Sprintf("processing relay %s", relay.RequestHashString()))

	servicerAddress, e := getServicerAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)

	if e != nil {
		return nil, sdk.ErrInternal("could not convert servicer hex to public key")
	}

	s, ok := servicerPkMap.Load(servicerAddress)
	if !ok {
		return nil, sdk.ErrInternal("failed to find correct servicer PK")
	}

	servicerNode := s.(*servicer)

	// attempt to execute
	respPayload, err := execute(relay, chains, &servicerNode.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not send relay %s with error: %s", relay.RequestHashString(), err.Error()))
		return nil, err
	}
	// generate response object
	resp := &pocketTypes.RelayResponse{
		Response: respPayload,
		Proof:    relay.Proof,
	}

	// sign the response
	sig, er := servicerNode.PrivateKey.Sign(resp.Hash())
	if er != nil {
		logger.Error(
			fmt.Sprintf("could not sign response for address: %s with hash: %v, with error: %s",
				servicerAddress, resp.HashString(), er.Error()),
		)
		return nil, pocketTypes.NewKeybaseError(pocketTypes.ModuleName, er)
	}
	// attach the signature in hex to the response
	resp.Signature = hex.EncodeToString(sig)
	// track the relay time
	relayTime := time.Since(relayTimeStart)
	// add to metrics
	addRelayMetricsFunc := func() {
		logger.Debug(fmt.Sprintf("adding metric for relay %s", relay.RequestHashString()))
		pocketTypes.GlobalServiceMetric().AddRelayTimingFor(relay.Proof.Blockchain, float64(relayTime.Milliseconds()), &servicerNode.Address)
		pocketTypes.GlobalServiceMetric().AddRelayFor(relay.Proof.Blockchain, &servicerNode.Address)
	}
	go addRelayMetricsFunc()
	return resp, nil
}

// handleRelay - evaluate node status, validate relay payload and call processRelay
func handleRelay(r *pocketTypes.Relay) (res *pocketTypes.RelayResponse, dispatch *dispatchResponse, err error) {
	servicerAddress, e := getServicerAddressFromPubKeyAsString(r.Proof.ServicerPubKey)

	if e != nil {
		return nil, nil, errors.New("could not convert servicer hex to public key")
	}

	s, ok := servicerPkMap.Load(servicerAddress)
	if !ok {
		return nil, nil, errors.New("failed to find correct servicer PK")
	}

	servicerNode := s.(*servicer)

	if servicerNode.Status == nil {
		return nil, nil, fmt.Errorf("pocket node is currently unavailable")
	}

	if servicerNode.Status.IsStarting {
		return nil, nil, fmt.Errorf("pocket node is unable to retrieve synced status from tendermint node, cannot service in this state")
	}

	if servicerNode.Status.IsCatchingUp {
		return nil, nil, fmt.Errorf("pocket node is currently syncing to the blockchain, cannot service in this state")
	}

	err = validate(r)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"could not validate relay %s for app: %s, for chainID %v on node %s, at session height: %v, with error: %s",
				r.RequestHashString(),
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
				r.Meta.BlockHeight,
				err.Error(),
			),
		)

		return
	}

	// store relay on cache; once we hit this point this relay will be processed so should be notified to servicer even
	// if process is shutdown
	storeRelay(r)

	res, err = processRelay(r)

	if err != nil && pocketTypes.ErrorWarrantsDispatch(err) {
		// TODO: check if for request header hash we have the dispatch
		header := pocketTypes.SessionHeader{
			ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
			Chain:              r.Proof.Blockchain,
			SessionBlockHeight: r.Proof.SessionBlockHeight,
		}

		hash := header.Hash()

		if appSession, ok := servicerNode.LoadAppSession(hash); !ok {
			response := meshRPCSessionResult{}
			err1 := getAppSession(r, &response)
			if err1 != nil {
				logger.Error(
					fmt.Sprintf(
						"error getting app %s session; hash %s",
						r.Proof.Token.ApplicationPublicKey,
						hash,
					),
				)
			} else {
				dispatch = response.Dispatch
			}
		} else {
			dispatch = appSession.Dispatch
		}
	}

	// add to task group pool
	if servicerNode.Worker.Stopped() {
		// this should not happen, but just in case avoid a panic here.
		logger.Error(fmt.Sprintf("Worker of servicer %s was already stopped", servicerNode.Address.String()))
		return
	}

	servicerNode.Worker.Submit(func() {
		notifyServicer(r)
	})

	return
}

// meshHealth - handle mesh health request
func meshHealth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res := meshHealthResponse{
		Version: AppVersion,
	}
	j, er := json.Marshal(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshNodeRelay - handle mesh node relay request, call handleRelay
func meshNodeRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if cors(&w, r) {
		return
	}
	var relay = pocketTypes.Relay{}

	if err := PopModel(w, r, ps, &relay); err != nil {
		response := meshRPCRelayResponse{
			Success: false,
			Error:   newSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	logger.Debug(fmt.Sprintf("handling relay %s", relay.RequestHashString()))
	res, dispatch, err := handleRelay(&relay)

	if err != nil {
		response := meshRPCRelayResponse{
			Success:  false,
			Error:    newSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
			Dispatch: dispatch,
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	response := RPCRelayResponse{
		Signature: res.Signature,
		Response:  res.Response,
	}

	j, er := json.Marshal(response)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
	logger.Debug(fmt.Sprintf("relay %s done", relay.RequestHashString()))
}

// meshSimulateRelay - handle a simulated relay to test connectivity to the chains that this should be serving.
func meshSimulateRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = simRelayParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	chain, err := chains.GetChain(params.RelayNetworkID)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	_url := strings.Trim(chain.URL, `/`)
	if len(params.Payload.Path) > 0 {
		_url = _url + "/" + strings.Trim(params.Payload.Path, `/`)
	}

	logger.Debug(
		fmt.Sprintf(
			"executing simulated relay of chain %s",
			chain.ID,
		),
	)
	// do basic http request on the relay
	res, er := executeMeshHTTPRequest(
		params.Payload.Data, _url, app.GlobalMeshConfig.UserAgent,
		chain.BasicAuth, params.Payload.Method, params.Payload.Headers,
	)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, res, r.URL.Path, r.Host)
}

// newSdkErrorFromPocketSdkError - return a mesh node sdkErrorResponse from a pocketcore sdk.Error
func newSdkErrorFromPocketSdkError(e sdk.Error) *sdkErrorResponse {
	return &sdkErrorResponse{
		Code:      e.Code(),
		Codespace: e.Codespace(),
		Error:     e.Error(),
	}
}

// newPocketSdkErrorFromSdkError - return a pocketcore sdk.Error from a mesh node sdkErrorResponse
func newPocketSdkErrorFromSdkError(e *sdkErrorResponse) sdk.Error {
	return sdk.NewError(e.Codespace, e.Code, errors.New(e.Error).Error())
}

// checkAddressIsSupported - use on pocket node side to verify if the address is handled by the running process.
func checkAddressIsSupported(address string) error {
	if address == "" {
		return errors.New("missing query param address")
	} else {
		if pocketTypes.GlobalPocketConfig.LeanPocket {
			// if lean pocket enabled, grab the targeted servicer through the relay proof
			nodeAddress, err := sdk.AddressFromHex(address)
			if err != nil {
				return errors.New("could not convert servicer hex")
			}
			_, err = pocketTypes.GetPocketNodeByAddress(&nodeAddress)
			if err != nil {
				return errors.New("failed to find correct servicer private key")
			}
		} else {
			// get self node (your validator) from the current state
			node := pocketTypes.GetPocketNode()
			nodeAddress := node.GetAddress()
			if nodeAddress.String() != address {
				return errors.New("failed to find correct servicer private key")
			}
		}
	}

	return nil
}

// meshServicerNodeRelay - receive relays that was processed by mesh node on /v1/client/relay
func meshServicerNodeRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var relay = pocketTypes.Relay{}

	if cors(&w, r) {
		return
	}

	token := r.URL.Query().Get("authtoken")
	if token != app.AuthToken.Value {
		WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	if verify == "true" {
		code := 200
		// useful just to test that mesh node is able to reach servicer
		response := meshRPCRelayResult{
			Success:  true,
			Error:    nil,
			Dispatch: nil,
		}

		address := r.URL.Query().Get("address")
		if err := checkAddressIsSupported(address); err != nil {
			response.Success = false
			response.Error = err
			code = 400
		}

		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
		return
	}

	if err := PopModel(w, r, ps, &relay); err != nil {
		response := RPCRelayErrorResponse{
			Error: err,
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	_, dispatch, err := app.PCA.HandleRelay(relay, true)
	if err != nil {
		response := meshRPCRelayResult{
			Success:  false,
			Error:    err,
			Dispatch: dispatch,
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	response := meshRPCRelayResult{
		Success:  true,
		Dispatch: dispatch,
	}
	j, er := json.Marshal(response)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNodeSession - receive requests from mesh node to validate a session for an app/servicer/blockchain on the servicer node data
func meshServicerNodeSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var session pocketTypes.MeshSession

	token := r.URL.Query().Get("authtoken")
	if token != app.AuthToken.Value {
		WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	if verify == "true" {
		code := 200
		// useful just to test that mesh node is able to reach servicer
		response := meshRPCSessionResult{
			Success:  true,
			Error:    nil,
			Dispatch: nil,
		}

		address := r.URL.Query().Get("address")
		if err := checkAddressIsSupported(address); err != nil {
			response.Success = false
			response.Error = newSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error()))
			code = 400
		}

		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
		return
	}

	if err := PopModel(w, r, ps, &session); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	res, err := app.PCA.HandleMeshSession(session)

	if err != nil {
		response := meshRPCSessionResult{
			Success: false,
			Error:   newSdkErrorFromPocketSdkError(err),
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	dispatch := dispatchResponse{
		BlockHeight: res.Session.BlockHeight,
		Session: dispatchSession{
			Header: res.Session.Session.SessionHeader,
			Key:    hex.EncodeToString(res.Session.Session.SessionKey),
			Nodes:  make([]dispatchSessionNode, 0),
		},
	}

	for i := range res.Session.Session.SessionNodes {
		sNode, ok := res.Session.Session.SessionNodes[i].(nodesTypes.Validator)
		if !ok {
			continue
		}
		dispatch.Session.Nodes = append(dispatch.Session.Nodes, dispatchSessionNode{
			Address:       sNode.Address.String(),
			Chains:        sNode.Chains,
			Jailed:        sNode.Jailed,
			OutputAddress: sNode.OutputAddress.String(),
			PublicKey:     sNode.PublicKey.String(),
			ServiceUrl:    sNode.ServiceURL,
			Status:        sNode.Status,
			Tokens:        sNode.GetTokens().String(),
			UnstakingTime: sNode.UnstakingCompletionTime,
		})
	}

	response := meshRPCSessionResult{
		Success:         true,
		Dispatch:        &dispatch,
		RemainingRelays: json.Number(strconv.FormatInt(res.RemainingRelays, 10)),
	}
	j, er := json.Marshal(response)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// isAuthorized - check if the request is authorized using authToken of the auth.json file
func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	token := r.URL.Query().Get("authtoken")
	if token == meshAuthToken.Value {
		return true
	} else {
		WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return false
	}
}

// meshStop - gracefully stop mesh rpc server. Also, this should stop new relays and wait/flush all pending relays, otherwise they will get loose.
func meshStop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}
	StopMeshRPC()
	fmt.Println("Stop Successful, PID:" + fmt.Sprint(os.Getpid()))
	os.Exit(0)
}

// meshChains - return load chains from app.GlobalMeshConfig.ChainsName file
func meshChains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	c := make([]pocketTypes.HostedBlockchain, 0)

	for _, chain := range chains.M {
		c = append(c, chain)
	}

	j, err := json.Marshal(c)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	WriteRaw(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNode - return servicer node configured by servicer_priv_key.json - return address
func meshServicerNode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	servicers := make([]types4.PublicPocketNode, 0)

	mutex.Lock()
	for _, a := range servicerList {
		servicers = append(servicers, types4.PublicPocketNode{
			Address: a,
		})
	}
	mutex.Unlock()

	j, err := json.Marshal(servicers)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	WriteRaw(w, string(j), r.URL.Path, r.Host)
}

// meshUpdateChains - update chains in memory and also chains.json file.
func meshUpdateChains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	var hostedChainsSlice []pocketTypes.HostedBlockchain
	if err := PopModel(w, r, ps, &hostedChainsSlice); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	m := make(map[string]pocketTypes.HostedBlockchain)
	for _, chain := range hostedChainsSlice {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			WriteErrorResponse(w, 400, fmt.Sprintf("invalid ID: %s in network identifier in json", chain.ID))
			return
		}
	}
	chains = &pocketTypes.HostedBlockchains{
		M: m,
		L: sync.RWMutex{},
	}

	j, er := json.Marshal(chains.M)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)

	updateChains(hostedChainsSlice)
}

// getMeshRoutes - return routes that will be handled/proxied by mesh rpc server
func getMeshRoutes(simulation bool) Routes {
	routes := Routes{
		// Proxy
		Route{Name: "AppVersion", Method: "GET", Path: "/v1", HandlerFunc: proxyRequest},
		Route{Name: "Health", Method: "GET", Path: "/v1/health", HandlerFunc: proxyRequest},
		Route{Name: "Challenge", Method: "POST", Path: "/v1/client/challenge", HandlerFunc: proxyRequest},
		Route{Name: "ChallengeCORS", Method: "OPTIONS", Path: "/v1/client/challenge", HandlerFunc: proxyRequest},
		Route{Name: "HandleDispatch", Method: "POST", Path: "/v1/client/dispatch", HandlerFunc: proxyRequest},
		Route{Name: "HandleDispatchCORS", Method: "OPTIONS", Path: "/v1/client/dispatch", HandlerFunc: proxyRequest},
		Route{Name: "SendRawTx", Method: "POST", Path: "/v1/client/rawtx", HandlerFunc: proxyRequest},
		Route{Name: "Stop", Method: "POST", Path: "/v1/private/stop", HandlerFunc: proxyRequest},
		Route{Name: "QueryChains", Method: "POST", Path: "/v1/private/chains", HandlerFunc: proxyRequest},
		Route{Name: "QueryAccount", Method: "POST", Path: "/v1/query/account", HandlerFunc: proxyRequest},
		Route{Name: "QueryAccounts", Method: "POST", Path: "/v1/query/accounts", HandlerFunc: proxyRequest},
		Route{Name: "QueryAccountTxs", Method: "POST", Path: "/v1/query/accounttxs", HandlerFunc: proxyRequest},
		Route{Name: "QueryACL", Method: "POST", Path: "/v1/query/acl", HandlerFunc: proxyRequest},
		Route{Name: "QueryAllParams", Method: "POST", Path: "/v1/query/allparams", HandlerFunc: proxyRequest},
		Route{Name: "QueryApp", Method: "POST", Path: "/v1/query/app", HandlerFunc: proxyRequest},
		Route{Name: "QueryAppParams", Method: "POST", Path: "/v1/query/appparams", HandlerFunc: proxyRequest},
		Route{Name: "QueryApps", Method: "POST", Path: "/v1/query/apps", HandlerFunc: proxyRequest},
		Route{Name: "QueryBalance", Method: "POST", Path: "/v1/query/balance", HandlerFunc: proxyRequest},
		Route{Name: "QueryBlock", Method: "POST", Path: "/v1/query/block", HandlerFunc: proxyRequest},
		Route{Name: "QueryBlockTxs", Method: "POST", Path: "/v1/query/blocktxs", HandlerFunc: proxyRequest},
		Route{Name: "QueryDAOOwner", Method: "POST", Path: "/v1/query/daoowner", HandlerFunc: proxyRequest},
		Route{Name: "QueryHeight", Method: "POST", Path: "/v1/query/height", HandlerFunc: proxyRequest},
		Route{Name: "QueryNode", Method: "POST", Path: "/v1/query/node", HandlerFunc: proxyRequest},
		Route{Name: "QueryNodeClaim", Method: "POST", Path: "/v1/query/nodeclaim", HandlerFunc: proxyRequest},
		Route{Name: "QueryNodeClaims", Method: "POST", Path: "/v1/query/nodeclaims", HandlerFunc: proxyRequest},
		Route{Name: "QueryNodeParams", Method: "POST", Path: "/v1/query/nodeparams", HandlerFunc: proxyRequest},
		Route{Name: "QueryNodes", Method: "POST", Path: "/v1/query/nodes", HandlerFunc: proxyRequest},
		Route{Name: "QueryParam", Method: "POST", Path: "/v1/query/param", HandlerFunc: proxyRequest},
		Route{Name: "QueryPocketParams", Method: "POST", Path: "/v1/query/pocketparams", HandlerFunc: proxyRequest},
		Route{Name: "QueryState", Method: "POST", Path: "/v1/query/state", HandlerFunc: proxyRequest},
		Route{Name: "QuerySupply", Method: "POST", Path: "/v1/query/supply", HandlerFunc: proxyRequest},
		Route{Name: "QuerySupportedChains", Method: "POST", Path: "/v1/query/supportedchains", HandlerFunc: proxyRequest},
		Route{Name: "QueryTX", Method: "POST", Path: "/v1/query/tx", HandlerFunc: proxyRequest},
		Route{Name: "QueryUpgrade", Method: "POST", Path: "/v1/query/upgrade", HandlerFunc: proxyRequest},
		Route{Name: "QuerySigningInfo", Method: "POST", Path: "/v1/query/signinginfo", HandlerFunc: proxyRequest},
		Route{Name: "LocalNodes", Method: "POST", Path: "/v1/private/nodes", HandlerFunc: proxyRequest},
		Route{Name: "QueryUnconfirmedTxs", Method: "POST", Path: "/v1/query/unconfirmedtxs", HandlerFunc: proxyRequest},
		Route{Name: "QueryUnconfirmedTx", Method: "POST", Path: "/v1/query/unconfirmedtx", HandlerFunc: proxyRequest},
		// start mesh things
		Route{Name: "MeshHealth", Method: "GET", Path: "/v1/mesh/health", HandlerFunc: meshHealth},
		Route{Name: "MeshService", Method: "POST", Path: "/v1/client/relay", HandlerFunc: reuseBody(meshNodeRelay)},
		Route{Name: "StopMeshNode", Method: "POST", Path: "/v1/private/mesh/stop", HandlerFunc: meshStop},
		Route{Name: "QueryMeshNodeChains", Method: "POST", Path: "/v1/private/mesh/chains", HandlerFunc: meshChains},
		Route{Name: "MeshNodeServicer", Method: "POST", Path: "/v1/private/mesh/servicer", HandlerFunc: meshServicerNode},
		Route{Name: "UpdateMeshNodeChains", Method: "POST", Path: "/v1/private/mesh/updatechains", HandlerFunc: meshUpdateChains},
	}

	// check if simulation is turn on
	if simulation {
		simRoute := Route{Name: "SimulateRequest", Method: "POST", Path: "/v1/client/sim", HandlerFunc: meshSimulateRelay}
		routes = append(routes, simRoute)
	}

	return routes
}

// checkServicerHealth - check server /v1/health endpoint
func checkServicerHealth(servicerNode *servicer) error {
	requestURL := fmt.Sprintf("%s/v1/health", servicerNode.ServicerURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := servicerClient.Do(req)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return // add log here
		}
	}(resp.Body)

	if resp.StatusCode != 200 || !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		if resp.StatusCode != 200 {
			err = errors.New(fmt.Sprintf("servicer %s is returning a non 200 code response from /v1/health", servicerNode.Address.String()))
		} else {
			err = errors.New(fmt.Sprintf("servicer %s is returning a non json response from /v1/health", servicerNode.Address.String()))
		}

		return err
	}

	res := &app.HealthResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}

	servicerNode.Status = res

	return nil
}

// checkNotifyServicerEndpoint - check servicer ServicerRelayEndpoint endpoint
func checkServicerEndpoint(servicerNode *servicer, endpoint string) error {
	requestURL := fmt.Sprintf(
		"%s%s?authtoken=%s&verify=true&address=%s",
		servicerNode.ServicerURL,
		endpoint,
		servicerAuthToken.Value,
		servicerNode.Address.String(),
	)
	req, err := http.NewRequest("POST", requestURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ServicerHeader, servicerNode.Address.String())
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := servicerClient.Do(req)

	if err != nil {
		return errors.New(
			fmt.Sprintf(
				"error verifying %s connectivity for servicer %s error=%s",
				endpoint,
				servicerNode.Address.String(),
				err.Error(),
			),
		)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	isSuccess := resp.StatusCode == 200

	if !isSuccess {
		return errors.New(
			fmt.Sprintf(
				"error verifying %s connectivity for servicer %s error=non 200 code code=%d",
				endpoint,
				servicerNode.Address.String(),
				resp.StatusCode,
			),
		)
	}

	return nil
}

// connectivityChecks - run check over critical endpoints that mesh node need to be able to reach on servicer
func connectivityChecks() {
	logger.Info("start connectivity checks")
	mutex.Lock()
	totalServicers := len(servicerList)
	mutex.Unlock()
	connectivityWorkerPool := pond.New(
		totalServicers, 0, pond.MinWorkers(10),
		pond.IdleTimeout(time.Duration(app.GlobalMeshConfig.WorkersIdleTimeout)*time.Millisecond),
		pond.Strategy(pond.Eager()),
	)

	var success uint32

	// check health for all the servicer nodes before start.
	servicerPkMap.Range(func(key, value any) bool {
		servicerNode := value.(*servicer)
		connectivityWorkerPool.Submit(func() {
			e1 := checkServicerEndpoint(servicerNode, ServicerRelayEndpoint)
			e2 := checkServicerEndpoint(servicerNode, ServicerSessionEndpoint)
			if e1 != nil {
				logger.Error(e1.Error())
			}
			if e2 != nil {
				logger.Error(e2.Error())
			}
			atomic.AddUint32(&success, 1)
			return
		})
		return true
	})

	// Wait for all HTTP requests to complete.
	connectivityWorkerPool.StopAndWait()

	if success == 0 {
		logger.Error(fmt.Sprintf("any servicer was able to be reach at %s", ServicerRelayEndpoint))
		log2.Fatal(fmt.Sprintf("servicers=%d; reachable=%d", totalServicers, success))
	}

	if int(success) != totalServicers {
		logger.Error(fmt.Sprintf("IMPORTANT!!! %d servicer are not reachable at %s", totalServicers-int(success), ServicerRelayEndpoint))
		logger.Error("you should stop this and fix the connectivity before continue")
	}
	logger.Info("connectivity check pass")
}

// nodeHealthStatusPooling - schedule a node heal pooling
func nodeHealthStatusPooling(c *cron.Cron, servicerNode *servicer) {
	e := checkServicerHealth(servicerNode)
	if e != nil {
		logger.Error(
			fmt.Sprintf(
				"servicer %s failed health check at: GET %s/v1/health error=%s",
				servicerNode.Address,
				servicerNode.ServicerURL,
				e.Error(),
			),
		)
	}

	_, err := c.AddFunc("@every 60s", func() {
		value, ok := servicerPkMap.Load(servicerNode.Address.String())
		if !ok {
			log2.Fatal(fmt.Sprintf("unable to load %s from servicers", servicerNode.Address.String()))
		}

		servicerNode := value.(*servicer)

		e := checkServicerHealth(servicerNode)
		if e != nil {
			logger.Error(
				fmt.Sprintf(
					"servicer %s failed health check at: GET %s/v1/health error=%s",
					servicerNode.Address,
					servicerNode.ServicerURL,
					e.Error(),
				),
			)
		}
	})

	if err != nil {
		log2.Fatal(err)
	}
}

// cleanOldSessions - clean up sessions that are longer than 50 blocks (just to be sure they are not needed)
func cleanOldSessions(c *cron.Cron) {
	_, err := c.AddFunc("@every 30m", func() {
		servicerPkMap.Range(func(_, sn any) bool {
			servicerNode := sn.(*servicer)
			servicerNode.SessionCache.Range(func(key any, ss any) bool {
				appSession := ss.(*appCache)
				hash, err := hex.DecodeString(key.(string))
				if err != nil {
					logger.Error("error decoding session hash to delete from cache " + err.Error())
					return true
				}

				if appSession.Dispatch == nil {
					servicerNode.DeleteAppSession(hash)
				} else if appSession.Dispatch.Session.Header.SessionBlockHeight < (servicerNode.Status.Height - 6) {
					servicerNode.DeleteAppSession(hash)
				}

				return true
			})
			return true
		})
	})

	if err != nil {
		log2.Fatal(err)
	}
}

// getChainsFilePath - return chains file path resolved by config.json
func getChainsFilePath() string {
	return app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ChainsName
}

// getServicersFilePath - return servicers file path resolved by config.json
func getServicersFilePath() string {
	return app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ServicerPrivateKeyFile
}

func getServicersFromFile() []servicerFile {
	path := getServicersFilePath()
	logger.Info("reading private key path=" + path)
	var readServicers []servicerFile
	data, err := os.ReadFile(path)
	if err != nil {
		log2.Fatal(fmt.Errorf("an error occurred attempting to read the servicer key file: %s", err.Error()))
	}

	if err := json.Unmarshal(data, &readServicers); err != nil {
		log2.Fatal(fmt.Errorf("an error occurred attempting to parse the servicer key file: %s", err.Error()))
	}

	return readServicers
}

// loadServicerNodes - read servicer address and cast to sdk.Address
func loadServicerNodes() int {
	servicersPath := getServicersFilePath()

	readServicers := getServicersFromFile()
	if len(readServicers) == 0 {
		log2.Fatal(fmt.Errorf("read 0 servicers from servicer key file: %s", servicersPath))
	}

	loadedServicerList := make([]string, 0)

	for i, s := range readServicers {
		pk, err := crypto.NewPrivateKey(s.PrivateKey)
		if err != nil {
			log2.Fatal(fmt.Errorf("error parsing private key at index=%d of the file %s", i, servicersPath))
		}

		address, err := sdk.AddressFromHex(pk.PubKey().Address().String())
		if err != nil {
			log2.Fatal(fmt.Errorf("error getting address from private key at index=%d of the file %s", i, servicersPath))
		}

		addressStr := address.String()

		if _, ok := servicerPkMap.Load(addressStr); !ok {
			logger.Info(fmt.Sprintf("initializing servicer %s health cron job", addressStr))
			servicerCronJobs := cron.New()

			sRecord := servicer{
				PrivateKey:  pk,
				Address:     address,
				ServicerURL: s.ServicerUrl,
				Status: &app.HealthResponse{
					IsStarting:   true,
					IsCatchingUp: true,
					Height:       1,
				},
				Crons:        servicerCronJobs,
				Worker:       newWorker(fmt.Sprintf("Notify Servicer %s", addressStr)),
				SessionCache: sync.Map{},
			}

			servicerPkMap.Store(addressStr, &sRecord)

			// check node status before start and schedule job
			nodeHealthStatusPooling(servicerCronJobs, &sRecord)

			servicerCronJobs.Start()
			logger.Info(fmt.Sprintf("servicer %s health cron job started", addressStr))
		}

		loadedServicerList = append(loadedServicerList, addressStr)
	}

	totalServicers := len(loadedServicerList)
	mutex.Lock()
	servicerList = loadedServicerList
	mutex.Unlock()

	return totalServicers
}

// removeServicers - stop receiving work for them and remove after that.
func removeServicers(servicers []string) {
	if len(servicers) > 0 {
		logger.Debug(
			fmt.Sprintf(
				"start drain of %d servicers after a hot reload of %s",
				len(servicers),
				getServicersFilePath(),
			),
		)
	}
	for _, address := range servicers {
		s, ok := servicerPkMap.LoadAndDelete(address)

		if !ok {
			// is not there
			continue
		}
		servicerNode := s.(*servicer)
		logger.Info(
			fmt.Sprintf(
				"removing servicer %s after a hot reload of %s",
				servicerNode.Address.String(),
				getServicersFilePath(),
			),
		)
		servicerNode.Crons.Stop()
		servicerNode.Worker.StopAndWait()
		logger.Info(
			fmt.Sprintf(
				"servicer %s successfuly drained and removed",
				servicerNode.Address.String(),
			),
		)
	}
}

// getAuthTokenFromFile - read from path a json that match sdk.AuthToken struct
func getAuthTokenFromFile(path string) sdk.AuthToken {
	logger.Info("reading authtoken from path=" + path)
	t := sdk.AuthToken{}

	var jsonFile *os.File
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			return
		}
	}(jsonFile)

	if _, err := os.Stat(path); err == nil {
		jsonFile, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log2.Fatalf("cannot open auth token json file: " + err.Error())
		}
		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log2.Fatalf("cannot read auth token json file: " + err.Error())
		}
		err = json.Unmarshal(b, &t)
		if err != nil {
			log2.Fatalf("cannot read auth token json file into json: " + err.Error())
		}
	}

	return t
}

// loadAuthTokens - load mesh node authtoken and servicer authtoken
func loadAuthTokens() {
	dataDir := app.GlobalMeshConfig.DataDir
	meshNodeAuthFile := dataDir + app.FS + app.GlobalMeshConfig.AuthTokenFile
	servicerAuthFile := dataDir + app.FS + app.GlobalMeshConfig.ServicerAuthTokenFile
	// used to authenticate request to mesh node on /v1/private paths
	meshAuthToken = getAuthTokenFromFile(meshNodeAuthFile)
	// used to call servicer node on private path to notify about relays
	servicerAuthToken = getAuthTokenFromFile(servicerAuthFile)
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

// initLogger - initialize logger
func initLogger() (logger log.Logger) {
	logger = log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			fmt.Printf("expected level key to be first, got %v", keyvals[0])
			log2.Fatal(1)
		}
		switch keyvals[1].(kitlevel.Value).String() {
		case "info":
			return term.FgBgColor{Fg: term.Green}
		case "debug":
			return term.FgBgColor{Fg: term.DarkBlue}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	})
	logger, err := flags.ParseLogLevel(app.GlobalMeshConfig.LogLevel, logger, "info")
	if err != nil {
		log2.Fatal(err)
	}
	return
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

// newWorker - generate a new worker.
func newWorker(handler string) *pond.WorkerPool {
	panicHandler := func(p interface{}) {
		logger.Error(fmt.Sprintf("%s Worker task throw panic error: %v", handler, p))
	}

	var strategy pond.ResizingStrategy

	switch app.GlobalMeshConfig.WorkerStrategy {
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
		log2.Fatal(fmt.Sprintf("strategy %s is not a valid option; allowed values are: lazy|eager|balanced", app.GlobalMeshConfig.WorkerStrategy))
	}

	return pond.New(
		app.GlobalMeshConfig.MaxWorkers, app.GlobalMeshConfig.MaxWorkersCapacity,
		pond.IdleTimeout(time.Duration(app.GlobalMeshConfig.WorkersIdleTimeout)*time.Millisecond),
		pond.PanicHandler(panicHandler),
		pond.Strategy(strategy),
	)
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
			servicerAddress, err := getServicerAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)
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

			s, ok := servicerPkMap.Load(servicerAddress)
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

			servicerNode, ok := s.(*servicer)
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

			servicerNode.Worker.Submit(func() {
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

// GetServicerMeshRoutes - return routes that need to be added to servicer to allow mesh node to communicate with.
func GetServicerMeshRoutes() Routes {
	routes := Routes{
		{Name: "MeshRelay", Method: "POST", Path: ServicerRelayEndpoint, HandlerFunc: meshServicerNodeRelay},
		{Name: "MeshSession", Method: "POST", Path: ServicerSessionEndpoint, HandlerFunc: meshServicerNodeSession},
	}

	return routes
}

// StopMeshRPC - stop http server
func StopMeshRPC() {
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
	servicerPkMap.Range(func(key, value any) bool {
		servicerNode := value.(*servicer)

		logger.Debug(fmt.Sprintf("stopping worker pool of servicer %s", servicerNode.Address.String()))
		servicerNode.Worker.Stop()
		logger.Debug(fmt.Sprintf("worker pool of servicer %s stopped!", servicerNode.Address.String()))

		logger.Debug(fmt.Sprintf("stopping health cron job of servicer %s", servicerNode.Address.String()))
		servicerNode.Crons.Stop()
		logger.Debug(fmt.Sprintf("health cron job of servicer %s stopped!", servicerNode.Address.String()))

		return true
	})
	interceptorPool.Stop()
	logger.Info("worker pools stopped!")

	logger.Info("stopping clean session cron job")
	cronJobs.Stop()
	logger.Info("clean session job stopped!")
}

// StartMeshRPC - Start mesh rpc server
func StartMeshRPC(simulation bool) {
	ctx, cancel := context.WithCancel(context.Background())
	finish = cancel
	defer cancel()
	logger = initLogger()
	// initialize pseudo random to choose servicer url
	rand.Seed(time.Now().Unix())
	// load auth token files (servicer and mesh node)
	loadAuthTokens()
	// retrieve the nonNative blockchains your node is hosting
	chains = loadHostedChains()
	// turn on chains hot reload
	go initHotReload()
	// initialize prometheus metrics
	pocketTypes.InitGlobalServiceMetric(chains, logger, app.GlobalMeshConfig.PrometheusAddr, app.GlobalMeshConfig.PrometheusMaxOpenfiles)
	// instantiate all the http clients used to call Chains and Servicer
	prepareHttpClients()
	// read mesh node routes
	routes := getMeshRoutes(simulation)
	// read servicer
	totalServicers := loadServicerNodes()
	// check servicers are reachable at required endpoints
	connectivityChecks()
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
			Router(routes),
			time.Duration(app.GlobalMeshConfig.RPCTimeout)*time.Millisecond,
			"Server Timeout Handling Request",
		),
	}

	go catchSignal()

	logger.Info(fmt.Sprintf("start serving relay as mesh node for %d servicer nodes", totalServicers))

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
