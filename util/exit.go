package util

import (
	"fmt"
	"os"
	
	"github.com/pokt-network/pocket-core/logs"
)

func ExitGracefully(message string) {
	logs.NewLog("Shutting down Pocket Core: "+message, logs.InfoLevel, logs.JSONLogFormat)
	fmt.Println("Shutting down Pocket Core: " + message)
	os.Exit(0)
}
