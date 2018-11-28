package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/net"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"io/ioutil"
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
	// hard code in some nodes
	var empty []string
	n1:= node.Node{"211057e8a7bbf340614b55fce0c481f3da8179b1",
	"","","","","","",empty}
	n2:= node.Node{"211057e8a7bbf340614b55fce0c481f3da8179b2",
		"","","","","","",empty}
	n3:= node.Node{"211057e8a7bbf340614b55fce0c481f3da8179b3",
		"","","","","","",empty}
	// add to peerList
	net.GetPeerList()
	net.AddNodePeerList(n1)
	net.AddNodePeerList(n2)
	net.AddNodePeerList(n3)
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
	body,_:=ioutil.ReadAll(resp.Body)
	fmt.Println("BODY: " + string(body))
	var data []node.Node
	err = json.Unmarshal(body,&data)
	fmt.Println(data)
	if err!=nil{
		t.Fatalf("Unable to unmarshall json node response 2: "+ err.Error())
	}
	if(data[0].GID!=n1.GID){
		t.Fatalf("Nodes are not in correct order")
	}
}
