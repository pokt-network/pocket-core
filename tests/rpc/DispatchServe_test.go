package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestDispatchServe(t *testing.T) {
	// create arbitrary blockchains
	ethereum := node.Blockchain{Name: "ethereum", NetID: "1", Version: "1.0"}
	rinkeby  := node.Blockchain{Name: "ethereum", NetID: "4", Version: "1.0"}
	bitcoin := node.Blockchain{Name: "bitcoin", NetID: "1", Version: "1.0"}
	bitcoinv1 := node.Blockchain{Name: "bitcoin", NetID: "1", Version: "1.1"}
	bitcoinCash := node.Blockchain{Name: "bitcoinCash", NetID: "1", Version: "1.0"}
	// create arbitrary nodes
	node1:= node.Node{
		GID:"node1",
		IP:"ip1",
		Blockchains:[]node.Blockchain{ethereum, rinkeby, bitcoin}}
	node2:= node.Node{
		GID:"node2",
		IP:"ip2",
		Blockchains:[]node.Blockchain{ethereum, bitcoin, bitcoinv1}}
	node3:= node.Node{
		GID:"node3",
		IP:"ip3",
		Blockchains:[]node.Blockchain{bitcoinCash, ethereum, bitcoinv1}}
	// add them to dispatchPeers
	node.NewDispatchPeer(node1)
	node.NewDispatchPeer(node2)
	node.NewDispatchPeer(node3)
	// get node lists
	// json call string for dispatch serve
	jsons:=[]byte("{\"DevID\": \"foo\", \"Blockchains\": [{\"name\":\"ethereum\",\"netid\":\"1\",\"version\":\"1.0\"}," +
		"{\"name\":\"bitcoin\",\"netid\":\"1\",\"version\":\"1.0\"}]}")
	// start relay server
	go http.ListenAndServe(":"+config.GetConfigInstance().Relayrpcport, shared.NewRouter(relay.RelayRoutes()))
	// url for the POST request
	u := "http://localhost:" + config.GetConfigInstance().Relayrpcport + "/v1/dispatch/serve"
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsons))
	if err != nil{
		t.Fatalf(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// create new http client
	client := &http.Client{}
	// Execute the request
	resp, err := client.Do(req)
	// Handle errors
	if err != nil {
		t.Errorf(err.Error())
	}
	// Deferred: close the body of the response
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string][]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("Unable to unmarshall json node response : " + err.Error())
	}
	fmt.Println(result)
	btcAns:= []string{"ip1","ip2"}
	btcKey := strings.ToUpper(bitcoin.Name)+"V"+bitcoin.Version+" | NetID "+bitcoin.NetID
	if !reflect.DeepEqual(result[btcKey],btcAns) {
		t.Fatalf("The resulting dispatchPeers is not as expected")
	}
}
