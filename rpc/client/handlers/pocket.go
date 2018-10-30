// This package is contains the handler functions needed for the Client API
package handlers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pocket_network/pocket-core/rpc/shared"
	"net/http"
)

// Define all API handlers that are under the 'pocket' category within this file.

/*
 "GetPocketBCInfo" handles the localhost:<client-port>/v1/pocket call.
 */
func GetPocketBCInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetLatestBlock" handles the localhost:<client-port>/v1/pocket/block call.
 */
func GetLatestBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetBlockByHash" handles the localhost:<client-port>/v1/pocket/block/hash call.
 */
func GetBlockByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetBlkTxCntByHash" handles the localhost:<client-port>/v1/pocket/block/hash/transaction_count call.
 */
func GetBlkTxCntByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetBlockByNum" handles the localhost:<client-port>/v1/pocket/block/number call.
 */
func GetBlockByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetBlkTxCntByNum" handles the localhost:<client-port>/v1/pocket/block/number/transaction_count call.
 */
func GetBlkTxCntByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetProtocolVersion" handles the localhost:<client-port>/v1/pocket/version_count call.
 */
func GetProtocolVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
