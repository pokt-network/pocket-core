// The relay forwarding plugin for rpc medium
package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// "ExecuteRequest" takes in the raw json string and forwards it to the port
func ExecuteRequest(jsonStr []byte, port string) (string, error) {
	fmt.Println("EXECUTING REQUEST")
	url := "http://localhost:" + port
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
