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
		shared.Route{"GetClientAPIVersion", "POST", "/v1/", handlers.GetClientAPIVersion},
		shared.Route{"GetAccount", "POST", "/v1/account", handlers.GetAccount},
		shared.Route{"IsAccountActive", "POST", "/v1/account/active/", handlers.IsAccountActive},
		shared.Route{"GetAccountBalance", "POST", "/v1/account/balance/", handlers.GetAccountBalance},
		shared.Route{"GetDateJoined", "POST", "/v1/account/joined/", handlers.GetDateJoined},
		shared.Route{"GetAccountKarma", "POST", "/v1/account/karma/", handlers.GetAccountKarma},
		shared.Route{"GetlastActive", "POST", "/v1/account/last_active/", handlers.GetLastActive},
		shared.Route{"GetAccTxCount", "POST", "/v1/account/transaction_count/", handlers.GetAccTxCount},
		shared.Route{"GetAccSessCount", "POST", "/v1/account/session_count/", handlers.GetAccSessCount},
		shared.Route{"GetAccStatus", "POST", "/v1/account/status/", handlers.GetAccStatus},
		shared.Route{"GetClientInfo", "POST", "/v1/client/", handlers.GetClientInfo},
		shared.Route{"GetClientID", "POST", "/v1/client/id/", handlers.GetClientID},
		shared.Route{"GetClientVersion", "POST", "/v1/client/version/", handlers.GetClientVersion},
		shared.Route{"GetCliSyncStatus", "POST", "/v1/client/syncing/", handlers.GetCliSyncStatus},
		shared.Route{"GetNetworkInfo", "POST", "/v1/network/", handlers.GetNetworkInfo},
		shared.Route{"GetNetworkID", "POST", "/v1/network/id/", handlers.GetNetworkID},
		shared.Route{"GetPeerCount", "POST", "/v1/network/peer_count/", handlers.GetPeerCount},
		shared.Route{"GetPeerList", "POST", "/v1/network/peer_list/", handlers.GetPeerList},
		shared.Route{"GetPeers", "POST", "/v1/network/peers/", handlers.GetPeers},
		shared.Route{"GetPersonalInfo", "POST", "/v1/personal/", handlers.GetPersonalInfo},
		shared.Route{"ListAccounts", "POST", "/v1/personal/list_accounts/", handlers.ListAccounts},
		shared.Route{"PersonalNetOptions", "POST", "/v1/personal/network/", handlers.PersonalNetOptions},
		shared.Route{"EnterNetwork", "POST", "/v1/personal/network/enter/", handlers.EnterNetwork},
		shared.Route{"ExitNetwork", "POST", "/v1/personal/network/exit/", handlers.ExitNetwork},
		shared.Route{"GetPrimaryAddr", "POST", "/v1/personal/primary_address/", handlers.GetPrimaryAddr},
		shared.Route{"SendPOKT", "POST", "/v1/personal/send/", handlers.SendPOKT},
		shared.Route{"SendPOKTRaw", "POST", "/v1/personal/send/raw/", handlers.SendPOKTRaw},
		shared.Route{"Sign", "POST", "/v1/personal/sign/", handlers.Sign},
		shared.Route{"StakeOptions", "POST", "/v1/personal/stake/", handlers.StakeOptions},
		shared.Route{"Stake", "POST", "/v1/personal/stake/add/", handlers.Stake},
		shared.Route{"UnStake", "POST", "/v1/personal/stake/remove/", handlers.UnStake},
		shared.Route{"GetPocketBCInfo", "POST", "/v1/pocket/", handlers.GetPocketBCInfo},
		shared.Route{"GetLatestBlock", "POST", "/v1/pocket/block/", handlers.GetLatestBlock},
		shared.Route{"GetBlockByHash", "POST", "/v1/pocket/block/hash/", handlers.GetBlockByHash},
		shared.Route{"GetBlkTxCntByHash", "POST", "/v1/pocket/block/hash/transaction_handlers.count/", handlers.GetBlkTxCntByHash},
		shared.Route{"GetBlockByNum", "POST", "/v1/pocket/block/number/", handlers.GetBlockByNum},
		shared.Route{"GetBlkTxCntByNum", "POST", "/v1/pocket/block/number/transaction_count/", handlers.GetBlkTxCntByNum},
		shared.Route{"GetProtocolVersion", "POST", "/v1/pocket/version/", handlers.GetProtocolVersion},
		shared.Route{"TxOptions", "POST", "/v1/pocket/transaction/", handlers.TxOptions},
		shared.Route{"GetTxByHash", "POST", "/v1/pocket/transaction/hash/", handlers.GetTxByHash},
	}
	return routes
}
