package integration

import (
	"fmt"
	"testing"
	
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/util"
)

/*
Integration Testing Assumptions:
1) Dispatcher is hosting a testrpc instance that is labeled as (Blockchain: 0 | NetworkID: 0 | Version: 0)
2) Dispatcher has white listed DEVID1 (Dev) and GID1 (SN)
 */

func TestDispatch(t *testing.T) {
	blockchain := node.Blockchain{Name: "ethereum", NetID: "0", Version: "0"}
	c := node.Chains()
	c[blockchain]=node.HostedChain{Blockchain: blockchain, Port: "8080", Medium: "rpc"}
	
	const ethereumRequest = "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":67}"
	const relayRequest = "{\"blockchain\":\"ethereum\", \"networkid\":\"0\", \"version\":\"0\", \"data\":\"" + ethereumRequest + "\", \"devid\":\"DEVID1\"}"
	fmt.Println(relayRequest)
	resp, _:= util.RPCRequest("http://"+config.GlobalConfig().DisIP+":"+config.GlobalConfig().DisRPort+"/v1/relay", relayRequest, util.POST)
	fmt.Println(resp)
}
