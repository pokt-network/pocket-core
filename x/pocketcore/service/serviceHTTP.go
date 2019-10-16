package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// "executeHTTPRequest" takes in the raw json string and forwards it to the RPC endpoint
// todo improved http responses
func executeHTTPRequest(payload []byte, url string, method string) (string, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewHTTPStatusCodeError(resp.StatusCode)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
