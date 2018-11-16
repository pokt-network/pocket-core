// This package contains files for the Client API
package client

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// Define all Client API routes in this file.

/*
The "Route" structure defines the generalization of an api route.
 */
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

/*
"Routes" is a slice that holds all of the routes within one structure.
 */
type Routes []Route

/*
"clientRoutes" is a function that returns all of the routes of the API.
 */
func ClientRoutes() shared.Routes {
	routes := shared.Routes{
		shared.Route{"GetClientAPIVersion", "POST", "/v1", GetClientAPIVersion},
		shared.Route{"GetAccount", "POST", "/v1/account", GetAccount},
		shared.Route{"IsAccountActive", "POST", "/v1/account/active", IsAccountActive},
		shared.Route{"GetAccountBalance", "POST", "/v1/account/balance", GetAccountBalance},
		shared.Route{"GetDateJoined", "POST", "/v1/account/joined", GetDateJoined},
		shared.Route{"GetAccountKarma", "POST", "/v1/account/karma", GetAccountKarma},
		shared.Route{"GetlastActive", "POST", "/v1/account/last_active", GetLastActive},
		shared.Route{"GetAccTxCount", "POST", "/v1/account/transaction_count", GetAccTxCount},
		shared.Route{"GetAccSessCount", "POST", "/v1/account/session_count", GetAccSessCount},
		shared.Route{"GetAccStatus", "POST", "/v1/account/status", GetAccStatus},
		shared.Route{"GetClientInfo", "POST", "/v1/client", GetClientInfo},
		shared.Route{"GetClientID", "POST", "/v1/client/id", GetClientID},
		shared.Route{"GetClientVersion", "POST", "/v1/client/version", GetClientVersion},
		shared.Route{"GetCliSyncStatus", "POST", "/v1/client/syncing", GetCliSyncStatus},
		shared.Route{"GetNetworkInfo", "POST", "/v1/network", GetNetworkInfo},
		shared.Route{"GetNetworkID", "POST", "/v1/network/id", GetNetworkID},
		shared.Route{"GetPeerCount", "POST", "/v1/network/peer_count", GetPeerCount},
		shared.Route{"GetPeerList", "POST", "/v1/network/peer_list", GetPeerList},
		shared.Route{"GetPeers", "POST", "/v1/network/peers", GetPeers},
		shared.Route{"GetPersonalInfo", "POST", "/v1/personal", GetPersonalInfo},
		shared.Route{"ListAccounts", "POST", "/v1/personal/list_accounts", ListAccounts},
		shared.Route{"PersonalNetOptions", "POST", "/v1/personal/network", PersonalNetOptions},
		shared.Route{"EnterNetwork", "POST", "/v1/personal/network/enter", EnterNetwork},
		shared.Route{"ExitNetwork", "POST", "/v1/personal/network/exit", ExitNetwork},
		shared.Route{"GetPrimaryAddr", "POST", "/v1/personal/primary_address", GetPrimaryAddr},
		shared.Route{"SendPOKT", "POST", "/v1/personal/send", SendPOKT},
		shared.Route{"SendPOKTRaw", "POST", "/v1/personal/send/raw", SendPOKTRaw},
		shared.Route{"Sign", "POST", "/v1/personal/sign", Sign},
		shared.Route{"StakeOptions", "POST", "/v1/personal/stake", StakeOptions},
		shared.Route{"Stake", "POST", "/v1/personal/stake/add", Stake},
		shared.Route{"UnStake", "POST", "/v1/personal/stake/remove", UnStake},
		shared.Route{"GetPocketBCInfo", "POST", "/v1/pocket", GetPocketBCInfo},
		shared.Route{"GetLatestBlock", "POST", "/v1/pocket/block", GetLatestBlock},
		shared.Route{"GetBlockByHash", "POST", "/v1/pocket/block/hash", GetBlockByHash},
		shared.Route{"GetBlkTxCntByHash", "POST", "/v1/pocket/block/hash/transaction_count", GetBlkTxCntByHash},
		shared.Route{"GetBlockByNum", "POST", "/v1/pocket/block/number", GetBlockByNum},
		shared.Route{"GetBlkTxCntByNum", "POST", "/v1/pocket/block/number/transaction_count", GetBlkTxCntByNum},
		shared.Route{"GetProtocolVersion", "POST", "/v1/pocket/version", GetProtocolVersion},
		shared.Route{"TxOptions", "POST", "/v1/pocket/transaction", TxOptions},
		shared.Route{"GetTxByHash", "POST", "/v1/pocket/transaction/hash", GetTxByHash},
	}
	return routes
}
