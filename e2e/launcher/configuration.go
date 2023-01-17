package launcher

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type NodeConfiguration struct {
	PrivateKey string `json:"private_key"`
	ConfigPath string `json:"config_path"`
}

type NetworkConfiguration struct {
	NetworkId          string               `json:"network_id"`
	GenesisPath        string               `json:"genesis_path"`
	NodeConfigurations []*NodeConfiguration `json:"node_configurations"`
}

const defaultNetworkConfigFilename = "network.json"

func loadNetworkConfiguration(networkConfigDir string) (networkConfig *NetworkConfiguration, err error) {
	networkConfigDir, err = filepath.Abs(networkConfigDir)
	if err != nil {
		return
	}

	var networkConfigFileContents []byte
	networkConfigFileContents, err = os.ReadFile(filepath.Join(networkConfigDir, defaultNetworkConfigFilename))
	if err != nil {
		return
	}

	err = json.Unmarshal(networkConfigFileContents, &networkConfig)
	if err != nil {
		return
	}

	if !filepath.IsAbs(networkConfig.GenesisPath) {
		networkConfig.GenesisPath = filepath.Join(networkConfigDir, networkConfig.GenesisPath)
	}

	for _, nodeConfig := range networkConfig.NodeConfigurations {
		if !filepath.IsAbs(nodeConfig.ConfigPath) {
			nodeConfig.ConfigPath = filepath.Join(networkConfigDir, nodeConfig.ConfigPath)
		}
	}

	return
}
