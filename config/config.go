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
	Clientid      string `json:"CLIENTID"` // This variable holds a client identifier string.
	Version       string `json:"VERSION"`  // This variable holds the client version string.
	Datadir       string `json:"DATADIR"`  // This variable holds the working directory string.
	Clientrpc     bool   `json:"CRPC"`     // This variable describes if the client rpc is running.
	Clientrpcport string `json:"CRPCPORT"` // This variable holds the client rpc port string.
	Relayrpc      bool   `json:"RRPC"`     // This variable describes if the relay rpc is running.
	Relayrpcport  string `json:"RRPCPORT"` // This variable holds the relay rpc port string.
}

var (
	instance *config
	once     sync.Once
	datadir  = flag.String("datadir", _const.DATADIR, "setup the data director for the DB and keystore")
	// A boolean variable derived from flags, that describes whether or not to print the version of the client.
	client_rpc = flag.Bool("clientrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the client rpc (default :8545)
	client_rpcport = flag.String("clientrpcport", "8545", "specified port to run rpc")
	// A boolean variable derived from flags, that describes whether or not to start the relay rpc server.
	relay_rpc = flag.Bool("relayrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the relay rpc (default :8546)
	relay_rpcport = flag.String("relayrpcport", "8546", "specified port to run rpc")
	// A boolean variable derived from flags, that describes whether or not to start the web sockets server
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
		*relay_rpcport}		// the port the relay rpc is running.
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
	return instance					// else return the configuration
}
