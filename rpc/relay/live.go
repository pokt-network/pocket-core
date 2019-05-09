package relay

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/db"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/service"
)

// "Register" handles the localhost:<relay-port>/v1/register call.
func Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// if not a dispatcher
	if !config.GlobalConfig().Dispatch {
		shared.WriteErrorResponse(w, 405, "Not a dispatch node")
		return
	}
	// if in deprecated mode
	if config.GlobalConfig().DisMode == _const.DISMODEDEPRECATED {
		shared.WriteErrorResponse(w, 410, "Deprecated, please upgrade software")
		return
	}
	n := node.Node{}
	// if cannot populate model
	if err := shared.PopModel(w, r, ps, &n); err != nil {
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	// if within white list
	if node.EnsureSNWL(node.SWL(), n.GID) {
		node.PeerList().Add(n)
		node.DispatchPeers().Add(n)
		if _, err := db.DB().Add(n); err != nil {
			fmt.Println(err.Error())
			shared.WriteErrorResponse(w, 500, "unable to write peer to database")
			return
		}
		// if within migrate mode
		if config.GlobalConfig().DisMode == _const.DISMODEMIGRATE {
			_, err := service.HandleReport(&service.Report{IP: n.IP, Message: "This node has not upgraded Pocket Core"})
			if err != nil {
				logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
			}
			shared.WriteJSONResponse(w, "WARNING: Pocket Core is now in the Migration Phase. Please upgrade your software as this version will soon be deprecated and not supported", r.Host)
			return
		}
		shared.WriteJSONResponse(w, "Success! Your node is now registered in the Pocket Network", r.Host)
		return
	}
	shared.WriteErrorResponse(w, 401, "Invalid credentials")
}

// "UnRegister" handles the localhost:<relay-port>/v1/unregister call.
func UnRegister(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	n := node.Node{}
	if err := shared.PopModel(w, r, ps, &n); err != nil {
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	node.PeerList().Remove(n)
	node.DispatchPeers().Delete(n)
	if _, err := db.DB().Remove(n); err != nil {
		shared.WriteErrorResponse(w, 500, "unable to remove peer from database")
		return
	}
	shared.WriteJSONResponse(w, "Success! Your node is now unregistered from the Pocket Network", r.Host)
}

func RegisterInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.InfoStruct(r, "Register", node.Node{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}

func UnRegisterInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.InfoStruct(r, "UnRegister", node.Node{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}
