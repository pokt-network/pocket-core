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
		shared.Route{Name: "Version", Method: "GET", Path: "/v1", HandlerFunc: Version},
		shared.Route{Name: "GetRoutes", Method: "GET", Path: "/v1/routes", HandlerFunc: GetRoutes},
		shared.Route{Name: "Report", Method: "POST", Path: "/v1/report", HandlerFunc: Report},
		shared.Route{Name: "ReportInfo", Method: "GET", Path: "/v1/report", HandlerFunc: ReportInfo},
		shared.Route{Name: "Dispatch", Method: "POST", Path: "/v1/dispatch", HandlerFunc: Dispatch},
		shared.Route{Name: "DispatchInfo", Method: "GET", Path: "/v1/dispatch", HandlerFunc: DispatchInfo},
		shared.Route{Name: "Relay", Method: "POST", Path: "/v1/relay/", HandlerFunc: Relay},
		shared.Route{Name: "RelayInfo", Method: "GET", Path: "/v1/relay/", HandlerFunc: RelayInfo},
	}
	return routes
}

// "GetRoutes" handles the localhost:<relay-port>/routes call.
func GetRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var paths []string
	for _, v := range Routes() {
		if v.Method != "GET" {
			paths = append(paths, v.Path)
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	j, err := json.MarshalIndent(paths, "", "    ")
	if err != nil {
		logs.NewLog("Unable to marshal GetRoutes to JSON", logs.ErrorLevel, logs.JSONLogFormat)
	}
	shared.WriteRawJSONResponse(w, j)
}
