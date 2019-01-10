// This package is contains the handler functions needed for the Relay API
package relay

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/plugin/rpc-plugin"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"net/http"
	"os"
)

// "relay.go" defines API handlers that are under the 'relay' category within this file.

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
	relay := &Relay{}                               // create empty relay structure
	shared.PopulateModelFromParams(w, r, ps, relay) // populate the relay struct from params //TODO handle error for populate model from params (in all cases within codebase!)
	response, err := RouteRelay(*relay)             // route the relay to the correct chain
	if err != nil {
		// TODO handle error
	}
	shared.WriteJSONResponse(w, response) // relay the response
}

/*
"RelayReadInfo" handles a get request to localhost:<relay-port>/v1/relay/read call.
And provides the developers with an in-client reference to the API call
*/
func RelayReadInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "RelayRead", Relay{}, "Response from hosted chain")
	shared.WriteInfoResponse(w, info)
}

/*
 "RelayWrite" handles the localhost:<relay-port>/v1/relay/write call.
*/
func RelayWrite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
"RouteRelay" routes the relay to the specified hosted chain
*/
func RouteRelay(relay Relay) (string, error) {
	port := node.GetHostedChainPort(relay.Blockchain, relay.NetworkID, relay.Version)
	if port == "" {
		logs.NewLog("Not a supported blockchain", logs.ErrorLevel, logs.JSONLogFormat)
		return "Error: not a supported blockchain", nil // TODO custom error here
	}
	return rpc_plugin.ExecuteRequest([]byte(relay.Data),port)
}

// NOTE: This is for the centralized dispatcher of Pocket core mvp, may be removed for production
func ReportServiceNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	report := &Report{}
	shared.PopulateModelFromParams(w,r,ps,report)
	response, err := HandleReport(report)
	if err != nil {
		//TODO handle errors
	}
	shared.WriteJSONResponse(w,response)
}

// NOTE: This is for the centralized dispatcher of Pocket core mvp, may be removed for production
func HandleReport(report *Report) (string,error){
	f, err := os.OpenFile(_const.REPORTFILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	text, err:= json.Marshal(report)
	if err != nil {
		return "500 ERROR", err
	}
	if _, err = f.WriteString(string(text)+"\n"); err != nil {
		return "500 ERROR", err
	}
	return "Okay! The node has been successfully reported to our servers and will be reviewed! Thank you!", err
}

func ReportServiceNodeInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "ReportServiceNode", Report{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}
