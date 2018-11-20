// This package maintains the client configuration.
package config

import (
	"github.com/pokt-network/pocket-core/const"
	"log"
	"os"
)

/*
This file is for building the Pocket Core configuration
 */

/*
This function builds the configuration needed for the client.
 */
func BuildConfiguration() {
	buildDataDir()
	buildLogsDir()
}

/*
This function builds the directory needed for DB and keystore etc.
 */
func buildDataDir() {
	err := os.MkdirAll(GetInstance().Datadir, os.ModePerm);
	if err != nil {
		// If unable to write the folder... Probably unable to write this log file
		//logs.LogConstructorAndLog("Unable to create ",logs.ErrorLevel,logs.JSONLogFormat)
		// So redundantly log with built in logger to print and quit
		log.Fatal(err.Error())
	}
}

func buildLogsDir() {
	err := os.MkdirAll(GetInstance().Datadir+_const.FILESEPARATOR+"logs", os.ModePerm);
	if err != nil {
		// If unable to write the folder... Probably unable to write this log file
		//logs.LogConstructorAndLog("Unable to create ",logs.ErrorLevel,logs.JSONLogFormat)
		// So redundantly log with built in logger to print and quit
		log.Fatal(err.Error())
	}
}
