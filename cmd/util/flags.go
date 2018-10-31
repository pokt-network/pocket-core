// This package is for command line specific utilities.
package util

import (
	"flag"
)

// flags.go specifies startup command flags for the client

var (
	// A string variable derived from flags, used to setup data directory of Pocket Core.
	datadir = flag.String("datadir", "/path/to/data/dir", "setup the data directory for the DB and keystore")
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
"StartConfig" reads in specific command line arguments passed to the client and executes events based
on the input.
 */
func StartConfig() {
	flag.Parse()
	NewConfiguration(*datadir, *client_rpc, *relay_rpc, *client_rpcport, *relay_rpcport)
}
