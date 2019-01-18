package shared

import (
	"github.com/julienschmidt/httprouter"
)

// The "Route" structure defines the generalization of an api route.
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

// "Routes" is a slice that holds all of the routes within one structure.
type Routes []Route
