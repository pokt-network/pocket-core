package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetPocketBCInfo" handles the localhost:<client-port>/v1/pocket call.
func GetPocketBCInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetLatestBlock" handles the localhost:<client-port>/v1/pocket/block call.
func GetLatestBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetBlockByHash" handles the localhost:<client-port>/v1/pocket/block/hash call.
func GetBlockByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetBlkTxCntByHash" handles the localhost:<client-port>/v1/pocket/block/hash/transaction_count call.
func GetBlkTxCntByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetBlockByNum" handles the localhost:<client-port>/v1/pocket/block/number call.
func GetBlockByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetBlkTxCntByNum" handles the localhost:<client-port>/v1/pocket/block/number/transaction_count call.
func GetBlkTxCntByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetProtocolVersion" handles the localhost:<client-port>/v1/pocket/version_count call.
func GetProtocolVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}
