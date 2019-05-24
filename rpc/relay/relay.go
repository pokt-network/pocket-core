package relay

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/service"
)

// "Relay" handles the localhost:<relay-port>/v1/relaycall.
func Relay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	relay := &service.Relay{}
	shared.SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	if err := shared.PopModel(w, r, ps, relay); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	if relay.Blockchain == "" || relay.NetworkID == "" || relay.DevID == "" || relay.Data == "" {
		shared.WriteErrorResponse(w, 400, "The request was not properly formatted")
		return
	}
	response, err := service.RouteRelay(*relay)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 500, err.Error())
		return
	}
	shared.WriteJSONResponse(w, response, r.URL.Path, r.Host) // relay the response
}

// "RelayInfo" handles a get request to localhost:<relay-port>/v1/relay call.
// And provides the developers with an in-client reference to the API call
func RelayInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.InfoStruct(r, "Relay", service.Relay{}, "Response from hosted chain")
	shared.WriteInfoResponse(w, info)
}

// "Report" is client side protection against a bad/faulty service node.
func Report(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	report := &service.Report{}
	if err := shared.PopModel(w, r, ps, report); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	if report.IP == "" || report.Message == "" {
		err := errors.New("empty field error")
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		shared.WriteErrorResponse(w, 401, err.Error())
		return
	}
	response, err := service.HandleReport(report)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
	shared.WriteJSONResponse(w, response, r.URL.Path, r.Host)
}

// "ReportInfo" provides an in-client refrence to the api
func ReportInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.InfoStruct(r, "Report", service.Report{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}
