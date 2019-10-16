package blockchain

import (
	"github.com/pokt-network/pocket-core/tests/fixtures"
	"github.com/pokt-network/pocket-core/types"
)

func GetLatestBlock() types.BlockID {
	// todo
	return fixtures.GenerateBlockHash()
}

func GetLatestSessionBlock() types.BlockID {
	// todo
	return fixtures.GenerateBlockHash()
}

func GetNodes() (*types.Nodes, error) { // this is essentially -> dispatchPeers()
	// todo
	return fixtures.GetNodes()
}

func GetApplications() (*types.Applications, error) {
	// todo
	return fixtures.GetApplications()
}

func GetMaxNumberOfRelaysForApp(applicationPubKey string) int {
	// todo
	return 5000
}
