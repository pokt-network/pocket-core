package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"time"
)

// HandleRelay handles an api (read/write) request to a non-native (external) blockchain
func (k Keeper) HandleRelay(ctx sdk.Ctx, relay pc.Relay) (*pc.RelayResponse, sdk.Error) {
	relayTimeStart := time.Now()

	sessionBlockHeight := relay.Proof.SessionBlockHeight

	if !k.IsProofSessionHeightWithinTolerance(ctx, sessionBlockHeight) {
		// For legacy support, we are intentionally returning the invalid block height error.
		return nil, pc.NewInvalidBlockHeightError(pc.ModuleName)
	}

	var node *pc.PocketNode
	// There is reference to node address so that way we don't have to recreate address twice for pre-leanpokt
	var nodeAddress sdk.Address

	if pc.GlobalPocketConfig.LeanPocket {
		// if lean pocket enabled, grab the targeted servicer through the relay proof
		servicerRelayPublicKey, err := crypto.NewPublicKey(relay.Proof.ServicerPubKey)
		if err != nil {
			return nil, sdk.ErrInternal("Could not convert servicer hex to public key")
		}
		nodeAddress = sdk.GetAddress(servicerRelayPublicKey)
		node, err = pc.GetPocketNodeByAddress(&nodeAddress)
		if err != nil {
			return nil, sdk.ErrInternal("Failed to find correct servicer PK")
		}
	} else {
		// get self node (your validator) from the current state
		node = pc.GetPocketNode()
		nodeAddress = node.GetAddress()
	}

	// retrieve the nonNative blockchains your node is hosting
	hostedBlockchains := k.GetHostedBlockchains()
	// ensure the validity of the relay
	maxPossibleRelays, err := relay.Validate(ctx, k.posKeeper, k.appKeeper, k, hostedBlockchains, sessionBlockHeight, node)
	if err != nil {
		if pc.GlobalPocketConfig.RelayErrors {
			ctx.Logger().Error(
				fmt.Sprintf("could not validate relay for app: %s for chainID: %v with error: %s",
					relay.Proof.ServicerPubKey,
					relay.Proof.Blockchain,
					err.Error(),
				),
			)
			ctx.Logger().Debug(
				fmt.Sprintf(
					"could not validate relay for app: %s, for chainID %v on node %s, at session height: %v, with error: %s",
					relay.Proof.ServicerPubKey,
					relay.Proof.Blockchain,
					nodeAddress.String(),
					sessionBlockHeight,
					err.Error(),
				),
			)
		}
		return nil, err
	}
	// move this to a worker that will insert this proof in a series style to avoid memory consumption and relay proof race conditions
	// https://github.com/pokt-network/pocket-core/issues/1457
	pc.GlobalEvidenceWorker.Submit(func() {
		// store the proof before execution, because the proof corresponds to the previous relay
		relay.Proof.Store(maxPossibleRelays, node.EvidenceStore)
	})
	// attempt to execute
	respPayload, err := relay.Execute(hostedBlockchains, &nodeAddress)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("could not send relay with error: %s", err.Error()))
		return nil, err
	}
	// generate response object
	resp := &pc.RelayResponse{
		Response: respPayload,
		Proof:    relay.Proof,
	}
	// sign the response
	sig, er := node.PrivateKey.Sign(resp.Hash())
	if er != nil {
		ctx.Logger().Error(
			fmt.Sprintf("could not sign response for address: %s with hash: %v, with error: %s",
				nodeAddress.String(), resp.HashString(), er.Error()),
		)
		return nil, pc.NewKeybaseError(pc.ModuleName, er)
	}
	// attach the signature in hex to the response
	resp.Signature = hex.EncodeToString(sig)
	// track the relay time
	relayTime := time.Since(relayTimeStart)
	// add to metrics
	addRelayMetricsFunc := func() {
		pc.GlobalServiceMetric().AddRelayTimingFor(relay.Proof.Blockchain, float64(relayTime.Milliseconds()), &nodeAddress)
		pc.GlobalServiceMetric().AddRelayFor(relay.Proof.Blockchain, &nodeAddress)
	}
	if pc.GlobalPocketConfig.LeanPocket {
		go addRelayMetricsFunc()
	} else {
		addRelayMetricsFunc()
	}
	return resp, nil
}

// "HandleChallenge" - Handles a client relay response challenge request
func (k Keeper) HandleChallenge(ctx sdk.Ctx, challenge pc.ChallengeProofInvalidData) (*pc.ChallengeResponse, sdk.Error) {

	var node *pc.PocketNode
	// There is reference to self address so that way we don't have to recreate address twice for pre-leanpokt
	var nodeAddress sdk.Address

	if pc.GlobalPocketConfig.LeanPocket {
		// try to retrieve a PocketNode that was part of session
		for _, r := range challenge.MajorityResponses {
			servicerRelayPublicKey, err := crypto.NewPublicKey(r.Proof.ServicerPubKey)
			if err != nil {
				continue
			}
			potentialNodeAddress := sdk.GetAddress(servicerRelayPublicKey)
			potentialNode, err := pc.GetPocketNodeByAddress(&nodeAddress)
			if err != nil || potentialNode == nil {
				continue
			}
			node = potentialNode
			nodeAddress = potentialNodeAddress
			break
		}
		if node == nil {
			return nil, pc.NewNodeNotInSessionError(pc.ModuleName)
		}
	} else {
		node = pc.GetPocketNode()
		nodeAddress = node.GetAddress()
	}

	sessionBlkHeight := k.GetLatestSessionBlockHeight(ctx)
	// get the session context
	sessionCtx, er := ctx.PrevCtx(sessionBlkHeight)
	if er != nil {
		return nil, sdk.ErrInternal(er.Error())
	}
	// get the application that staked on behalf of the client
	app, found := k.GetAppFromPublicKey(sessionCtx, challenge.MinorityResponse.Proof.Token.ApplicationPublicKey)
	if !found {
		return nil, pc.NewAppNotFoundError(pc.ModuleName)
	}
	// generate header
	header := pc.SessionHeader{
		ApplicationPubKey:  challenge.MinorityResponse.Proof.Token.ApplicationPublicKey,
		Chain:              challenge.MinorityResponse.Proof.Blockchain,
		SessionBlockHeight: sessionCtx.BlockHeight(),
	}
	// check cache
	session, found := pc.GetSession(header, node.SessionStore)
	// if not found generate the session
	if !found {
		var err sdk.Error
		blockHashBz, er := sessionCtx.BlockHash(k.Cdc, sessionCtx.BlockHeight())
		if er != nil {
			return nil, sdk.ErrInternal(er.Error())
		}
		session, err = pc.NewSession(sessionCtx, ctx, k.posKeeper, header, hex.EncodeToString(blockHashBz), int(k.SessionNodeCount(sessionCtx)))
		if err != nil {
			return nil, err
		}
		// add to cache
		pc.SetSession(session, node.SessionStore)
	}
	// validate the challenge
	err := challenge.ValidateLocal(header, app.GetMaxRelays(), app.GetChains(), int(k.SessionNodeCount(sessionCtx)), session.SessionNodes, nodeAddress, node.EvidenceStore)
	if err != nil {
		return nil, err
	}
	// store the challenge in memory
	challenge.Store(app.GetMaxRelays(), node.EvidenceStore)
	// update metric

	if pc.GlobalPocketConfig.LeanPocket {
		go pc.GlobalServiceMetric().AddChallengeFor(header.Chain, &nodeAddress)
	} else {
		pc.GlobalServiceMetric().AddChallengeFor(header.Chain, &nodeAddress)
	}

	return &pc.ChallengeResponse{Response: fmt.Sprintf("successfully stored challenge proof for %s", challenge.MinorityResponse.Proof.ServicerPubKey)}, nil
}
