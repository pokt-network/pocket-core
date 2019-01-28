package client

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "Routes" is a function that returns all of the routes of the API.
func Routes() shared.Routes {
	routes := shared.Routes{
		shared.Route{Name: "Routes", Method: "GET", Path: "/v1/routes", HandlerFunc: GetRoutes},
		shared.Route{Name: "Register", Method: "POST", Path: "/v1/register", HandlerFunc: Register},
		shared.Route{Name: "UnRegister", Method: "POST", Path: "/v1/unregister", HandlerFunc: UnRegister},
		shared.Route{Name: "RegisterInfo", Method: "GET", Path: "/v1/register", HandlerFunc: RegisterInfo},
		shared.Route{Name: "UnRegisterInfo", Method: "GET", Path: "/v1/unregister", HandlerFunc: UnRegisterInfo},
		shared.Route{Name: "GetClientAPIVersion", Method: "POST", Path: "/v1", HandlerFunc: GetClientAPIVersion},
		shared.Route{Name: "GetAccount", Method: "POST", Path: "/v1/account", HandlerFunc: GetAccount},
		shared.Route{Name: "IsAccountActive", Method: "POST", Path: "/v1/account/active", HandlerFunc: IsAccountActive},
		shared.Route{Name: "GetAccountBalance", Method: "POST", Path: "/v1/account/balance", HandlerFunc: GetAccountBalance},
		shared.Route{Name: "GetDateJoined", Method: "POST", Path: "/v1/account/joined", HandlerFunc: GetDateJoined},
		shared.Route{Name: "GetAccountKarma", Method: "POST", Path: "/v1/account/karma", HandlerFunc: GetAccountKarma},
		shared.Route{Name: "GetlastActive", Method: "POST", Path: "/v1/account/last_active", HandlerFunc: GetLastActive},
		shared.Route{Name: "GetAccTxCount", Method: "POST", Path: "/v1/account/transaction_count", HandlerFunc: GetAccTxCount},
		shared.Route{Name: "GetAccSessCount", Method: "POST", Path: "/v1/account/session_count", HandlerFunc: GetAccSessCount},
		shared.Route{Name: "GetAccStatus", Method: "POST", Path: "/v1/account/status", HandlerFunc: GetAccStatus},
		shared.Route{Name: "GetClientInfo", Method: "POST", Path: "/v1/client", HandlerFunc: GetClientInfo},
		shared.Route{Name: "GetClientID", Method: "POST", Path: "/v1/client/id", HandlerFunc: GetClientID},
		shared.Route{Name: "GetClientVersion", Method: "POST", Path: "/v1/client/version", HandlerFunc: GetClientVersion},
		shared.Route{Name: "GetCliSyncStatus", Method: "POST", Path: "/v1/client/syncing", HandlerFunc: GetCliSyncStatus},
		shared.Route{Name: "GetNetworkInfo", Method: "POST", Path: "/v1/network", HandlerFunc: GetNetworkInfo},
		shared.Route{Name: "GetNetworkID", Method: "POST", Path: "/v1/network/id", HandlerFunc: GetNetworkID},
		shared.Route{Name: "GetPeerCount", Method: "POST", Path: "/v1/network/peer_count", HandlerFunc: GetPeerCount},
		shared.Route{Name: "GetPeerList", Method: "POST", Path: "/v1/network/peer_list", HandlerFunc: GetPeerList},
		shared.Route{Name: "GetPeers", Method: "POST", Path: "/v1/network/peers", HandlerFunc: GetPeers},
		shared.Route{Name: "GetPersonalInfo", Method: "POST", Path: "/v1/personal", HandlerFunc: GetPersonalInfo},
		shared.Route{Name: "ListAccounts", Method: "POST", Path: "/v1/personal/list_accounts", HandlerFunc: ListAccounts},
		shared.Route{Name: "PersonalNetOptions", Method: "POST", Path: "/v1/personal/network", HandlerFunc: PersonalNetOptions},
		shared.Route{Name: "EnterNetwork", Method: "POST", Path: "/v1/personal/network/enter", HandlerFunc: EnterNetwork},
		shared.Route{Name: "ExitNetwork", Method: "POST", Path: "/v1/personal/network/exit", HandlerFunc: ExitNetwork},
		shared.Route{Name: "GetPrimaryAddr", Method: "POST", Path: "/v1/personal/primary_address", HandlerFunc: GetPrimaryAddr},
		shared.Route{Name: "SendPOKT", Method: "POST", Path: "/v1/personal/send", HandlerFunc: SendPOKT},
		shared.Route{Name: "SendPOKTRaw", Method: "POST", Path: "/v1/personal/send/raw", HandlerFunc: SendPOKTRaw},
		shared.Route{Name: "Sign", Method: "POST", Path: "/v1/personal/sign", HandlerFunc: Sign},
		shared.Route{Name: "StakeOptions", Method: "POST", Path: "/v1/personal/stake", HandlerFunc: StakeOptions},
		shared.Route{Name: "Stake", Method: "POST", Path: "/v1/personal/stake/add", HandlerFunc: Stake},
		shared.Route{Name: "UnStake", Method: "POST", Path: "/v1/personal/stake/remove", HandlerFunc: UnStake},
		shared.Route{Name: "GetPocketBCInfo", Method: "POST", Path: "/v1/pocket", HandlerFunc: GetPocketBCInfo},
		shared.Route{Name: "GetLatestBlock", Method: "POST", Path: "/v1/pocket/block", HandlerFunc: GetLatestBlock},
		shared.Route{Name: "GetBlockByHash", Method: "POST", Path: "/v1/pocket/block/hash", HandlerFunc: GetBlockByHash},
		shared.Route{Name: "GetBlkTxCntByHash", Method: "POST", Path: "/v1/pocket/block/hash/transaction_count", HandlerFunc: GetBlkTxCntByHash},
		shared.Route{Name: "GetBlockByNum", Method: "POST", Path: "/v1/pocket/block/number", HandlerFunc: GetBlockByNum},
		shared.Route{Name: "GetBlkTxCntByNum", Method: "POST", Path: "/v1/pocket/block/number/transaction_count", HandlerFunc: GetBlkTxCntByNum},
		shared.Route{Name: "GetProtocolVersion", Method: "POST", Path: "/v1/pocket/version", HandlerFunc: GetProtocolVersion},
		shared.Route{Name: "TxOptions", Method: "POST", Path: "/v1/pocket/transaction", HandlerFunc: TxOptions},
		shared.Route{Name: "GetTxByHash", Method: "POST", Path: "/v1/pocket/transaction/hash", HandlerFunc: GetTxByHash},
	}
	return routes
}

// "GetRoutes" handles the localhost:<relay-port>/routes call.
func GetRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var paths []string
	for _, v := range Routes() {
		paths = append(paths, v.Path)
	}
	j, err := json.MarshalIndent(paths, "", "    ")
	if err != nil {
		logs.NewLog("Unable to marshal GetRoutes to JSON", logs.ErrorLevel, logs.JSONLogFormat)
	}
	shared.WriteRawJSONResponse(w, j)
}
