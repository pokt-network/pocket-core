// This package is for the RPC/REST API
package rpc

import "github.com/julienschmidt/httprouter"

/*
"NewClientRouter" creates a new httprouter from all of the routes and corresponding functions dealing with local calls.
 */
func NewClientRouter(routes Routes) *httprouter.Router {
	// Declare a new http router.
	router := httprouter.New()
	// For each 'route' within 'routes'
	for _, route := range routes {
		// Setup the router for this route.
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	// Return the Router
	return router
}
