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
// TODO configuration updating
type config struct {
	GID				string `json:"GID"`				// This variable holds the nodes GID
	Clientid      	string `json:"CLIENTID"`    	// This variable holds a client identifier string.
	Version       	string `json:"VERSION"`       	// This variable holds the client version string.
	Datadir        	string `json:"DATADIR"`      	// This variable holds the working directory string.
	Clientrpc      	bool   `json:"CRPC"`         	// This variable describes if the client rpc is running.
	Clientrpcport  	string `json:"CRPCPORT"`     	// This variable holds the client rpc port string.
	Relayrpc       	bool   `json:"RRPC"`         	// This variable describes if the relay rpc is running.
	Relayrpcport   	string `json:"RRPCPORT"`     	// This variable holds the relay rpc port string.
	ChainsFilepath 	string `json:"HOSTEDCHAINS"` 	// This variable holds the filepath to the chains.json
	PeerFile       	string `json:"PEERFILE"`     	// This variable holds the filepath to the peerFile.json
	Whitelist	   string `json:"WHITELIST"`	// This variable holds the filepath to the whitelist.json
}

var (
	instance *config
	once     sync.Once
	gid = flag.String("gid","0", "set the selfNode.GID for pocket core mvp")
	datadir  = flag.String("datadir", _const.DATADIR, "setup the data directory for the DB and keystore")
	// A boolean variable derived from flags, that describes whether or not to print the version of the client.
	clientRpc = flag.Bool("clientrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the client rpc (default :8080)
	clientRpcport = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	// A boolean variable derived from flags, that describes whether or not to start the relay rpc server.
	relayRpc = flag.Bool("relayrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the relay rpc (default :8081)
	relayRpcport = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	// A string variable derived from flags, that specifies a custom path for the hosted chains
	hostedChains = flag.String("hostedchains", _const.CHAINSFILENAME, "specifies the filepath for hosted chains")
	// A string variable derived from flags, that specifies the filepath for peerFile.json
	peerFile = flag.String("peerFile", _const.PEERFILENAME, "specifies the filepath for peers.json")
	// A string variable derived from flags, that specifies the fielpath for whitelist.json
	whitelistFIle = flag.String("whitelist", _const.WHITELISTFILENAME, "specifies the filepath for whitelist.json")
)

func InitializeConfiguration() {
	flag.Parse()        				// built in function to parse the flags above.
	GetConfigInstance() 				// returns the thread safe instance of the client configuration.
}

/*
"NewConfiguration() is a Constructor function of the configuration type.
*/
func newConfiguration() {
	instance = &config{
		*gid,						// the global identifier of this node
		_const.CLIENTID, 		// client identifier is set in global constants.
		_const.VERSION,  		// client version is set in global constants.
		*datadir,        		// data directory path specified by the flag.
		*clientRpc,      		// the client rpc is running.
		*clientRpcport,  	// the port the client rpc is running.
		*relayRpc,       		// the relay rpc is running.
		*relayRpcport,   	// the port the relay rpc is running.
		*hostedChains,		// the filepath for the chains.json
		*peerFile,				// the filepath of the peers.json
		*whitelistFIle}       	// the filepath of the whitelist.json
}

/*
"PrintConfiguration()" prints the client configuration information to the CLI.
*/
func PrintConfiguration() {
	data, _ := json.MarshalIndent(instance, "", "    ")       	// pretty configure the json data
	fmt.Println("Pocket Core Configuration:\n", string(data)) 			// pretty print the pocket configuration
}

/*
"GetConfigInstance()" returns the configuration object in a thread safe manner.
*/
func GetConfigInstance() *config { 						// singleton structure to return the configuration object
	once.Do(func() { 									// thread safety.
			newConfiguration()
	})
	return instance 									// return the configuration
}
