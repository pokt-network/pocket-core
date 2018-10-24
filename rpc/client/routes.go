package client

import(
	"github.com/pocket_network/pocket-core/rpc"
	"github.com/pocket_network/pocket-core/rpc/client/handlers"
)
/*
"clientRoutes" is a function that returns all of the routes of the API.
 */
func ClientRoutes() rpc.Routes {
	routes := rpc.Routes{
		rpc.Route{"GetClientAPIVersion", "POST", "/v1/", handlers.GetClientAPIVersion},
		rpc.Route{"GetAccount", "POST", "/v1/account", handlers.GetAccount},
		rpc.Route{"IsAccountActive", "POST", "/v1/account/active/", handlers.IsAccountActive},
		rpc.Route{"GetAccountBalance", "POST", "/v1/account/balance/", handlers.GetAccountBalance},
		rpc.Route{"GetDateJoined", "POST", "/v1/account/joined/", handlers.GetDateJoined},
		rpc.Route{"GetAccountKarma", "POST", "/v1/account/karma/", handlers.GetAccountKarma},
		rpc.Route{"GetlastActive", "POST", "/v1/account/last_active/", handlers.GetLastActive},
		rpc.Route{"GetAccTxCount", "POST", "/v1/account/transaction_count/", handlers.GetAccTxCount},
		rpc.Route{"GetAccSessCount", "POST", "/v1/account/session_count/", handlers.GetAccSessCount},
		rpc.Route{"GetAccStatus", "POST", "/v1/account/status/", handlers.GetAccStatus},
		rpc.Route{"GetClientInfo", "POST", "/v1/client/", handlers.GetClientInfo},
		rpc.Route{"GetClientID", "POST", "/v1/client/id/", handlers.GetClientID},
		rpc.Route{"GetClientVersion", "POST", "/v1/client/version/", handlers.GetClientVersion},
		rpc.Route{"GetCliSyncStatus", "POST", "/v1/client/syncing/", handlers.GetCliSyncStatus},
		rpc.Route{"GetNetworkInfo", "POST", "/v1/network/", handlers.GetNetworkInfo},
		rpc.Route{"GetNetworkID", "POST", "/v1/network/id/", handlers.GetNetworkID},
		rpc.Route{"GetPeerCount", "POST", "/v1/network/peer_count/", handlers.GetPeerCount},
		rpc.Route{"GetPeerList", "POST", "/v1/network/peer_list/", handlers.GetPeerList},
		rpc.Route{"GetPeers", "POST", "/v1/network/peers/", handlers.GetPeers},
		rpc.Route{"GetPersonalInfo", "POST", "/v1/personal/", handlers.GetPersonalInfo},
		rpc.Route{"ListAccounts", "POST", "/v1/personal/list_accounts/", handlers.ListAccounts},
		rpc.Route{"PersonalNetOptions", "POST", "/v1/personal/network/", handlers.PersonalNetOptions},
		rpc.Route{"EnterNetwork", "POST", "/v1/personal/network/enter/", handlers.EnterNetwork},
		rpc.Route{"ExitNetwork", "POST", "/v1/personal/network/exit/", handlers.ExitNetwork},
		rpc.Route{"GetPrimaryAddr", "POST", "/v1/personal/primary_address/", handlers.GetPrimaryAddr},
		rpc.Route{"SendPOKT", "POST", "/v1/personal/send/", handlers.SendPOKT},
		rpc.Route{"SendPOKTRaw", "POST", "/v1/personal/send/raw/", handlers.SendPOKTRaw},
		rpc.Route{"Sign", "POST", "/v1/personal/sign/", handlers.Sign},
		rpc.Route{"StakeOptions", "POST", "/v1/personal/stake/", handlers.StakeOptions},
		rpc.Route{"Stake", "POST", "/v1/personal/stake/add/", handlers.Stake},
		rpc.Route{"UnStake", "POST", "/v1/personal/stake/remove/", handlers.UnStake},
		rpc.Route{"GetPocketBCInfo", "POST", "/v1/pocket/", handlers.GetPocketBCInfo},
		rpc.Route{"GetLatestBlock", "POST", "/v1/pocket/block/", handlers.GetLatestBlock},
		rpc.Route{"GetBlockByHash", "POST", "/v1/pocket/block/hash/", handlers.GetBlockByHash},
		rpc.Route{"GetBlkTxCntByHash", "POST", "/v1/pocket/block/hash/transaction_handlers.count/", handlers.GetBlkTxCntByHash},
		rpc.Route{"GetBlockByNum", "POST", "/v1/pocket/block/number/", handlers.GetBlockByNum},
		rpc.Route{"GetBlkTxCntByNum", "POST", "/v1/pocket/block/number/transaction_count/", handlers.GetBlkTxCntByNum},
		rpc.Route{"GetProtocolVersion", "POST", "/v1/pocket/version/", handlers.GetProtocolVersion},
		rpc.Route{"TxOptions", "POST", "/v1/pocket/transaction/", handlers.TxOptions},
		rpc.Route{"GetTxByHash", "POST", "/v1/pocket/transaction/hash/", handlers.GetTxByHash},
	}
	return routes
}
