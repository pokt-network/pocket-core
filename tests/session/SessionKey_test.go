package session

import (
	"bytes"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSessionKey(t *testing.T) {
	// Start server instance
	go http.ListenAndServe(":"+config.GetInstance().Relayrpcport, shared.NewRouter(relay.RelayRoutes()))
	// @ Url
	u := "http://localhost:" + config.GetInstance().Relayrpcport + "/v1/dispatch/serve"
	// Create json string
	jsonString := []byte(`{"devid":"testing"}`)
	// Create post request
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonString))
	// Handle errors
	if err != nil {
		t.Errorf(err.Error())
	}
	// Set header of the request
	req.Header.Set("Content-Type","application/json")
	// Create a new http client
	client := &http.Client{}
	// Execute the request
	resp, err := client.Do(req)
	// Handle errors
	if err != nil {
		t.Errorf(err.Error())
	}
	// Deferred: close the body of the response
	defer resp.Body.Close()
	// Read the body from the response using ioutil
	body, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(body))
}
