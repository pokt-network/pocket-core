package mesh

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"io"
	"io/ioutil"
	"time"
)

// storeRelay - persist relay to disk
func storeRelay(relay *pocketTypes.Relay) {
	hash := relay.RequestHash()
	logger.Debug(fmt.Sprintf("storing relay %s", relay.RequestHashString()))
	rb, err := json.Marshal(relay)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = relaysCacheDb.Put(hash, rb)
	if err != nil {
		logger.Error(fmt.Sprintf("error adding relay %s to cache. %s", relay.RequestHashString(), err.Error()))
	}

	return
}

// decodeCacheRelay - decode []byte relay from cache to pocketTypes.Relay
func decodeCacheRelay(body []byte) (relay *pocketTypes.Relay) {
	if err := json.Unmarshal(body, &relay); err != nil {
		logger.Error("error decoding cache relay")
		// todo: delete key from cache?
		return nil
	}
	return
}

// deleteCacheRelay - delete a key from relay cache
func deleteCacheRelay(relay *pocketTypes.Relay) {
	hash := relay.RequestHash()
	err := relaysCacheDb.Delete(hash)
	if err != nil {
		logger.Error("error deleting relay from cache %s", hex.EncodeToString(hash))
		return
	}

	return
}

// evaluateServicerError - this will change internalCache[hash].IsValid bool depending on the result of the evaluation
func evaluateServicerError(r *pocketTypes.Relay, err *SdkErrorResponse) (isSessionStillValid bool) {
	isSessionStillValid = !ShouldInvalidateSession(err.Code) // we should not retry if is invalid

	if isSessionStillValid {
		return isSessionStillValid
	}

	session, e1 := sessionStorage.GetSession(r)

	if e1 != nil {
		servicerAddress, _ := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)
		logger.Error(
			fmt.Sprintf(
				"failure getting session from storage app=%s chain=%s servicer=%s err=%s",
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
				e1.Error,
			),
		)
		// todo: should we invalidate here if we are not sure what is going on?
		return
	}

	session.InvalidateNodeSession(r.Proof.ServicerPubKey, err)

	return
}

// notifyServicer - call servicer to ack about the processed relay.
func notifyServicer(r *pocketTypes.Relay) {
	relayTimeStart := time.Now()
	// discard this relay at the end of this function, to end this function the servicer will be retried N times
	defer deleteCacheRelay(r)

	result := RPCRelayResponse{}
	ctx := context.WithValue(context.Background(), "result", &result)
	jsonData, err := json.Marshal(r)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error encoding relay hash=%s err=%s",
				r.RequestHashString(),
				err.Error(),
			),
		)
		return
	}

	servicerAddress, err := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"unable to decode servicer publicKey=%s to address from relay of app=%s chain=%s sessionHeight=%d",
				r.Proof.ServicerPubKey,
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				r.Proof.SessionBlockHeight,
			),
		)
		return
	}

	logger.Debug(
		fmt.Sprintf(
			"delivering relay notification of app=%s chain=%s sessionHeight=%d servicer=%s",
			r.Proof.Token.ApplicationPublicKey,
			r.Proof.Blockchain,
			r.Proof.SessionBlockHeight,
			servicerAddress,
		),
	)

	servicerNode, ok := servicerMap.Load(servicerAddress)

	if !ok {
		logger.Error(
			fmt.Sprintf(
				"unable to find servicer=%s to notify relay of app=%s chain=%s sessionHeight=%d",
				servicerAddress,
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				r.Proof.SessionBlockHeight,
			),
		)
		return
	}

	nodeSession, e := sessionStorage.GetNodeSession(r)

	if e != nil || (nodeSession.Validated && !nodeSession.IsValid) {
		if e != nil {
			// just to have a different text that provide better understanding of what is going on on the logs.
			logger.Error(
				fmt.Sprintf(
					"relay for app=%s chain=%s servicer=%s was not able to be delivered because we encounter an error getting session from session storage error=%s",
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
					servicerAddress,
					e.Error,
				),
			)
			return
		}

		logger.Error(
			fmt.Sprintf(
				"relay for app=%s chain=%s servicer=%s was not able to be delivered because session is already invalidated",
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
			),
		)
		return
	} else if !nodeSession.Validated {
		// we should re-schedule it to later on the session is validated, or even validated but invalid to be discarded.
		servicerNode.Node.Worker.Submit(func() {
			notifyServicer(r)
		})
		return
	}

	requestURL := fmt.Sprintf(
		"%s%s?chain=%s&app=%s",
		servicerNode.Node.URL,
		ServicerRelayEndpoint,
		r.Proof.Blockchain,
		r.Proof.Token.ApplicationPublicKey,
	)
	req, err := retryablehttp.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonData))
	req.Header.Set(AuthorizationHeader, servicerAuthToken.Value)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error formatting url to notify servicer=%s for relay of app=%s chain=%s sessionHeight=%d at url=%s with err=%s",
				servicerAddress,
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				r.Proof.SessionBlockHeight,
				requestURL,
				err.Error(),
			),
		)
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, "500")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ServicerHeader, servicerAddress)
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := relaysClient.Do(req)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error dispatching relay to app=%s chain=%s servicer=%s err=%s",
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
				err.Error(),
			),
		)
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, "500")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	// read the body just to allow http 1.x be able to reuse the connection
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"couldn't parse response body app=%s chain=%s servicer=%s err=%s",
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
				err.Error(),
			),
		)
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, fmt.Sprintf("%d", resp.StatusCode))
		return
	}

	isSuccess := resp.StatusCode == 200

	if result.Dispatch != nil && result.Dispatch.BlockHeight > servicerNode.Node.Status.Height {
		servicerNode.Node.Status.Height = result.Dispatch.BlockHeight
	}

	if !isSuccess {
		logger.Debug(
			fmt.Sprintf(
				"servicer %s reject relay %s\n: CODE=%d\nCODESPACE=%s\nMESSAGE=%s",
				servicerAddress, r.RequestHashString(),
				result.Error.Code, result.Error.Codespace, result.Error.Error,
			),
		)
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, fmt.Sprintf("%d", resp.StatusCode))
		evaluateServicerError(r, result.Error)
	} else {
		logger.Debug(fmt.Sprintf("servicer processed relay %s successfully", r.RequestHashString()))

		ns, e1 := sessionStorage.GetNodeSession(r)

		if e1 != nil {
			logger.Error(
				fmt.Sprintf(
					"error getting session from storage for app=%s chain=%s servicer=%s err=%s",
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
					servicerAddress,
					e1.Error,
				),
			)
			servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, "500")
			return
		}

		exhausted := ns.CountRelay()

		if exhausted {
			logger.Debug(
				fmt.Sprintf(
					"servicer %s exhaust relays for app %s at blockchain %s",
					servicerNode.Address.String(),
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
				),
			)
		} else {
			logger.Debug(
				fmt.Sprintf(
					"servicer %s has %d remaining relays to process for app %s at blockchain %s",
					servicerNode.Address.String(),
					ns.RemainingRelays,
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
				),
			)
		}

		// track the notify relay time
		relayDuration := time.Since(relayTimeStart)
		servicerNode.Node.MetricsWorker.AddServiceMetricRelayFor(r, &servicerNode.Address, relayDuration, true)
	}

	return
}

// execute - Attempts to do a request on the non-native blockchain specified
func execute(r *pocketTypes.Relay, hostedBlockchains *pocketTypes.HostedBlockchains, address *sdk.Address) (string, sdk.Error) {
	node := GetNodeFromAddress(address.String())

	if node == nil {
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ServicerNotFoundStatusType, "500")
		return "", sdk.ErrInternal("Failed to find correct servicer PK")
	}

	// retrieve the hosted blockchain url requested
	chain, err := hostedBlockchains.GetChain(r.Proof.Blockchain)
	if err != nil {
		// metric track
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ChainNotFoundStatusType, "500")
		return "", err
	}

	// do basic http request on the relay
	res, er, code := ExecuteBlockchainHTTPRequest(r.Payload, chain)
	if er != nil {
		// metric track
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ChainStatusType, fmt.Sprintf("%d", code))
		return res, pocketTypes.NewHTTPExecutionError(ModuleName, er)
	}

	if code >= 400 {
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ChainStatusType, fmt.Sprintf("%d", code))
	}

	return res, nil
}

// processRelay - call execute and create RelayResponse or Error in case. Also trigger relay metrics.
func processRelay(relay *pocketTypes.Relay) (*pocketTypes.RelayResponse, sdk.Error) {
	relayTimeStart := time.Now()
	logger.Debug(fmt.Sprintf("processing relay %s", relay.RequestHashString()))

	servicerAddress, e := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)

	if e != nil {
		return nil, sdk.ErrInternal("could not convert servicer hex to public key")
	}

	servicerNode, ok := servicerMap.Load(servicerAddress)
	if !ok {
		return nil, sdk.ErrInternal("failed to find correct servicer PK")
	}

	// attempt to execute
	respPayload, err := execute(relay, chains, &servicerNode.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not send relay %s with error: %s", relay.RequestHashString(), err.Error()))
		return nil, err
	}
	// generate response object
	resp := &pocketTypes.RelayResponse{
		Response: respPayload,
		Proof:    relay.Proof,
	}

	// sign the response
	sig, er := servicerNode.PrivateKey.Sign(resp.Hash())
	if er != nil {
		logger.Error(
			fmt.Sprintf("could not sign response for address: %s with hash: %v, with error: %s",
				servicerAddress, resp.HashString(), er.Error()),
		)
		return nil, pocketTypes.NewKeybaseError(pocketTypes.ModuleName, er)
	}
	// attach the signature in hex to the response
	resp.Signature = hex.EncodeToString(sig)
	// track the relay time
	relayTimeDuration := time.Since(relayTimeStart)
	// queue metric of relay
	servicerNode.Node.MetricsWorker.AddServiceMetricRelayFor(relay, &servicerNode.Address, relayTimeDuration, false)
	return resp, nil
}

// validate - evaluate relay to understand if should or not processed.
func validate(r *pocketTypes.Relay) sdk.Error {
	logger.Debug(fmt.Sprintf("validating relay %s", r.RequestHashString()))
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return pocketTypes.NewEmptyPayloadDataError(ModuleName)
	}
	// validate appPubKey
	if err := pocketTypes.PubKeyVerification(r.Proof.Token.ApplicationPublicKey); err != nil {
		return err
	}
	// validate chain
	if err := pocketTypes.NetworkIdentifierVerification(r.Proof.Blockchain); err != nil {
		return pocketTypes.NewEmptyChainError(ModuleName)
	}
	// validate the relay merkleHash = request merkleHash
	if r.Proof.RequestHash != r.RequestHashString() {
		return pocketTypes.NewRequestHashError(ModuleName)
	}
	// validate servicer public key
	servicerAddress, e := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)
	if e != nil {
		return sdk.ErrInternal("could not convert servicer hex to public key")
	}
	// load servicer from servicer map, if not there maybe the servicer is pk is not loaded
	if _, ok := servicerMap.Load(servicerAddress); !ok {
		return pocketTypes.NewInvalidSessionError(ModuleName)
	}

	session, e1 := sessionStorage.GetSession(r)

	if e1 != nil {
		logger.Error(
			fmt.Sprintf(
				"failure getting session information for app=%s chain=%s servicer=%s err=%s",
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
				e1.Error,
			),
		)
		return pocketTypes.NewInvalidSessionKeyError(ModuleName, errors.New(e1.Error))
	}

	if session != nil {
		nodeSession, e2 := session.GetNodeSessionByPubKey(r.Proof.ServicerPubKey)

		if e2 != nil {
			logger.Error(
				fmt.Sprintf(
					"failure getting session information for app=%s chain=%s servicer=%s err=%s",
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
					servicerAddress,
					e2.Error(),
				),
			)
			return pocketTypes.NewInvalidSessionKeyError(ModuleName, e2)
		}

		if !nodeSession.IsValid {
			if nodeSession.Error != nil {
				return NewPocketSdkErrorFromSdkError(nodeSession.Error)
			} else {
				err := errors.New(fmt.Sprintf(
					"invalid session for app=%s chain=%s servicer=%s",
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
					servicerAddress,
				))
				logger.Error(err.Error())
				return pocketTypes.NewInvalidSessionKeyError(ModuleName, err)
			}
		}
	} else {
		err := errors.New(fmt.Sprintf(
			"session not found for app=%s chain=%s servicer=%s ",
			r.Proof.Token.ApplicationPublicKey,
			r.Proof.Blockchain,
			servicerAddress,
		))
		logger.Error(err.Error())
		return pocketTypes.NewInvalidSessionKeyError(ModuleName, err)
	}

	// is needed we call the node and validate if there is not a validation already in place get done by the cron?
	return nil
}

// HandleRelay - evaluate node status, validate relay payload and call processRelay
func HandleRelay(r *pocketTypes.Relay) (res *pocketTypes.RelayResponse, dispatch *DispatchResponse, err error) {
	servicerAddress, e := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)

	if e != nil {
		return nil, nil, errors.New("could not convert servicer hex to public key")
	}

	servicerNode, ok := servicerMap.Load(servicerAddress)
	if !ok {
		return nil, nil, errors.New("failed to find correct servicer PK")
	}

	if servicerNode.Node.Status == nil {
		return nil, nil, fmt.Errorf("pocket node is currently unavailable")
	}

	if servicerNode.Node.Status.IsStarting {
		return nil, nil, fmt.Errorf("pocket node is unable to retrieve synced status from tendermint node, cannot service in this state")
	}

	if servicerNode.Node.Status.IsCatchingUp {
		return nil, nil, fmt.Errorf("pocket node is currently syncing to the blockchain, cannot service in this state")
	}

	err = validate(r)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"could not validate relay %s for app: %s, for chainID %v on node %s, at session height: %v, with error: %s",
				r.RequestHashString(),
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
				r.Meta.BlockHeight,
				err.Error(),
			),
		)

		return
	}

	// store relay on cache; once we hit this point this relay will be processed so should be notified to servicer even
	// if process is shutdown
	storeRelay(r)

	res, err = processRelay(r)

	if err != nil && pocketTypes.ErrorWarrantsDispatch(err) {

		session, e1 := sessionStorage.GetSession(r)

		if e1 != nil {
			logger.Error(
				fmt.Sprintf(
					"error getting session for app=%s chain=%s servicer=%s sessionHeight=%d",
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
					servicerAddress,
					r.Proof.SessionBlockHeight,
				),
			)
		} else {
			dispatch = session.Dispatch
		}
	}

	// add to task group pool
	if servicerNode.Node.Worker.Stopped() {
		// this should not happen, but just in case avoid a panic here.
		logger.Error(fmt.Sprintf("worker of node=%s was already stopped", servicerNode.Node.URL))
		return
	}

	// 50k - 50
	servicerNode.Node.Worker.Submit(func() {
		notifyServicer(r)
	})

	return
}
