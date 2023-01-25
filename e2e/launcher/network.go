package launcher

import (
	"fmt"
	"log"
	"os"
)

type Network struct {
	Nodes         []*Node
	NetworkConfig *NetworkConfiguration
}

const (
	networkConfigsLauncherPath = "./launcher/network_configs/"
)

func LaunchNetwork(networkConfigDir string, executablePath string) (network *Network, err error) {
	log.Printf("loading network from config: ./%s%s\n", networkConfigsLauncherPath, networkConfigDir)

	networkConfig, err := loadNetworkConfiguration(fmt.Sprintf("../%s/%s", networkConfigsLauncherPath, networkConfigDir))
	if err != nil {
		return
	}

	networkRootDir, err := os.MkdirTemp("", networkConfig.NetworkId+"-")
	if err != nil {
		return
	}
	log.Printf("root network directory: %v", networkRootDir)

	var nodes []*Node
	for _, nodeConfig := range networkConfig.NodeConfigurations {
		var node *Node
		node, err = newNode(nodeConfig, networkRootDir, networkConfig.GenesisPath, executablePath)
		if err != nil {
			return
		}

		err = node.PocketServer.Start("--datadir="+node.DataDir, "--keybase=false")
		if err != nil {
			return
		}

		nodes = append(nodes, node)
	}

	network = &Network{
		Nodes:         nodes,
		NetworkConfig: networkConfig,
	}

	return
}
