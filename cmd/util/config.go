// This package is for command line specific utilities.
package util

import (
	"encoding/json"
	"fmt"
	"github.com/pocket_network/pocket-core/const"
	"github.com/pocket_network/pocket-core/rpc"
)

// This file describes all of the configuration properties of the client (set by startup flags)

type configuration struct {
	Clientid      string `json:"CLIENTID"` // This variable holds a client identifier string.
	Version       string `json:"VERSION"`  // This variable holds the client version string.
	Datadir       string `json:"DATADIR"`  // This variable holds the working directory string.
	Clientrpc     bool   `json:"CRPC"`     // This variable describes if the client rpc is running.
	Clientrpcport string `json:"CRPCPORT"` // This variable holds the client rpc port string.
	Relayrpc      bool   `json:"RRPC"`     // This variable describes if the relay rpc is running.
	Relayrpcport  string `json:"RRPCPORT"` // This variable holds the relay rpc port string.
}

// "config" is a pointer to the configuration object.
var config configuration

/*
"NewConfiguration() is a Constructor function of the configuration type.
 */
func NewConfiguration(dd string, crpc bool, rrpc bool, crpcport string, rrpcport string) {
	if (config == configuration{}) {		// This checks for empty configuration.
		config = configuration{
			_const.CLIENTID,
			_const.VERSION,
			dd,
			crpc,
			crpcport,
			rrpc,
			rrpcport}
	}
}

/*
"runConfiguration" executes the specified configuration for the client.
 */
func runConfiguration(){
	if config.Clientrpc{
		go rpc.StartClientRPC(config.Clientrpcport)
	}
	if config.Relayrpc {
		rpc.StartRelayRPC(config.Relayrpcport) // TODO convert to go routine
	}
}

/*
"PrintConfiguration()" prints the client configuration information to the CLI.
 */
func printConfiguration() {
	data, _:= json.MarshalIndent(config, "","    ")
	fmt.Println("Pocket Core Configuration:\n",string(data))
}

/*
"GetConfig" returns the client configuration
 */
func GetConfig() configuration {
	return config
}

// TODO edit configuration settings ex: SetRRPCPort(port string){}
