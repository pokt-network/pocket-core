// The relay forwarding plugin for rpc medium
package rpc

import (
	"bytes"
	"errors"
	"github.com/pokt-network/pocket-core/util"
	"io/ioutil"
	"net/http"
)

// "ExecuteRequest" takes in the raw json string and forwards it to the port
func ExecuteRequest(jsonStr []byte, host string, port string) (string, error) {
	url, err := util.URLProto(host + ":" + port)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("500: no response error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
