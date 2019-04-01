package relay

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"net/http"
)

// "WhiteList" handles the localhost:<relay-port>/v1/whitelist call.
func WhiteList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	nd := &node.Node{}
	err := shared.PopModel(w, r, ps, nd)
	if err != nil {
		shared.WriteErrorResponse(w, 500, err.Error())
		return
	}
	if !node.EnsureWL(node.SNWL, nd.GID) {
		shared.WriteErrorResponse(w, 401, "invalid authentication")
		return
	}
	b, err := json.MarshalIndent(node.DevWL.ToSlice(), "", "")
	if err != nil {
		shared.WriteErrorResponse(w, 500, err.Error())
		return
	}
	shared.WriteRawJSONResponse(w, b)
}
