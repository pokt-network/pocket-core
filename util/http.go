package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Method int

const (
	POST Method = iota + 1
	GET
)

// "String" converts a Method Iota to a string
func (m Method) String() string {
	return [...]string{"GET", "POST"}[m]
}

// "RPCRequ" executes an RPC request
func RPCRequ(url string, data []byte, m Method) (string, error) {
	req, err := http.NewRequest(m.String(), url, bytes.NewBuffer(data))
	if err != nil {
		return "", errors.New("Cannot convert struct to json " + err.Error())
	}
	defer req.Body.Close()
	return rpcRequ(url, req)
}

// "StructRPCReq" sends an RPC request and returns the response
func StructRPCReq(url string, data interface{}, m Method) (string, error) {
	// convert structure to json
	j, err := json.Marshal(data)
	// handle error
	if err != nil {
		return "", errors.New("Cannot convert struct to json " + err.Error())
	}
	// create new post request
	req, err := http.NewRequest(m.String(), url, bytes.NewBuffer(j))
	// hanlde error
	if err != nil {
		return "", errors.New("Cannot create request " + err.Error())
	}
	defer req.Body.Close()
	return rpcRequ(url, req)
}

func rpcRequ(url string, req *http.Request) (string, error) {
	if req == nil {
		return "", errors.New("nil request passed to function")
	}
	if req.Body == nil {
		return "", errors.New("nil request passed to function")
	}
	defer req.Body.Close()
	// setup header for json data
	req.Header.Set("Content-Type", "application/json")
	// setup http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("unable to do request " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(body))
	}
	// get body of response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("unable to unmarshal response: " + err.Error())
	}
	defer resp.Body.Close()
	return string(body), nil
}
