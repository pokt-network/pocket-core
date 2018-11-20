// This package maintains the client configuration.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pokt-network/pocket-core/const"
	"sync"
)

// This file describes all of the configuration properties of the client (set by startup flags)

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
	datadir=flag.String("datadir",_const.DATADIR, "setup the data director for the DB and keystore")
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

/*
"parseFlags" reads in specific command line arguments passed to the client and executes events based
on the input.
 */
 /*
 The default value is `%APPDATA%\Pocket` for Windows, `~/.pocket` for Linux, `~/Library/Pocket` for Mac
  */
func parseFlags() {
	flag.Parse()
}

func InitializeConfiguration(){
	parseFlags()
	GetInstance()
}

/*
"NewConfiguration() is a Constructor function of the configuration type.
 */
func newConfiguration() {
		instance = &config{
			_const.CLIENTID,
			_const.VERSION,
			*datadir,
			*client_rpc,
			*client_rpcport,
			*relay_rpc,
			*relay_rpcport}
}

/*
"PrintConfiguration()" prints the client configuration information to the CLI.
 */
func PrintConfiguration() {
	data, _ := json.MarshalIndent(instance, "", "    ")
	fmt.Println("Pocket Core Configuration:\n", string(data))
}

/*
"GetInstance()" returns the configuration object in a thread safe manner.
 */
func GetInstance() *config {
	once.Do(func() {
		if instance==nil {
			newConfiguration()
		}
	})
	return instance
}
