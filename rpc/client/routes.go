package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
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
		shared.Route{Name: "Version", Method: "GET", Path: "/v1", HandlerFunc: Version},
		// shared.Route{Name: "Account", Method: "POST", Path: "/v1/account", HandlerFunc: Account},
		// shared.Route{Name: "IsAccountActive", Method: "POST", Path: "/v1/account/active", HandlerFunc: IsAccountActive},
		// shared.Route{Name: "Balance", Method: "POST", Path: "/v1/account/balance", HandlerFunc: Balance},
		// shared.Route{Name: "DateJoined", Method: "POST", Path: "/v1/account/joined", HandlerFunc: DateJoined},
		// shared.Route{Name: "AcountKarma", Method: "POST", Path: "/v1/account/karma", HandlerFunc: AcountKarma},
		// shared.Route{Name: "lastActive", Method: "POST", Path: "/v1/account/last_active", HandlerFunc: LastActive},
		// shared.Route{Name: "AcctTXCount", Method: "POST", Path: "/v1/account/transaction_count", HandlerFunc: AcctTXCount},
		// shared.Route{Name: "AccSessCount", Method: "POST", Path: "/v1/account/session_count", HandlerFunc: AccSessCount},
		// shared.Route{Name: "AccStatus", Method: "POST", Path: "/v1/account/status", HandlerFunc: AccStatus},
		// shared.Route{Name: "CliInfo", Method: "POST", Path: "/v1/client", HandlerFunc: CliInfo},
		// shared.Route{Name: "CliID", Method: "POST", Path: "/v1/client/id", HandlerFunc: CliID},
		// shared.Route{Name: "CliVersion", Method: "POST", Path: "/v1/client/version", HandlerFunc: CliVersion},
		// shared.Route{Name: "CliSyncStatus", Method: "POST", Path: "/v1/client/syncing", HandlerFunc: CliSyncStatus},
		// shared.Route{Name: "NetInfo", Method: "POST", Path: "/v1/network", HandlerFunc: NetInfo},
		// shared.Route{Name: "NetID", Method: "POST", Path: "/v1/network/id", HandlerFunc: NetID},
		// shared.Route{Name: "PeerCount", Method: "POST", Path: "/v1/network/peer_count", HandlerFunc: PeerCount},
		// shared.Route{Name: "PeerList", Method: "POST", Path: "/v1/network/peer_list", HandlerFunc: PeerList},
		// shared.Route{Name: "PL", Method: "POST", Path: "/v1/network/peers", HandlerFunc: PL},
		// shared.Route{Name: "PersonalInfo", Method: "POST", Path: "/v1/personal", HandlerFunc: PersonalInfo},
		// shared.Route{Name: "Accounts", Method: "POST", Path: "/v1/personal/list_accounts", HandlerFunc: Accounts},
		// shared.Route{Name: "EnterNetwork", Method: "POST", Path: "/v1/personal/network/enter", HandlerFunc: EnterNetwork},
		// shared.Route{Name: "ExitNetwork", Method: "POST", Path: "/v1/personal/network/exit", HandlerFunc: ExitNetwork},
		// shared.Route{Name: "PrimaryAddr", Method: "POST", Path: "/v1/personal/primary_address", HandlerFunc: PrimaryAddr},
		// shared.Route{Name: "SendPOKT", Method: "POST", Path: "/v1/personal/send", HandlerFunc: SendPOKT},
		// shared.Route{Name: "SendPOKTRaw", Method: "POST", Path: "/v1/personal/send/raw", HandlerFunc: SendPOKTRaw},
		// shared.Route{Name: "Sign", Method: "POST", Path: "/v1/personal/sign", HandlerFunc: Sign},
		// shared.Route{Name: "Stake", Method: "POST", Path: "/v1/personal/stake/add", HandlerFunc: Stake},
		// shared.Route{Name: "UnStake", Method: "POST", Path: "/v1/personal/stake/remove", HandlerFunc: UnStake},
		// shared.Route{Name: "BCInfo", Method: "POST", Path: "/v1/pocket", HandlerFunc: BCInfo},
		// shared.Route{Name: "LatestBlock", Method: "POST", Path: "/v1/pocket/block", HandlerFunc: LatestBlock},
		// shared.Route{Name: "BlockByHash", Method: "POST", Path: "/v1/pocket/block/hash", HandlerFunc: BlockByHash},
		// shared.Route{Name: "BlkTXCountByHash", Method: "POST", Path: "/v1/pocket/block/hash/transaction_count", HandlerFunc: BlkTXCountByHash},
		// shared.Route{Name: "BlkByNum", Method: "POST", Path: "/v1/pocket/block/number", HandlerFunc: BlkByNum},
		// shared.Route{Name: "BlkCntByNum", Method: "POST", Path: "/v1/pocket/block/number/transaction_count", HandlerFunc: BlkCntByNum},
		// shared.Route{Name: "ProtVersion", Method: "POST", Path: "/v1/pocket/version", HandlerFunc: ProtVersion},
		// shared.Route{Name: "TxByHash", Method: "POST", Path: "/v1/pocket/transaction/hash", HandlerFunc: TxByHash},
	}
	return routes
}

// "GetRoutes" handles the localhost:<client-port>/routes call.
func GetRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.GetRoutes(w, r, ps, Routes())
}
