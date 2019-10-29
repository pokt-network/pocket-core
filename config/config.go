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
	CID       string `json:"clientid"`     // This variable holds a client identifier string.
	Ver       string `json:"version"`      // This variable holds the client version string.
	DD        string `json:"datadir"`      // This variable holds the working directory string.
	CRPC      bool   `json:"crpc"`         // This variable describes if the client rpc is running.
	CRPCPort  string `json:"crpcport"`     // This variable holds the client rpc port string.
	RRPC      bool   `json:"rrpc"`         // This variable describes if the relay rpc is running.
	RRPCPort  string `json:"rrpcport"`     // This variable holds the relay rpc port string.
	CFile     string `json:"hostedchains"` // This variable holds the filepath to the chains.json.
	LogFormat string `json:"logformat"`    // This variable changes the log storage format.
	LogLevel  string `json:"loglevel"`     // This variable changes the log level.
	LogDir    string `json:"logdir"`       // This variable changes the log storage dir.

}

var (
	c         *config
	once      sync.Once
	dd        = flag.String("datadirectory", _const.DATADIR, "setup the data directory for the DB and keystore")
	rRpcPort  = flag.String("relayrpcport", "8081", "specified port to run relay rpc")
	cFile     = flag.String("cfile", _const.CHAINSFILENAME, "specifies the filepath for chains.json")
	cRpcPort  = flag.String("clientrpcport", "8080", "specified port to run client rpc")
	cRpc      = flag.Bool("clientrpc", true, "whether or not to start the rpc server")
	rRpc      = flag.Bool("relayrpc", true, "whether or not to start the rpc server")
	logFormat = flag.String("logformat", "", "Log format for storing logs.  ex.:[.json, .log] ('.json' is used by default)")
	logLevel  = flag.String("loglevel", "INFO", "Log level.  ex.:[TRACE, DEBUG, INFO, ERROR, FATAL, PANIC] ('INFO' is used by default)")
	logDir    = flag.String("logdir", _const.DATADIR+_const.FILESEPARATOR+"logs"+_const.FILESEPARATOR, "setup the log directory.  ex.:['/var/log/'] ('~/.pocket/logs/' is used by default)")
)

// "Init" initializes the configuration object.
func Init() {
	// built in function to parse the flags above.
	flag.Parse()

	// In case user sets logdir flag without logformat, we assume logformat to be .json by default
	if isFlagPassed("logdir") == true {
		if isFlagPassed("logformat") == false {
			flag.Set("logformat", ".json")
		}
	}
	// Check if valid loglevel is passed 
	if isFlagPassed("loglevel") == true {
		loglevel := fmt.Sprintf("%s", flag.Lookup("loglevel").Value)
		_, err := log.ParseLevel(loglevel)

		// Throw fatal error and exit if invalid loglevel is received
		if err != nil {
			logger("loglevel flag not valid", log.FatalLevel)
		}
	}
	// sets up filepaths for config files
	filePaths()
	// returns the thread safe c of the client configuration.
	GlobalConfig()
}

// A simple log for showing the pocket configuration
func logger(output string, level log.Level) {

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true

	log.SetFormatter(Formatter)

	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	switch level {
	case log.InfoLevel:
		log.Info(output)
	case log.DebugLevel:
		log.Debug(output)
	case log.FatalLevel:
		log.Fatal(output)
	case log.PanicLevel:
		log.Panic(output)
	default:
		log.Info(output)
	}


}

// "Print()" prints the client configuration information to the CLI.
func Print() {
	data, _ := json.MarshalIndent(c, "", "    ")
	var output = fmt.Sprintf("%s", string(data))
	logger(output, log.InfoLevel)
}

// Validate if a flag has value
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
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
		*logLevel,
		*logDir}
}
