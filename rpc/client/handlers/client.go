package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 "GetClient" handles the localhost:<client-port>/v1/client call.
 */
func GetClientInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetClientID" handles the localhost:<client-port>/v1/client/id call.
 */
func GetClientID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetClientVersion" handles the localhost:<client-port>/v1/client/version call.
 */
func GetClientVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetCliSyncStatus" handles the localhost:<client-port>/v1/client/syncing call.
 */
func GetCliSyncStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
