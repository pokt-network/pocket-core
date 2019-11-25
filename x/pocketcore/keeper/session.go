package keeper

import (
	"encoding/hex"
	"errors"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// verifies the session for a node
func (k Keeper) SessionVerification(ctx sdk.Context, nodeVerify nodeexported.ValidatorI, app appexported.ApplicationI, blockchain string, blockHash string, blockHeight int64, allActiveNodes []nodeexported.ValidatorI) (session *pc.Session, er error) {
	// validate that app staked for chain
	if _, found := app.GetChains()[blockchain]; !found { // todo is it possible for application information to change while this is being called?
		return nil, errors.New("app not staked for chain")
	}

	sess, err := pc.NewSession(hex.EncodeToString(app.GetConsPubKey().Bytes()), blockchain, blockHash, blockHeight, allActiveNodes, int(k.SessionNodeCount(ctx)))
	if err != nil {
		return sess, pc.NewSessionGenerationError(err)
	}
	if !contains(sess.Nodes, nodeVerify) {
		return sess, pc.InvalidSessionError
	}
	return sess, nil
}

func (k Keeper) IsSessionBlock(ctx sdk.Context) bool {
	frequency := k.posKeeper.SessionBlockFrequency(ctx)
	return ctx.BlockHeight()%frequency == 1
}

func contains(nodes []nodeexported.ValidatorI, node nodeexported.ValidatorI) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

func (k Keeper) GetLatestSessionBlock(ctx sdk.Context) sdk.Context{
	sessionBlockHeight := (ctx.BlockHeight() % int64(k.posKeeper.SessionBlockFrequency(ctx))) + 1
	return ctx.WithBlockHeight(sessionBlockHeight)
}
