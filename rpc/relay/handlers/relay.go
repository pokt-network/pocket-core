package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 "RelayOptions" handles the localhost:<relay-port>/v1/relay call.
 */
func RelayOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "RelayRead" handles the localhost:<relay-port>/v1/relay/read call.
 */
func RelayRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "RelayWrite" handles the localhost:<relay-port>/v1/relay/write call.
 */
func RelayWrite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
