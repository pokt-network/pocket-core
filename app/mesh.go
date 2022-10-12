package app

import (
	"encoding/json"
	sdk "github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	log2 "log"
	"os"
)

var (
	GlobalMeshConfig MeshConfig
)

type MeshConfig struct {
	// Mesh Node
	DataDir                string `json:"data_dir"`
	RPCPort                string `json:"rpc_port"`
	ChainsName             string `json:"chains_name"`
	RPCTimeout             int64  `json:"rpc_timeout"`
	LogLevel               string `json:"log_level"`
	UserAgent              string `json:"user_agent"`
	AuthTokenFile          string `json:"auth_token_file"`
	JSONSortRelayResponses bool   `json:"json_sort_relay_responses"`
	// Prometheus
	PrometheusAddr         string `json:"pocket_prometheus_port"`
	PrometheusMaxOpenfiles int    `json:"prometheus_max_open_files"`
	// Cache
	RelayCacheFile   string `json:"relay_cache_file"`
	SessionCacheFile string `json:"session_cache_file"`
	// Workers
	WorkerStrategy     string `json:"worker_strategy"`
	MaxWorkers         int    `json:"max_workers"`
	MaxWorkersCapacity int    `json:"max_workers_capacity"`
	WorkersIdleTimeout int    `json:"workers_idle_timeout"`
	// Servicer
	ServicerPrivateKeyFile string `json:"servicer_private_key_file"`
	ServicerURL            string `json:"servicer_url"`
	ServicerRPCTimeout     int64  `json:"servicer_rpc_timeout"`
	ServicerAuthTokenFile  string `json:"servicer_auth_token_file"`
	ServicerRetryMaxTimes  int    `json:"servicer_retry_max_times"`
	ServicerRetryWaitMin   int    `json:"servicer_retry_wait_min"`
	ServicerRetryWaitMax   int    `json:"servicer_retry_wait_max"`
}

func defaultMeshConfig(dataDir string) MeshConfig {
	c := MeshConfig{
		// Mesh Node
		DataDir:                dataDir,
		RPCPort:                sdk.DefaultRPCPort,
		ChainsName:             sdk.DefaultChainsName,
		RPCTimeout:             sdk.DefaultRPCTimeout,
		LogLevel:               "*:info, *:error",
		UserAgent:              sdk.DefaultUserAgent,
		AuthTokenFile:          "auth" + FS + "mesh.json",
		JSONSortRelayResponses: sdk.DefaultJSONSortRelayResponses,
		// Prometheus
		PrometheusAddr:         sdk.DefaultPocketPrometheusListenAddr,
		PrometheusMaxOpenfiles: sdk.DefaultPrometheusMaxOpenFile,
		// Cache
		RelayCacheFile:   "data" + FS + "relays.pkt",
		SessionCacheFile: "data" + FS + "session.pkt",
		// Worker
		WorkerStrategy:     "balanced",
		MaxWorkers:         10,
		MaxWorkersCapacity: 1000,
		WorkersIdleTimeout: 100,
		// Servicer
		ServicerPrivateKeyFile: "key" + FS + "key.json",
		ServicerURL:            sdk.DefaultRemoteCLIURL,
		ServicerRPCTimeout:     sdk.DefaultRPCTimeout,
		ServicerAuthTokenFile:  "auth" + FS + "servicer.json",
		ServicerRetryMaxTimes:  10,
		ServicerRetryWaitMin:   5,
		ServicerRetryWaitMax:   180,
	}

	return c
}

func InitMeshConfig(datadir string) {
	log2.Println("Initializing Pocket Datadir")
	// set up the codec
	MakeCodec()
	if datadir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log2.Fatal("could not get home directory for data dir creation: " + err.Error())
		}
		datadir = home + FS + sdk.DefaultDDName
		log2.Println("datadir = " + datadir)
	}
	c := defaultMeshConfig(datadir)
	// read from ccnfig file
	configFilepath := datadir + FS + sdk.ConfigDirName + FS + sdk.ConfigFileName
	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		log2.Println("no config file found... creating the datadir @ "+c.DataDir+FS+sdk.ConfigDirName, os.ModePerm)
		// ensure directory path made
		err = os.MkdirAll(c.DataDir+FS+sdk.ConfigDirName, os.ModePerm)
		if err != nil {
			log2.Fatal(err)
		}
	}
	var jsonFile *os.File
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			return
		}
	}(jsonFile)
	// if file exists open, else create and open
	if _, err := os.Stat(configFilepath); err == nil {
		jsonFile, err = os.OpenFile(configFilepath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log2.Fatalf("cannot open config json file: " + err.Error())
		}
		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log2.Fatalf("cannot read config file: " + err.Error())
		}
		err = json.Unmarshal(b, &c)
		if err != nil {
			log2.Fatalf("cannot read config file into json: " + err.Error())
		}
	} else if os.IsNotExist(err) {
		// if does not exist create one
		jsonFile, err = os.OpenFile(configFilepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			log2.Fatalf("canot open/create config json file: " + err.Error())
		}
		b, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			log2.Fatalf("cannot marshal default config into json: " + err.Error())
		}
		// write to the file
		_, err = jsonFile.Write(b)
		if err != nil {
			log2.Fatalf("cannot write default config to json file: " + err.Error())
		}
	}

	GlobalMeshConfig = c
}
