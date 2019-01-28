package relay

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "Routes" is a function that returns all of the routes of the API.
func Routes() shared.Routes {
	routes := shared.Routes{
		shared.Route{Name: "GetRelayAPIVersion", Method: "GET", Path: "/v1", HandlerFunc: GetRelayAPIVersion},
		shared.Route{Name: "GetRoutes", Method: "GET", Path: "/v1/routes", HandlerFunc: GetRoutes},
		shared.Route{Name: "ReportServiceNode", Method: "POST", Path: "/v1/report", HandlerFunc: ReportServiceNode},
		shared.Route{Name: "ReportServiceNodeInfo", Method: "GET", Path: "/v1/report", HandlerFunc: ReportServiceNodeInfo},
		shared.Route{Name: "DispatchOptions", Method: "POST", Path: "/v1/dispatch", HandlerFunc: DispatchOptions},
		shared.Route{Name: "DispatchServe", Method: "POST", Path: "/v1/dispatch/serve", HandlerFunc: DispatchServe},
		shared.Route{Name: "DispatchServeInfo", Method: "GET", Path: "/v1/dispatch/serve", HandlerFunc: DispatchServeInfo},
		shared.Route{Name: "Relay", Method: "POST", Path: "/v1/relay/", HandlerFunc: Forward},
		shared.Route{Name: "RelayReadInfo", Method: "GET", Path: "/v1/relay/", HandlerFunc: RelayReadInfo},
	}
	return routes
}

// "GetRoutes" handles the localhost:<relay-port>/routes call.
func GetRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var paths []string
	for _, v := range Routes() {
		paths = append(paths, v.Path)
	}
	j, err := json.MarshalIndent(paths, "", "    ")
	if err != nil {
		logs.NewLog("Unable to marshal GetRoutes to JSON", logs.ErrorLevel, logs.JSONLogFormat)
	}
	shared.WriteRawJSONResponse(w, j)
}
