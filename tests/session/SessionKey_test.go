package session

import (
	"bytes"
	"encoding/json"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/session"
	"github.com/pokt-network/pocket-core/util"
	"net/http"
	"testing"
)

/* DEPRECATED
func TestSessionKey(t *testing.T) {
	// Start server instance
	go http.ListenAndServe(":"+config.GetConfigInstance().Relayrpcport, shared.NewRouter(relay.RelayRoutes()))
	// @ Url
	u := "http://localhost:" + config.GetConfigInstance().Relayrpcport + "/v1/dispatch/serve"
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
	response := new(shared.JSONResponse)
	json.NewDecoder(resp.Body).Decode(response)
	expectedKey := util.BytesToHex(session.GenerateSessionKey("testing"))
	t.Log("Expected Key: "+expectedKey )
	t.Log("Generated Key: "+response.Data)
	if response.Data!=expectedKey {
		t.Errorf("Response does not contain expected key...")
	}
}
*/

func TestSessionKey(t *testing.T) {
	// Start server instance
	go http.ListenAndServe(":"+config.GetConfigInstance().Relayrpcport, shared.NewRouter(relay.RelayRoutes()))
	// @ Url
	u := "http://localhost:" + config.GetConfigInstance().Relayrpcport + "/v1/dispatch/serve"
	// Create json string
	jsonString := []byte(`{"devid":"asdf"}`)
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
	response := new(shared.JSONResponse)
	json.NewDecoder(resp.Body).Decode(response)
	expectedKey := util.BytesToUInt32(session.GenerateSessionKey("asdf"))
	t.Log(expectedKey)
	t.Log("Generated Key: "+response.Data)
}
