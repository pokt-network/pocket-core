// This package is contains the handler functions needed for the Relay API
package handlers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pocket_network/pocket-core/rpc/shared"
	"net/http"
)

// Define all API handlers that are under the 'relay' category within this file.

/*
 "RelayOptions" handles the localhost:<relay-port>/v1/relay call.
 */
func RelayOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "RelayRead" handles the localhost:<relay-port>/v1/relay/read call.
 */
func RelayRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "RelayWrite" handles the localhost:<relay-port>/v1/relay/write call.
 */
func RelayWrite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
