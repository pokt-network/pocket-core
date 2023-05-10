package mesh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/puzpuzpuz/xsync"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"io/ioutil"
	log2 "log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

type RichChain struct {
	Label string `json:"label"`
}

var (
	ChainNameMap = xsync.NewMapOf[string]()
)

// getChainsFilePath - return chains file path resolved by config.json
func getChainsFilePath() string {
	return app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ChainsName
}

// loadLocalChainsNameMap - load local chain name map
func loadLocalChainsNameMap() ([]byte, error) {
	var chainsPath = app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ChainsNameMap

	if _, err := os.Stat(chainsPath); err != nil && os.IsNotExist(err) {
		e := errors.New(fmt.Sprintf("chains_name_map not found @ %s", chainsPath))
		logger.Error(e.Error())
		return []byte{}, e
	}

	return os.ReadFile(chainsPath)
}

// loadRemoteChainsNameMap - load remote chain name map
func loadRemoteChainsNameMap() ([]byte, error) {
	req, err := http.NewRequest("GET", app.GlobalMeshConfig.RemoteChainsNameMap, nil)
	req.Header.Set("Content-Type", "application/json")
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := chainsClient.Do(req)

	if err != nil {
		return []byte{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	// read the body just to allow http 1.x be able to reuse the connection
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := errors.New(fmt.Sprintf("Couldn't parse response body. Error: %s", err.Error()))
		return []byte{}, e
	}

	isSuccess := resp.StatusCode == 200

	if !isSuccess {
		return []byte{}, errors.New(
			fmt.Sprintf(
				"error=StatusCode != 200 code=%d",
				resp.StatusCode,
			),
		)
	}

	return body, nil
}

// loadChainsNameMap - load chain name map from local or remote depending on config.json file
func loadChainsNameMap() {
	var data []byte

	shouldSkip := app.GlobalMeshConfig.RemoteChainsNameMap == "" && app.GlobalMeshConfig.ChainsNameMap == ""

	if shouldSkip {
		return
	}

	// check which one get data
	if app.GlobalMeshConfig.ChainsNameMap != "" {
		logger.Debug(fmt.Sprintf("loading chains name map from local %s", app.GlobalMeshConfig.ChainsNameMap))
		d, err := loadLocalChainsNameMap()
		if err != nil {
			logger.Error("unable to load local chains map", "err", err)
			return
		}

		data = d
	}

	// remote has precedence over local because could be shared among multiple nodes.
	if app.GlobalMeshConfig.RemoteChainsNameMap != "" {
		logger.Debug(fmt.Sprintf("loading chains name map from remote %s", app.GlobalMeshConfig.RemoteChainsNameMap))
		d, err := loadRemoteChainsNameMap()
		if err != nil {
			logger.Error("unable to load remote chains map", "err", err)
			return
		}

		data = d
	}

	if len(data) == 0 {
		logger.Error("there is no data to be read/load from chains map (local or remote)")
		return
	}

	plainSchemaLoader := gojsonschema.NewSchemaLoader()
	plainSchemaStringLoader := gojsonschema.NewStringLoader(plainChainsMapSchema)
	plainSchema, plainSchemaError := plainSchemaLoader.Compile(plainSchemaStringLoader)
	if plainSchemaError != nil {
		log2.Fatal(fmt.Errorf("an error occurred loading plain chains name map json schema: %s", plainSchemaError.Error()))
	}

	richSchemaLoader := gojsonschema.NewSchemaLoader()
	richSchemaStringLoader := gojsonschema.NewStringLoader(richChainsMapSchema)
	richSchema, richSchemaError := richSchemaLoader.Compile(richSchemaStringLoader)
	if richSchemaError != nil {
		log2.Fatal(fmt.Errorf("an error occurred loading rich chains name map json schema: %s", richSchemaError.Error()))
	}

	strData := gojsonschema.NewStringLoader(string(data[:]))
	if r, e := plainSchema.Validate(strData); e != nil || len(r.Errors()) > 0 {
		if r2, e2 := richSchema.Validate(strData); e2 != nil || len(r2.Errors()) > 0 {
			logger.Error("unable to validate json from chains name map sources (local or remote)")
			logger.Error(fmt.Sprintf("chains name rich schema errors: %v", r.Errors()))
			logger.Error(fmt.Sprintf("chains name plain schema errors: %v", r2.Errors()))
			return
		} else {
			// rich chain format
			richRecords := map[string]RichChain{}
			if err := json.Unmarshal(data, &richRecords); err != nil {
				logger.Error("an error occurred attempting to parse rich chains name map", "err", err)
				return
			}

			for chain, richChain := range richRecords {
				ChainNameMap.Store(chain, richChain.Label)
			}
		}
	} else {
		// fallback
		// plain chain format
		plainRecords := map[string]string{}
		if err := json.Unmarshal(data, &plainRecords); err != nil {
			logger.Error("an error occurred attempting to parse plain chains name map", "err", err)
			return
		}

		for chain, name := range plainRecords {
			ChainNameMap.Store(chain, name)
		}
	}
}

// UpdateChains - update chainName file with the retrieve chains value.
func UpdateChains(c []pocketTypes.HostedBlockchain) (*pocketTypes.HostedBlockchains, error) {
	m := make(map[string]pocketTypes.HostedBlockchain)
	for _, chain := range c {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			return chains, errors.New(fmt.Sprintf("invalid ID: %s in network identifier in json", chain.ID))
		}

		m[chain.ID] = chain
	}

	chains.L.Lock()
	chains.M = m
	chains.L.Unlock()

	var chainsPath = app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ChainsName
	var jsonFile *os.File
	if _, err := os.Stat(chainsPath); err != nil && os.IsNotExist(err) {
		e := errors.New(fmt.Sprintf("no chains.json found @ %s", chainsPath))
		logger.Error(e.Error())
		return chains, e
	}
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(chainsPath, os.O_WRONLY, os.ModePerm)
	if err != nil {
		log2.Fatal(app.NewInvalidChainsError(err))
	}
	// create dummy input for the file
	res, err := json.MarshalIndent(c, "", "  ")
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
	return chains, nil
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

// reloadChains - reload chainsName file
func reloadChains() {
	logger.Debug("initializing chains reload process")
	chainsPath := getChainsFilePath()
	// if file exists open, else create and open
	var jsonFile *os.File
	var bz []byte
	if !FileExist(chainsPath) {
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
	logger.Debug("reloading chains name map")
	loadChainsNameMap()
	logger.Debug("chains reload process done")
}

// initHotReload - initialize keys and chains file change detection
func initChainsHotReload() {
	if app.GlobalMeshConfig.ChainsHotReloadInterval <= 0 {
		logger.Info("skipping hot reload due to chains_hot_reload_interval is less or equal to 0")
		return
	}

	logger.Info(fmt.Sprintf("chains hot reload set to run every %d milliseconds", app.GlobalMeshConfig.ChainsHotReloadInterval))

	for {
		time.Sleep(time.Duration(app.GlobalMeshConfig.ChainsHotReloadInterval) * time.Millisecond)
		reloadChains()
	}
}

// ExecuteBlockchainHTTPRequest - run the non-native blockchain http request reusing chains http client.
func ExecuteBlockchainHTTPRequest(payload pocketTypes.Payload, chain pocketTypes.HostedBlockchain) (string, error, int) {
	data := payload.Data
	_url := strings.Trim(chain.URL, `/`)

	if len(payload.Path) > 0 {
		_url = _url + "/" + strings.Trim(payload.Path, `/`)
	}

	var m string
	if payload.Method == "" {
		m = pocketTypes.DEFAULTHTTPMETHOD
	} else {
		m = payload.Method
	}

	if app.GlobalMeshConfig.ChainRequestPathCleanup {
		_url = strings.Map(func(r rune) rune {
			if unicode.IsGraphic(r) {
				return r
			}
			return -1
		}, _url)
	}

	// generate an http request
	req, e1 := http.NewRequest(m, _url, bytes.NewBuffer([]byte(data)))
	if e1 != nil {
		logger.Error("unable to create blockchain request ", "err", e1.Error())
		return "", errors.New("unable to process blockchain request"), 500
	}
	if chain.BasicAuth.Username != "" {
		req.SetBasicAuth(chain.BasicAuth.Username, chain.BasicAuth.Password)
	}
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	// add headers if needed
	if len(payload.Headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range payload.Headers {
			req.Header.Set(k, v)
		}
	}

	// some users report lots of EOF due to connections trying to behind reused but net.Http fails to understand it.
	if app.GlobalMeshConfig.ChainDropConnections {
		req.Header.Set("Connection", "close")
		req.Close = true
	}

	// execute the request
	resp, e2 := chainsClient.Do(req)
	if e2 != nil {
		if os.IsTimeout(e2) {
			logStr := fmt.Sprintf(
				"request timeout for CHAIN=%s PATH=%s METHOD=%s",
				chain.ID, payload.Path, m,
			)
			logger.Error("blockchain call timeout: ", e2.Error())
			return "", errors.New(logStr), 504
		}

		logStr := fmt.Sprintf(
			"blockchain request for CHAIN=%s PATH=%s METHOD=%s",
			chain.ID, payload.Path, m,
		)

		if app.GlobalMeshConfig.LogChainRequest {
			logStr = logStr + " " + fmt.Sprintf("REQ=%s", data)
		}

		logger.Error(fmt.Sprintf("blockchain call error: %v", e2))
		return "", errors.New(logStr), 500
	}
	defer func(Body io.ReadCloser) {
		e3 := Body.Close()
		if e3 != nil {

		}
	}(resp.Body)

	// this log remains local.
	logStr := fmt.Sprintf(
		"executing blockchain request for CHAIN=%s CHAIN_URL=%s PATH=%s METHOD=%s STATUS=%d PARSE_URL=%s",
		chain.ID, chain.URL, payload.Path, m, resp.StatusCode, _url,
	)

	// read all bz
	body, e4 := ioutil.ReadAll(resp.Body)
	if e4 != nil {
		logger.Error("unable to process blockchain response payload ", " error ", e4.Error())
		return "", errors.New("unable to process blockchain response payload"), 500
	}

	if app.GlobalMeshConfig.JSONSortRelayResponses {
		body = []byte(sortJSONResponse(string(body)))
	}

	if app.GlobalMeshConfig.LogChainRequest {
		logStr = logStr + " " + fmt.Sprintf("REQ=%s", data)
	}

	if app.GlobalMeshConfig.LogChainResponse {
		logStr = logStr + " " + fmt.Sprintf("RES=%s", string(body))
	}

	if resp.StatusCode >= 400 {
		logger.Error(logStr)
	} else {
		logger.Debug(logStr)
	}

	return string(body), nil, resp.StatusCode
}

// GetChains - return current chains list.
func GetChains() *pocketTypes.HostedBlockchains {
	return chains
}
