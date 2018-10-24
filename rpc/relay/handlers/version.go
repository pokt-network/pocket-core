package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 "GetRelayAPIVersion" handles the localhost:<relay-port>/v1 call.
 */
func GetRelayAPIVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

