package util

import (
  "fmt"
	"github.com/pokt-network/pocket-core/logs"
	"os"
)

func ExitGracefully(message string){ //TODO exit gracefully function
  logs.NewLog("Shutting down Pocket Core because " + message, logs.InfoLevel, logs.JSONLogFormat)
  fmt.Println("Shutting down Pocket Core because " + message)
  os.Exit(0)
}
