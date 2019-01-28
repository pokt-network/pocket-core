package relay

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetRelayAPIVersion" handles the localhost:<relay-port>/v1 call.
func GetRelayAPIVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, _const.RAPIVERSION)
}

// "GetRoutes" handles the localhost:<relay-port>/routes call.
func GetRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// expose routes in json format todo
	shared.WriteResponse(w, _const.RAPIVERSION)
}