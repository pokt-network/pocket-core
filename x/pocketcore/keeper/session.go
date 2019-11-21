package keeper

import (
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// verifies the session for a node
func (k Keeper) SessionVerification(ctx sdk.Context, nodeVerify nodeexported.ValidatorI, appPubKey string, blockchain string, blockHash string, allActiveNodes []nodeexported.ValidatorI) error {
	sess, err := pc.NewSession(appPubKey, blockchain, blockHash, allActiveNodes, int(k.SessionNodeCount(ctx)))
	if err != nil {
		return pc.NewServiceSessionGenerationError(err)
	}
	if !contains(sess.Nodes, nodeVerify) {
		return pc.InvalidSessionError
	}
	return nil
}

func contains(nodes []nodeexported.ValidatorI, node nodeexported.ValidatorI) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}
