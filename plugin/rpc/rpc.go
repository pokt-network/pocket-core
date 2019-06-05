// The relay forwarding plugin for rpc medium
package rpc

import (
	"bytes"
	"github.com/pokt-network/pocket-core/util"
	"io/ioutil"
	"net/http"
)

const (
	POST = "POST"
)

// "ExecuteHTTPRequest" takes in the raw json string and forwards it to the HTTP endpoint
func ExecuteHTTPRequest(payload []byte, u string, method string) (string, error) {
	if method == "" {
		method = POST
	}
	ur, err := util.URLProto(u)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, ur, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	if req != nil {
		if req.Body != nil {
			defer req.Body.Close()
		}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	if resp != nil {
		if resp.Body != nil {
			defer resp.Body.Close()
		}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
