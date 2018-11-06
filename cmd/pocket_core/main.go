// Pocket Core: This is the starting point of the CLI.
package main

import (
	"github.com/pocket_network/pocket-core/config"
	"github.com/pocket_network/pocket-core/rpc"
)

//TODO add logging

/*
"main" is the starting function of the client.
 Keep main as light as possible by calling accessory functions.
*/
func main() {
	startClient()
}

/*
"startClient" Starts the client with the given initial configuration.
 */
func startClient(){
	config.InitializeConfiguration()
	config.PrintConfiguration()
	rpc.RunAPIEndpoints()
}
