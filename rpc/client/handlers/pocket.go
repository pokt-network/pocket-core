package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 "GetPocketBCInfo" handles the localhost:<client-port>/v1/pocket call.
 */
func GetPocketBCInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetLatestBlock" handles the localhost:<client-port>/v1/pocket/block call.
 */
func GetLatestBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetBlockByHash" handles the localhost:<client-port>/v1/pocket/block/hash call.
 */
func GetBlockByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetBlkTxCntByHash" handles the localhost:<client-port>/v1/pocket/block/hash/transaction_count call.
 */
func GetBlkTxCntByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetBlockByNum" handles the localhost:<client-port>/v1/pocket/block/number call.
 */
func GetBlockByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetBlkTxCntByNum" handles the localhost:<client-port>/v1/pocket/block/number/transaction_count call.
 */
func GetBlkTxCntByNum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetProtocolVersion" handles the localhost:<client-port>/v1/pocket/version_count call.
 */
func GetProtocolVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
