// This package is contains the handler functions needed for the Relay API
package relay

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/plugin/pcp-bitcoin"
	"github.com/pokt-network/pocket-core/plugin/pcp-ethereum"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"net/http"
)

const (
	relayReadMethod = "relayRead()"
)

// "relay.go" defines API handlers that are under the 'relay' category within this file.

/*
 "RelayOptions" handles the localhost:<relay-port>/v1/relay call.
*/
func RelayOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "RelayRead" handles the localhost:<relay-port>/v1/relay/read call.
*/
func RelayRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	relay := &Relay{}                               // create empty relay structure
	shared.PopulateModelFromParams(w, r, ps, relay) // populate the relay struct from params
	response := RouteRelay(*relay)                  // route the relay to the correct chain
	shared.WriteJSONResponse(w, response)           // relay the response
}

/*
"RelayReadInfo" handles a get request to localhost:<relay-port>/v1/relay/read call.
And provides the developers with an in-client reference to the API call
*/
func RelayReadInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "RelayRead", Relay{}, "Response from hosted chain")
	shared.WriteInfoResponse(w, info)
}

/*
 "RelayWrite" handles the localhost:<relay-port>/v1/relay/write call.
*/
func RelayWrite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
"RouteRelay" routes the relay to the specified hosted chain
*/
func RouteRelay(relay Relay) string {
	switch relay.Blockchain {
	case "ethereum":
		return pcp_ethereum.ExecuteRequest([]byte(relay.Data), config.GetConfigInstance().Ethrpcport)
	case "bitcoin":
		return pcp_bitcoin.ExecuteRequest([]byte(relay.Data), config.GetConfigInstance().Btcrpcport)
	}
	logs.NewLog("Not a supported blockchain", logs.ErrorLevel, logs.JSONLogFormat)
	return "Error: not a supported blockchain"
}
