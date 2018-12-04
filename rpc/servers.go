// This package is for the RPC/REST API
package rpc

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/client"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"log"
	"net/http"
)

// "servers.go" defines RPC/REST API serving functions within this file.

/*
"RunAPIEndpoints" executes the specified configuration for the client.
*/
func RunAPIEndpoints() {
	if config.GetConfigInstance().Clientrpc { // if flag set
		go StartClientRPC(config.GetConfigInstance().Clientrpcport) // run the client rpc in a goroutine
	}
	if config.GetConfigInstance().Relayrpc { // if flag set
		go StartRelayRPC(config.GetConfigInstance().Relayrpcport) // run the relay rpc in a goroutine
	}
}

/*
"startClientRPC" starts the client RPC/REST API server at a specific port.
*/
func StartClientRPC(port string) {
	log.Fatal(http.ListenAndServe(":"+port, shared.NewRouter(client.ClientRoutes()))) // This starts the client RPC API.
	logs.NewLog("Started client server", logs.InfoLevel, logs.JSONLogFormat)
}

/*
"startRelayRPC" starts the client RPC/REST API server at a specific port.
*/
func StartRelayRPC(port string) {
	log.Fatal(http.ListenAndServe(":"+port, shared.NewRouter(relay.RelayRoutes()))) // This starts the relay RPC API.
	logs.NewLog("Started relay server", logs.InfoLevel, logs.JSONLogFormat)
}
