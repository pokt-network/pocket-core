package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/pokt-network/pocket-core/const"
)

// TODO configuration updating
type config struct {
	GID      string `json:"GID"`          // This variable holds self.GID.
	CID      string `json:"CLIENTID"`     // This variable holds a client identifier string.
	Ver      string `json:"VERSION"`      // This variable holds the client version string.
	DD       string `json:"DATADIR"`      // This variable holds the working directory string.
	CRPC     bool   `json:"CRPC"`         // This variable describes if the client rpc is running.
	CRPCPort string `json:"CRPCPORT"`     // This variable holds the client rpc port string.
	RRPC     bool   `json:"RRPC"`         // This variable describes if the relay rpc is running.
	RRPCPort string `json:"RRPCPort"`     // This variable holds the relay rpc port string.
	CFile    string `json:"HOSTEDCHAINS"` // This variable holds the filepath to the chains.json.
	PFile    string `json:"PEERFILE"`     // This variable holds the filepath to the peerFile.json.
	SNWL     string `json:"SNWL"`         // This variable holds the filepath to the service_whitelist.json.
	DWL      string `json:"DWL"`          // This variable holds the filepath to the developer_whitelist.json
}

var (
	c        *config
	once     sync.Once
	gid      = flag.String("gid", "GID1", "set the selfNode.GID for pocket core mvp")
	dd       = flag.String("dd", _const.DATADIR, "setup the data directory for the DB and keystore")
	rRpcPort = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	cFile    = flag.String("cfile", _const.CHAINSFILENAME, "specifies the filepath for chains.json")
	pFile    = flag.String("pfile", _const.PEERFILENAME, "specifies the filepath for peers.json")
	snwl     = flag.String("sFile", _const.SNWLFILENAME, "specifies the filepath for service_whitelist.json")
	dwl      = flag.String("dFile", _const.DWLFILENAME, "specifies the filepath for developer_whitelist.json")
	cRpcPort = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	cRpc     = flag.Bool("clientrpc", true, "whether or not to start the rpc server")
	rRpc     = flag.Bool("relayrpc", true, "whether or not to start the rpc server")
)

// "Init" initializes the configuration object.
func Init() {
	// built in function to parse the flags above.
	flag.Parse()
	// returns the thread safe c of the client configuration.
	GetInstance()
}

// "newConfiguration() is a constructor function of the configuration type.
func newConfiguration() {
	c = &config{
		*gid,
		_const.CLIENTID,
		_const.VERSION,
		*dd,
		*cRpc,
		*cRpcPort,
		*rRpc,
		*rRpcPort,
		*cFile,
		*pFile,
		*snwl,
		*dwl}
}

// "Print()" prints the client configuration information to the CLI.
func Print() {
	data, _ := json.MarshalIndent(c, "", "    ")              // pretty configure the json data
	fmt.Println("Pocket Core Configuration:\n", string(data)) // pretty print the pocket configuration
}

// "GetInstance()" returns the configuration object in a thread safe manner.
func GetInstance() *config { // singleton structure to return the configuration object
	once.Do(func() { // thread safety.
		newConfiguration()
	})
	return c // return the configuration
}
