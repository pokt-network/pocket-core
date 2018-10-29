// This package is contains the handler functions needed for the Client API
package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Define all API handlers that are under the 'version' category within this file.

/*
 "getClientAPIVersion" handles the localhost:<client-port>/v1 call.
 */
func GetClientAPIVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
