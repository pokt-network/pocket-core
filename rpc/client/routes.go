// This package contains files for the Client API
package client

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pocket_network/pocket-core/rpc/client/handlers"
	"github.com/pocket_network/pocket-core/rpc/shared"
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
		shared.Route{"GetClientAPIVersion", "GET", "/v1", handlers.GetClientAPIVersion},
		shared.Route{"GetAccount", "GET", "/v1/account", handlers.GetAccount},
		shared.Route{"IsAccountActive", "GET", "/v1/account/active", handlers.IsAccountActive},
		shared.Route{"GetAccountBalance", "GET", "/v1/account/balance", handlers.GetAccountBalance},
		shared.Route{"GetDateJoined", "GET", "/v1/account/joined", handlers.GetDateJoined},
		shared.Route{"GetAccountKarma", "GET", "/v1/account/karma", handlers.GetAccountKarma},
		shared.Route{"GetlastActive", "GET", "/v1/account/last_active", handlers.GetLastActive},
		shared.Route{"GetAccTxCount", "GET", "/v1/account/transaction_count", handlers.GetAccTxCount},
		shared.Route{"GetAccSessCount", "GET", "/v1/account/session_count", handlers.GetAccSessCount},
		shared.Route{"GetAccStatus", "GET", "/v1/account/status", handlers.GetAccStatus},
		shared.Route{"GetClientInfo", "GET", "/v1/client", handlers.GetClientInfo},
		shared.Route{"GetClientID", "GET", "/v1/client/id", handlers.GetClientID},
		shared.Route{"GetClientVersion", "GET", "/v1/client/version", handlers.GetClientVersion},
		shared.Route{"GetCliSyncStatus", "GET", "/v1/client/syncing", handlers.GetCliSyncStatus},
		shared.Route{"GetNetworkInfo", "GET", "/v1/network", handlers.GetNetworkInfo},
		shared.Route{"GetNetworkID", "GET", "/v1/network/id", handlers.GetNetworkID},
		shared.Route{"GetPeerCount", "GET", "/v1/network/peer_count", handlers.GetPeerCount},
		shared.Route{"GetPeerList", "GET", "/v1/network/peer_list", handlers.GetPeerList},
		shared.Route{"GetPeers", "GET", "/v1/network/peers", handlers.GetPeers},
		shared.Route{"GetPersonalInfo", "GET", "/v1/personal", handlers.GetPersonalInfo},
		shared.Route{"ListAccounts", "GET", "/v1/personal/list_accounts", handlers.ListAccounts},
		shared.Route{"PersonalNetOptions", "GET", "/v1/personal/network", handlers.PersonalNetOptions},
		shared.Route{"EnterNetwork", "GET", "/v1/personal/network/enter", handlers.EnterNetwork},
		shared.Route{"ExitNetwork", "GET", "/v1/personal/network/exit", handlers.ExitNetwork},
		shared.Route{"GetPrimaryAddr", "GET", "/v1/personal/primary_address", handlers.GetPrimaryAddr},
		shared.Route{"SendPOKT", "GET", "/v1/personal/send", handlers.SendPOKT},
		shared.Route{"SendPOKTRaw", "GET", "/v1/personal/send/raw", handlers.SendPOKTRaw},
		shared.Route{"Sign", "GET", "/v1/personal/sign", handlers.Sign},
		shared.Route{"StakeOptions", "GET", "/v1/personal/stake", handlers.StakeOptions},
		shared.Route{"Stake", "GET", "/v1/personal/stake/add", handlers.Stake},
		shared.Route{"UnStake", "GET", "/v1/personal/stake/remove", handlers.UnStake},
		shared.Route{"GetPocketBCInfo", "GET", "/v1/pocket", handlers.GetPocketBCInfo},
		shared.Route{"GetLatestBlock", "GET", "/v1/pocket/block", handlers.GetLatestBlock},
		shared.Route{"GetBlockByHash", "GET", "/v1/pocket/block/hash", handlers.GetBlockByHash},
		shared.Route{"GetBlkTxCntByHash", "GET", "/v1/pocket/block/hash/transaction_handlers.count", handlers.GetBlkTxCntByHash},
		shared.Route{"GetBlockByNum", "GET", "/v1/pocket/block/number", handlers.GetBlockByNum},
		shared.Route{"GetBlkTxCntByNum", "GET", "/v1/pocket/block/number/transaction_count", handlers.GetBlkTxCntByNum},
		shared.Route{"GetProtocolVersion", "GET", "/v1/pocket/version", handlers.GetProtocolVersion},
		shared.Route{"TxOptions", "GET", "/v1/pocket/transaction", handlers.TxOptions},
		shared.Route{"GetTxByHash", "GET", "/v1/pocket/transaction/hash", handlers.GetTxByHash},
	}
	return routes
}
