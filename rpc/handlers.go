// This package is for the RPC/REST API
package rpc

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
)

/*
"handlers.go holds all of the handler functions for the specified API routes.
 */

func mockAPIFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	example := &Example{}
	populateModelFromParams(w, r, ps, example)
	writeResponse(w, example)
}

//Populates a model from the params
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
