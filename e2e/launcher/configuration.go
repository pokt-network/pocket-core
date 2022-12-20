package launcher

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type NetworkConfiguration struct {
	NetworkId          string               `json:"network_id"`
	GenesisPath        string               `json:"genesis_path"`
	NodeConfigurations []*NodeConfiguration `json:"node_configurations"`
}

type NodeConfiguration struct {
	PrivateKey string `json:"private_key"`
	ConfigPath string `json:"config_path"`
}

const networkConfigurationFileName = "network.json"

func loadNetworkConfiguration(networkConfigDirectory string) (networkConfiguration *NetworkConfiguration, err error) {
	networkConfigDirectory, err = filepath.Abs(networkConfigDirectory)
	if err != nil {
		return
	}

	var networkConfigurationFileContents []byte
	networkConfigurationFileContents, err = os.ReadFile(filepath.Join(networkConfigDirectory, networkConfigurationFileName))
	if err != nil {
		return
	}

	err = json.Unmarshal(networkConfigurationFileContents, &networkConfiguration)
	if err != nil {
		return
	}

	if !filepath.IsAbs(networkConfiguration.GenesisPath) {
		networkConfiguration.GenesisPath = filepath.Join(networkConfigDirectory, networkConfiguration.GenesisPath)
	}

	for _, nodeConfig := range networkConfiguration.NodeConfigurations {
		if !filepath.IsAbs(nodeConfig.ConfigPath) {
			nodeConfig.ConfigPath = filepath.Join(networkConfigDirectory, nodeConfig.ConfigPath)
		}
	}

	return
}
