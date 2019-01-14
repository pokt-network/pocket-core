// This package contains files for the Relay API
package relay

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "routes.go" defines all of the relay routes within this file.

/*
The "Route" structure defines the generalization of an api route.
*/
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

/*
"Routes" is a slice that holds all of the routes within one structure.
*/
type Routes []Route

/*
"relayRoutes" is a function that returns all of the routes of the API.
*/
func RelayRoutes() shared.Routes {
	routes := shared.Routes{
		shared.Route{"GetRelayAPIVersion", "GET", "/v1", GetRelayAPIVersion},
		shared.Route{"ReportServiceNode","POST","/v1/report", ReportServiceNode},
		shared.Route{"ReportServiceNodeInfo","GET","/v1/report",ReportServiceNodeInfo},
		shared.Route{"DispatchOptions", "POST", "/v1/dispatch", DispatchOptions},
		shared.Route{"DispatchServe", "POST", "/v1/dispatch/serve", DispatchServe},
		shared.Route{"DispatchServeInfo", "GET", "/v1/dispatch/serve", DispatchServeInfo},
		shared.Route{"RelayOptions", "POST", "/v1/relay", RelayOptions},
		shared.Route{"RelayRead", "POST", "/v1/relay/read", RelayRead},
		shared.Route{"RelayReadInfo", "GET", "/v1/relay/read", RelayReadInfo},
		shared.Route{"RelayWrite", "POST", "/v1/relay/write", RelayWrite},
	}
	return routes
}
