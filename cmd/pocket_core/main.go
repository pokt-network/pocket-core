// This package is the starting point of Pocket Core.
package main

import (
	"fmt"
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

	defer func() {
		sig := <-signalChannel
		message := fmt.Sprintf("Exit signal %s received\n", sig)
		fmt.Println(message)
		os.Exit(3)
	}()
}
