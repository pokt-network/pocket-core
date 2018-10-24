package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 "DispatchOptions" handles the localhost:<relay-port>/v1/dispatch call.
 */
func DispatchOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "DispatchServe" handles the localhost:<relay-port>/v1/dispatch/serve call.
 */
func DispatchServe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
