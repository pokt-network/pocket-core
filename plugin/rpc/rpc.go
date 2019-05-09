// The relay forwarding plugin for rpc medium
package rpc

import (
	"bytes"
	"github.com/pokt-network/pocket-core/util"
	"io/ioutil"
	"net/http"
	"net/url"
)

// "ExecuteRequest" takes in the raw json string and forwards it to the port
func ExecuteRequest(jsonStr []byte, u *url.URL) (string, error) {
	ur, err := util.URLProto(u.String() + u.Path)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", ur, bytes.NewBuffer(jsonStr))
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
