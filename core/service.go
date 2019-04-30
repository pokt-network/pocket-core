package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

// "RouteRelay" routes the relay to the specified hosted chain
// This call handles REST and traditional JSON RPC
func RouteRelay(relay Relay) (string, error) {
	chain := GetHostedChains().GetChainFromBytes(relay.Blockchain)
	chainURL := chain.URL
	if relay.Path != nil && len(relay.Path) != 0 {
		resturl := string(relay.Path)
		strings.TrimSuffix(chainURL, "/")
		strings.TrimPrefix(resturl, "/")
		chainURL += "/" + resturl
	}
	return executeHTTPRequest(relay.Payload, chainURL, string(relay.Method))
}

// "executeHTTPRequest" takes in the raw json string and forwards it to the HTTP endpoint
func executeHTTPRequest(payload []byte, url string, method string) (string, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
