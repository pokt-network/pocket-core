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
	DataDir                    string `json:"data_dir"`
	RPCPort                    string `json:"rpc_port"`
	ClientRPCTimeout           int64  `json:"client_rpc_timeout"`
	ClientRPCReadTimeout       int64  `json:"client_rpc_read_timeout"`
	ClientRPCReadHeaderTimeout int64  `json:"client_rpc_read_header_timeout"`
	ClientRPCWriteTimeout      int64  `json:"client_rpc_write_timeout"`
	LogLevel                   string `json:"log_level"`
	LogChainRequest            bool   `json:"log_chain_request"`
	LogChainResponse           bool   `json:"log_chain_response"`
	UserAgent                  string `json:"user_agent"`
	AuthTokenFile              string `json:"auth_token_file"`
	JSONSortRelayResponses     bool   `json:"json_sort_relay_responses"`

	// Chains
	ChainsName                  string `json:"chains_name"`
	ChainsNameMap               string `json:"chains_name_map"`
	RemoteChainsNameMap         string `json:"remote_chains_name_map"`
	ChainRPCTimeout             int64  `json:"chain_rpc_timeout"`
	ChainRPCMaxIdleConnections  int    `json:"chain_rpc_max_idle_connections"`
	ChainRPCMaxConnsPerHost     int    `json:"chain_rpc_max_conns_per_host"`
	ChainRPCMaxIdleConnsPerHost int    `json:"chain_rpc_max_idle_conns_per_host"`
	ChainDropConnections        bool   `json:"chain_drop_connections"`

	// Relay Cache
	RelayCacheFile                         string `json:"relay_cache_file"`
	RelayCacheBackgroundSyncInterval       int    `json:"relay_cache_background_sync_interval"`
	RelayCacheBackgroundCompactionInterval int    `json:"relay_cache_background_compaction_interval"`

	// Hot Reload Interval in milliseconds
	KeysHotReloadInterval   int `json:"keys_hot_reload_interval"`
	ChainsHotReloadInterval int `json:"chains_hot_reload_interval"`

	// Workers
	ServicerWorkerStrategy     string `json:"servicer_worker_strategy"`
	ServicerMaxWorkers         int    `json:"servicer_max_workers"`
	ServicerMaxWorkersCapacity int    `json:"servicer_max_workers_capacity"`
	ServicerWorkersIdleTimeout int    `json:"servicer_workers_idle_timeout"`

	// Servicer
	ServicerPrivateKeyFile         string `json:"servicer_private_key_file"`
	ServicerRPCTimeout             int64  `json:"servicer_rpc_timeout"`
	ServicerRPCMaxIdleConnections  int    `json:"servicer_rpc_max_idle_connections"`
	ServicerRPCMaxConnsPerHost     int    `json:"servicer_rpc_max_conns_per_host"`
	ServicerRPCMaxIdleConnsPerHost int    `json:"servicer_rpc_max_idle_conns_per_host"`
	ServicerAuthTokenFile          string `json:"servicer_auth_token_file"`
	ServicerRetryMaxTimes          int    `json:"servicer_retry_max_times"`
	ServicerRetryWaitMin           int    `json:"servicer_retry_wait_min"`
	ServicerRetryWaitMax           int    `json:"servicer_retry_wait_max"`

	// Node Health check interval (seconds)
	NodeCheckInterval int `json:"node_check_interval"`

	// Session cache (in-memory) clean up interval (seconds)
	SessionCacheCleanUpInterval int `json:"session_cache_clean_up_interval"`

	// Prometheus
	PrometheusAddr         string `json:"pocket_prometheus_port"`
	PrometheusMaxOpenfiles int    `json:"prometheus_max_open_files"`
	// Metrics uniq moniker name
	MetricsMoniker string `json:"metrics_moniker"`
	// Metrics Workers
	MetricsWorkerStrategy      string `json:"metrics_worker_strategy"`
	MetricsMaxWorkers          int    `json:"metrics_max_workers"`
	MetricsMaxWorkersCapacity  int    `json:"metrics_max_workers_capacity"`
	MetricsWorkersIdleTimeout  int    `json:"metrics_workers_idle_timeout"`
	MetricsAttachServicerLabel bool   `json:"metrics_attach_servicer_label"`
	// Metrics report interval in seconds
	MetricsReportInterval int `json:"metrics_report_interval"`
}

func defaultMeshConfig(dataDir string) MeshConfig {
	c := MeshConfig{
		// Mesh Node
		DataDir: dataDir,
		RPCPort: sdk.DefaultRPCPort,
		// following values are to be able to handle very big response from blockchains.
		ClientRPCTimeout:           120000,
		ClientRPCReadTimeout:       60000,
		ClientRPCReadHeaderTimeout: 50000,
		ClientRPCWriteTimeout:      90000,
		LogLevel:                   "*:info, *:error",
		LogChainRequest:            false,
		LogChainResponse:           false,
		UserAgent:                  sdk.DefaultUserAgent,
		AuthTokenFile:              "auth" + FS + "mesh.json",
		JSONSortRelayResponses:     sdk.DefaultJSONSortRelayResponses,
		// Chains
		ChainsName:                  sdk.DefaultChainsName,
		ChainsNameMap:               "",
		RemoteChainsNameMap:         "",
		ChainRPCTimeout:             sdk.DefaultRPCTimeout,
		ChainRPCMaxIdleConnections:  2500,
		ChainRPCMaxConnsPerHost:     2500,
		ChainRPCMaxIdleConnsPerHost: 2500,
		ChainDropConnections:        false,
		// Relay Cache
		RelayCacheFile:                         "data" + FS + "relays.pkt",
		RelayCacheBackgroundSyncInterval:       3600,
		RelayCacheBackgroundCompactionInterval: 18000,
		// Hot Reload
		KeysHotReloadInterval:   180000,
		ChainsHotReloadInterval: 180000,
		// Servicer
		ServicerPrivateKeyFile:         "key" + FS + "key.json",
		ServicerRPCTimeout:             sdk.DefaultRPCTimeout,
		ServicerRPCMaxIdleConnections:  2500,
		ServicerRPCMaxConnsPerHost:     2500,
		ServicerRPCMaxIdleConnsPerHost: 2500,
		ServicerAuthTokenFile:          "auth" + FS + "servicer.json",
		ServicerRetryMaxTimes:          10,
		ServicerRetryWaitMin:           5,
		ServicerRetryWaitMax:           180,
		ServicerWorkerStrategy:         "balanced",
		ServicerMaxWorkers:             50,
		ServicerMaxWorkersCapacity:     50000,
		ServicerWorkersIdleTimeout:     10000,
		// Node Check
		NodeCheckInterval: 60,
		// Session cache (in-memory) clean up interval (seconds)
		SessionCacheCleanUpInterval: 1800,
		// Metrics
		// Prometheus
		PrometheusAddr:             sdk.DefaultPocketPrometheusListenAddr,
		PrometheusMaxOpenfiles:     sdk.DefaultPrometheusMaxOpenFile,
		MetricsMoniker:             "geo-mesh-node",
		MetricsAttachServicerLabel: false,
		MetricsWorkerStrategy:      "lazy",
		MetricsMaxWorkers:          50,
		MetricsMaxWorkersCapacity:  50000,
		MetricsWorkersIdleTimeout:  10000,
		MetricsReportInterval:      10,
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
