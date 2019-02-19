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
	Dispatch bool   `json:"DISPATCH"`     // This variable describes whether or not this node is a dispatcher
	DisMode  int    `json:"DISMODE"`      // The mode by which the dispatch runs in (NORM, MIGRATE, DEPCRECATED)
	DBEND    string `json:"DBENDPOINT"`   // The endpoint of the centralized database for dispatch configuration
	DisIP    string `json:"DISIP"`        // The IP address of the centralized dispatcher
	DisCPort string `json:"DISCPort"`     // The client port of the centralized dispatcher
	DisRPort string `json:"DISCPort"`     // The relay port of the centralized dispatcher
}

var (
	c        *config
	once     sync.Once
	gid      = flag.String("gid", "GID1", "set the selfNode.GID for pocket core mvp")
	dd       = flag.String("dd", _const.DATADIR, "setup the data directory for the DB and keystore")
	rRpcPort = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	cFile    = flag.String("cfile", _const.CHAINSFILENAME, "specifies the filepath for chains.json")
	pFile    = flag.String("pfile", _const.PEERFILENAME, "specifies the filepath for peers.json")
	snwl     = flag.String("sfile", _const.SNWLFILENAME, "specifies the filepath for service_whitelist.json")
	dwl      = flag.String("dfile", _const.DWLFILENAME, "specifies the filepath for developer_whitelist.json")
	cRpcPort = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	cRpc     = flag.Bool("clientrpc", true, "whether or not to start the rpc server")
	rRpc     = flag.Bool("relayrpc", true, "whether or not to start the rpc server")
	dispatch = flag.Bool("dispatch", false, "specifies if this node is operating as a dispatcher")
	dismode  = flag.Int("dismode", _const.DISMODENORMAL, "specifies the mode by which the dispatcher is operating (0) Normal, (1) Migrate, (2) Deprecated")
	dbend    = flag.String("dbend", _const.DBENDPOINT, "specifies the database endpoint for the centralized dispatcher")
	disip    = flag.String("disip", _const.DISPATCHIP, "specifies the address of the centralized dispatcher")
	discport = flag.String("discport", _const.DISPATCHCLIENTPORT, "specifies the client port of the centralized dispatcher")
	disrport = flag.String("disrport", _const.DISPATCHRELAYPORT, "specifies the relay port of the centralized dispatcher")
)

// "Init" initializes the configuration object.
func Init() {
	// built in function to parse the flags above.
	flag.Parse()
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
		*dwl,
		*dispatch,
		*dismode,
		*dbend,
		*disip,
		*discport,
		*disrport}
}
