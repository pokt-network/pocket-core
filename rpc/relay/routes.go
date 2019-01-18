package relay

import (
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "Routes" is a function that returns all of the routes of the API.
func Routes() shared.Routes {
	routes := shared.Routes{
		shared.Route{Name: "GetRelayAPIVersion", Method: "GET", Path: "/v1", HandlerFunc: GetRelayAPIVersion},
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
