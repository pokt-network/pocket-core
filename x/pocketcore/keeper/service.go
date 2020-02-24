package keeper

import (
	"fmt"
	"encoding/hex"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// this is the main call for a service node handling a relay request
func (k Keeper) HandleRelay(ctx sdk.Context, relay pc.Relay) (*pc.RelayResponse, sdk.Error) {
	ctx.Logger().Info(fmt.Sprintf("HandleRelay(Relay = %+v)", relay))
	// get the latest session block height because this relay will correspond with the latest session
	sessionBlockHeight := k.GetLatestSessionBlock(ctx).BlockHeight()
	// retrieve all service nodes available from world state to do session generation (the session data is needed to service)
	allNodes := k.GetAllNodes(ctx)
	// get self node (your validator) from the current state
	selfNode, err := k.GetSelfNode(ctx)
	if err != nil {
		return nil, err
	}
	// retrieve the nonNative blockchains your node is hosting
	hostedBlockchains := k.GetHostedBlockchains()
	// get the application that staked on behalf of the client
	app, found := k.GetAppFromPublicKey(ctx, relay.Proof.Token.ApplicationPublicKey)
	if !found {
		return nil, pc.NewAppNotFoundError(pc.ModuleName)
	}
	// ensure the validity of the relay
	if err := relay.Validate(ctx, selfNode, hostedBlockchains, sessionBlockHeight, int(k.SessionNodeCount(ctx)), allNodes, app); err != nil {
		ctx.Logger().Error(fmt.Errorf("could not validate for %v, %v, %v %v, %v, %v \n", selfNode, hostedBlockchains, sessionBlockHeight, int(k.SessionNodeCount(ctx)), allNodes, app).Error())
		return nil, err
	}
	// store the proof before execution, because the proof corresponds to the previous relay
	if err := relay.HandleProof(ctx, sessionBlockHeight); err != nil {
		ctx.Logger().Error(fmt.Errorf("could not handle proof for %v", sessionBlockHeight).Error())
		return nil, err
	}
	// attempt to execute
	respPayload, err := relay.Execute(hostedBlockchains)
	if err != nil {
		return nil, err
	}
	// generate response object
	resp := &pc.RelayResponse{
		RequestHash: relay.HashString(),
		Response:    respPayload,
		Proof: pc.RelayProof{
			Blockchain:         relay.Proof.Blockchain,
			SessionBlockHeight: sessionBlockHeight,
			ServicerPubKey:     selfNode.GetPublicKey().RawString(),
			Token:              relay.Proof.Token,
		},
	}
	// sign the response
	sig, _, er := (k.Keybase).Sign(selfNode.GetAddress(), k.coinbasePassphrase, resp.Hash())
	if er != nil {
		ctx.Logger().Error(fmt.Errorf("could not sign response for address: %v with hash: %v \n", selfNode.GetAddress().String(), resp.Hash()).Error())
		return nil, pc.NewKeybaseError(pc.ModuleName, er)
	}
	resp.Signature = hex.EncodeToString(sig)
	return resp, nil
}
