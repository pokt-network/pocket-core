package relay

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetRelayAPIVersion" handles the localhost:<relay-port>/v1 call.
func GetRelayAPIVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, _const.RAPIVERSION)
}
