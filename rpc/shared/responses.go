// This package is shared between the different RPC packages
package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// "responses.go is the interface for the types of responses for the API.

/*
"JSONResponse" is a metadata and data response in JSON format.
 */
type JSONResponse struct {
	Data interface{} `json:"data"`
}

/*
"JSONErrorResponse" is an error response in JSON format.
 */
type JSONErrorResponse struct {
	Error *APIError `json:"error"`
}

/*
"APIError" is an error feedback structure containing a title and a status.
 */
type APIError struct {
	Status int    `json:"error"`
	Title  string `json:"title"`
}

/*
"WriteResponse" writes a normal JSON response.
 */
func WriteResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&JSONResponse{Data: m}); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
/*
"WriteInfo" provides useful information about the api URL when get is called
 */
 // TODO implement and test
func WriteInfoResponse(w http.ResponseWriter, information APIReference) {
	b, err := json.Marshal(information)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _,err := io.WriteString(w,string(b)); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
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
