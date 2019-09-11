package util

import (
	"os"
	"os/signal"

	"github.com/pokt-network/pocket-core/logs"
)

func ExitGracefully(message string) {
	logs.Log("Shutting down Pocket Core: "+message, logs.InfoLevel, logs.JSONLogFormat)
	logs.Log("Shutting down Pocket Core: "+message, logs.InfoLevel, logs.TextLogFormatter)

	os.Exit(0)
}

func WaitForExit() {
	// Catches OS system interrupt signal and calls unregister
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	select {
	case sig := <-c:
		ExitGracefully(sig.String() + " command executed.")
	}
}
