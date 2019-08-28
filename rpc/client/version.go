package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "getClientAPIVersion" handles the localhost:<client-port>/v1 call.
func Version(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, _const.APIVERSION)
}
