// This package is shared between the different RPC packages
package shared

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
)

// Define all shared API handlers in this file.

/*
Unused mock api function for example.
 */
func mockAPIFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	example := &Example{}
	populateModelFromParams(w, r, ps, example)
	WriteResponse(w, example)
}

/*
Populate the model from the parameters of the POST call.
 */
func populateModelFromParams(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}
