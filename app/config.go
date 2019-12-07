package app

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto/keys"
	"github.com/tendermint/tendermint/node"
)

func GetKeybase() keys.Keybase { // todo
	return nil
}

func GetHostedChains() types.HostedBlockchains {
	return types.HostedBlockchains{} // todo
}

func GetTendermintNode() *node.Node {
	return nil // todo
}

func GetCoinbasePassphrase() string {
	return "" // todo
}

func GetGenesisFile() string {
	return "" // todo
}
