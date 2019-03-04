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
	GID         string `json:"GID"`          // This variable holds self.GID.
	IP          string `json:"IP"`           // This variable holds the ip of the client
	Port        string `json:"PORT"`         // The public service port that will be displayed to the clients
	CID         string `json:"CLIENTID"`     // This variable holds a client identifier string.
	Ver         string `json:"VERSION"`      // This variable holds the client version string.
	DD          string `json:"DATADIR"`      // This variable holds the working directory string.
	CRPC        bool   `json:"CRPC"`         // This variable describes if the client rpc is running.
	CRPCPort    string `json:"CRPCPORT"`     // This variable holds the client rpc port string.
	RRPC        bool   `json:"RRPC"`         // This variable describes if the relay rpc is running.
	RRPCPort    string `json:"RRPCPort"`     // This variable holds the relay rpc port string.
	CFile       string `json:"HOSTEDCHAINS"` // This variable holds the filepath to the chains.json.
	PFile       string `json:"PEERFILE"`     // This variable holds the filepath to the peerFile.json.
	SNWL        string `json:"SNWL"`         // This variable holds the filepath to the service_whitelist.json.
	DWL         string `json:"DWL"`          // This variable holds the filepath to the developer_whitelist.json
	Dispatch    bool   `json:"DISPATCH"`     // This variable describes whether or not this node is a dispatcher
	DisMode     int    `json:"DISMODE"`      // The mode by which the dispatch runs in (NORM, MIGRATE, DEPCRECATED)
	DBEndpoint  string `json:"DBENDPOINT"`   // The endpoint of the centralized database for dispatch configuration
	DBTableName string `json:"DBTABLE"`      // The table name of the centralized dispatcher database
	DBRegion    string `json:"DBREGION"`     // The aws reigion for the centralized datablase for dispatch configuration
	DisIP       string `json:"DISIP"`        // The IP address of the centralized dispatcher
	DisCPort    string `json:"DISCPort"`     // The client port of the centralized dispatcher
	DisRPort    string `json:"DISCPort"`     // The relay port of the centralized dispatcher
	PRefresh    int    `json:"PREFRESH"`     // The peer refresh time for the centralized dispatcher in seconds
}

var (
	c           *config
	once        sync.Once
	gid         = flag.String("gid", "GID1", "set the self GID prefix for pocket core mvp node")
	ip          = flag.String("ip", _const.DEFAULTIP, "set the IP address of the pocket core mvp node, if not set, uses public ip")
	port        = flag.String("port", _const.DEFAULTPORT, "set the publicly displayed servicing port")
	dd          = flag.String("datadirectory", _const.DATADIR, "setup the data directory for the DB and keystore")
	rRpcPort    = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	cFile       = flag.String("cfile", _const.CHAINSFILENAME, "specifies the filepath for chains.json")
	pFile       = flag.String("pfile", _const.PEERFILENAME, "specifies the filepath for peers.json")
	snwl        = flag.String("sfile", _const.SNWLFILENAME, "specifies the filepath for service_whitelist.json")
	dwl         = flag.String("dfile", _const.DWLFILENAME, "specifies the filepath for developer_whitelist.json")
	cRpcPort    = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	cRpc        = flag.Bool("clientrpc", true, "whether or not to start the rpc server")
	rRpc        = flag.Bool("relayrpc", true, "whether or not to start the rpc server")
	dispatch    = flag.Bool("dispatch", false, "specifies if this node is operating as a dispatcher")
	dismode     = flag.Int("dismode", _const.DISMODENORMAL, "specifies the mode by which the dispatcher is operating (0) Normal, (1) Migrate, (2) Deprecated")
	dbend       = flag.String("dbend", _const.DBENDPOINT, "specifies the database endpoint for the centralized dispatcher")
	dbtable     = flag.String("dbtable", _const.DBTABLENAME, "specifies the database tablename for the centralized dispatcher")
	dbregion    = flag.String("dbregion", _const.DBREGION, "specifies the region of the db for the centralized dispatcher")
	disip       = flag.String("disip", _const.DISPATCHIP, "specifies the address of the centralized dispatcher")
	discport    = flag.String("discport", _const.DISPATCHCLIENTPORT, "specifies the client port of the centralized dispatcher")
	disrport    = flag.String("disrport", _const.DISPATCHRELAYPORT, "specifies the relay port of the centralized dispatcher")
	peerrefresh = flag.Int("peerrefresh", _const.DBREFRESH, "specifies the peer refresh time for the centralized dispatcher liveness checks")
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
		*ip,
		*port,
		_const.CLIENTID,
		_const.VERSION,
		*dd,
		*cRpc,
		*cRpcPort,
		*rRpc,
		*rRpcPort,
		*dd + _const.FILESEPARATOR + *cFile,
		*dd + _const.FILESEPARATOR + *pFile,
		*dd + _const.FILESEPARATOR + *snwl,
		*dd + _const.FILESEPARATOR + *dwl,
		*dispatch,
		*dismode,
		*dbend,
		*dbtable,
		*dbregion,
		*disip,
		*discport,
		*disrport,
		*peerrefresh}
}
