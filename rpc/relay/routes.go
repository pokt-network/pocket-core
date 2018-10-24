package relay

import (
	"github.com/pocket_network/pocket-core/rpc"
	"github.com/pocket_network/pocket-core/rpc/relay/handlers"
)

/*
"relayRoutes" is a function that returns all of the routes of the API.
 */
func RelayRoutes() rpc.Routes {
	routes := rpc.Routes{
		rpc.Route{"GetRelayAPIVersion", "POST", "/v1/", handlers.GetRelayAPIVersion},
		rpc.Route{"DispatchOptions", "POST", "/v1/dispatch/", handlers.DispatchOptions},
		rpc.Route{"DispatchServe", "POST", "/v1/dispatch/serve/", handlers.DispatchServe},
		rpc.Route{"RelayOptions", "POST", "/v1/relay/", handlers.RelayOptions},
		rpc.Route{"RelayRead", "POST", "/v1/relay/read/", handlers.RelayRead},
		rpc.Route{"RelayWrite", "POST", "/v1/relay/write/", handlers.RelayWrite},
	}
	return routes
}

