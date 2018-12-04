// This package is the native pocket core plugin for ethereum.
package pcp_ethereum

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

/*
"ExecuteRequest" takes in the raw json string and forwards it to the ethereum port
*/
func ExecuteRequest(jsonStr []byte, ethPort string) string {
	url := "http://localhost:" + ethPort                               // create a url for ethereum port
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr)) // call POST request to forward string
	req.Header.Set("Content-Type", "application/json")                 // specify json header
	resp, err := (&http.Client{}).Do(req)                              // execute request
	if err != nil {                                                    // handle error
		panic(err)
	}
	defer resp.Body.Close()              // close body after function completes
	body, _ := ioutil.ReadAll(resp.Body) // get the body from the response
	return string(body)                  // returns the response
}
