package node

import (
	"fmt"
	"github.com/pokt-network/pocket-core/logs"
	"os"
)

// "ExitGracefully" is the shutdown sequece of Pocket Core
func ExitGracefully(message string) {
	// unregister from the network
	if err := UnRegister(0); err != nil {
		logs.NewLog("Shutting down Pocket Core: "+err.Error(), logs.InfoLevel, logs.JSONLogFormat)
		fmt.Fprint(os.Stderr, "\nShutting down Pocket Core: "+err.Error())
		os.Exit(1)
	}
	logs.NewLog("Shutting down Pocket Core: "+message, logs.InfoLevel, logs.JSONLogFormat)
	fmt.Fprint(os.Stdout, "Shutting down Pocket Core: "+message)
	// Exit honoring deferred calls
	os.Exit(3)
}

// "WaitForExit" listens for interrupt signal and calls unregister
func WaitForExit(message string) {
	ExitGracefully(message)
}
