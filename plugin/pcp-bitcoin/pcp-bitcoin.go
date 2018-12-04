// This package is the native pocket core plugin for bitcoin.
package pcp_bitcoin

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

/*
"ExecuteRequest" takes in the raw json string and forwards it to the bitcoin port
*/
func ExecuteRequest(jsonStr []byte, btcport string) string {
	url := "http://localhost:" + btcport                               // create a url for btc port
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
