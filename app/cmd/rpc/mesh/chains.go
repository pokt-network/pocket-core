package mesh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"io"
	"io/ioutil"
	log2 "log"
	"net/http"
	"os"
	"sync"
	"time"
)

// getChainsFilePath - return chains file path resolved by config.json
func getChainsFilePath() string {
	return app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ChainsName
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
	chainsPath := getChainsFilePath()
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

// executeBlockchainHTTPRequest - run the non-native blockchain http request reusing chains http client.
func executeBlockchainHTTPRequest(payload, url, userAgent string, basicAuth pocketTypes.BasicAuth, method string, headers map[string]string) (string, error) {
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
