// This package is shared between the different RPC packages
package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// "responses.go is the interface for the types of responses for the API.

/*
"WriteResponse" writes a normal JSON response.
 */
func WriteResponse(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	b,err := json.MarshalIndent(&JSONResponse{m},"","\t");
	if err!= nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	} else{
		w.Write(b)
	}
}
/*
"WriteInfo" provides useful information about the api URL when get is called
 */
func WriteInfoResponse(w http.ResponseWriter, information APIReference) {
	b,err := json.MarshalIndent(information,"","\t");
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err!= nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	} else{
		w.Write(b)
	}
}

/*
"WriteErrorResponse" writes an error JSON response.
 */
func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.
		NewEncoder(w).
		Encode(&JSONErrorResponse{Error: &APIError{Status: errorCode, Title: errorMsg}})
}
