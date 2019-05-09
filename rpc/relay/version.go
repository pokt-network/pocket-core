package relay

import (
	"net/http"
	
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "Version" handles the localhost:<relay-port>/v1 call.
func Version(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, _const.RAPIVERSION, r.URL.Path, r.Host)
}
