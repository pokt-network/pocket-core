// This package is for the RPC/REST API
package rpc

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/rpc/client"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"log"
	"net/http"
)

// Define RPC/REST API serving functions within this file.

/*
"RunAPIEndpoints" executes the specified configuration for the client.
 */
func RunAPIEndpoints() {
	if config.GetInstance().Clientrpc {
		go startClientRPC(config.GetInstance().Clientrpcport)
	}
	if config.GetInstance().Relayrpc {
		startRelayRPC(config.GetInstance().Relayrpcport) // TODO convert to go routine
	}
}

/*
"startClientRPC" starts the client RPC/REST API server at a specific port.
 */
func startClientRPC(port string) {
	// This starts the client RPC API.
	log.Fatal(http.ListenAndServe(":"+port, shared.NewRouter(client.ClientRoutes())))
}

/*
"startRelayRPC" starts the client RPC/REST API server at a specific port.
 */
func startRelayRPC(port string) {
	// This starts the relay RPC API.
	log.Fatal(http.ListenAndServe(":"+port, shared.NewRouter(relay.RelayRoutes())))
}
