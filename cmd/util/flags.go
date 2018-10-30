// This package is for command line specific utilities.
package util

import (
	"flag"
	"fmt"

	"github.com/pocket_network/pocket-core/rpc"
)

// flags.go specifies startup command flags for the client

var (
	// A string variable derived from flags, used to setup data directory of Pocket Core.
	datadir = flag.String("datadir", "/path/to/data/dir", "setup the data directory for the DB and keystore")
	// A boolean variable derived from flags, that describes whether or not to print the version of the client.
	print_version = flag.Bool("version", false, "whether or not to print the version of the client")
	// A boolean variable derived from flags, that describes whether or not to start the client rpc server.
	client_rpc = flag.Bool("clientrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the client rpc (default :8545)
	client_rpcport = flag.String("clientrpcport", "8545", "specified port to run rpc")
	// A boolean variable derived from flags, that describes whether or not to start the relay rpc server.
	relay_rpc = flag.Bool("relayrpc", false, "whether or not to start the rpc server")
	// A string variable derived from flags, that specifies which port to run the listener for the relay rpc (default :8546)
	relay_rpcport = flag.String("relayrpcport", "8546", "specified port to run rpc")
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

//TODO move tests
func CLI_Test() {
	fmt.Println("CLIENT ID:", clientIdentifier)
	if *print_version {
		fmt.Println("VERSION:", version)
	}
	fmt.Println("DATADIR:", *datadir)
	fmt.Println("CLIENT RPC:", *client_rpc)
	if *client_rpc {
		fmt.Println("CLIENT RPCPORT:", *client_rpcport)
		go rpc.StartClientRPC(*client_rpcport)
	}
	fmt.Println("RELAY RPC:", *relay_rpc)
	if *relay_rpc {
		fmt.Println("RELAY RPCPORT:", *relay_rpcport)
		rpc.StartRelayRPC(*relay_rpcport) // TODO should be changed to goroutine
	}
	fmt.Println("WS:", *ws)
	if *ws {
		fmt.Println("WSADDR:", *wsaddr)
		fmt.Println("WSPORT:", *wsport)
	}
}
