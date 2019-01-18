// This package is shared between the different RPC packages
package shared

import "github.com/julienschmidt/httprouter"

// "NewRouter" creates a new httprouter from all of the routes and corresponding functions dealing with local calls.
func NewRouter(routes Routes) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	return router
}
