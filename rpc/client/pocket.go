package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "BCInfo" handles the localhost:<client-port>/v1/pocket call.
func BCInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "LatestBlock" handles the localhost:<client-port>/v1/pocket/block call.
func LatestBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "BlockByHash" handles the localhost:<client-port>/v1/pocket/block/hash call.
func BlockByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "BlkTXCountByHash" handles the localhost:<client-port>/v1/pocket/block/hash/transaction_count call.
func BlkTXCountByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "BlkByNum" handles the localhost:<client-port>/v1/pocket/block/number call.
func BlkByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "BlkCntByNum" handles the localhost:<client-port>/v1/pocket/block/number/transaction_count call.
func BlkCntByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "ProtVersion" handles the localhost:<client-port>/v1/pocket/version_count call.
func ProtVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}
