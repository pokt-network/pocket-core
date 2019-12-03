package keeper

import (
	"encoding/hex"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// this is the main call for a service node handling a relay request
func (k Keeper) HandleRelay(ctx sdk.Context, relay pc.Relay) (*pc.RelayResponse, sdk.Error) {
	// get the latest session block height because this relay will correspond with the latest session
	sessionBlockHeight := k.GetLatestSessionBlock(ctx).BlockHeight()
	// retrieve all service nodes available from world state to do session generation
	allNodes := k.GetAllNodes(ctx)
	// get self node (your validator) from world state
	selfNode, err := k.GetSelfNode(ctx)
	if err != nil {
		return nil, err
	}
	// get hosted blockchains
	hostedBlockchains := k.GetHostedBlockchains()
	// get application that staked for client
	app, found := k.GetAppFromPublicKey(ctx, relay.Proof.Token.ApplicationPublicKey)
	if !found {
		return nil, pc.NewAppNotFoundError(pc.ModuleName)
	}
	// validate the relay
	if err := relay.Validate(ctx, selfNode, hostedBlockchains, sessionBlockHeight, int(k.SessionNodeCount(ctx)), allNodes, app); err != nil {
		return nil, err
	}
	// store the previous proof
	if err := relay.HandleProof(ctx, sessionBlockHeight, int(app.GetMaxRelays().Int64())); err != nil {
		return nil, err
	}
	header := pc.PORHeader{
		ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
		Chain:              relay.Blockchain,
		SessionBlockHeight: sessionBlockHeight,
	}
	// attempt to execute
	respPayload, err := relay.Execute(hostedBlockchains)
	if err != nil {
		return nil, err
	}
	// generate response object
	resp := &pc.RelayResponse{
		Response: respPayload,
		Proof: pc.Proof{
			Index:              pc.GetAllTix().GetNextTicket(header, int(app.GetMaxRelays().Int64())),
			SessionBlockHeight: sessionBlockHeight,
			ServicerPubKey:     hex.EncodeToString(selfNode.GetConsPubKey().Bytes()),
			Token:              relay.Proof.Token,
		},
	}
	// sign the response
	sig, _, er := k.keybase.Sign(sdk.AccAddress(selfNode.GetAddress()), k.coinbasePassphrase, resp.Hash())
	if er != nil {
		return nil, pc.NewKeybaseError(pc.ModuleName, er)
	}
	resp.Signature = hex.EncodeToString(sig)
	return resp, nil
}
