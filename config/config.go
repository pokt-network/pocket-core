/*
Package config builds, maintains, and serves as the source of the client's configuration

"config.go" describes all of the configuration properties of the client (set by startup flags)
"build.go" is for building the Pocket Core configuration
*/
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/pokt-network/pocket-core/const"
)

// TODO configuration updating
type config struct {
	Clientid      string `json:"CLIENTID"`   // This variable holds a client identifier string.
	Version       string `json:"VERSION"`    // This variable holds the client version string.
	Datadir       string `json:"DATADIR"`    // This variable holds the working directory string.
	Clientrpc     bool   `json:"CRPC"`       // This variable describes if the client rpc is running.
	Clientrpcport string `json:"CRPCPORT"`   // This variable holds the client rpc port string.
	Relayrpc      bool   `json:"RRPC"`       // This variable describes if the relay rpc is running.
	Relayrpcport  string `json:"RRPCPORT"`   // This variable holds the relay rpc port string.
	Ethereum      bool   `json:"ETHEREUM"`   // This variable describes if Ethereum is hosted.
	Ethrpcport    string `json:"ETHRPCPORT"` // This variable holds the port the ETH rpc is running on.
	Bitcoin       bool   `json:"BITCOIN"`    // This variable describes if Bitcoin is hosted.
	Btcrpcport    string `json:"BTCRPCPORT"` // This variable holds the port the BTC rpc is running on.
	PeerFile      string `json:"PEERFILE"`   // This variable holds the filepath to the peerFile.json
	ManPeers      bool   `json:"MANPEERS"`   // This variable specifies if manual peers are being used
}

var (
	instance       *config
	once           sync.Once
	datadir        = flag.String("datadir", _const.DATADIR, "setup the data directory for the DB and keystore")
	runClientRpc   = flag.Bool("clientrpc", false, "whether or not to start the rpc server")
	clientRpcport  = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	runRelayRpc    = flag.Bool("relayrpc", false, "whether or not to start the rpc server")
	relayRpcport   = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	ethereumHosted = flag.Bool("ethereum", false, "whether or not ethereum is hosted")
	ethRpcport     = flag.String("ethrpcport", "8545", "specified port to run ethereum rpc")
	bitcoinHosted  = flag.Bool("bitcoin", false, "whether or not bitcoin is hosted")
	btcRpcport     = flag.String("btcrpcport", "8333", "specified port to run bitcoin rpc")
	manPeers       = flag.Bool("manpeers", false, "specifies if peers are manually added")
	peerFile       = flag.String("peerFile", _const.DATADIR+_const.FILESEPARATOR+"peers.json", "specifies the filepath for peers.json")
)

func constructConfiguration() {
	instance = &config{
		// Client identifier; set in global constants.
		_const.CLIENTID,
		// Client version; set in global constants.
		_const.VERSION,
		// Data directory path; specified by flag.
		*datadir,
		// Client RPC .
		*runClientRpc,
		// Port for the client RPC.
		*clientRpcport,
		// Whether the relay rpc should run.
		*runRelayRpc,
		// Port the relay RPC.
		*relayRpcport,
		// ethereum is hosted
		*ethereumHosted,
		// the port Ethereum's rpc is on
		*ethRpcport,
		// bitcoin is hosted
		*bitcoinHosted,
		// the port Bitcoin's rpc is on
		*btcRpcport,
		// the filepath of the peers.json
		*peerFile,
		// using manual peers
		*manPeers}
}

/*
"PrintConfiguration()" prints the client configuration information to the CLI.
*/
func PrintConfiguration() {
	data, _ := json.MarshalIndent(instance, "", "    ")       // pretty configure the json data
	fmt.Println("Pocket Core Configuration:\n", string(data)) // pretty print the pocket configuration
}

/*
"GetConfigInstance()" returns the configuration object.

Because there is only one configuration per instance of the program, initialization is guarded by a sync.Once instance.
This invocation ensures the flags are parsed before the configuration is instantiated.
*/
func GetConfigInstance() *config {
	once.Do(func() {
		flag.Parse()
		constructConfiguration()
	})
	return instance
}
