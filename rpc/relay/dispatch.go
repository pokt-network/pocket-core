// This package is contains the handler functions needed for the Relay API
package relay

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/dispatch"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "Dispatch" handles the localhost:<relay-port>/v1/dispatch/serve call.
func Dispatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	d := &dispatch.Dispatch{}
	if err := shared.PopModel(w, r, ps, d); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	if d.DevID == "" || len(d.Blockchains) == 0 {
		shared.WriteErrorResponse(w, 400, "Request was not formatted properly")
	}
	res, err, code := dispatch.Serve(d)
	if err != nil {
		shared.WriteErrorResponse(w, code, err.Error())
	}
	shared.WriteRawJSONResponse(w, res, r.Host)
}

// "DispatchInfo" handles a get request to localhost:<relay-port>/v1/dispatch/serve call.
// And provides the developers with an in-client reference to the API call
func DispatchInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := shared.InfoStruct(r, "Dispatch", dispatch.Dispatch{}, "zero or more service nodes")
	shared.WriteInfoResponse(w, info)
}
