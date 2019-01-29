// This package maintains the client configuration.
package config

import (
	"log"
	"os"

	"github.com/pokt-network/pocket-core/const"
)

// "Build" builds the configuration structure needed for the client.
func Build() {
	// builds the data directory on the local machine
	dataDir()
	// builds the logs directory within the data directory
	logsDir()
}

// "dataDir" builds the directory for program files.
func dataDir() {
	// attempts to make the data directory.
	if err := os.MkdirAll(Get().DD, os.ModePerm); err != nil {
		// doesn't use custom logs, because they may or may not be available at this point
		log.Fatalf(err.Error())
	}
}

// "logsDir" builds the directory for logs.
func logsDir() {
	// attempts to make the logs directory
	if err := os.MkdirAll(Get().DD+_const.FILESEPARATOR+"logs", os.ModePerm); err != nil {
		log.Fatal(err.Error())
	}
}
