// This package is for the RPC/REST API
package rpc

import (
	"github.com/julienschmidt/httprouter"
)

/*
"routes.go" is responsible for declaring all of the possible routes for the API.
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
