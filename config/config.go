// This package maintains the client configuration.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pokt-network/pocket-core/const"
	"sync"
)

// "config.go" describes all of the configuration properties of the client (set by startup flags)

type config struct {
	Clientid      string `json:"CLIENTID"` 		// This variable holds a client identifier string.
	Version       string `json:"VERSION"`  		// This variable holds the client version string.
	Datadir       string `json:"DATADIR"`  		// This variable holds the working directory string.
	Clientrpc     bool   `json:"CRPC"`     		// This variable describes if the client rpc is running.
	Clientrpcport string `json:"CRPCPORT"`		// This variable holds the client rpc port string.
	Relayrpc      bool   `json:"RRPC"`    		// This variable describes if the relay rpc is running.
	Relayrpcport  string `json:"RRPCPORT"` 		// This variable holds the relay rpc port string.
	Ethereum	  bool	 `json:"ETHEREUM"` 		// This variable describes if Ethereum is hosted.
	Ethrpcport	  string `json:"ETHRPCPORT"` 	// This variable holds the port the ETH rpc is running on.
	Bitcoin		  bool 	 `json:"BITCOIN"`		// This variable describes if Bitcoin is hosted.
	Btcrpcport	  string `json:"BTCRPCPORT"`	// This variable holds the port the BTC rpc is running on.
}

var (
	instance *config
	once     sync.Once
	datadir  = flag.String("datadir", _const.DATADIR, "setup the data director for the DB and keystore")
	// A boolean variable derived from flags, that describes whether or not to print the version of the client.
	client_rpc = flag.Bool("clientrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the client rpc (default :8080)
	client_rpcport = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	// A boolean variable derived from flags, that describes whether or not to start the relay rpc server.
	relay_rpc = flag.Bool("relayrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the relay rpc (default :8081)
	relay_rpcport = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	// A boolean variable derived from flags, that specifies if Ethereum is hosted.
	ethereum = flag.Bool("ethereum", false, "whether or not ethereum is hosted")
	// A string variable derived from flags, that specifies which port Ethereum's json rpc is running.
	eth_rpcport = flag.String("ethrpcport", "8545", "specified port to run ethereum rpc")
	// A boolean variable derived from flags, that specifies if Bitcoin is hosted.
	bitcoin = flag.Bool("bitcoin", false, "whether or not bitcoin is hosted")
	// A string variable derived from flags, that specifies which port Bitcoin's json rpc is running.
	btc_rpcport = flag.String("btcrpcport", "8333", "specified port to run bitcoin rpc")
)

func InitializeConfiguration() {
	flag.Parse()        // built in function to parse the flags above.
	GetConfigInstance() // returns the thread safe instance of the client configuration.
}

/*
"NewConfiguration() is a Constructor function of the configuration type.
 */
func newConfiguration() {
	instance = &config{
		_const.CLIENTID,			// client identifier is set in global constants.
		_const.VERSION,			// client version is set in global constants.
		*datadir,				// data directory path specified by the flag.
		*client_rpc,			// the client rpc is running.
		*client_rpcport,		// the port the client rpc is running.
		*relay_rpc,				// the relay rpc is running.
		*relay_rpcport,		// the port the relay rpc is running.
		*ethereum,				// ethereum is hosted
		*eth_rpcport,			// the port Ethereum's rpc is on
		*bitcoin,				// bitcoin is hosted
		*btc_rpcport}			// the port Bitcoin's rpc is on
}

/*
"PrintConfiguration()" prints the client configuration information to the CLI.
 */
func PrintConfiguration() {
	data, _ := json.MarshalIndent(instance, "", "    ")           	// pretty configure the json data
	fmt.Println("Pocket Core Configuration:\n", string(data))     			// pretty print the pocket configuration
}

/*
"GetConfigInstance()" returns the configuration object in a thread safe manner.
 */
func GetConfigInstance() *config { 	// singleton structure to return the configuration object
	once.Do(func() {				// thread safety.
		if instance == nil {		// if nil make a new configuration
			newConfiguration()
		}
	})
	return instance					// return the configuration
}
