package config

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"

	"github.com/pokt-network/pocket-core/const"
)

type config struct {
	CID         string `json:"clientid"`     // This variable holds a client identifier string.
	Ver         string `json:"version"`      // This variable holds the client version string.
	DD          string `json:"datadir"`      // This variable holds the working directory string.
	CRPC        bool   `json:"crpc"`         // This variable describes if the client rpc is running.
	CRPCPort    string `json:"crpcport"`     // This variable holds the client rpc port string.
	RRPC        bool   `json:"rrpc"`         // This variable describes if the relay rpc is running.
	RRPCPort    string `json:"rrpcport"`     // This variable holds the relay rpc port string.
	CFile       string `json:"hostedchains"` // This variable holds the filepath to the chains.json.
	LogFormat   string `json:"logformat"`    // This variable changes the log storage format
	LogSize     int    `json:"logsize"`      // This variable sets the log file size specified in megabytes.
	LogAge      int    `json:"logage"`       // This variable sets the log file age specified in days.
	LogBackups  int    `json:"logbackups"`   // This variable sets the log file backup quantity.
	LogCompress bool   `json:"logcompress"`  // This variable enabled the log compression.

}

var (
	c           *config
	once        sync.Once
	dd          = flag.String("datadirectory", _const.DATADIR, "setup the data directory for the DB and keystore")
	rRpcPort    = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	cFile       = flag.String("cfile", _const.CHAINSFILENAME, "specifies the filepath for chains.json")
	cRpcPort    = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	cRpc        = flag.Bool("clientrpc", true, "whether or not to start the rpc server")
	rRpc        = flag.Bool("relayrpc", true, "whether or not to start the rpc server")
	logFormat   = flag.String("logformat", ".json", "Log format for storing logs.  ex.:[.json, .log] ('.json' is used by default)")
	logSize     = flag.Int("logsize", 500, "Log size of logs before rotating. Size is given in Megabytes. ex.: [500 = 500MB] ('500' used by default)")
	logAge      = flag.Int("logage", 15, "Log rotate age. Age is given in days. ex.: [20 = 20 days] ('15' used by default)")
	logBackups  = flag.Int("logbackups", 10, "Log backup quantity. Quantity of logfiles to keep before deleting. ex.: [20 = 20 backups files] ('15' used by default)")
	logCompress = flag.Bool("logcompress", true, "Enables log compress. Quantity of logfiles to keep before deleting. ex.: [20 = 20 backups files] ('15' used by default)")
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

// A simple log for showing the pocket configuration
func logger(output string) {

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true

	log.SetFormatter(Formatter)

	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	log.Info(output)

}

// "Print()" prints the client configuration information to the CLI.
func Print() {
	data, _ := json.MarshalIndent(c, "", "    ")
	var output = fmt.Sprintf("%s", string(data))
	logger(output)
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
		*cFile,
		*logFormat,
		*logSize,
		*logAge,
		*logBackups,
		*logCompress}
}
