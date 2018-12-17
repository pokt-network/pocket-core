// This package is the starting point of the CLI.
package main

import (
	"bufio"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/net"
	"github.com/pokt-network/pocket-core/net/session"
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
"manualPeers" checks if manual peers are specified and adds them to the peerlist.
 */
func manualPeers(){
	if config.GetConfigInstance().ManPeers {					// if flag enabled
		net.Manualpeers(config.GetConfigInstance().PeerFile)	// add peers from file
	}
}

/*
"startClient" Starts the client with the given initial configuration.
 */
func startClient(){
	config.InitializeConfiguration()                	// initializes the configuration from flags and defaults.
	config.BuildConfiguration()                     	// builds the proper structure on pc for core client to operate.
	config.PrintConfiguration()                     	// print the configuration the the cmd.
	manualPeers()										// check for manual peers
	logs.NewLog("Started Client ", logs.InfoLevel,logs.JSONLogFormat) 	// log start message
	rpc.RunAPIEndpoints()                           	// runs the server endpoints for client and relay api.
	session.ServeAndListen("3333","localhost")
	fmt.Print("Press any key + 'Return' to quit: ") 	// prompt user to exit
	input := bufio.NewScanner(os.Stdin)             	// unnecessary temporary entry
	input.Scan()                                    	// wait
}
