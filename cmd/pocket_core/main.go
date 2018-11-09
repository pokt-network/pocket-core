// Pocket Core: This is the starting point of the CLI.
package main

import (
	"github.com/logmatic/logmatic-go"
	"github.com/pocket_network/pocket-core/config"
	"github.com/pocket_network/pocket-core/logs"
	"github.com/pocket_network/pocket-core/rpc"
	"github.com/sirupsen/logrus"
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
	logs.LogConstructorAndLog("test.log","pocket-core/cmd/pocket_core/main.go",
		"27", "Testing testing 123", logrus.DebugLevel, &logmatic.JSONFormatter{})
	config.InitializeConfiguration()
	config.PrintConfiguration()
	rpc.RunAPIEndpoints()
}
