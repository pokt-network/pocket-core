// This package is contains the handler functions needed for the Relay API
package relay

import (
	"net/http"
	
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "Dispatch" handles the localhost:<relay-port>/v1/dispatch/serve call.
func Dispatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "DispatchInfo" handles a get request to localhost:<relay-port>/v1/dispatch/serve call.
// And provides the developers with an in-client reference to the API call
func DispatchInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}
