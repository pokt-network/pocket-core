// Pocket Core: This is the starting point of the CLI.
package main

import (
	"bufio"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc"
	"os"
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
func startClient() {
	config.InitializeConfiguration()
	config.BuildConfiguration()
	config.PrintConfiguration()
	logs.LogConstructorAndLog("TESTING TESTING", logs.InfoLevel, logs.JSONLogFormat)
	rpc.RunAPIEndpoints()
	fmt.Print("Press any key + 'Return' to quit: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
