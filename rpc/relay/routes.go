// This package contains files for the Relay API
package relay

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pocket_network/pocket-core/rpc/shared"
)

/*
"relayRoutes" is a function that returns all of the routes of the API.
 */

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

func RelayRoutes() shared.Routes {
	routes := shared.Routes{
		shared.Route{"GetRelayAPIVersion", "POST", "/v1", GetRelayAPIVersion},
		shared.Route{"DispatchOptions", "POST", "/v1/dispatch", DispatchOptions},
		shared.Route{"DispatchServe", "POST", "/v1/dispatch/serve", DispatchServe},
		shared.Route{"RelayOptions", "POST", "/v1/relay", RelayOptions},
		shared.Route{"RelayRead", "POST", "/v1/relay/read", RelayRead},
		shared.Route{"RelayWrite", "POST", "/v1/relay/write", RelayWrite},
	}
	return routes
}
