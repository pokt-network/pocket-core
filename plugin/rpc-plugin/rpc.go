// The relay forwarding plugin for rpc medium
package rpc_plugin

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

/*
"ExecuteRequest" takes in the raw json string and forwards it to the port
*/
func ExecuteRequest(jsonStr []byte, port string) (string, error) {
	url := "http://localhost:" + port                                  			// create a url for the port
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr)) 	// call POST request to forward string
	req.Header.Set("Content-Type", "application/json")               // specify json header
	resp, err := (&http.Client{}).Do(req)                              			// execute request
	if err != nil {                                                    			// handle error
		return "", err
	}
	defer resp.Body.Close()              										// close body after function completes
	body, _ := ioutil.ReadAll(resp.Body) 										// get the body from the response
	return string(body), nil             										// returns the response
}
