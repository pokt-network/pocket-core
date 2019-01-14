// This package is the starting point of the CLI.
package main

import (
	"bufio"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc"
	"os"
)

// "main.go" is the entry point of the client.

/*
"init" is a built in function that is automatically called before main.
*/
func init() {
	crypto.GenerateSeed() // generates seed for randomization
}

/*
"main" is the starting function of the client.
 Keep main as light as possible by calling accessory functions.
*/
func main() {
	startClient() // see function call below
}

/*
"manualPeers" parses peers from peers.json file
*/
func peersFromFile() {
	if err := node.ManualPeersFile(config.GetConfigInstance().PeerFile); err!=nil { // add peers from file
		// TODO handle error (note: if file doesn't exist this still should work)
	}
}

/*
"chainsFromFile" parses hosted chains from chains.json file
 */
func chainsFromFile(){
	if err := node.HostedChainsFile(config.GetConfigInstance().ChainsFilepath); err!=nil {
		// TODO handle error (note: if hosted chains file doesn't exist how to proceed?"
	}
}

/*
"startClient" Starts the client with the given initial configuration.
*/
func startClient() {
	config.InitializeConfiguration()                                  // initializes the configuration from flags and defaults.
	config.BuildConfiguration()                                       // builds the proper structure on pc for core client to operate.
	config.PrintConfiguration()                                       // print the configuration the the cmd.
	peersFromFile()                                                   // check for manual peers
	node.GetPeerList().AddPeersToDispatchStructure()                  // add peers to dispatch structure
	chainsFromFile()                                                  // check for chains.json file
	logs.NewLog("Started client", logs.InfoLevel, logs.JSONLogFormat) // log start message
	rpc.StartAPIServers()                                             // runs the server endpoints for client and relay api.
	message.RunMessageServers()                                       // runs servers for messages
	fmt.Print("Press any key + 'Return' to quit: ")                   // prompt user to exit
	input := bufio.NewScanner(os.Stdin)                               // unnecessary temporary entry
	input.Scan()                                                      // wait
}
