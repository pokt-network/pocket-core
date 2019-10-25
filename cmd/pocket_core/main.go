// This package is the starting point of Pocket Core.
package main

import (
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc"
	"github.com/pokt-network/pocket-core/util"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// "init" is a built in function that is automatically called before main.
func init() {
	// generates seed for randomization
	rand.Seed(time.Now().UTC().UnixNano())
}

// "main" is the starting function of the client.
func main() {
	startClient()
}

// "startClient" Starts the client with the given initial configuration.
func startClient() {

	// We trap kill signals (2,3,15,9)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Kill,
		os.Interrupt)

	// initializes the configuration from flags and defaults
	config.Init()
	// builds the proper structure on pc for core client to operate
	config.Build()
	// runs the server endpoints for client and relay api
	rpc.StartServers()
	// logs the client starting
	logs.Log("Started Pocket Core", logs.InfoLevel, logs.JSONLogFormat)
	// logs start message to stdout
	logs.Log("Started Pocket Core", logs.InfoLevel, logs.TextLogFormatter)
	// print the configuration
	config.Print()
	// Exits gracefully printing the exit signal
	defer func() {
		sig := <-signalChannel

		message := fmt.Sprintf("Exit signal %s received\n", sig)
		util.ExitGracefully(message)
	}()
}
