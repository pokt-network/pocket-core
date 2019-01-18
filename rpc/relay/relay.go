package relay

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/service"
)

// "Forward" handles the localhost:<relay-port>/v1/relaycall.
func Forward(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	relay := &service.Relay{}
	shared.PopulateModelFromParams(w, r, ps, relay) // TODO handle error for populate model from params (in all cases within codebase!)
	response, err := service.RouteRelay(*relay)
	if err != nil {
		// TODO handle error
	}
	shared.WriteJSONResponse(w, response) // relay the response
}

// "RelayReadInfo" handles a get request to localhost:<relay-port>/v1/relay call.
// And provides the developers with an in-client reference to the API call
func RelayReadInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "Forward", service.Relay{}, "Response from hosted chain")
	shared.WriteInfoResponse(w, info)
}

// DISCLAIMER: This is for the centralized dispatcher of Pocket core mvp, may be removed for production
// "ReportServiceNode" is client side protection against a bad/faulty service node.
func ReportServiceNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	report := &service.Report{}
	shared.PopulateModelFromParams(w, r, ps, report)
	response, err := service.HandleReport(report)
	if err != nil {
		// TODO handle errors
	}
	shared.WriteJSONResponse(w, response)
}

// "ReportServiceNodeInfo" provides an in-client refrence to the api
func ReportServiceNodeInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "ReportServiceNode", service.Report{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}
