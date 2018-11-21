// This package is contains the handler functions needed for the Relay API
package relay

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"net/http"
)

// "version.go" defines API handlers that are under the 'version' category within this file.

/*
 "GetRelayAPIVersion" handles the localhost:<relay-port>/v1 call.
 */
func GetRelayAPIVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

