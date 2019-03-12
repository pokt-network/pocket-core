// The relay forwarding plugin for rpc medium
package rpc

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// "ExecuteRequest" takes in the raw json string and forwards it to the port
func ExecuteRequest(jsonStr []byte, host string, port string) (string, error) {
	if !strings.Contains(host, "http") {
		host = "http://" + host
	}
	url := host + ":" + port
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
