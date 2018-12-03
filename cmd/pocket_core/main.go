// This package is the starting point of the CLI.
package main

import (
	"bufio"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/net"
	"github.com/pokt-network/pocket-core/rpc"
	"os"
)

// "main.go" is the entry point of the client.

/*
"init" is a built in function that is automatically called before main.
 */
func init() {
	crypto.GenerateSeed() 								// generates seed for randomization
}

/*
"main" is the starting function of the client.
 Keep main as light as possible by calling accessory functions.
*/
func main() {
	startClient() 										// see function call below
}

/*
"startClient" Starts the client with the given initial configuration.
 */
func startClient() {
	config.InitializeConfiguration()                	// initializes the configuration from flags and defaults.
	config.BuildConfiguration()                     	// builds the proper structure on pc for core client to operate.
	config.PrintConfiguration()                     	// print the configuration the the cmd.
	net.DummyList()										// feed the peerlist with dummy data
	rpc.RunAPIEndpoints()                           	// runs the server endpoints for client and relay api.
	fmt.Print("Press any key + 'Return' to quit: ") 	// prompt user to exit
	input := bufio.NewScanner(os.Stdin)             	// unnecessary temporary entry
	input.Scan()                                    	// wait
}
