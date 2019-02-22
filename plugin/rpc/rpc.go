// The relay forwarding plugin for rpc medium
package rpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// "ExecuteRequest" takes in the raw json string and forwards it to the port
func ExecuteRequest(jsonStr []byte, host string,port string) (string, error) {
	url := host+":"+ port
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
