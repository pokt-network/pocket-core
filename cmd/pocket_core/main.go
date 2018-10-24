// Pocket Core: This is the starting point of the CLI.
package main

import (
	"github.com/pocket_network/pocket-core/cmd/util"
	"github.com/pocket_network/pocket-core/rpc"
)

//TODO add logging

// "Main" is the starting function of the client.
// Keep main as light as possible by calling accessory functions.
func main() {
	util.ParseFlags()
	util.PrintClientInfo()
	// TODO change to flag function
	rpc.StartRPC("8080")
}
