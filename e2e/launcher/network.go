package launcher

import (
	"fmt"
	"log"
	"os"
)

type Network struct {
	Nodes                []*Node
	NetworkConfiguration *NetworkConfiguration
}

func LaunchNetwork(networkConfigDirectory string, executablePath string) (network Network, err error) {
	fmt.Println("./launcher/network_configs/" + networkConfigDirectory)
	networkConfiguration, err := loadNetworkConfiguration("../launcher/network_configs/" + networkConfigDirectory)
	if err != nil {
		return
	}
	networkRootDirectory, err := os.MkdirTemp("", networkConfiguration.NetworkId+"-")
	if err != nil {
		return
	}
	log.Printf("root network directory: %v", networkRootDirectory)
	var nodes []*Node
	for _, nodeConfiguration := range networkConfiguration.NodeConfigurations {
		var node *Node
		node, err = newNode(nodeConfiguration, networkRootDirectory, networkConfiguration.GenesisPath, executablePath)
		if err != nil {
			return
		}
		node.PocketServer.Start("--datadir="+node.DataDir, "--keybase=false")
		nodes = append(nodes, node)
	}
	network = Network{
		Nodes:                nodes,
		NetworkConfiguration: networkConfiguration,
	}
	return
}
