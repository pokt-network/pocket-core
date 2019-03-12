// This package is for the RPC/REST API
package rpc

import (
	"log"
	"net/http"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "StartServers" executes the specified configuration for the client.
func StartServers() {
	if config.GlobalConfig().RRPC { // if flag set
		go StartRelayRPC(config.GlobalConfig().RRPCPort) // run the relay rpc in a goroutine
	}
}

// "startRelayRPC" starts the client RPC/REST API server at a specific port.
func StartRelayRPC(port string) {
	log.Fatal(http.ListenAndServe(":"+port, shared.Router(relay.Routes()))) // This starts the relay RPC API.
	logs.NewLog("Started relay server", logs.InfoLevel, logs.JSONLogFormat)
}
