// This package is for command line specific utilities.
package util

import (
	"flag"
	"fmt"
)
/*
flags.go specifies startup command flags for the client
 */

var (
	// A string variable derived from flags, used to setup data directory of Pocket Core.
	datadir = flag.String("datadir", "/path/to/data/dir", "setup the data directory for the DB and keystore")
	// A boolean variable derived from flags, that describes whether or not to print the version of the client.
	print_version = flag.Bool("version", false, "whether or not to print the version of the client")
	// A boolean variable derived from flags, that describes whether or not to start the rpc server.
	rpc = flag.Bool("rpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the rpc (default :8545)
	rpcport = flag.String("rpcport", "8545", "specified port to run rpc")
	// A boolean variable derived from flags, that describes whether or not to start the web sockets server
	ws = flag.Bool("ws", false, "whether or not to start the ws server")
	// A string variable derived from flags, that specifies the websocket listening interface (default :localhost).
	wsaddr = flag.String("wsaddr", "localhost", "specified ws listening addr")
	// A string variable derived from flags, that specifies the websocket listening port (default :8546).
	wsport = flag.String("wsport", "8546", "specifying ws listening port")

)

/*
"ParseFlags" reads in specific command line arguments passed to the client and executes events based
on the input.
 */
func ParseFlags() {
	flag.Parse()
}

/*
"PrintClientInfo" prints startup information about the client
 */
func PrintClientInfo(){
	fmt.Println("CLIENT ID:", clientIdentifier)
	if *print_version{
		fmt.Println("VERSION:", version)
	}
	fmt.Println("DATADIR:", *datadir)
	fmt.Println("RPC:", *rpc)
	if *rpc{
		fmt.Println("RPCPORT:",*rpcport)
	}
	fmt.Println("WS:", *ws);
	if *ws{
		fmt.Println("WSADDR:",*wsaddr)
		fmt.Println("WSPORT:", *wsport)
	}
}
