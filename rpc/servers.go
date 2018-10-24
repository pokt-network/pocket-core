// This package is for the RPC/REST API
package rpc

import (
	"log"
	"net/http"
)

/*
TODO need separate endpoints for custom facing (relay) vs. local facing APIs
 */
func StartEndpoints() {

}

/*
"StartRPC" starts an RPC/REST API server at a specific port.
 */
func StartRPC(port string) {
	log.Fatal(http.ListenAndServe(":"+port, NewClientRouter(AllRoutes())))
}
