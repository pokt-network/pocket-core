package util

import (
	"os"

	"github.com/pokt-network/pocket-core/logs"
)

func ExitGracefully(message string) {
	logs.Log("Shutting down Pocket Core: "+message, logs.InfoLevel, logs.JSONLogFormat)
	logs.Log("Shutting down Pocket Core: "+message, logs.InfoLevel, logs.TextLogFormatter)

	os.Exit(3)
}
