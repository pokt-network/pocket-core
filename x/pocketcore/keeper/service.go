package keeper

import (
	"encoding/hex"
	"fmt"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// this is the main call for a service node handling a relay request
func (k Keeper) HandleRelay(ctx sdk.Ctx, relay pc.Relay) (*pc.RelayResponse, sdk.Error) {
	// get the latest session block height because this relay will correspond with the latest session
	sessionBlockHeight := k.GetLatestSessionBlockHeight(ctx)
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
	// get the session context
	sessionCtx, er := ctx.PrevCtx(sessionBlockHeight)
	if er != nil {
		return nil, sdk.ErrInternal(er.Error())
	}
	// ensure the validity of the relay
	if err := relay.Validate(ctx, selfNode, hostedBlockchains, sessionBlockHeight, int(k.SessionNodeCount(sessionCtx)), allNodes, app); err != nil {
		ctx.Logger().Error(fmt.Errorf("could not validate for %v, %v, %v %v, %v, %v \n", selfNode, hostedBlockchains, sessionBlockHeight, int(k.SessionNodeCount(sessionCtx)), allNodes, app).Error())
		return nil, err
	}
	// store the proof before execution, because the proof corresponds to the previous relay
	relay.Proof.Handle()
	// attempt to execute
	respPayload, err := relay.Execute(hostedBlockchains)
	if err != nil {
		return nil, err
	}
	// generate response object
	resp := &pc.RelayResponse{
		Response: respPayload,
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

func (k Keeper) HandleChallenge(ctx sdk.Ctx, challenge pc.ChallengeProofInvalidData) (*pc.ChallengeResponse, sdk.Error) {
	// get self node (your validator) from the current state
	selfNode, err := k.GetSelfNode(ctx)
	if err != nil {
		return nil, err
	}
	sessionBlkHeight := k.GetLatestSessionBlockHeight(ctx)
	// get the session context
	sessionCtx, er := ctx.PrevCtx(sessionBlkHeight)
	if er != nil {
		return nil, sdk.ErrInternal(er.Error())
	}
	// get the application that staked on behalf of the client
	app, found := k.GetAppFromPublicKey(ctx, challenge.MinorityResponse.Proof.Token.ApplicationPublicKey)
	// generate header
	header := pc.SessionHeader{
		ApplicationPubKey:  app.GetPublicKey().RawString(),
		Chain:              challenge.MinorityResponse.Proof.Blockchain,
		SessionBlockHeight: sessionCtx.BlockHeight(),
	}
	// check cache
	session, found := pc.GetSession(header)
	// if not found generate the session
	if !found {
		var err sdk.Error
		nodes := k.GetAllNodes(ctx)
		session, err = pc.NewSession(header, pc.BlockHash(sessionCtx), nodes, int(k.SessionNodeCount(sessionCtx)))
		if err != nil {
			return nil, err
		}
		// add to cache
		pc.SetSession(session)
	}
	if !found {
		return nil, pc.NewAppNotFoundError(pc.ModuleName)
	}
	// validate the challenge
	err = challenge.ValidateLocal(app.GetMaxRelays().Int64(), sessionBlkHeight, app.GetChains(), int(k.SessionNodeCount(sessionCtx)), session.SessionNodes, selfNode.GetAddress())
	if err != nil {
		return nil, err
	}
	// store the challenge in memory
	challenge.Handle()
	return nil, nil
}
