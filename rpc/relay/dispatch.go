// This package is contains the handler functions needed for the Relay API
package relay

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pocket_network/pocket-core/rpc/shared"
	"io/ioutil"
	"log"
	"net/http"
)

// Define all API handlers that are under the 'dispatch' category within this file.

/*
 "DispatchOptions" handles the localhost:<relay-port>/v1/dispatch call.
 */
func DispatchOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "DispatchServe" handles the localhost:<relay-port>/v1/dispatch/serve call.
 */
func DispatchServe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", reqBody)
	//TODO unmarshal this into a model
}
