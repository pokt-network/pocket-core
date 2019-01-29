package relay

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/service"
)

// "Forward" handles the localhost:<relay-port>/v1/relaycall.
func Forward(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	relay := &service.Relay{}
	if err:=shared.PopulateModelFromParams(w, r, ps, relay); err!=nil{
		logs.NewLog(err.Error(),logs.ErrorLevel, logs.JSONLogFormat)
	}
	response, err := service.RouteRelay(*relay)
	if err != nil {
		logs.NewLog(err.Error(),logs.ErrorLevel, logs.JSONLogFormat)
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
	if err:=shared.PopulateModelFromParams(w, r, ps, report); err!=nil{
		logs.NewLog(err.Error(),logs.ErrorLevel, logs.JSONLogFormat)
	}
	response, err := service.HandleReport(report)
	if err != nil {
		logs.NewLog(err.Error(),logs.ErrorLevel, logs.JSONLogFormat)
	}
	shared.WriteJSONResponse(w, response)
}

// "ReportServiceNodeInfo" provides an in-client refrence to the api
func ReportServiceNodeInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "ReportServiceNode", service.Report{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}
