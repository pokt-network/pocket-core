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
	GID            string `json:"GID"`            // This variable holds self.GID.
	IP             string `json:"IP"`             // This variable holds the ip of the client
	Port           string `json:"PORT"`           // The public service port that will be displayed to the clients
	CID            string `json:"CLIENTID"`       // This variable holds a client identifier string.
	Ver            string `json:"VERSION"`        // This variable holds the client version string.
	DD             string `json:"DATADIR"`        // This variable holds the working directory string.
	RRPC           bool   `json:"RRPC"`           // This variable describes if the relay rpc is running.
	RRPCPort       string `json:"RRPCPort"`       // This variable holds the relay rpc port string.
	CFile          string `json:"HOSTEDCHAINS"`   // This variable holds the filepath to the chains.json.
	SNWL           string `json:"SNWL"`           // This variable holds the filepath to the service_whitelist.json.
	DWL            string `json:"DWL"`            // This variable holds the filepath to the developer_whitelist.json
	Dispatch       bool   `json:"DISPATCH"`       // This variable describes whether or not this node is a dispatcher
	DisMode        int    `json:"DISMODE"`        // The mode by which the dispatch runs in (NORM, MIGRATE, DEPCRECATED)
	DBEndpoint     string `json:"DBENDPOINT"`     // The endpoint of the centralized database for dispatch configuration
	DBTableName    string `json:"DBTABLE"`        // The table name of the centralized dispatcher database
	DBRegion       string `json:"DBREGION"`       // The aws reigion for the centralized datablase for dispatch configuration
	DisIP          string `json:"DISIP"`          // The IP address of the centralized dispatcher
	DisRPort       string `json:"DISRPort"`       // The relay port of the centralized dispatcher
	PRefresh       int    `json:"PREFRESH"`       // The peer refresh time for the centralized dispatcher in seconds
	RequestTimeout int    `json:"REQUESTTIMEOUT"` // The timeout for http requests
}

var (
	c              *config
	once           sync.Once
	gid            = flag.String("gid", "GID1", "set the self GID prefix for pocket core mvp node")
	ip             = flag.String("ip", _const.DEFAULTIP, "set the IP address of the pocket core mvp node, if not set, uses public ip")
	port           = flag.String("port", _const.DEFAULTPORT, "set the publicly displayed servicing port")
	dd             = flag.String("datadirectory", _const.DATADIR, "setup the data directory for the DB and keystore")
	rRpcPort       = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	cFile          = flag.String("cfile", _const.CHAINFILEPLACEHOLDER, "specifies the filepath for chains.json")
	snwl           = flag.String("sfile", _const.SNWLFILENAMEPLACEHOLDER, "specifies the filepath for service_whitelist.json")
	dwl            = flag.String("dfile", _const.DWLFILENAMEPLACEHOLDER, "specifies the filepath for developer_whitelist.json")
	rRpc           = flag.Bool("relayrpc", true, "whether or not to start the rpc server")
	dispatch       = flag.Bool("dispatch", false, "specifies if this node is operating as a dispatcher")
	dismode        = flag.Int("dismode", _const.DISMODENORMAL, "specifies the mode by which the dispatcher is operating (0) Normal, (1) Migrate, (2) Deprecated")
	dbend          = flag.String("dbend", _const.DBENDPOINT, "specifies the database endpoint for the centralized dispatcher")
	dbtable        = flag.String("dbtable", _const.DBTABLENAME, "specifies the database tablename for the centralized dispatcher")
	dbregion       = flag.String("dbregion", _const.DBREGION, "specifies the region of the db for the centralized dispatcher")
	disip          = flag.String("disip", _const.DISPATCHIP, "specifies the address of the centralized dispatcher")
	disrport       = flag.String("disrport", _const.DISPATCHRELAYPORT, "specifies the relay port of the centralized dispatcher")
	peerrefresh    = flag.Int("peerrefresh", _const.DBREFRESH, "specifies the peer refresh time for the centralized dispatcher liveness checks")
	requestTimeout = flag.Int("requestTimeout", _const.TIMEOUT, "specifies the timeout for http requests (ms)")
)

// "Init" initializes the configuration object.
func Init() {
	// built in function to parse the flags above.
	flag.Parse()
	// generates filepaths from data directory flag
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
		*gid,
		*ip,
		*port,
		_const.CLIENTID,
		_const.VERSION,
		*dd,
		*rRpc,
		*rRpcPort,
		*cFile,
		*snwl,
		*dwl,
		*dispatch,
		*dismode,
		*dbend,
		*dbtable,
		*dbregion,
		*disip,
		*disrport,
		*peerrefresh,
		*requestTimeout}
}
