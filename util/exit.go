package util

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pokt-network/pocket-core/logs"
)

func ExitGracefully(message string) {
	logs.NewLog("Shutting down Pocket Core: "+message, logs.InfoLevel, logs.JSONLogFormat)
	fmt.Fprint(os.Stdout, "Shutting down Pocket Core: "+message)
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
