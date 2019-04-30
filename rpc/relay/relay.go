package relay

import (
	"github.com/pokt-network/pocket-core/core"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// TODO E2E test on all rpc calls
// "Relay" handles the localhost:<relay-port>/v1/relaycall.
func Relay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	relay := &core.Relay{}
	if err := shared.PopModel(w, r, ps, relay); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	if err := relay.ErrorCheck(); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	response, err := core.RouteRelay(*relay)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 500, err.Error())
		return
	}
	shared.WriteJSONResponse(w, response)
}

// "RelayInfo" handles a get request to localhost:<relay-port>/v1/relay call.
// And provides the developers with an in-client reference to the API call
func RelayInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.InfoStruct(r, "Relay", core.Relay{}, "Response from hosted chain")
	shared.WriteInfoResponse(w, info)
}
