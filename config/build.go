// This package maintains the client configuration.
package config

import (
	"github.com/pokt-network/pocket-core/const"
	"log"
	"os"
)

// "build.go" is for building the Pocket Core configuration

/*
This function builds the configuration needed for the client.
 */
func BuildConfiguration() {
	buildDataDir()													// builds the data directory on the local machine
	buildLogsDir()													// builds the logs directory within the datadirectory
}

/*
This function builds the directory needed for DB and keystore etc.
 */
func buildDataDir() {
	err := os.MkdirAll(GetConfigInstance().Datadir, os.ModePerm) 	// attempts to make the data directory.
	if err != nil {													// if unable to create custom data directory
		log.Fatal(err.Error())										// notice use of built in log constructor
	}                                                            	// if the data directory isn't built then no use
}																	// using custom logs.

func buildLogsDir() {
	err := os.MkdirAll(GetConfigInstance().Datadir+ 			// attempts to make the logs directory
		_const.FILESEPARATOR+"logs", os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())										// log if error
	}
}
