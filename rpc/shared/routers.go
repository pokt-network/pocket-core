// This package is shared between the different RPC packages
package shared

import "github.com/julienschmidt/httprouter"

// "routers.go" handles the shared router call

/*
"NewRouter" creates a new httprouter from all of the routes and corresponding functions dealing with local calls.
 */
func NewRouter(routes Routes) *httprouter.Router {
	router := httprouter.New()										// Declare a new http router.
	for _, route := range routes {									// For each 'route' within 'routes'
		router.Handle(route.Method, route.Path, route.HandlerFunc) 	// Setup the router for this route.
	}																// Return the Router
	return router
}

