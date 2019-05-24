// This package is shared between the different RPC packages
package shared

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/logs"
	"net/http"
)

// "WriteJSONResponse" writes a JSON response.
func WriteJSONResponse(w http.ResponseWriter, m string, path, ip string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		if err := logs.NewRPCLog(false, path, ip, err.Error()); err != nil {
			logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		}
	} else {
		if err := logs.NewRPCLog(false, path, ip, string(m)); err != nil {
			logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		}
		w.Write(b)
	}
}

// "WriteRawJSON" writes a byte array.
func WriteRawJSONResponse(w http.ResponseWriter, b []byte, path, ip string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if err := logs.NewRPCLog(false, path, ip, string(b)); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
	}
	w.Write(b)
}

// "WriteInfo" provides useful information about the api URL when get is called
func WriteInfoResponse(w http.ResponseWriter, information APIReference) {
	b, err := json.MarshalIndent(information, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	} else {
		w.Write(b)
	}
}

// "WriteErrorResponse" writes an error JSON response.
func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&JSONErrorResponse{Error: &APIError{Status: errorCode, Title: errorMsg}})
}
