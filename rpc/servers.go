// This package is for the RPC/REST API
package rpc

import (
	"github.com/pocket_network/pocket-core/rpc/client"
	"github.com/pocket_network/pocket-core/rpc/relay"
	"log"
	"net/http"
)

func StartEndpoints() {

}

/*
"StartClientRPC" starts the client RPC/REST API server at a specific port.
 */
func StartClientRPC(port string) {
	// This starts the client RPC API.
	log.Fatal(http.ListenAndServe(":"+port, NewRouter(client.ClientRoutes())))
}

/*
"StartRelayRPC" starts the client RPC/REST API server at a specific port.
 */
func StartRelayRPC(port string) {
	// This starts the relay RPC API.
	log.Fatal(http.ListenAndServe(":"+port, NewRouter(relay.RelayRoutes())))
}
