package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/pokt-network/pocket-core/const"
)

// TODO configuration updating through CLI
type config struct {
	CID      string `json:"clientid"`     // This variable holds a client identifier string.
	Ver      string `json:"version"`      // This variable holds the client version string.
	DD       string `json:"datadir"`      // This variable holds the working directory string.
	CRPC     bool   `json:"crpc"`         // This variable describes if the client rpc is running.
	CRPCPort string `json:"crpcport"`     // This variable holds the client rpc port string.
	RRPC     bool   `json:"rrpc"`         // This variable describes if the relay rpc is running.
	RRPCPort string `json:"rrpcport"`     // This variable holds the relay rpc port string.
	CFile    string `json:"hostedchains"` // This variable holds the filepath to the chains.json.
}

var (
	c        *config
	once     sync.Once
	dd       = flag.String("datadirectory", _const.DATADIR, "setup the data directory for the DB and keystore")
	rRpcPort = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	cFile    = flag.String("cfile", _const.CHAINSFILENAME, "specifies the filepath for chains.json")
	cRpcPort = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	cRpc     = flag.Bool("clientrpc", true, "whether or not to start the rpc server")
	rRpc     = flag.Bool("relayrpc", true, "whether or not to start the rpc server")
)

// "Init" initializes the configuration object.
func Init() {
	// built in function to parse the flags above.
	flag.Parse()
	// sets up filepaths for config files
	filePaths()
	// returns the thread safe c of the client configuration.
	GlobalConfig()
}

// "Print()" prints the client configuration information to the CLI.
func Print() {
	data, _ := json.MarshalIndent(c, "", "    ")
	fmt.Println("Pocket Core Configuration:\n", string(data))
}

// "GlobalConfig()" returns the configuration object in a thread safe manner.
func GlobalConfig() *config { // singleton structure to return the configuration object
	once.Do(func() { // thread safety.
		newConfiguration()
	})
	return c // return the configuration
}

// "newConfiguration() is a constructor function of the configuration type.
func newConfiguration() {
	c = &config{
		_const.CLIENTID,
		_const.VERSION,
		*dd,
		*cRpc,
		*cRpcPort,
		*rRpc,
		*rRpcPort,
		*cFile}
}
