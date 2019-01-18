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
	"testing"
)

func TestDispatchServe(t *testing.T) {
	// create arbitrary blockchains
	ethereum := node.Blockchain{Name: "ethereum", NetID: "1", Version: "1.0"}
	rinkeby := node.Blockchain{Name: "ethereum", NetID: "4", Version: "1.0"}
	bitcoin := node.Blockchain{Name: "bitcoin", NetID: "1", Version: "1.0"}
	bitcoinv1 := node.Blockchain{Name: "bitcoin", NetID: "1", Version: "1.1"}
	bitcoinCash := node.Blockchain{Name: "bitcoinCash", NetID: "1", Version: "1.0"}
	// create arbitrary nodes
	node1 := node.Node{
		GID:         "node1",
		IP:          "ip1",
		RelayPort:   "0",
		Blockchains: []node.Blockchain{ethereum, rinkeby, bitcoin}}
	node2 := node.Node{
		GID:         "node2",
		IP:          "ip2",
		RelayPort:   "0",
		Blockchains: []node.Blockchain{rinkeby, bitcoin, bitcoinv1}}
	node3 := node.Node{
		GID:         "node3",
		IP:          "ip3",
		RelayPort:   "0",
		Blockchains: []node.Blockchain{bitcoinCash, rinkeby, bitcoinv1}}
	// add them to dispatchPeers
	dp := node.GetDispatchPeers()
	dp.Add(node1)
	dp.Add(node2)
	dp.Add(node3)
	// add foo to the whitelist
	node.GetDWL().Add("foo")
	// json call string for dispatch serve
	requestJSON := []byte("{\"DevID\": \"foo\", \"Blockchains\": [{\"name\":\"ethereum\",\"netid\":\"1\",\"version\":\"1.0\"}]}")
	// start relay server
	go http.ListenAndServe(":"+config.GetInstance().RRPCPort, shared.NewRouter(relay.Routes()))
	// url for the POST request
	u := "http://localhost:" + config.GetInstance().RRPCPort + "/v1/dispatch/serve"
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(requestJSON))
	if err != nil {
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
	fmt.Println(string(body))
	var result map[string][]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("Unable to unmarshall json node response : " + err.Error())
	}
	expectedBody := map[string][]string{"ETHEREUMV1.0 | NetID 1": {"ip1:0"}}
	fmt.Println("EXPECTED BODY:", expectedBody)
	fmt.Println("RECEIVED BODY:", result)
	if !reflect.DeepEqual(result, expectedBody) {
		t.Fatalf("The resulting dispatchPeers is not as expected")
	}
}
