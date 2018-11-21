// This package is shared between the different RPC packages
package shared

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
)

// "handlers.go" defines shared API handlers in this file.

/*
Populate the model from the parameters of the POST call.
 */
func PopulateModelFromParams(_ http.ResponseWriter, r *http.Request, _ httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))	// get the request body
	if err != nil {													// handle error
		return err
	}
	if err := r.Body.Close(); err != nil {							// try to close the body
		return err													// handle error
	}
	if err := json.Unmarshal(body, model); err != nil {				// unmarshal the body into a model
		return err													// handle error
	}
	return nil														// return null pointer
}
