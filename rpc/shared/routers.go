// This package is shared between the different RPC packages
package shared

import "github.com/julienschmidt/httprouter"

// "Router" creates a new httprouter from all of the routes and corresponding functions dealing with local calls.
func Router(routes Routes) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	return router
}
