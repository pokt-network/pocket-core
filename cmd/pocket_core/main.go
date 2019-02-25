// This package is the starting point of Pocket Core.
package main

import (
	"os"
	"os/signal"

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
	node.ConfigFiles()
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
	// runs a check on all service nodes periodically
	db.CheckPeers()
	// sends an entry message to the centralized dispatcher
	node.Register()
	// logs the client starting
	logs.NewLog("Started Pocket Core", logs.InfoLevel, logs.JSONLogFormat)
	// Catches OS system interrupt signal and calls unregister
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	select {
	case sig := <-c:
		// Call util.ExitGracefully
	}
}
