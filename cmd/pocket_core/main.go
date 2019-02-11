// This package is the starting point of Pocket Core.
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/db"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc"
)

// "init" is a built in function that is automatically called before main.
func init() {
	// generates seed for randomization
	crypto.GenerateSeed()
}

// "main" is the starting function of the client.
func main() {
	startClient()
}

// "startClient" Starts the client with the given initial configuration.
func startClient() {
	// initializes the configuration from flags and defaults
	config.Init()
	// builds the proper structure on pc for core client to operate
	config.Build()
	// builds node structures from files
	node.Files()
	// print the configuration the the cmd
	config.Print()
	// add peers to dispatch structure
	node.PeerList().CopyToDP()
	// check for hosted chains
	node.TestChains()
	// runs the server endpoints for client and relay api
	rpc.StartServers()
	// runs servers for messages
	message.StartServers()
	// run db refresh on peers (if dispatch node)
	db.PeersRefresh()
	// sends an entry message to the centralized dispatcher
	node.Register()
	// logs the client starting
	logs.NewLog("Started Pocket Core", logs.InfoLevel, logs.JSONLogFormat)
	// prompt user to exit
	fmt.Print("Press any key + 'Return' to quit: ")
	// wait for input
	bufio.NewScanner(os.Stdin).Scan()
	// send exit message to dispatcher
	node.UnRegister(0)
}
