package rpc

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
Unit test for APIReference
 */
func TestApiReference(t *testing.T) {
	// Start server instance
	go http.ListenAndServe(":"+config.GetConfigInstance().Relayrpcport, shared.NewRouter(relay.RelayRoutes()))
	// @ Url
	u := "http://localhost:" + config.GetConfigInstance().Relayrpcport + "/v1/dispatch/serve"
	// Send get request
	resp, err := http.Get(u)
	if err != nil {
		t.Errorf("Unable to get request at " + u + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			t.Errorf(err2.Error())
		}
		bodyString := string(bodyBytes)
		// Log response
		t.Log(bodyString)
	} else {
		t.Errorf("Failed at " + resp.Status)
	}
}
