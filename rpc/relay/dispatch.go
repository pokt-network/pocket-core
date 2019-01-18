// This package is contains the handler functions needed for the Relay API
package relay

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/dispatch"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// TODO fix dispatch serve example APIInformation
// "DispatchOptions" handles the localhost:<relay-port>/v1/dispatch call.
func DispatchOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "DispatchServe" handles the localhost:<relay-port>/v1/dispatch/serve call.
func DispatchServe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	d := &dispatch.Dispatch{}
	shared.PopulateModelFromParams(w, r, ps, d)
	shared.WriteRawJSONResponse(w, dispatch.Serve(d))
}

// "DispatchServeInfo" handles a get request to localhost:<relay-port>/v1/dispatch/serve call.
// And provides the developers with an in-client reference to the API call
func DispatchServeInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := shared.CreateInfoStruct(r, "Serve", dispatch.Dispatch{}, "list of service nodes")
	shared.WriteInfoResponse(w, info)
}
